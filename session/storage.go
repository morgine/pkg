package session

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Storage interface {
	// SaveToken 保存 token
	SaveToken(userID, token string, expires int64) error
	// CheckAndRefreshToken 检测 token，如果 token 合法则刷新过期时间
	CheckAndRefreshToken(userID, token string, expires int64) (ok bool, err error)
	// RemoveToken 移除 token
	RemoveToken(userID, token string) error
	// RemoveUser 移除用户所有 token
	RemoveUser(userID string) error
}

type RedisStorage struct {
	client    *redis.Client
	keyPrefix string
}

func (rs *RedisStorage) RemoveToken(userID, token string) error {
	return rs.client.Del(noCtx, rs.tokenKey(userID, token)).Err()
}

func (rs *RedisStorage) RemoveUser(userID string) error {
	keys, err := rs.client.Keys(noCtx, rs.userKey(userID)).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if len(keys) > 0 {
		return rs.client.Del(noCtx, keys...).Err()
	}
	return nil
}

func NewRedisStorage(keyPrefix string, client *redis.Client) Storage {
	return &RedisStorage{keyPrefix: keyPrefix, client: client}
}

func (rs *RedisStorage) userKey(id string) string {
	return rs.keyPrefix + id
}

func (rs *RedisStorage) tokenKey(id, token string) string {
	return rs.userKey(id) + "_" + token
}

var noCtx = context.Background()

func (rs *RedisStorage) SaveToken(id, token string, expires int64) error {
	return rs.client.Set(noCtx, rs.tokenKey(id, token), "1", time.Duration(expires)*time.Second).Err()
}

func (rs *RedisStorage) CheckAndRefreshToken(id, token string, expires int64) (ok bool, err error) {
	key := rs.tokenKey(id, token)
	savedToken, err := rs.client.Get(noCtx, key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if savedToken != "1" {
		return false, nil
	} else {
		err = rs.client.Expire(noCtx, key, time.Duration(expires)*time.Second).Err()
		if err != nil {
			return false, nil
		} else {
			return true, nil
		}
	}
}
