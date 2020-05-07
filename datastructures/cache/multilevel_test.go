package cache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/splitio/go-toolkit/logging"
)

type LayerMock struct {
	getCall func(key string) (interface{}, error)
	setCall func(key string, value interface{}) error
}

func (m *LayerMock) Get(key string) (interface{}, error) {
	return m.getCall(key)
}

func (m *LayerMock) Set(key string, value interface{}) error {
	return m.setCall(key, value)
}

type callTracker struct {
	calls map[string]int
	t     *testing.T
}

func newCallTracker(t *testing.T) *callTracker {
	return &callTracker{calls: make(map[string]int), t: t}
}

func (c *callTracker) track(name string) { c.calls[name]++ }

func (c *callTracker) reset() { c.calls = make(map[string]int) }

func (c *callTracker) checkCall(name string, count int) {
	if c.calls[name] != count {
		c.t.Errorf("calls for '%s' should be %d", name, count)
	}
}

func (c *callTracker) checkTotalCalls(count int) {
	if len(c.calls) != count {
		c.t.Errorf("The nomber of total calls should be '%d' and is '%d'", count, len(c.calls))
	}
}

func TestMultiLevelCache(t *testing.T) {
	// To test this we setup 3 layers of caching in order of querying: top -> mid -> bottom
	// Top layer has key1, doesn't have key2 (returns Miss), has key3 expired and errors out when requesting any other Key
	// Mid layer has key 2, returns a Miss on any other key, and fails the test if key1 is fetched (because it was present on top layer)
	// Bottom layer fails if key1 or 2 are requested, has key 3. returns Miss if any other key is requested
	calls := newCallTracker(t)
	topLayer := &LayerMock{
		getCall: func(key string) (interface{}, error) {
			calls.track(fmt.Sprintf("top:get:%s", key))
			switch key {
			case "key1":
				return "value1", nil
			case "key2":
				return nil, &Miss{Where: "layer1", Key: "key2"}
			case "key3":
				return nil, &Expired{Key: "key3", Value: "someOtherValue"}
			default:
				return nil, errors.New("someError")
			}
		},
		setCall: func(key string, value interface{}) error {
			calls.track(fmt.Sprintf("top:set:%s", key))
			switch key {
			case "key1":
				t.Error("Set should not be called on the top layer for key1")
				break
			case "key2":
				break
			case "key3":
				break
			default:
				return errors.New("someError")
			}
			return nil
		},
	}

	midLayer := &LayerMock{
		getCall: func(key string) (interface{}, error) {
			calls.track(fmt.Sprintf("mid:get:%s", key))
			switch key {
			case "key1":
				t.Error("Get should not be called on the mid layer for key1")
				return nil, nil
			case "key2":
				return "value2", nil
			default:
				return nil, &Miss{Where: "layer2", Key: key}
			}
		},
		setCall: func(key string, value interface{}) error {
			calls.track(fmt.Sprintf("mid:set:%s", key))
			switch key {
			case "key1":
				t.Error("Set should not be called on the mid layer for key1")
			case "key2":
				t.Error("Set should not be called on the mid layer for key2")
			case "key3":
			default:
				return errors.New("someError")
			}
			return nil
		},
	}

	bottomLayer := &LayerMock{
		getCall: func(key string) (interface{}, error) {
			calls.track(fmt.Sprintf("bot:get:%s", key))
			switch key {
			case "key1":
				t.Error("Get should not be called on the mid layer for key1")
				return nil, nil
			case "key2":
				t.Error("Get should not be called on the mid layer for key1")
				return nil, nil
			case "key3":
				return "value3", nil
			default:
				return nil, &Miss{Where: "layer3", Key: key}
			}
		},
		setCall: func(key string, value interface{}) error {
			calls.track(fmt.Sprintf("bot:set:%s", key))
			switch key {
			case "key1":
				t.Error("Set should not be called on the mid layer for key1")
			case "key2":
				t.Error("Set should not be called on the mid layer for key2")
			default:
				return errors.New("someError")
			}
			return nil
		},
	}

	cache := MultiLevelCacheImpl{
		logger: logging.NewLogger(nil),
		layers: []Layer{topLayer, midLayer, bottomLayer},
	}

	value1, err := cache.Get("key1")
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}
	if value1 != "value1" {
		t.Error("Get 'key1' should return 'value1'. Got: ", value1)
	}
	calls.checkCall("top:get:key1", 1)
	calls.checkTotalCalls(1)

	calls.reset()
	value2, err := cache.Get("key2")
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}
	if value2 != "value2" {
		t.Error("Get 'key2' should return 'value2'. Got: ", value2)
	}
	calls.checkCall("top:get:key2", 1)
	calls.checkCall("mid:get:key2", 1)
	calls.checkCall("top:set:key2", 1)
	calls.checkTotalCalls(3)

	calls.reset()
	value3, err := cache.Get("key3")
	if err != nil {
		t.Error("Error should be nil. Was: ", err)
	}

	if value3 != "value3" {
		t.Error("Get 'key3' should return 'value3'. Got: ", value3)
	}
	calls.checkCall("top:get:key3", 1)
	calls.checkCall("mid:get:key3", 1)
	calls.checkCall("bot:get:key3", 1)
	calls.checkCall("mid:set:key3", 1)
	calls.checkCall("top:set:key3", 1)
	calls.checkTotalCalls(5)

	calls.reset()
	value4, err := cache.Get("key4")
	if err == nil {
		t.Error("Error should be returned when getting nonexistant key.")
	}

	asMiss, ok := err.(*Miss)
	if !ok {
		t.Errorf("Error should be of Miss type. Is %T", err)
	}

	if asMiss.Where != "ALL_LEVELS" || asMiss.Key != "key4" {
		t.Errorf("Incorrect 'Where' or 'Key'. Got: %+v", asMiss)
	}

	if value4 != nil {
		t.Errorf("Value returned for GET 'key4' should be nil. Is: %+v", value4)
	}
	calls.checkCall("top:get:key4", 1)
	calls.checkCall("top:get:key4", 1)
	calls.checkCall("top:get:key4", 1)
	calls.checkTotalCalls(3)
}
