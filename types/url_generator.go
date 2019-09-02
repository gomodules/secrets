package types

import (
	"sync"
)

type URLGenerator func() string

var (
	urlFn     URLGenerator
	rotateURL bool

	m sync.RWMutex
)

// Config should be called once at the start of a program to configure
// URLGenerator and secret rotation policy
func Config(fn URLGenerator, autoRotate bool) {
	m.Lock()
	urlFn = fn
	rotateURL = autoRotate
	m.Unlock()
}
