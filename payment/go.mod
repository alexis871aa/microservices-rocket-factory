module github.com/alexis871aa/microservices-rocket-factory/payment

go 1.24.4

require (
	github.com/alexis871aa/microservices-rocket-factory/platform v0.0.0-00010101000000-000000000000
	github.com/alexis871aa/microservices-rocket-factory/shared v0.0.0-20250808212652-d1eb6bf8f836
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.74.2
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250811230008-5f3141c8851a // indirect
	google.golang.org/protobuf v1.36.7 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alexis871aa/microservices-rocket-factory/shared => ../shared

replace github.com/alexis871aa/microservices-rocket-factory/platform => ../platform
