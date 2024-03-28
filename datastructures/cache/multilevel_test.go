package cache

import (
	"context"
	"testing"

	"github.com/splitio/go-toolkit/v6/datastructures/cache/mocks"
	"github.com/splitio/go-toolkit/v6/logging"
	"github.com/stretchr/testify/assert"
)

func TestMultiLevelCache(t *testing.T) {
	// To test this we setup 3 layers of caching in order of querying: top -> mid -> bottom
	// Top layer has key1, doesn't have key2 (returns Miss), has key3 expired and errors out when requesting any other Key
	// Mid layer has key 2, returns a Miss on any other key, and fails the test if key1 is fetched (because it was present on top layer)
	// Bottom layer fails if key1 or 2 are requested, has key 3. returns Miss if any other key is requested

	ctx := context.Background()

	topLayer := &mocks.LayerMock{}
	topLayer.On("Get", ctx, "key1").Once().Return("value1", nil)
	topLayer.On("Get", ctx, "key2").Once().Return("", &Miss{Where: "layer1", Key: "key2"})
	topLayer.On("Get", ctx, "key3").Once().Return("value1", &Expired{Key: "key3", Value: "someOtherValue"})
	topLayer.On("Get", ctx, "key4").Once().Return("", &Miss{Where: "layer1", Key: "key4"})
	topLayer.On("Set", ctx, "key2", "value2").Once().Return(nil)
	topLayer.On("Set", ctx, "key3", "value3").Once().Return(nil)

	midLayer := &mocks.LayerMock{}
	midLayer.On("Get", ctx, "key2").Once().Return("value2", nil)
	midLayer.On("Get", ctx, "key3").Once().Return("", &Miss{Where: "layer2", Key: "key3"}, nil)
	midLayer.On("Get", ctx, "key4").Once().Return("", &Miss{Where: "layer2", Key: "key4"})
	midLayer.On("Set", ctx, "key3", "value3").Once().Return(nil)

	bottomLayer := &mocks.LayerMock{}
	bottomLayer.On("Get", ctx, "key3").Once().Return("value3", nil)
	bottomLayer.On("Get", ctx, "key4").Once().Return("", &Miss{Where: "layer3", Key: "key4"})

	cacheML := MultiLevelCacheImpl[string, string]{
		logger: logging.NewLogger(nil),
		layers: []MLCLayer[string, string]{topLayer, midLayer, bottomLayer},
	}

	value1, err := cacheML.Get(ctx, "key1")
	assert.Nil(t, err)
	assert.Equal(t, "value1", value1)

	value2, err := cacheML.Get(ctx, "key2")
	assert.Nil(t, err)
	assert.Equal(t, "value2", value2)

	value3, err := cacheML.Get(ctx, "key3")
	assert.Nil(t, err)
	assert.Equal(t, "value3", value3)

	value4, err := cacheML.Get(ctx, "key4")
	assert.NotNil(t, err)
	asMiss, ok := err.(*Miss)
	assert.True(t, ok)
	assert.Equal(t, "ALL_LEVELS", asMiss.Where)
	assert.Equal(t, "key4", asMiss.Key)
	assert.Equal(t, "", value4)
}
