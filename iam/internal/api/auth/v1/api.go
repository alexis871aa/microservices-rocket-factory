package v1

import (
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/service"
	authV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
