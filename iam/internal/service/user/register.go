package user

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
)

func (s *service) Register(ctx context.Context, userInfo model.UserInfo, password string) (userUUID string, err error) {
	userUUID, err = s.userRepository.Create(ctx, userInfo, password)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}
