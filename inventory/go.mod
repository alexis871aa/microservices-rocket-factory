module github.com/alexis871aa/microservices-rocket-factory/inventory

go 1.24.4

replace github.com/alexis871aa/microservices-rocket-factory/shared => ../shared

require (
	github.com/alexis871aa/microservices-rocket-factory/shared v0.0.0-20250714133457-49d7d66272ee
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/go-faster/errors v0.7.1
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.51.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.74.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250715232539-7130f93afb79 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
