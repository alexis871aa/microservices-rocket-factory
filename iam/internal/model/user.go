package model

import "time"

type NotificationMethod struct {
	ProviderName string // Провайдер: telegram, email, push и т.д.
	Target       string // Адрес/идентификатор назначения (email, чат-id)
}

type UserInfo struct {
	Login               string               // Логин
	Email               string               // Email
	NotificationMethods []NotificationMethod //	Каналы уведомлений
}

type User struct {
	UserUUID  string     // UserUUID пользователя
	Info      UserInfo   // Базовая информация
	CreatedAt time.Time  // Дата создания
	UpdatedAt *time.Time // Дата обновления
}
