package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisWrapperKeysAndScan(t *testing.T) {
	rc := redis.NewUniversalClient(&redis.UniversalOptions{})
	client := &ClientImpl{wrapped: rc}

	for i := 0; i < 10; i++ {
		client.Set(fmt.Sprintf("utest.key-del%d", i), 0, 1*time.Hour)
	}

	keys, err := client.Keys("utest*").Multi()
	assert.Nil(t, err)
	assert.Equal(t, 10, len(keys))
	var cursor uint64

	scanKeys := make([]string, 0)
	for {
		result := client.Scan(cursor, "utest*", 10)
		assert.Nil(t, result.Err())

		cursor = uint64(result.Int())

		keys, err := result.Multi()
		assert.Nil(t, err)

		scanKeys = append(scanKeys, keys...)

		if cursor == 0 {
			break
		}
	}

	assert.Equal(t, 10, len(scanKeys))
	for i := 0; i < 10; i++ {
		client.Del(fmt.Sprintf("utest.key-del%d", i))
	}
}

func TestRedisWrapperPipeline(t *testing.T) {
	rc := redis.NewUniversalClient(&redis.UniversalOptions{})
	client := &ClientImpl{wrapped: rc}

	client.RPush("key1", "e1", "e2", "e3")
	client.Set("key-del1", 0, 1*time.Hour)
	client.Set("key-del2", 0, 1*time.Hour)
	res := client.SetNX("key-setnx-cient", "field-test-1", 1*time.Hour)
	assert.True(t, res.Bool(), "setnx should be executed successfully")
	client.Del("key-setnx-cient")

	pipe := client.Pipeline()
	pipe.LRange("key1", 0, 5)
	pipe.LLen("key1")
	pipe.LTrim("key1", 2, -1)
	pipe.HIncrBy("key-test", "field-test", 5)
	pipe.HIncrBy("key-test", "field-test-2", 4)
	pipe.HIncrBy("key-test", "field-test-2", 3)
	pipe.HLen("key-test")
	pipe.Set("key-set", "field-test-1", 1*time.Hour)
	pipe.SAdd("key-sadd", []interface{}{"field-test-1", "field-test-2"})
	pipe.SMembers("key-sadd")
	pipe.SRem("key-sadd", []interface{}{"field-test-1", "field-test-2"})
	pipe.Incr("key-incr")
	pipe.Decr("key-incr")
	pipe.SetNX("key-setnx", "field-test-1", 30*time.Minute)
	pipe.Del([]string{"key-del1", "key-del2", "key-setnx"}...)
	result, err := pipe.Exec()
	assert.Nil(t, err)
	assert.Equal(t, 15, len(result))

	items, _ := result[0].Multi()
	assert.Equal(t, []string{"e1", "e2", "e3"}, items)
	assert.Equal(t, int64(3), result[1].Int())
	assert.Equal(t, int64(1), client.LLen("key1").Int())
	assert.Equal(t, int64(5), result[3].Int())
	assert.Equal(t, int64(4), result[4].Int())
	assert.Equal(t, int64(7), result[5].Int())
	assert.Equal(t, int64(2), result[6].Int())
	assert.Equal(t, int64(6), client.HIncrBy("key-test", "field-test", 1).Int())
	assert.Equal(t, "field-test-1", client.Get("key-set").String())
	assert.Equal(t, int64(2), result[8].Int())
	d, _ := result[9].Multi()
	assert.Equal(t, 2, len(d))
	assert.Equal(t, int64(2), result[10].Int())
	assert.Equal(t, int64(1), result[11].Int())
	assert.Equal(t, int64(0), result[12].Int())
	assert.True(t, result[13].Bool(), "setnx should be executed successfully")
	assert.Equal(t, int64(3), result[14].Int())

	client.Del([]string{"key1", "key-test", "key-set", "key-incr"}...)
}
