package v1

import (
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/service"
	userV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
