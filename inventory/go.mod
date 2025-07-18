module github.com/alexis871aa/microservices-rocket-factory/inventory

go 1.24.4

replace github.com/alexis871aa/microservices-rocket-factory/shared => ../shared

require (
	github.com/alexis871aa/microservices-rocket-factory/shared v0.0.0-20250714133457-49d7d66272ee
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/go-faster/errors v0.7.1
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.51.0
	google.golang.org/grpc v1.74.0
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250715232539-7130f93afb79 // indirect
)
