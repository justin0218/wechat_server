package services

import (
	"wechat_server/store"
)

type baseService struct {
	Redis  store.Redis
	Config store.Config
}
