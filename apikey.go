package api

import (
	"github.com/labstack/echo"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
)

var (
	ErrFieldDoesntExistOrNotType = errors.New("Field does not exist or is incorrect type")
)

type KeyStore interface {
	Check(string, *echo.Context) (bool)
}

type APIKey struct {
	store		KeyStore
}

func NewKey(ks KeyStore) (*APIKey, error) {
	return &APIKey{store: ks}, nil
}

func (k *APIKey) Middleware() echo.Middleware {
	return func(ctx *echo.Context) error {
		return k.authorize(ctx)
	}
}

func (k *APIKey) authorize(ctx *echo.Context) error {
	return nil
}

type Key struct {
	Key		string
	Paths	[]string
	extra	map[string]interface{}
}

func (k *Key) GetField(field string) (interface{}, error) {
	if v, ok := k.extra[field]; ok {
		return v, nil
	}
	return false, ErrFieldDoesntExistOrNotType
}

func (k *Key) GetBool(field string) (bool, error) {
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

type JSONKeyStore struct {
	filename		string
	keys			map[string]*Key
}

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

func (j *JSONKeyStore) Check(key string, ctx *echo.Context) bool {
	return false
}
