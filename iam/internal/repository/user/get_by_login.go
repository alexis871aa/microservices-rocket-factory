package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
	repoConverter "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/model"
)

func (r *repository) GetUserByLogin(ctx context.Context, login, password string) (*model.User, error) {
	builder := sq.Select("user_uuid", "info", "created_at", "updated_at", "password_hash").
		From("users").
		Where("info->>'login' = ?", login).
		PlaceholderFormat(sq.Dollar).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return nil, err
	}

	var user repoModel.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.UserUUID,
		&user.Info,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Password,
	)
	if err != nil {
		log.Printf("failed to select user: %v\n", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, repoModel.ErrUserNotFound
	}

	return repoConverter.RepoToUser(user), nil
}
