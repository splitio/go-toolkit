package cache

import (
	"errors"
	"github.com/splitio/go-toolkit/logging"
	"testing"
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

func TestMultiLevelCache(t *testing.T) {
	calls := map[string]struct{}{}
	topLayer := &LayerMock{
		getCall: func(key string) (interface{}, error) {
			switch key {
			case "key1":
				return "value1", nil
			case "key2":
				return nil, &Miss{Where: "layer1", Key: "key2"}
			default:
				return nil, errors.New("someError")
			}
		},
		setCall: func(key string, value interface{}) error {
			switch key {
			case "key1":
				t.Error("Set should not be called on the top layer for key1")
			case "key2":
				calls["top:set:key2"] = struct{}{}
			case "key3":
				t.Error("Set should not be called on the top layer for key3")
			default:
				return errors.New("someError")

			}
			return nil
		},
	}

	midLayer := &LayerMock{
		getCall: func(key string) (interface{}, error) {
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
			switch key {
			case "key1":
				t.Error("Set should not be called on the mid layer for key1")
			case "key2":
				t.Error("Set should not be called on the mid layer for key2")
			case "key3":
				calls["mid:set:key3"] = struct{}{}
			default:
				return errors.New("someError")
			}
			return nil
		},
	}

	bottomLayer := &LayerMock{
		getCall: func(key string) (interface{}, error) {
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

	value2, err := cache.Get("key2")
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}
	if value2 != "value2" {
		t.Error("Get 'key2' should return 'value2'. Got: ", value2)
	}

	value3, err := cache.Get("key3")
	if err != nil {
		t.Error("Error should be nil. Was: ", err)
	}

	if value3 != "value3" {
		t.Error("Get 'key3' should return 'value3'. Got: ", value3)
	}

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

	if _, ok := calls["top:set:key2"]; !ok {
		t.Error("Top layer should have executed Set operation for key2")
	}

	if _, ok := calls["mid:set:key3"]; !ok {
		t.Error("Mid layer should have executed Set operation for key3")
	}

	if len(calls) != 2 {
		t.Errorf("There should only be 2 calls of set operations. Got: %+v", calls)
	}
}
