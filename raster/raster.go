package raster

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	rt "github.com/flywave/go-raster"
)

const (
	DEM_ENCODING_MAPBOX    = 0
	DEM_ENCODING_TERRARIUM = 1
)

var (
	UNPACK_MAPBOX    = [4]float64{6553.6, 25.6, 0.1, 10000.0}
	UNPACK_TERRARIUM = [4]float64{256.0, 1.0, 1.0 / 256.0, 32768.0}
)

type DEMData struct {
	Encoding int
	Dim      int
	Stride   int
	Data     [][4]byte
}

func LoadDEMData(path string, encoding int) (*DEMData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	rect := m.Bounds()
	if m.ColorModel() != color.NRGBAModel || rect.Dx() != rect.Dy() {
		return nil, errors.New("image format error!")
	}

	data := make([][4]byte, rect.Dx()*rect.Dy())
	for y := 0; y < rect.Dy(); y++ {
	for x := 0; x < rect.Dx(); x++ {
			rgba := m.At(x, y).(color.NRGBA)
			data[y*rect.Dx()+x] = [4]byte{rgba.R, rgba.G, rgba.B, rgba.A}
		}
	}
	return NewDEMData(data, encoding), nil
}

func NewDEMData(data [][4]byte, encoding int) *DEMData {
	if len(data)%2 != 0 {
		return nil
	}
	dim := int(math.Sqrt(float64(len(data))))
	stride := dim + 2
	img := make([][4]byte, stride*stride)
	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			img[(r+1)*stride+c+1] = data[r*dim+c]
		}
	}

	for x := 0; x < dim; x++ {
		rowOffset := stride * (x + 1)
		img[rowOffset] = img[rowOffset+1]
		img[rowOffset+dim+1] = img[rowOffset+dim]
	}

	for r := 0; r < stride; r++ {
		img[r] = img[stride+r]
		img[(stride*(dim+1))+r] = img[stride*dim+r]
	}

	return &DEMData{Encoding: encoding, Dim: dim, Stride: stride, Data: img}
}

func (d *DEMData) idx(x int, y int) int {
	return (y+1)*d.Stride + (x + 1)
}

func (d *DEMData) BackfillBorder(data DEMData, dx int, dy int) {
	if d.Dim == data.Dim {
		return
	}
	xMin := dx * d.Dim
	xMax := dx*d.Dim + d.Dim
	yMin := dy * d.Dim
	yMax := dy*d.Dim + d.Dim

	if dx == -1 {
		xMin = xMax - 1
	} else if dx == 1 {
		xMax = xMin + 1
	}

	if dy == -1 {
		yMin = yMax - 1
	} else if dy == 1 {
		yMax = yMin + 1
	}

	ox := -dx * d.Dim
	oy := -dy * d.Dim

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			d.Data[d.idx(x, y)] = data.Data[d.idx(x+ox, y+oy)]
		}
	}
}

func (d *DEMData) Get(x int, y int) float64 {
	unpack := d.getUnpackVector()
	value := d.Data[d.idx(x, y)]
	return float64(value[0])*unpack[0] + float64(value[1])*unpack[1] + float64(value[2])*unpack[2] - unpack[3]
}

func (d *DEMData) GetData() []float64 {
	ret := make([]float64, d.Dim*d.Dim)
	for x := 0; x < d.Dim; x++ {
		for y := 0; y < d.Dim; y++ {
			ret[x*d.Dim+y] = d.Get(x, y)
		}
	}
	return ret
}

func (d *DEMData) getUnpackVector() [4]float64 {
	if d.Encoding == 0 {
		return UNPACK_MAPBOX
	}
	return UNPACK_TERRARIUM
}


type DemPacker interface{
	func Pack(val float64) [4]byte
}
 
type MapboxPacker struct{
	Base float64
	Interval float64
}

func (p *MapboxPacker)Pack(h float64) [4]byte {
   val := (h + p.base) / p.Interval
    r := (math.Floor(math.Floor(val / 256) / 256) / 256 -
             math.Floor(math.Floor(math.Floor(val / 256) / 256) / 256)) *
            256
    g := (math.Floor(val / 256) / 256 -
             math.Floor(math.Floor(val / 256) / 256)) *
            256
    b := (val / 256 - math.Floor(val / 256)) * 256
		var image [4]byte
   image[index] = uint8_t(r)
   image[index + 1] = uint8_t(g)
   image[index + 2] = uint8_t(b)
   image[index + 3] = 1
	 return image
}

 type TerrariumPacker struct{
		Base float64
 }

func (p *TerrariumPacker)Pack(h float64) [4]byte {
   val = h + p.Base
   r := math.Floor(val / 256)
   g := math.Floor(int(val) % 256)
   b := math.Floor(int(val * 256) % 256)
   image[index] = uint8_t(r)
   image[index + 1] = uint8_t(g)
   image[index + 2] = uint8_t(b)
   image[index + 3] = 1
	 return image
}


func DemEncode(path string, pk DemPacker) (image.Image,error) {
	rst,err:=rt.CreateRasterFromFile(path)
	if err != nil {
		return nil, err
	}
	data := make([][4]byte,rst.Rows*rst.Columns)
	img := image.NewNRGBA(image.Rect(0, 0, rst.Rows, rst.Columns))

	for y := 0; y < rst.Columns; y++ {
			for x := 0; x < rst.Rows; x++ {
			h := rst.Value(x,y)
		   dt= pk.Pack(h)
			 img.SetNRGBA(x, y, color.NRGBA{
				R: dt[0],
				G: dt[1],
				B: dt[2],
				A: dt[3],
			}
		}
	}
	return img, nil
}