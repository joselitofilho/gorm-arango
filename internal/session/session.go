package session

import "sync"

var once sync.Once

// type global
type singleton map[string]string

var (
	instance singleton
)

func Session() singleton {
	once.Do(func() {
		instance = make(singleton)
	})
	return instance
}
