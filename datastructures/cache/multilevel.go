package cache

import (
	"github.com/splitio/go-toolkit/logging"
)

// Layer is the interface that should be implemented for all caching structs to be used with this piece of code.
type Layer interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

// MultiLevelCache bundles a list of ordered cache layers (upper -> lower)
type MultiLevelCache interface {
	Get(key string) (interface{}, error)
}

// MultiLevelCacheImpl implements the MultiLevelCache interface
type MultiLevelCacheImpl struct {
	layers []Layer
	logger logging.LoggerInterface
}

// Get returns the value of the requested key (if found) and populates upper levels with it
func (c *MultiLevelCacheImpl) Get(key string) (interface{}, error) {
	toUpdate := make([]int, 0, len(c.layers))
	var toReturn interface{}
	for index, layer := range c.layers {
		item, err := layer.Get(key)
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
			toReturn = item
			break
		}
	}

	if toReturn == nil {
		return nil, &Miss{Where: "ALL_LEVELS", Key: key}
	}

	// Update upper layers if any
	for _, index := range toUpdate {
		if index < len(c.layers) { // Ignore any awkward index (if any)
			err := c.layers[index].Set(key, toReturn)
			if err != nil {
				c.logger.Error(err)
			}
		}
	}
	return toReturn, nil
}

// NewMultiLevel creates and returns a new MultiLevelCache instance
func NewMultiLevel(layers []Layer, logger logging.LoggerInterface) (*MultiLevelCacheImpl, error) {
	if logger == nil {
		logger = logging.NewLogger(nil)
	}

	return &MultiLevelCacheImpl{layers: layers, logger: logger}, nil
}
