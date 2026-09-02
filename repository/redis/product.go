package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"goproject/domain"
	"goproject/repository"

	redisClient "github.com/redis/go-redis/v9"
)

type ProductCache struct {
	client *redisClient.Client
	ttl    time.Duration
}

func NewProductCache(client *redisClient.Client, ttl time.Duration) *ProductCache {
	return &ProductCache{client: client, ttl: ttl}
}

func (cache *ProductCache) Get(ctx context.Context, id uint) (*domain.Product, error) {
	value, err := cache.client.Get(ctx, productKey(id)).Bytes()
	if err == redisClient.Nil {
		return nil, repository.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	var product domain.Product
	if err := json.Unmarshal(value, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (cache *ProductCache) Set(ctx context.Context, product *domain.Product) error {
	value, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, productKey(product.ID), value, cache.ttl).Err()
}

func (cache *ProductCache) Delete(ctx context.Context, id uint) error {
	return cache.client.Del(ctx, productKey(id)).Err()
}

func (cache *ProductCache) Ping(ctx context.Context) error {
	return cache.client.Ping(ctx).Err()
}

func productKey(id uint) string {
	return fmt.Sprintf("product:%d", id)
}
