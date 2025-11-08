module github.com/alexis871aa/microservices-rocket-factory/iam

go 1.24.4

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/alexis871aa/microservices-rocket-factory/platform v0.0.0-20250908115142-da6ebaef04a0
	github.com/alexis871aa/microservices-rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v11 v11.3.1
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/ogen-go/ogen v1.16.0
	github.com/pressly/goose/v3 v3.26.0
	github.com/samber/lo v1.52.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.43.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20251009144603-d2f985daa21b // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
)

replace github.com/alexis871aa/microservices-rocket-factory/shared => ../shared

replace github.com/alexis871aa/microservices-rocket-factory/platform => ../platform
