module github.com/alexis871aa/microservices-rocket-factory/inventory

go 1.24.4

require (
	github.com/alexis871aa/microservices-rocket-factory/shared v0.0.0-20250808212652-d1eb6bf8f836
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/go-faster/errors v0.7.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/samber/lo v1.51.0
	github.com/stretchr/testify v1.10.0
	go.mongodb.org/mongo-driver v1.17.4
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250811230008-5f3141c8851a // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alexis871aa/microservices-rocket-factory/shared => ../shared

replace github.com/alexis871aa/microservices-rocket-factory/platform => ../platform
