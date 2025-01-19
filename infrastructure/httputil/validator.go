package httputil

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	v    *validator.Validate
	once sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		v = validator.New()
	})
	return v
}
