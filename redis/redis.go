package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisCli Redis客户端
var RedisCli *RedisClient

const RedisAddress = "192.168.220.128:6379"
const RedisPassword = "78KJtyjg0928abc"

func init() {
	//利用选项模式进行初始化Redis客户端
	RedisCli = NewRedisClient(
		WithAddress(RedisAddress),
		WithPassword(RedisPassword),
		WithPrefix("Li:"),
	)
}

type RedisClient struct {
	Client  *redis.Client
	prefix  string         //prefix参数是当前结构体特有的，可通过选项模式动态添加次参数默认值
	options *redis.Options //redis参数太多可通过RedisClientOption选项模式动态添加初始值
}

// NewRedisClient 创建一个Redis客户端
func NewRedisClient(options ...RedisClientOption) *RedisClient {
	//初始化空结构体
	c := &RedisClient{
		Client:  nil,
		prefix:  "",
		options: &redis.Options{},
	}

	//由于redis参数非常多，这里提供选项模式对redis参数进行选择性初始化，提供扩展效率
	for _, option := range options {
		option(c) //这里可以考虑调整为c，可以让自定义的prefix也可以使用选项模式
	}

	//初始化client参数
	c.Client = redis.NewClient(c.options)

	return c
}

// RedisClientOption 是一个redis配置选项函数类型
type RedisClientOption func(options *RedisClient)

// WithPassword 返回一个设置 password 的选项函数
func WithPassword(password string) RedisClientOption {
	return func(c *RedisClient) {
		c.options.Password = password
	}
}

// WithAddress 返回一个设置 address 的选项函数
func WithAddress(address string) RedisClientOption {
	return func(c *RedisClient) {
		c.options.Addr = address
	}
}

// WithPrefix 返回一个设置 prefix 的选项函数
func WithPrefix(prefix string) RedisClientOption {
	return func(c *RedisClient) {
		c.prefix = prefix
	}
}

func (c *RedisClient) GetFullKey(key string) string {
	return c.prefix + key
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, c.GetFullKey(key), value, expiration).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.Client.Get(ctx, c.GetFullKey(key))
}

func (c *RedisClient) GetLock(ctx context.Context, key string, expiration time.Duration) (err error) {
	success, err := c.Client.SetNX(ctx, c.GetFullKey(key), 1, expiration).Result()
	if err != nil {
		return err
	}
	if success {
		return nil
	}
	return errors.New("redis get lock error: key already exists")
}

func (c *RedisClient) ReleaseLock(ctx context.Context, key string) error {
	res := c.Client.Del(ctx, c.GetFullKey(key))
	if res.Err() != nil {
		return fmt.Errorf("redis del error: %s", res.Err().Error())
	}
	return nil
}
