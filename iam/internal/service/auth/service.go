package auth

import (
	"time"

	iamService "github.com/alexis871aa/microservices-rocket-factory/iam/internal/service"
)

var _ iamService.AuthService = (*service)(nil)

type service struct {
	sessionRepository iamService.SessionRepository
	userRepository    iamService.UserRepository
	cacheTTL          time.Duration
}

func NewService(sessionRepository iamService.SessionRepository, userRepository iamService.UserRepository, cacheTTL time.Duration) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		cacheTTL:          cacheTTL,
	}
}
