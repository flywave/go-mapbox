module github.com/flywave/go-mapbox

go 1.23.0

toolchain go1.24.4

require (
	github.com/flywave/freetype v0.0.0-20200612054648-f2aab071ba59
	github.com/flywave/go-cog v0.0.0-20250314092301-4673589220b8
	github.com/flywave/go-geom v0.0.0-20250607125323-f685bf20f12c
	github.com/flywave/go-pbf v0.0.0-20230306063816-5e5b0da27bbd
	github.com/flywave/go-quantized-mesh v0.0.0-20210525134750-cb854922974d
	github.com/flywave/imaging v1.6.5
	github.com/flywave/webp v1.1.2
	github.com/go-courier/reflectx v1.3.4
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/paulmach/go.geojson v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.14.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/flywave/go-geo v0.0.0-20250314091853-e818cb9de299 // indirect
	github.com/flywave/go-geoid v0.0.0-20210705014121-cd8f70cb88bb // indirect
	github.com/flywave/go-geos v0.0.0-20210924031454-d16b758e2026 // indirect
	github.com/flywave/go-proj v0.0.0-20211220121303-46dc797a5cd0 // indirect
	github.com/flywave/go3d v0.0.0-20231213061711-48d3c5834480 // indirect
	github.com/google/tiff v0.0.0-20161109161721-4b31f3041d9a // indirect
	github.com/hhrutter/lzw v0.0.0-20190829144645-6f07a24e8650 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9

replace github.com/flywave/go-geom => ../go-geom

replace github.com/flywave/go-geos => ../go-geos

replace github.com/flywave/go-proj => ../go-proj

replace github.com/flywave/go-geoid => ../go-geoid
