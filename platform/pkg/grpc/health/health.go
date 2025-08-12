package health

import "google.golang.org/grpc/health/grpc_health_v1"

type Server struct {
	grpc_health_v1.UnimplementedHealthServer
}
