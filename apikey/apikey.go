package api

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
)

var (
	ErrFieldDoesntExistOrNotType = errors.New("Field does not exist or is incorrect type")
	ErrKeyRequired               = errors.New("Key is required")
	ErrKeyNotFound               = errors.New("Key not found")
	ErrKeyNotAuthorised          = errors.New("Key not authorised")
)

// KeyStore is an interface for the APIKey to a backend
type KeyStore interface {
	// Get returns a pointer to a Key type from the keystore
	Get(string) (*Key, error)
}

// APIKey is a object to manage the authorization of API requests
type APIKey struct {
	store      KeyStore
	pathRegexp *regexp.Regexp
	paths      []string
	// Required decided if the requestor MUST supply an apikey
	Required bool
}

// NewKey creates a new API Key management object using the provided
// KeyStore as a backend for storage
func NewAPIKey(ks KeyStore) (*APIKey, error) {
	return &APIKey{store: ks}, nil
}

// Middleware returns a function that can be passed to echo.Use
func (k *APIKey) Middleware() echo.Middleware {
	return func(ctx *echo.Context) error {
		var kv *Key
		var err error
		if key, ok := ctx.Request().Header["X-API-Key"]; ok {
			kv, err = k.store.Get(key[0])
			if err != nil {
				return nil
			}
			ctx.Set("apikey", kv)
			return nil
		}
		if key := ctx.Param("apikey"); key != "" {
			kv, err = k.store.Get(key)
			if err != nil {
				return nil
			}
			ctx.Set("apikey", kv)
			return nil
		}
		if k.pathRegexp != nil {
			if k.pathRegexp.MatchString(ctx.Request().URL.Path) {
				return nil
			}
		}
		if k.Required {
			ctx.JSON(416, ErrKeyRequired)
			return ErrKeyRequired
		}
		return nil
	}
}

func (k *APIKey) PublicPath(path string) error {
	k.paths = append(k.paths, path)
	r, err := regexp.Compile("^(" + strings.Replace(strings.Join(k.paths, "|"), "/", "\\/", -1) + ")")
	if err != nil {
		return err
	}
	k.pathRegexp = r
	return nil
}
