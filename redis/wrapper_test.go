package redis

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/splitio/go-toolkit/v5/testhelpers"
)

func TestRedisWrapperPipeline(t *testing.T) {
	rc := redis.NewUniversalClient(&redis.UniversalOptions{})
	client := &ClientImpl{wrapped: rc}

	client.Del("key1")
	client.Del("key-test")
	client.RPush("key1", "e1", "e2", "e3")

	pipe := client.Pipeline()
	pipe.LRange("key1", 0, 5)
	pipe.LLen("key1")
	pipe.LTrim("key1", 2, -1)
	pipe.HIncrBy("key-test", "field-test", 5)
	result, err := pipe.Exec()
	if err != nil {
		t.Error("there should not be any error. Got: ", err)
	}

	if len(result) != 4 {
		t.Error("there should be 4 elements")
	}

	items, _ := result[0].Multi()
	testhelpers.AssertStringSliceEquals(t, items, []string{"e1", "e2", "e3"}, "result of lrange should be e1,e2,e3")
	if l := result[1].Int(); l != 3 {
		t.Error("length should be 3. is: ", l)
	}

	if i := client.LLen("key1").Int(); i != 1 {
		t.Error("new length should be 1. Is: ", i)
	}

	if ib := client.HIncrBy("key-test", "field-test", 1); ib.Int() != 6 {
		t.Error("new count should be 6")
	}
}
