package mocks

import (
	"time"

	"github.com/splitio/go-toolkit/v5/redis"
)

// MockResultOutput mocks struct
type MockResultOutput struct {
	ErrCall             func() error
	IntCall             func() int64
	StringCall          func() string
	BoolCall            func() bool
	DurationCall        func() time.Duration
	ResultCall          func() (int64, error)
	ResultStringCall    func() (string, error)
	MultiCall           func() ([]string, error)
	MultiInterfaceCall  func() ([]interface{}, error)
	MapStringStringCall func() (map[string]string, error)
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

// MapStringString mocks MapStringString
func (m *MockResultOutput) MapStringString() (map[string]string, error) {
	return m.MapStringStringCall()
}

// MpockPipeline  impl
type MockPipeline struct {
	LRangeCall   func(key string, start, stop int64)
	LTrimCall    func(key string, start, stop int64)
	LLenCall     func(key string)
	HIncrByCall  func(key string, field string, value int64)
	HLenCall     func(key string)
	SetCall      func(key string, value interface{}, expiration time.Duration)
	IncrCall     func(key string)
	DecrCall     func(key string)
	SAddCall     func(key string, members ...interface{})
	SRemCall     func(key string, members ...interface{})
	SMembersCall func(key string)
	DelCall      func(keys ...string)
	ExecCall     func() ([]redis.Result, error)
}

func (m *MockPipeline) LRange(key string, start, stop int64) {
	m.LRangeCall(key, start, stop)
}

func (m *MockPipeline) LTrim(key string, start, stop int64) {
	m.LTrimCall(key, start, stop)
}

func (m *MockPipeline) LLen(key string) {
	m.LLenCall(key)
}

func (m *MockPipeline) HIncrBy(key string, field string, value int64) {
	m.HIncrByCall(key, field, value)
}

func (m *MockPipeline) HLen(key string) {
	m.HLenCall(key)
}

func (m *MockPipeline) Set(key string, value interface{}, expiration time.Duration) {
	m.SetCall(key, value, expiration)
}

func (m *MockPipeline) Incr(key string) {
	m.IncrCall(key)
}

func (m *MockPipeline) Decr(key string) {
	m.DecrCall(key)
}

func (m *MockPipeline) SAdd(key string, members ...interface{}) {
	m.SAddCall(key, members...)
}

func (m *MockPipeline) SRem(key string, members ...interface{}) {
	m.SRemCall(key, members...)
}

func (m *MockPipeline) SMembers(key string) {
	m.SMembers(key)
}

func (m *MockPipeline) Del(keys ...string) {
	m.DelCall(keys...)
}

func (m *MockPipeline) Exec() ([]redis.Result, error) {
	return m.ExecCall()
}

// MockClient mocks for testing purposes
type MockClient struct {
	ClusterModeCall            func() bool
	ClusterCountKeysInSlotCall func(slot int) redis.Result
	ClusterSlotForKeyCall      func(key string) redis.Result
	ClusterKeysInSlotCall      func(slot int, count int) redis.Result
	DelCall                    func(keys ...string) redis.Result
	GetCall                    func(key string) redis.Result
	SetCall                    func(key string, value interface{}, expiration time.Duration) redis.Result
	PingCall                   func() redis.Result
	ExistsCall                 func(keys ...string) redis.Result
	KeysCall                   func(pattern string) redis.Result
	SMembersCall               func(key string) redis.Result
	SIsMemberCall              func(key string, member interface{}) redis.Result
	SAddCall                   func(key string, members ...interface{}) redis.Result
	SRemCall                   func(key string, members ...interface{}) redis.Result
	IncrCall                   func(key string) redis.Result
	DecrCall                   func(key string) redis.Result
	RPushCall                  func(key string, values ...interface{}) redis.Result
	LRangeCall                 func(key string, start, stop int64) redis.Result
	LTrimCall                  func(key string, start, stop int64) redis.Result
	LLenCall                   func(key string) redis.Result
	ExpireCall                 func(key string, value time.Duration) redis.Result
	TTLCall                    func(key string) redis.Result
	MGetCall                   func(keys []string) redis.Result
	SCardCall                  func(key string) redis.Result
	EvalCall                   func(script string, keys []string, args ...interface{}) redis.Result
	HIncrByCall                func(key string, field string, value int64) redis.Result
	HGetAllCall                func(key string) redis.Result
	HSetCall                   func(key string, hashKey string, value interface{}) redis.Result
	TypeCall                   func(key string) redis.Result
	PipelineCall               func() redis.Pipeline
}

func (m *MockClient) ClusterMode() bool {
	return m.ClusterModeCall()
}

func (m *MockClient) ClusterCountKeysInSlot(slot int) redis.Result {
	return m.ClusterCountKeysInSlotCall(slot)
}

func (m *MockClient) ClusterSlotForKey(key string) redis.Result {
	return m.ClusterSlotForKeyCall(key)
}

func (m *MockClient) ClusterKeysInSlot(slot int, count int) redis.Result {
	return m.ClusterKeysInSlotCall(slot, count)
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
	return m.SAddCall(key, members...)
}

// SRem mocks SRem
func (m *MockClient) SRem(key string, members ...interface{}) redis.Result {
	return m.SRemCall(key, members...)
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

// HIncrBy mocks HIncrByCall
func (m *MockClient) HIncrBy(key string, field string, value int64) redis.Result {
	return m.HIncrByCall(key, field, value)
}

// HGetAll mocks HGetAll
func (m *MockClient) HGetAll(key string) redis.Result {
	return m.HGetAllCall(key)
}

// HSet implements HGetAll wrapper for redis
func (m *MockClient) HSet(key string, hashKey string, value interface{}) redis.Result {
	return m.HSetCall(key, hashKey, value)
}

// Type implements Type wrapper for redis with prefix
func (m *MockClient) Type(key string) redis.Result {
	return m.TypeCall(key)
}

// Pipeline mock
func (m *MockClient) Pipeline() redis.Pipeline {
	return m.PipelineCall()
}
