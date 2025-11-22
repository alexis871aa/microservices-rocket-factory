package converter

import (
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/model"
)

func UserInfoToRepo(info model.UserInfo) model.UserInfo {
	return info
}

func RepoToUserInfo(info model.UserInfo) model.UserInfo {
	return info
}

func RepoToUser(user repoModel.User) *model.User {
	return &model.User{
		UserUUID:  user.UserUUID,
		Info:      user.Info,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToRepo(user model.User, password []byte) *repoModel.User {
	return &repoModel.User{
		UserUUID:  user.UserUUID,
		Info:      user.Info,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Password:  string(password),
	}
}
