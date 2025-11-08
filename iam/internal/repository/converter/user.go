package converter

import (
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/model"
)

func UserInfoToRepo(info model.UserInfo) *repoModel.UserInfo {
	notificationMethods := make([]repoModel.NotificationMethod, 0, len(info.NotificationMethods))
	for _, method := range info.NotificationMethods {
		notificationMethods = append(notificationMethods, repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &repoModel.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: notificationMethods,
	}
}

func RepoToUserInfo(info repoModel.UserInfo) *model.UserInfo {
	notificationMethods := make([]model.NotificationMethod, 0, len(info.NotificationMethods))
	for _, method := range info.NotificationMethods {
		notificationMethods = append(notificationMethods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &model.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: notificationMethods,
	}
}

func RepoToUser(user repoModel.User) *model.User {
	notificationMethods := make([]model.NotificationMethod, 0, len(user.Info.NotificationMethods))
	for _, method := range user.Info.NotificationMethods {
		notificationMethods = append(notificationMethods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	return &model.User{
		UserUUID: user.UserUUID,
		Info: model.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: notificationMethods,
		},
	}
}

func UserToRepo(user model.User, password []byte) *repoModel.User {
	notificationMethods := make([]repoModel.NotificationMethod, 0, len(user.Info.NotificationMethods))
	for _, method := range user.Info.NotificationMethods {
		notificationMethods = append(notificationMethods, repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &repoModel.User{
		UserUUID: user.UserUUID,
		Info: repoModel.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: notificationMethods,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Password:  password,
	}
}
