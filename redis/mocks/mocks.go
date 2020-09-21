package mocks

import (
	"time"

	"github.com/splitio/go-toolkit/v3/redis"
)

// MockResultOutput mocks struct
type MockResultOutput struct {
	ErrCall            func() error
	IntCall            func() int64
	StringCall         func() string
	BoolCall           func() bool
	DurationCall       func() time.Duration
	ResultCall         func() (int64, error)
	ResultStringCall   func() (string, error)
	MultiCall          func() ([]string, error)
	MultiInterfaceCall func() ([]interface{}, error)
}

// Int mocks Int
func (m *MockResultOutput) Int() int64 {
	return m.IntCall()
}

// Err mocks Err
func (m *MockResultOutput) Err() error {
	return m.ErrCall()
}

// String mocks String
func (m *MockResultOutput) String() string {
	return m.StringCall()
}

// Bool mocks Bool
func (m *MockResultOutput) Bool() bool {
	return m.BoolCall()
}

// Duration mocks Duration
func (m *MockResultOutput) Duration() time.Duration {
	return m.DurationCall()
}

// Result mocks Result
func (m *MockResultOutput) Result() (int64, error) {
	return m.ResultCall()
}

// ResultString mocks ResultString
func (m *MockResultOutput) ResultString() (string, error) {
	return m.ResultStringCall()
}

// Multi mocks Multi
func (m *MockResultOutput) Multi() ([]string, error) {
	return m.MultiCall()
}

// MultiInterface mocks MultiInterface
func (m *MockResultOutput) MultiInterface() ([]interface{}, error) {
	return m.MultiInterfaceCall()
}

// MockClient mocks for testing purposes
type MockClient struct {
	DelCall       func(keys ...string) redis.Result
	GetCall       func(key string) redis.Result
	SetCall       func(key string, value interface{}, expiration time.Duration) redis.Result
	PingCall      func() redis.Result
	ExistsCall    func(keys ...string) redis.Result
	KeysCall      func(pattern string) redis.Result
	SMembersCall  func(key string) redis.Result
	SIsMemberCall func(key string, member interface{}) redis.Result
	SAddCall      func(key string, members ...interface{}) redis.Result
	SRemCall      func(key string, members ...interface{}) redis.Result
	IncrCall      func(key string) redis.Result
	DecrCall      func(key string) redis.Result
	RPushCall     func(key string, values ...interface{}) redis.Result
	LRangeCall    func(key string, start, stop int64) redis.Result
	LTrimCall     func(key string, start, stop int64) redis.Result
	LLenCall      func(key string) redis.Result
	ExpireCall    func(key string, value time.Duration) redis.Result
	TTLCall       func(key string) redis.Result
	MGetCall      func(keys []string) redis.Result
	SCardCall     func(key string) redis.Result
	EvalCall      func(script string, keys []string, args ...interface{}) redis.Result
}

// Del mocks get
func (m *MockClient) Del(keys ...string) redis.Result {
	return m.DelCall(keys...)
}

// Get mocks get
func (m *MockClient) Get(key string) redis.Result {
	return m.GetCall(key)
}

// Set mocks set
func (m *MockClient) Set(key string, value interface{}, expiration time.Duration) redis.Result {
	return m.SetCall(key, value, expiration)
}

// Exists mocks set
func (m *MockClient) Exists(keys ...string) redis.Result {
	return m.ExistsCall(keys...)
}

// Ping mocks ping
func (m *MockClient) Ping() redis.Result {
	return m.PingCall()
}

// Keys mocks keys
func (m *MockClient) Keys(pattern string) redis.Result {
	return m.KeysCall(pattern)
}

// SMembers mocks SMembers
func (m *MockClient) SMembers(key string) redis.Result {
	return m.SMembersCall(key)
}

// SIsMember mocks SIsMember
func (m *MockClient) SIsMember(key string, member interface{}) redis.Result {
	return m.SIsMemberCall(key, member)
}

// SAdd mocks SAdd
func (m *MockClient) SAdd(key string, members ...interface{}) redis.Result {
	return m.SAddCall(key)
}

// SRem mocks SRem
func (m *MockClient) SRem(key string, members ...interface{}) redis.Result {
	return m.SRemCall(key)
}

// Incr mocks Incr
func (m *MockClient) Incr(key string) redis.Result {
	return m.IncrCall(key)
}

// Decr mocks Decr
func (m *MockClient) Decr(key string) redis.Result {
	return m.DecrCall(key)
}

// RPush mocks RPush
func (m *MockClient) RPush(key string, values ...interface{}) redis.Result {
	return m.RPushCall(key, values...)
}

// LRange mocks LRange
func (m *MockClient) LRange(key string, start, stop int64) redis.Result {
	return m.LRangeCall(key, start, stop)
}

// LTrim mocks LTrim
func (m *MockClient) LTrim(key string, start, stop int64) redis.Result {
	return m.LTrimCall(key, start, stop)
}

// LLen mocks LLen
func (m *MockClient) LLen(key string) redis.Result {
	return m.LLenCall(key)
}

// Expire mocks Expire
func (m *MockClient) Expire(key string, value time.Duration) redis.Result {
	return m.ExpireCall(key, value)
}

// TTL mocks TTL
func (m *MockClient) TTL(key string) redis.Result {
	return m.TTLCall(key)
}

// MGet mocks MGet
func (m *MockClient) MGet(keys []string) redis.Result {
	return m.MGetCall(keys)
}

// SCard mocks SCard
func (m *MockClient) SCard(key string) redis.Result {
	return m.SCardCall(key)
}

// Eval mocks Eval
func (m *MockClient) Eval(script string, keys []string, args ...interface{}) redis.Result {
	return m.EvalCall(script, keys, args...)
}
