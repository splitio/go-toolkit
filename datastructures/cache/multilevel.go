package cache

import (
	"context"

	"github.com/splitio/go-toolkit/v6/logging"
)

// MLCLayer is the interface that should be implemented for all caching structs to be used with this piece of code.
type MLCLayer[K comparable, V comparable] interface {
	Get(ctx context.Context, key K) (V, error)
	Set(ctx context.Context, key K, value V) error
}

// MultiLevelCache bundles a list of ordered cache layers (upper -> lower)
type MultiLevelCache[K comparable, V comparable] interface {
	Get(ctx context.Context, key K) (V, error)
}

// MultiLevelCacheImpl implements the MultiLevelCache interface
type MultiLevelCacheImpl[K comparable, V comparable] struct {
	layers []MLCLayer[K, V]
	logger logging.LoggerInterface
}

// Get returns the value of the requested key (if found) and populates upper levels with it
func (c *MultiLevelCacheImpl[K, V]) Get(ctx context.Context, key K) (V, error) {
	toUpdate := make([]int, 0, len(c.layers))
	var item V
	var err error
	for index, layer := range c.layers {
		item, err = layer.Get(ctx, key)
		if err != nil {
			switch err.(type) {
			case *Miss:
				// If it's a miss, push the index to the stack, if we eventually find the item,
				// upper layers will be updated.
				toUpdate = append(toUpdate, index)
			case *Expired:
				// If the key is expired, push the index to the stack, if we eventually find the item,
				// upper layers will be updated.
				toUpdate = append(toUpdate, index)
			default:
				// Any other error implies simply skipping this layer.
				c.logger.Error(err)
			}
		} else {
			break
		}
	}

    var empty V
	if item == empty || err != nil {
		return empty, &Miss{Where: "ALL_LEVELS", Key: key}
	}

	// Update upper layers if any
	for _, index := range toUpdate {
		if index < len(c.layers) { // Ignore any awkward index (if any)
			err := c.layers[index].Set(ctx, key, item)
			if err != nil {
				c.logger.Error(err)
			}
		}
	}
	return item, nil
}

// NewMultiLevel creates and returns a new MultiLevelCache instance
func NewMultiLevel[K comparable, V comparable](layers []MLCLayer[K, V], logger logging.LoggerInterface) (*MultiLevelCacheImpl[K, V], error) {
	if logger == nil {
		logger = logging.NewLogger(nil)
	}

	return &MultiLevelCacheImpl[K, V]{layers: layers, logger: logger}, nil
}
