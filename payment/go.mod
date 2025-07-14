module github.com/alexis871aa/microservices-rocket-factory/payment

go 1.24.4

require (
	github.com/alexis871aa/microservices-rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/alexis871aa/microservices-rocket-factory/shared => ../shared
