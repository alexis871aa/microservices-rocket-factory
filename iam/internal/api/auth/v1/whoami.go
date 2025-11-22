package v1

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/converter"
	authV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	session, user, err := a.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &authV1.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
