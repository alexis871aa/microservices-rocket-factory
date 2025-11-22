package model

import (
	"time"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
)

type User struct {
	UserUUID  string         `json:"user_uuid"`
	Info      model.UserInfo `json:"info"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	Password  string         `json:"password_hash"`
}
