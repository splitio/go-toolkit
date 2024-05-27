package mocks

import (
	"time"

	"github.com/splitio/go-toolkit/v6/common"
	"github.com/splitio/go-toolkit/v6/redis"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

// ClusterCountKeysInSlot implements redis.Client.
func (m *MockClient) ClusterCountKeysInSlot(slot int) redis.Result {
	return m.Called(slot).Get(0).(redis.Result)
}

// ClusterKeysInSlot implements redis.Client.
func (m *MockClient) ClusterKeysInSlot(slot int, count int) redis.Result {
	return m.Called(slot, count).Get(0).(redis.Result)
}

// ClusterMode implements redis.Client.
func (m *MockClient) ClusterMode() bool {
	return m.Called().Bool(0)
}

// ClusterSlotForKey implements redis.Client.
func (m *MockClient) ClusterSlotForKey(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// Decr implements redis.Client.
func (m *MockClient) Decr(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// Del implements redis.Client.
func (m *MockClient) Del(keys ...string) redis.Result {
	return m.Called(common.AsInterfaceSlice(keys)...).Get(0).(redis.Result)
}

// Eval implements redis.Client.
func (m *MockClient) Eval(script string, keys []string, args ...interface{}) redis.Result {
	return m.Called(append([]interface{}{script, keys}, args...)...).Get(0).(redis.Result)
}

// Exists implements redis.Client.
func (m *MockClient) Exists(keys ...string) redis.Result {
	return m.Called(common.AsInterfaceSlice(keys)...).Get(0).(redis.Result)
}

// Expire implements redis.Client.
func (m *MockClient) Expire(key string, value time.Duration) redis.Result {
	return m.Called(key, value).Get(0).(redis.Result)
}

// Get implements redis.Client.
func (m *MockClient) Get(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// HGetAll implements redis.Client.
func (m *MockClient) HGetAll(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// HIncrBy implements redis.Client.
func (m *MockClient) HIncrBy(key string, field string, value int64) redis.Result {
	return m.Called(key, field, value).Get(0).(redis.Result)
}

// HSet implements redis.Client.
func (m *MockClient) HSet(key string, hashKey string, value interface{}) redis.Result {
	return m.Called(key, hashKey, value).Get(0).(redis.Result)
}

// Incr implements redis.Client.
func (m *MockClient) Incr(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// Keys implements redis.Client.
func (m *MockClient) Keys(pattern string) redis.Result {
	return m.Called(pattern).Get(0).(redis.Result)

}

// LLen implements redis.Client.
func (m *MockClient) LLen(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// LRange implements redis.Client.
func (m *MockClient) LRange(key string, start int64, stop int64) redis.Result {
	return m.Called(key, start, stop).Get(0).(redis.Result)
}

// LTrim implements redis.Client.
func (m *MockClient) LTrim(key string, start int64, stop int64) redis.Result {
	return m.Called(key, start, stop).Get(0).(redis.Result)
}

// MGet implements redis.Client.
func (m *MockClient) MGet(keys []string) redis.Result {
	return m.Called(keys).Get(0).(redis.Result)
}

// Ping implements redis.Client.
func (m *MockClient) Ping() redis.Result {
	return m.Called().Get(0).(redis.Result)
}

// Pipeline implements redis.Client.
func (m *MockClient) Pipeline() redis.Pipeline {
	return m.Called().Get(0).(redis.Pipeline)
}

// RPush implements redis.Client.
func (m *MockClient) RPush(key string, values ...interface{}) redis.Result {
	return m.Called(append([]interface{}{key}, values...)...).Get(0).(redis.Result)
}

// SAdd implements redis.Client.
func (m *MockClient) SAdd(key string, members ...interface{}) redis.Result {
	return m.Called(append([]interface{}{key}, members...)...).Get(0).(redis.Result)
}

// SCard implements redis.Client.
func (m *MockClient) SCard(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// SIsMember implements redis.Client.
func (m *MockClient) SIsMember(key string, member interface{}) redis.Result {
	return m.Called(key, member).Get(0).(redis.Result)
}

// SMembers implements redis.Client.
func (m *MockClient) SMembers(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// SRem implements redis.Client.
func (m *MockClient) SRem(key string, members ...interface{}) redis.Result {
	return m.Called(append([]interface{}{key}, members...)...).Get(0).(redis.Result)
}

// Scan implements redis.Client.
func (m *MockClient) Scan(cursor uint64, match string, count int64) redis.Result {
	return m.Called(cursor, match, count).Get(0).(redis.Result)
}

// Set implements redis.Client.
func (m *MockClient) Set(key string, value interface{}, expiration time.Duration) redis.Result {
	return m.Called(key, value, expiration).Get(0).(redis.Result)
}

// TTL implements redis.Client.
func (m *MockClient) TTL(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

// SetNX implements redis.Client.
func (m *MockClient) SetNX(key string, value interface{}, expiration time.Duration) redis.Result {
	return m.Called(key, value, expiration).Get(0).(redis.Result)
}

// Type implements redis.Client.
func (m *MockClient) Type(key string) redis.Result {
	return m.Called(key).Get(0).(redis.Result)
}

type MockPipeline struct {
	mock.Mock
}

// Decr implements redis.Pipeline.
func (m *MockPipeline) Decr(key string) {
	m.Called(key)
}

// Del implements redis.Pipeline.
func (m *MockPipeline) Del(keys ...string) {
	m.Called(common.AsInterfaceSlice(keys))
}

// Exec implements redis.Pipeline.
func (m *MockPipeline) Exec() ([]redis.Result, error) {
	args := m.Called()
	return args.Get(0).([]redis.Result), args.Error(1)
}

// HIncrBy implements redis.Pipeline.
func (m *MockPipeline) HIncrBy(key string, field string, value int64) {
	m.Called(key, field, value)
}

// HLen implements redis.Pipeline.
func (m *MockPipeline) HLen(key string) {
	m.Called(key)
}

// Incr implements redis.Pipeline.
func (m *MockPipeline) Incr(key string) {
	m.Called(key)
}

// LLen implements redis.Pipeline.
func (m *MockPipeline) LLen(key string) {
	m.Called(key)
}

// LRange implements redis.Pipeline.
func (m *MockPipeline) LRange(key string, start int64, stop int64) {
	m.Called(key, start, stop)
}

// LTrim implements redis.Pipeline.
func (m *MockPipeline) LTrim(key string, start int64, stop int64) {
	m.Called(key, start, stop)
}

// SAdd implements redis.Pipeline.
func (m *MockPipeline) SAdd(key string, members ...interface{}) {
	m.Called(append([]interface{}{key}, members...)...)
}

// SMembers implements redis.Pipeline.
func (m *MockPipeline) SMembers(key string) {
	m.Called(key)
}

// SRem implements redis.Pipeline.
func (m *MockPipeline) SRem(key string, members ...interface{}) {
	m.Called(append([]interface{}{key}, members...)...)
}

// Set implements redis.Pipeline.
func (m *MockPipeline) Set(key string, value interface{}, expiration time.Duration) {
	m.Called(key, value)
}

// SetNX implements redis.Pipeline.
func (m *MockPipeline) SetNX(key string, value interface{}, expiration time.Duration) {
	m.Called(key, value)
}

type MockResultOutput struct {
	mock.Mock
}

// Bool implements redis.Result.
func (m *MockResultOutput) Bool() bool {
	return m.Called().Bool(0)
}

// Duration implements redis.Result.
func (m *MockResultOutput) Duration() time.Duration {
	return m.Called().Get(0).(time.Duration)
}

// Err implements redis.Result.
func (m *MockResultOutput) Err() error {
	return m.Called().Error(0)
}

// Int implements redis.Result.
func (m *MockResultOutput) Int() int64 {
	return m.Called().Get(0).(int64)
}

func (m *MockResultOutput) String() string {
	return m.Called().Get(0).(string)
}

// MapStringString implements redis.Result.
func (m *MockResultOutput) MapStringString() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

// Multi implements redis.Result.
func (m *MockResultOutput) Multi() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// MultiInterface implements redis.Result.
func (m *MockResultOutput) MultiInterface() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

// Result implements redis.Result.
func (m *MockResultOutput) Result() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// ResultString implements redis.Result.
func (m *MockResultOutput) ResultString() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Val implements redis.Result.
func (m *MockResultOutput) Val() interface{} {
	return m.Called().Get(0)
}
