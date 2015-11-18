package api

import (
	"github.com/labstack/echo"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"sync"
)

var (
	ErrFieldDoesntExistOrNotType = errors.New("Field does not exist or is incorrect type")
	ErrKeyRequired = errors.New("Key is required")
	ErrKeyNotFound = errors.New("Key not found")
	ErrKeyNotAuthorised = errors.New("Key not authorised")
)

// KeyStore is an interface for the APIKey to a backend
type KeyStore interface {
	// Get returns a pointer to a Key type from the keystore
	Get(string) (*Key, error)
}

// APIKey is a object to manage the authorization of API requests
type APIKey struct {
	store		KeyStore
	// Required decided if the requestor MUST supply an apikey
	Required	bool
}

// NewKey creates a new API Key management object using the provided
// KeyStore as a backend for storage
func NewAPIKey(ks KeyStore) (*APIKey, error) {
	return &APIKey{store: ks}, nil
}

// Middleware returns a function that can be passed to echo.Use
func (k *APIKey) Middleware() echo.Middleware {
	return func(ctx *echo.Context) error {
		if key, ok := ctx.Request().Header["X-API-Key"]; ok {
			kv, err := k.store.Get(key[0])
			if err != nil {
				return nil
			}
			ctx.Set("apikey", kv)
		}
		if key := ctx.Param("apikey"); key != "" {
			kv, err := k.store.Get(key)
			if err != nil {
				return nil
			}
			ctx.Set("apikey", kv)
		}
		return ErrKeyRequired
	}
}

// Key is the representation of an individual API Key, it supports having extra data
// stored in it which can be collected out using the various methods
type Key struct {
	Key		string
	Paths	[]string
	extra	map[string]interface{}
}

func (k *Key) Field(field string) (interface{}, error) {
	if v, ok := k.extra[field]; ok {
		return v, nil
	}
	return false, ErrFieldDoesntExistOrNotType
}

func (k *Key) Bool(field string) (bool, error) {
	if v, ok := k.extra[field].(bool); ok {
		return v, nil
	}
	if v, ok := k.extra[field].(string); ok {
		switch (strings.ToLower(v)) {
		case "true": return true, nil
		case "false": return false, nil
		}
	}
	return false, ErrFieldDoesntExistOrNotType
}

func (k *Key) String(field string) (string, error) {
	if v, ok := k.extra[field].(string); ok {
		return v, nil
	}
	return "", ErrFieldDoesntExistOrNotType
}

// JSONKeyStore implements the KeyStore interface and loads its data from the file
// provided when you instantate the type
type JSONKeyStore struct {
	filename		string
	keys			map[string]*Key
	lock			sync.RWMutex
}

// NewJSONKeyStore creates a new instance of JSONKeyStore and attempts to load the data
// from the filename provided. If it fails to load key data it returns an error
func NewJSONKeyStore(filename string) (KeyStore, error) {
	jks := &JSONKeyStore{filename: filename, keys: make(map[string]*Key)}

	err := jks.loadData()
	if err != nil {
		return nil, err
	}

	return jks, nil
}

func (j *JSONKeyStore) loadData() error {
	buf, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, j.keys)
	if err != nil {
		return err
	}
	ex := make(map[string]map[string]interface{})
	err = json.Unmarshal(buf, ex)
	if err != nil {
		return err
	}
	for key, val := range ex {
		for field, dat := range val {
			if field == "paths" { continue }
			j.keys[key].extra[field] = dat
		}
	}
	return nil
}

func (j *JSONKeyStore) Get(key string) (*Key, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	if v, ok := j.keys[key]; ok {
		return v, nil
	}
	
	return nil, ErrKeyNotFound
}
