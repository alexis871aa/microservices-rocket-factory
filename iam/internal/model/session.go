package model

import "time"

type Session struct {
	Uuid      string     // UserUUID сессии
	CreatedAt time.Time  // Время создания
	UpdatedAt *time.Time // Время последнего обновления
	ExpiresAt time.Time  // Время истечения
}
