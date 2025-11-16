package storage

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore() *RedisStore {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "redis"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":6379",
		DB:   0,
	})

	return &RedisStore{
		Client: rdb,
	}
}

func (s *RedisStore) generateShortURL(originalURL string) string {
	hash := fnv.New32a()
	hash.Write([]byte(originalURL))
	return fmt.Sprintf("%x", hash.Sum32())
}

func (s *RedisStore) getDomain(originalURL string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(parsedURL.Host, "www."), nil
}

func (s *RedisStore) GetOriginalURL(shortURL string) (string, error) {
	originalURL, err := s.Client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("URL not found")
	} else if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s *RedisStore) GetDomainCounts() (map[string]int, error) {
	keys, err := s.Client.Keys(ctx, "domain:*").Result()
	if err != nil {
		return nil, err
	}
	domainCounts := make(map[string]int)

	for _, key := range keys {
		count, err := s.Client.Get(ctx, key).Int()
		if err != nil {
			return nil, err
		}
		domain := strings.TrimPrefix(key, "domain:")
		domainCounts[domain] = count
	}
	return domainCounts, nil
}

func (s *RedisStore) SaveURL(originalURL string) (string, error) {
	shortURL, err := s.Client.Get(ctx, originalURL).Result()
	if err == redis.Nil {
		shortURL = s.generateShortURL(originalURL)
		err = s.Client.Set(ctx, shortURL, originalURL, 0).Err()
		if err != nil {
			return "", err
		}
		err = s.Client.Set(ctx, shortURL, originalURL, 0).Err()
		if err != nil {
			return "", err
		}
		domain, err := s.getDomain(originalURL)
		if err != nil {
			return "", err
		}

		err = s.Client.Incr(ctx, fmt.Sprintf("domain:%s", domain)).Err()
		if err != nil {
			return "", nil
		}
	} else if err != nil {
		return "", nil
	}
	return shortURL, nil
}

func (s *RedisStore) FlushDB() {
	s.Client.FlushDB(ctx)
}
