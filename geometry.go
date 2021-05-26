package mapbox

type GeometryCoordinate []float64
type GeometryCoordinates []GeometryCoordinate
type GeometryCollection []GeometryCoordinates

const (
	FeatureTypeUnknown    = 0
	FeatureTypePoint      = 1
	FeatureTypeLineString = 2
	FeatureTypePolygon    = 3
)

type GeometryTileFeature interface {
	GetType() int
	GetValue(key string) interface{}
	GetProperties() map[string]interface{}
	GetID() int
	GetGeometries() GeometryCollection
}

type GeometryTileLayer interface {
	GetFeatureCount() int32
	GetFeature(int32) *GeometryTileFeature
	GetName() string
}

type GeometryTileData interface {
	GetLayer(key string) *GeometryTileLayer
	GetLayerNames() []string
}
