package internal

import "errors"

type Cache interface {
	Save(key string, value string) error
	Get(key string) (string, error)
}

type KeyValueCache struct {
	cache map[string]string
}

func NewKeyValueCache() *KeyValueCache {
	return &KeyValueCache{
		cache: make(map[string]string),
	}
}

func (c *KeyValueCache) Save(key string, value string) error {
	c.cache[key] = value
	return nil
}

func (c *KeyValueCache) Get(key string) (string, error) {
	v, ok := c.cache[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return v, nil
}
