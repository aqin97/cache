package cacheclient

import (
	"github.com/go-redis/redis"
)

type redisClient struct {
	*redis.Client
}

func (c *redisClient) get(key string) (string, error) {
	res, err := c.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func (c *redisClient) set(key, value string) error {
	return c.Set(key, value, 0).Err()
}

func (c *redisClient) del(key string) error {
	return c.Del(key).Err()
}

func (c *redisClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value, cmd.Error = c.get(cmd.Key)
	}
	if cmd.Name == "set" {
		cmd.Error = c.set(cmd.Key, cmd.Value)
	}
	if cmd.Name == "del" {
		cmd.Error = c.del(cmd.Key)
	}
	panic("unknow cmd name " + cmd.Name)
}

func (c *redisClient) PipelineRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	pipe := c.Pipeline()
	cmderrs := make([]redis.Cmder, len(cmds))
	for i, cmd := range cmds {
		if cmd.Name == "get" {
			cmderrs[i] = pipe.Get(cmd.Key)
		} else if cmd.Name == "set" {
			cmderrs[i] = pipe.Set(cmd.Key, cmd.Value, 0)
		} else if cmd.Name == "del" {
			cmderrs[i] = pipe.Del(cmd.Key)
		} else {
			panic("unknow cmd name " + cmd.Name)
		}
	}
	_, err := pipe.Exec()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	for i, cmd := range cmds {
		if cmd.Name == "get" {
			value, err := cmderrs[i].(*redis.StringCmd).Result()
			if err == redis.Nil {
				value, err = "", nil
			}
			cmd.Value, cmd.Error = value, err
		} else {
			cmd.Error = cmderrs[i].Err()
		}

	}
}

func newRedisClient(host string) *redisClient {
	return &redisClient{redis.NewClient(&redis.Options{
		Addr:        host + ":6379",
		ReadTimeout: -1,
	})}
}
