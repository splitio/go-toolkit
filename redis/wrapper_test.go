package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/splitio/go-toolkit/v5/testhelpers"
)

func TestRedisWrapperPipeline(t *testing.T) {
	rc := redis.NewUniversalClient(&redis.UniversalOptions{})
	client := &ClientImpl{wrapped: rc}

	client.Del("key1")
	client.Del("key-test")
	client.Del("key-set")
	client.RPush("key1", "e1", "e2", "e3")
	client.Set("key-del1", 0, 1*time.Hour)
	client.Set("key-del2", 0, 1*time.Hour)

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
	pipe.Del([]string{"key-del1", "key-del2"}...)
	result, err := pipe.Exec()
	if err != nil {
		t.Error("there should not be any error. Got: ", err)
	}

	if len(result) != 14 {
		t.Error("there should be 13 elements")
	}

	items, _ := result[0].Multi()
	testhelpers.AssertStringSliceEquals(t, items, []string{"e1", "e2", "e3"}, "result of lrange should be e1,e2,e3")
	if l := result[1].Int(); l != 3 {
		t.Error("length should be 3. is: ", l)
	}

	if i := client.LLen("key1").Int(); i != 1 {
		t.Error("new length should be 1. Is: ", i)
	}

	if c := result[3].Int(); c != 5 {
		t.Error("count should be 5. Is: ", c)
	}

	if c := result[4].Int(); c != 4 {
		t.Error("count should be 5. Is: ", c)
	}

	if c := result[5].Int(); c != 7 {
		t.Error("count should be 5. Is: ", c)
	}

	if l := result[6].Int(); l != 2 {
		t.Error("hlen should be 2. is: ", l)
	}

	if ib := client.HIncrBy("key-test", "field-test", 1); ib.Int() != 6 {
		t.Error("new count should be 6")
	}

	if ib := client.Get("key-set"); ib.String() != "field-test-1" {
		t.Error("it should be field-test-1")
	}

	if c := result[8].Int(); c != 2 {
		t.Error("count should be 2. Is: ", c)
	}
	if d, _ := result[9].Multi(); len(d) != 2 {
		t.Error("count should be 2. Is: ", len(d))
	}
	if c := result[10].Int(); c != 2 {
		t.Error("count should be 2. Is: ", c)
	}
	if c := result[11].Int(); c != 1 {
		t.Error("count should be 1. Is: ", c)
	}
	if c := result[12].Int(); c != 0 {
		t.Error("count should be zero. Is: ", c)
	}
	if c := result[13].Int(); c != 2 {
		t.Error("count should be 2. Is: ", c)
	}
}
