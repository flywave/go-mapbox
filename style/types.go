package style

type ObjectType int

const (
	ObjectTypeUnknown  ObjectType = 0
	ObjectTypeNode     ObjectType = 1
	ObjectTypeWay      ObjectType = 2
	ObjectTypeRelation ObjectType = 3
)

type ZoomLevel float64

const (
	MinZoomLevel ZoomLevel = 0
	MaxZoomLevel ZoomLevel = 255
)
