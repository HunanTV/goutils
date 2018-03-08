package goutils

import (
	"errors"
	"strings"
	"time"

	"gopkg.in/redis.v3"
)

type RedisConfig interface {
	GetAddr() string
	GetPoolNum() int
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetPoolTimeout() time.Duration
	GetDialTimeout() time.Duration
}

type RedisClient interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	MGet(keys ...string) *redis.SliceCmd
	IncrBy(key string, value int64) *redis.IntCmd
	TTL(key string) *redis.DurationCmd
	Del(keys ...string) *redis.IntCmd
	Keys(pattern string) *redis.StringSliceCmd
	HGet(key, field string) *redis.StringCmd
	HGetAllMap(key string) *redis.StringStringMapCmd
	HMSet(key, field, value string, pairs ...string) *redis.StatusCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Exists(key string) *redis.BoolCmd
}

type Redis struct {
	client RedisClient
}

func NewRedis(redisConfig RedisConfig) *Redis {
	return &Redis{client: initRedis(redisConfig)}
}

func initRedisNormal(redisConfig RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConfig.GetAddr(),
		PoolSize:     redisConfig.GetPoolNum(),
		ReadTimeout:  redisConfig.GetReadTimeout(),
		WriteTimeout: redisConfig.GetWriteTimeout(),
		PoolTimeout:  redisConfig.GetPoolTimeout(),
		DialTimeout:  redisConfig.GetDialTimeout(),
	})
	_, err := client.Ping().Result()
	if err != nil {
		Log.Error("init redis ping error:%s", err.Error())
	}
	return client, err
}

func initRedisCluster(redisConfig RedisConfig) (*redis.ClusterClient, error) {
	if len(redisConfig.GetAddr()) == 0 {
		Log.Fatal("null redis addr")
	}
	addrSegs := strings.Split(redisConfig.GetAddr(), ",")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrSegs,
		PoolSize:     redisConfig.GetPoolNum(),
		ReadTimeout:  redisConfig.GetReadTimeout(),
		WriteTimeout: redisConfig.GetWriteTimeout(),
		PoolTimeout:  redisConfig.GetPoolTimeout(),
		DialTimeout:  redisConfig.GetDialTimeout(),
	})
	_, err := client.Ping().Result()
	if err != nil {
		Log.Error("init redis ping error")
	}
	return client, err
}

func initRedis(redisConfig RedisConfig) RedisClient {
	var client RedisClient
	if strings.Contains(redisConfig.GetAddr(), ",") {
		client, _ = initRedisCluster(redisConfig)
	} else {
		client, _ = initRedisNormal(redisConfig)
	}
	return client
}

//get string key just for freq condition
func (r *Redis) Get(key string) (string, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return "", errors.New("not initied")
	}
	cmd := r.client.Get(key)
	err := cmd.Err()
	if err == redis.Nil {
		Log.Debug("key:(%s) is not exists", key)
		return "", nil
	} else if err != nil {
		Log.Error("get key:(%s), but err:(%s)", key, err)
		return "", err
	}
	value := cmd.Val()
	return value, err
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return errors.New("not initied")
	}
	cmd := r.client.Set(key, value, expiration)
	return cmd.Err()
}

func (r *Redis) Mget(keys ...string) ([]interface{}, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return nil, errors.New("not initied")
	}
	result := r.client.MGet(keys...)
	return result.Result()
}

func (r *Redis) HGet(key, field string) (string, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return "", errors.New("not initied")
	}
	cmd := r.client.HGet(key, field)
	err := cmd.Err()
	if err == redis.Nil {
		Log.Debug("key:(%s) is not exists", key)
		return "", nil
	} else if err != nil {
		Log.Error("get key:(%s), but err:(%s)", key, err)
		return "", err
	}
	value := cmd.Val()
	return value, err
}

func (r *Redis) HGetAllMap(key string) (map[string]string, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return nil, errors.New("not inited")
	}
	cmd := r.client.HGetAllMap(key)
	err := cmd.Err()
	if err == redis.Nil {
		Log.Debug("key:(%s) is not exists", key)
		return nil, nil
	}
	if err != nil {
		Log.Error("get key:(%s), but err:(%s)", key, err)
		return nil, err
	}
	value := cmd.Val()
	return value, nil
}

func (r *Redis) HMSet(key, feild, value string, pairs ...string) error {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return nil
	}
	cmd := r.client.HMSet(key, feild, value, pairs...)
	return cmd.Err()
}

func (r *Redis) SetExpire(key string, expiration time.Duration) error {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return nil
	}
	cmd := r.client.Expire(key, expiration)
	return cmd.Err()
}

func (r *Redis) Exists(key string) (bool, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return false, nil
	}
	cmd := r.client.Exists(key)
	return cmd.Val(), cmd.Err()
}

//IncrBy(key string, value int64) *IntCmd
func (r *Redis) Incby(key string, value int64) (int64, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return 0, errors.New("not initied")
	}
	return r.client.IncrBy(key, value).Result()
}

func (r *Redis) TTL(key string) (time.Duration, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return time.Nanosecond, errors.New("not initied")
	}
	return r.client.TTL(key).Result()
}
func (r *Redis) Del(key string) error {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return errors.New("not initied")
	}
	return r.client.Del(key).Err()
}

func (r *Redis) Keys(pattern string) ([]string, error) {
	if r.client == nil {
		Log.Error("Get redis value, but redis r.client is nil!")
		return nil, nil
	}
	stringSliceCmd := r.client.Keys(pattern)
	err := stringSliceCmd.Err()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		Log.Error("get all keys, but err:(%s)", err)
		return nil, err
	}
	return stringSliceCmd.Val(), nil
}
