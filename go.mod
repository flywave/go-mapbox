module github.com/flywave/go-mapbox

go 1.13

require (
	github.com/flywave/freetype v0.0.0-20200612054648-f2aab071ba59
	github.com/flywave/go-cog v0.0.0-20220109113741-c4b8ae8c49be
	github.com/flywave/go-geom v0.0.0-20211230100258-27b9a5f30082
	github.com/flywave/go-pbf v0.0.0-20210701015929-a3bdb1f6728e
	github.com/flywave/go-quantized-mesh v0.0.0-20210525134750-cb854922974d
	github.com/flywave/imaging v1.6.5
	github.com/flywave/webp v1.1.1
	github.com/go-courier/reflectx v1.3.4
	github.com/golang/protobuf v1.5.2
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/paulmach/go.geojson v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.14.0
)

replace golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9

replace github.com/flywave/go-geom => ../go-geom
