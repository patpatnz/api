// package apikey is an API key manager middleware for echo
package apikey

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/patpatnz/api/result"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	ErrFieldDoesntExistOrNotType = errors.New("Field does not exist or is incorrect type")
	ErrKeyRequired               = errors.New("Key is required")
	ErrKeyNotFound               = errors.New("Key not found")
	ErrKeyNotAuthorised          = errors.New("Key not authorised")
	ErrKeyPathNotAllowed         = errors.New("Key not authorised for this path")

	InvalidKey = &Key{Valid: false}
)

// KeyStore is an interface for the APIKey to a backend
type KeyStore interface {
	// Get returns a pointer to a Key type from the keystore
	Get(string) (*Key, error)
}

// APIKey is a object to manage the authorization of API requests
type APIKey struct {
	store            KeyStore
	pathMutex        sync.RWMutex
	publicPathRegexp *regexp.Regexp
	publicPaths      []string
	// Required decided if the requestor MUST supply an apikey
	Required bool
}

// NewKey creates a new API Key management object using the provided
// KeyStore as a backend for storage
func NewAPIKey(ks KeyStore) (*APIKey, error) {
	return &APIKey{store: ks}, nil
}

// Middleware returns a function that can be passed to, eg: echo.Use(ak.Middleware())
func (k *APIKey) Middleware() echo.Middleware {
	return func(ctx *echo.Context) error {
		var kv *Key
		var err error

		// Get key from X-API-Key header
		if key, ok := ctx.Request().Header["X-Api-Key"]; ok {
			kv, err = k.store.Get(key[0])
			if err != nil && err != ErrKeyNotFound {
				return ctx.JSON(500, "Error")
			}
		}

		// Get key from APIKEY query param
		if key := ctx.Query("apikey"); kv == nil && key != "" {
			log.Printf("Wat?")
			kv, err = k.store.Get(key)
			if err != nil && err != ErrKeyNotFound {
				return ctx.JSON(500, "Error")
			}
		}

		if kv != nil {
			ctx.Set("apikey", kv)
		} else {
			ctx.Set("apikey", InvalidKey)
		}

		k.pathMutex.RLock()

		// if this has publicly allowed paths then match and return neg if no match
		if kv == nil && k.publicPathRegexp != nil {
			if k.publicPathRegexp.MatchString(ctx.Request().URL.Path) && k.Required {
				k.pathMutex.RUnlock()
				return ctx.JSON(403, ErrKeyRequired)
			}
		}

		k.pathMutex.RUnlock()

		// if key is nil and keys are required then negatory it
		if kv == nil && k.Required {
			return result.Error(ctx, 403, ErrKeyRequired)
		}

		// skip key checks
		if kv == nil {
			return nil
		}

		// if key has paths set, check them
		if kv.pathRegexp != nil {
			if !kv.pathRegexp.MatchString(ctx.Request().URL.Path) {
				return ctx.JSON(412, ErrKeyPathNotAllowed)
			}
		}

		return nil
	}
}

func (k *APIKey) PublicPath(path string) error {
	k.pathMutex.Lock()
	defer k.pathMutex.Unlock()

	k.publicPaths = append(k.publicPaths, path)
	r, err := regexp.Compile("^(" + strings.Replace(strings.Join(k.publicPaths, "|"), "/", "\\/", -1) + ")")
	if err != nil {
		return err
	}
	k.publicPathRegexp = r
	return nil
}
