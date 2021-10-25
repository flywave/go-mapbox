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
	"strconv"
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

var formatStrings = [...]string{
	"",
	"pbf",
	"png",
	"jpg",
	"webp",
	"gzib",
	"zlib",
}

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

type LayerType int

const (
	BaseLayer LayerType = iota
	Overlay
)

var layerTypeStrings = [...]string{
	"baselayer",
	"overlay",
}

func (t LayerType) String() string {
	return layerTypeStrings[t]
}

func (t *LayerType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *LayerType) UnmarshalJSON(b []byte) error {
	var lt string

	if err := json.Unmarshal(b, &lt); err != nil {
		return err
	}

	for i, k := range layerTypeStrings {
		if k == lt {
			*t = LayerType(i)
			return nil
		}
	}

	return errors.New("invalid or unknown tile format")
}

func stringToLayerType(s string) LayerType {
	for i, k := range layerTypeStrings {
		if k == s {
			return LayerType(i)
		}
	}

	return BaseLayer
}

func layerTypeToString(s LayerType) string {
	switch s {
	case BaseLayer:
		return "baselayer"
	case Overlay:
		return "overlay"
	}
	return ""
}

type DB struct {
	db                 *sql.DB
	tileformat         TileFormat
	timestamp          time.Time
	hasUTFGrid         bool
	utfgridCompression TileFormat
}

func NewDB(filename string) (*DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	fileStat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file stats for mbtiles file: %s", filename)
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

func CreateDB(filename string, format TileFormat, md *Metadata) (*DB, error) {
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

	fileStat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file stats for mbtiles file: %s", filename)
	}

	out := DB{
		db:         db,
		tileformat: format,
		timestamp:  fileStat.ModTime().Round(time.Second),
	}

	if md != nil {
		out.UpdateMetadata(md)
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
		return errors.New("tileset does not contain UTFgrids")
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

func (ts *DB) GetMetadata() (*Metadata, error) {
	md := &Metadata{}

	rows, err := ts.db.Query("SELECT name, value FROM metadata WHERE value is not ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var key, value string
	for rows.Next() {
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}

		switch key {
		case "name":
			md.Name = value
		case "description":
			md.Description = value
		case "attribution":
			md.Attribution = value
		case "version":
			md.Version = value
		case "format":
			md.Format = stringToTileFormat(value)
		case "minzoom":
			md.MinZoom, err = strconv.Atoi(value)
		case "maxzoom":
			md.MaxZoom, err = strconv.Atoi(value)
		case "center":
			md.Center, err = stringToCenter(value)
		case "bounds":
			md.Bounds, err = stringToBounds(value)
		case "type":
			md.Type = stringToLayerType(value)
		case "json":
			err = json.Unmarshal([]byte(value), &md.LayerData)
		case "directory_layout":
			md.DirectoryLayout = value
		case "origin":
			md.Origin = value
		case "srs":
			md.Srs = value
		case "bounds_srs":
			md.BoundsSrs = value
		case "res_factor":
			md.ResFactor = stringToResFactor(value)
		case "tile_size":
			md.TileSize, err = stringToTileSize(value)
		}

		if err != nil {
			return nil, err
		}
	}

	if md.MaxZoom == 0 {
		var min, max string
		if err := ts.db.QueryRow("SELECT min(zoom_level), max(zoom_level) FROM tiles").Scan(&min, &max); err != nil {
			return nil, err
		}

		md.MinZoom, err = strconv.Atoi(min)
		if err != nil {
			return nil, err
		}

		md.MaxZoom, err = strconv.Atoi(max)
		if err != nil {
			return nil, err
		}
	}

	return md, nil
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

	for k, v := range values {
		_, err = stmt.Exec(k, v)
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

	return UNKNOWN, errors.New("could not detect tile format")
}
