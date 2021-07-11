module github.com/flywave/go-mapbox

go 1.13

require (
	github.com/flywave/freetype v0.0.0-20200612054648-f2aab071ba59
	github.com/flywave/go-geom v0.0.0-20210705081559-eee15cf4b503
	github.com/flywave/go-pbf v0.0.0-20210527131326-8e27970d0076
	github.com/flywave/go-raster v0.0.0-20210526065301-f50e348f662e
	github.com/flywave/imaging v1.6.5 // indirect
	github.com/go-courier/reflectx v1.3.4
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.5.2
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/paulmach/go.geojson v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
)

replace golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9
