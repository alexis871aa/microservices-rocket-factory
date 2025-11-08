package model

import (
	"time"
)

type User struct {
	UserUUID  string     `json:"user_uuid"`
	Info      UserInfo   `json:"info"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Password  []byte     `json:"password_hash"`
}

type UserInfo struct {
	Login               string               `json:"login"`                // Логин
	Email               string               `json:"email"`                // Email
	NotificationMethods []NotificationMethod `json:"notification_methods"` // Каналы уведомлений
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"` // Провайдер: telegram, email, push и т.д.
	Target       string `json:"target"`        // Адрес/идентификатор назначения (email, чат-id)
}
