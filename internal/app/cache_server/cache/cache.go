package cache

import (
	"log"
)

type Cache interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Del(string) error
	GetStatus() Status
}

func New(typ string) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemCache()
	}

	if c == nil {
		panic("unknown cache type.")
	}
	log.Println(typ, "cache is ready.")
	return c
}
