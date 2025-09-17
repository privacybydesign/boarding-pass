package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type InMemoryTokenStorage struct {
	TokenMap map[string]string
	mutex    sync.Mutex
}

func NewInMemoryTokenStorage() *InMemoryTokenStorage {
	return &InMemoryTokenStorage{
		TokenMap: make(map[string]string),
	}
}

type RedisTokenStorage struct {
	client    *redis.Client
	namespace string
}

func NewRedisTokenStorage(client *redis.Client, namespace string) *RedisTokenStorage {
	return &RedisTokenStorage{client: client, namespace: namespace}
}

// Should be safe to use in concurreny
type TokenStorage interface {
	// Store given token for the given email address,
	// returns an error when it somehow fails to store the value.
	// Should not return an error when the value already exists,
	// it should just update in that case.
	StoreToken(sessionPtr, token string) error

	// Should retrieve the token for the given email address
	// and return an error in any case where it fails to do so.
	RetrieveToken(sessionPtr string) (string, error)

	// Should remove the token and return an error if it fails to do so.
	// The value not being there should also be considered an error.
	RemoveToken(sessionPtr string) error
}

// ------------------------------------------------------------------------------

func createKey(namespace, sessionPtr string) string {
	return fmt.Sprintf("%s:token:%s", namespace, sessionPtr)
}

const Timeout time.Duration = 24 * time.Hour

func (s *RedisTokenStorage) StoreToken(sessionPtr, token string) error {
	ctx := context.Background()
	return s.client.Set(ctx, createKey(s.namespace, sessionPtr), token, Timeout).Err()
}

func (s *RedisTokenStorage) RetrieveToken(sessionPtr string) (string, error) {
	ctx := context.Background()
	return s.client.Get(ctx, createKey(s.namespace, sessionPtr)).Result()
}

func (s *RedisTokenStorage) RemoveToken(sessionPtr string) error {
	ctx := context.Background()
	return s.client.Del(ctx, createKey(s.namespace, sessionPtr)).Err()
}

func NewTokenStorage(storage_cfg *StorageConfig) TokenStorage {
	switch storage_cfg.Type {
	case "redis":
		return NewRedisTokenStorage(redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", storage_cfg.RedisConfig.Host, storage_cfg.RedisConfig.Port),
			Password: storage_cfg.RedisConfig.Password,
			DB:       0, // use default DB
		}), storage_cfg.RedisConfig.Namespace)
	case "inmemory":
		return NewInMemoryTokenStorage()
	default:
		log.Fatalf("unknown storage type: %s", storage_cfg.Type)
		return nil
	}
}

// ------------------------------------------------------------------------------

func (s *InMemoryTokenStorage) StoreToken(sessionPtr, token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.TokenMap[sessionPtr] = token
	return nil
}

func (s *InMemoryTokenStorage) RetrieveToken(sessionPtr string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if token, ok := s.TokenMap[sessionPtr]; ok {
		return token, nil
	} else {
		return "", fmt.Errorf("failed to find token for %s", sessionPtr)
	}
}

func (s *InMemoryTokenStorage) RemoveToken(sessionPtr string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.TokenMap[sessionPtr]; ok {
		delete(s.TokenMap, sessionPtr)
		return nil
	} else {
		return fmt.Errorf("failed to remove token for %s, because it wasn't there", sessionPtr)
	}
}
