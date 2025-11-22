package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type NotificationMethod struct {
	ProviderName string `json:"provider_name"` // Провайдер: telegram, email, push и т.д.
	Target       string `json:"target"`        // Адрес/идентификатор назначения (email, чат-id)
}

type UserInfo struct {
	Login               string               `json:"login"`                // Логин
	Email               string               `json:"email"`                // Email
	NotificationMethods []NotificationMethod `json:"notification_methods"` // Каналы уведомлений
}

// Scan реализует интерфейс sql.Scanner для UserInfo
func (ui *UserInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan UserInfo: value is not []byte")
	}

	return json.Unmarshal(bytes, ui)
}

// Value реализует интерфейс driver.Valuer для UserInfo
func (ui UserInfo) Value() (driver.Value, error) {
	return json.Marshal(ui)
}

type User struct {
	UserUUID  string     // UserUUID пользователя
	Info      UserInfo   // Базовая информация
	CreatedAt time.Time  // Дата создания
	UpdatedAt *time.Time // Дата обновления
}
