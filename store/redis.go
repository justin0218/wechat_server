package store

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Redis struct {
	conf *Config
}

func (s *Redis) Get() *redis.Client {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", s.conf.Get().Redis.Host, s.conf.Get().Redis.Port),
			Password: s.conf.Get().Redis.Pass,
			DB:       0,
		})
	})
	return redisClient
}

func (s *Redis) GetAccessTokenKey(appid string) string {
	return fmt.Sprintf("%s_access_token", appid)
}

func (s *Redis) GetAccessTicketKey(appid string) string {
	return fmt.Sprintf("%s_ticket", appid)
}

func (s *Redis) GetShortUrlKey(lurl string) string {
	return lurl
}
