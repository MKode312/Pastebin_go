package redis

import (
	"context"
	"fmt"
	"text_sharing/internal/config"
	"text_sharing/internal/storage"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

func NewClientForCache(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.CfgRedis.Address,
		WriteTimeout: config.CfgRedis.Timeout,
		ReadTimeout:  config.CfgRedis.Timeout,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return rdb, nil
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *Cache) SetLinksCache(linkID string, link string) error {
	status := c.client.Set(c.ctx, linkID, link, 0)
	if status.Err() != nil {
		return fmt.Errorf("failed to set link in the cache: %w", status.Err())
	}
	return nil
}

func (c *Cache) GetLinkFromCache(linkID string) (string, error) {
	const op = "storage.redis.GetLinkFromCache"
	link, err := c.client.Get(c.ctx, linkID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", storage.ErrCacheMiss
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return link, nil
}
