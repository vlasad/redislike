package cache

type CacheError string

func (e CacheError) Error() string {
	return string(e)
}

const (
	ErrorInvalidTTL  = CacheError("invalid ttl value")
	ErrorKeyNotFound = CacheError("key not found")
	ErrorWrongType   = CacheError("wrong type")
	ErrorNoItems     = CacheError("no items")
)
