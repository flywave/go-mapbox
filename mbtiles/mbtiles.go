package mbtiles

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TileFormat uint8

const (
	UNKNOWN TileFormat = iota
	GZIP
	ZLIB
	PNG
	JPG
	PBF
	WEBP
)

func (t TileFormat) String() string {
	switch t {
	case PNG:
		return "png"
	case JPG:
		return "jpg"
	case PBF:
		return "pbf"
	case WEBP:
		return "webp"
	default:
		return ""
	}
}

func (t TileFormat) ContentType() string {
	switch t {
	case PNG:
		return "image/png"
	case JPG:
		return "image/jpeg"
	case PBF:
		return "application/x-protobuf"
	case WEBP:
		return "image/webp"
	default:
		return ""
	}
}

type DB struct {
	filename           string
	db                 *sql.DB
	tileformat         TileFormat
	timestamp          time.Time
	hasUTFGrid         bool
	utfgridCompression TileFormat
}

func NewDB(filename string) (*DB, error) {
	_, id := filepath.Split(filename)
	id = strings.Split(id, ".")[0]

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	fileStat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file stats for mbtiles file: %s\n", filename)
	}

	var data []byte
	err = db.QueryRow("select tile_data from tiles limit 1").Scan(&data)
	if err != nil {
		return nil, err
	}
	tileformat, err := detectTileFormat(&data)
	if err != nil {
		return nil, err
	}
	if tileformat == GZIP {
		tileformat = PBF
	}
	out := DB{
		db:         db,
		tileformat: tileformat,
		timestamp:  fileStat.ModTime().Round(time.Second),
	}

	// UTFGrids
	var count int
	err = db.QueryRow("select count(*) from sqlite_master where name in ('grids', 'grid_data', 'grid_utfgrid', 'keymap', 'grid_key')").Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 5 {
		var gridData []byte
		err = db.QueryRow("select grid_utfgrid from grid_utfgrid limit 1").Scan(&gridData)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("could not read sample grid to determine type: %v", err)
			}
		} else {
			out.hasUTFGrid = true
			out.utfgridCompression, err = detectTileFormat(&gridData)
			if err != nil {
				return nil, fmt.Errorf("could not determine UTF Grid compression type: %v", err)
			}
		}
	}

	return &out, nil

}

func CreateDB(filename string, format TileFormat, description string, tilejson string) (*DB, error) {
	_, id := filepath.Split(filename)
	id = strings.Split(id, ".")[0]

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE tiles (zoom_level integer, tile_column integer, tile_row integer, tile_data blob);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = `CREATE UNIQUE INDEX idx_tile on tiles
	(zoom_level, tile_column, tile_row);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = `
	CREATE TABLE metadata (name text, value text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	values := [][]string{{"name", filename},
		{"type", "overlay"},
		{"version", "2"},
		{"description", description},
		{"format", format.String()},
		{"json", tilejson},
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("insert into metadata(value, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range values {
		_, err = stmt.Exec(i[1], i[0])
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()

	fileStat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file stats for mbtiles file: %s\n", filename)
	}

	out := DB{
		db:         db,
		tileformat: format,
		timestamp:  fileStat.ModTime().Round(time.Second),
	}

	return &out, nil
}

func (tileset *DB) StoreTile(z uint8, x uint64, y uint64, data []byte) error {
	stmt := "INSERT OR REPLACE INTO tiles (zoom_level, tile_column, tile_row, tile_data) VALUES (?,?,?,?);"

	_, err := tileset.db.Exec(stmt, z, x, y, data)
	if err != nil {
		return err
	}

	return nil
}

func (tileset *DB) ReadTile(z uint8, x uint64, y uint64, data *[]byte) error {
	err := tileset.db.QueryRow("select tile_data from tiles where zoom_level = ? and tile_column = ? and tile_row = ?", z, x, y).Scan(data)
	if err != nil {
		if err == sql.ErrNoRows {
			*data = nil
			return nil
		}
		return err
	}
	return nil
}

func (tileset *DB) ReadGrid(z uint8, x uint64, y uint64, data *[]byte) error {
	if !tileset.hasUTFGrid {
		return errors.New("Tileset does not contain UTFgrids")
	}

	err := tileset.db.QueryRow("select grid from grids where zoom_level = ? and tile_column = ? and tile_row = ?", z, x, y).Scan(data)
	if err != nil {
		if err == sql.ErrNoRows {
			*data = nil
			return nil
		}
		return err
	}

	keydata := make(map[string]interface{})
	var (
		key   string
		value []byte
	)

	rows, err := tileset.db.Query("select key_name, key_json FROM grid_data where zoom_level = ? and tile_column = ? and tile_row = ?", z, x, y)
	if err != nil {
		return fmt.Errorf("cannot fetch grid data: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			return fmt.Errorf("could not fetch grid data: %v", err)
		}
		valuejson := make(map[string]interface{})
		json.Unmarshal(value, &valuejson)
		keydata[key] = valuejson
	}

	if len(keydata) == 0 {
		return nil
	}

	var (
		zreader io.ReadCloser
		zwriter io.WriteCloser
		buf     bytes.Buffer
	)
	reader := bytes.NewReader(*data)

	if tileset.utfgridCompression == ZLIB {
		zreader, err = zlib.NewReader(reader)
		if err != nil {
			return err
		}
		zwriter = zlib.NewWriter(&buf)
	} else {
		zreader, err = gzip.NewReader(reader)
		if err != nil {
			return err
		}
		zwriter = gzip.NewWriter(&buf)
	}

	var utfjson map[string]interface{}
	jsonDecoder := json.NewDecoder(zreader)
	jsonDecoder.Decode(&utfjson)
	zreader.Close()

	utfjson["data"] = keydata
	if err != nil {
		return err
	}

	jsonEncoder := json.NewEncoder(zwriter)
	err = jsonEncoder.Encode(utfjson)
	if err != nil {
		return err
	}
	zwriter.Close()
	*data = buf.Bytes()

	return nil
}

func (tileset *DB) ReadMetadata() (*Metadata, error) {
	var (
		key   string
		value string
	)
	metadata := make(map[string]string)

	rows, err := tileset.db.Query("select * from metadata where value is not ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&key, &value)
		metadata[key] = value
	}

	_, hasMinZoom := metadata["minzoom"]
	_, hasMaxZoom := metadata["maxzoom"]
	if !(hasMinZoom && hasMaxZoom) {
		var minZoom, maxZoom int
		err := tileset.db.QueryRow("select min(zoom_level), max(zoom_level) from tiles").Scan(&minZoom, &maxZoom)
		if err != nil {
			return nil, err
		}
		metadata["minzoom"] = strconv.Itoa(minZoom)
		metadata["maxzoom"] = strconv.Itoa(maxZoom)
	}
	return NewMetadata(metadata), nil
}

func (tileset *DB) UpdateMetadata(md *Metadata) error {
	sqlStmt := `
	DELETE FROM metadata;
	`
	_, err := tileset.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	tx, err := tileset.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into metadata(value, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	values := md.ToMap()

	for _, i := range values {
		_, err = stmt.Exec(i[1], i[0])
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (tileset *DB) TileFormat() TileFormat {
	return tileset.tileformat
}

func (tileset *DB) TileFormatString() string {
	return tileset.tileformat.String()
}

func (tileset *DB) ContentType() string {
	return tileset.tileformat.ContentType()
}

func (tileset *DB) HasUTFGrid() bool {
	return tileset.hasUTFGrid
}

func (tileset *DB) UTFGridCompression() TileFormat {
	return tileset.utfgridCompression
}

func (tileset *DB) TimeStamp() time.Time {
	return tileset.timestamp
}

func (tileset *DB) Close() error {
	return tileset.db.Close()
}

func detectTileFormat(data *[]byte) (TileFormat, error) {
	patterns := map[TileFormat][]byte{
		GZIP: []byte("\x1f\x8b"),
		ZLIB: []byte("\x78\x9c"),
		PNG:  []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"),
		JPG:  []byte("\xFF\xD8\xFF"),
		WEBP: []byte("\x52\x49\x46\x46\xc0\x00\x00\x00\x57\x45\x42\x50\x56\x50"),
	}

	for format, pattern := range patterns {
		if bytes.HasPrefix(*data, pattern) {
			return format, nil
		}
	}

	return UNKNOWN, errors.New("Could not detect tile format")
}

func stringToFloats(str string) ([]float64, error) {
	split := strings.Split(str, ",")
	var out []float64
	for _, v := range split {
		value, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return out, fmt.Errorf("could not parse %q to floats: %v", str, err)
		}
		out = append(out, value)
	}
	return out, nil
}
