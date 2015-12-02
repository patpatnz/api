package apikey

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Key is the representation of an individual API Key, it supports having extra data
// stored in it which can be collected out using the various methods
type Key struct {
	Valid      bool
	Key        string
	pathLock   sync.RWMutex
	pathRegexp *regexp.Regexp
	paths      []string
	extra      map[string]interface{}
}

func (k *Key) CalculateRegexp() error {
	k.pathLock.Lock()
	defer k.pathLock.Unlock()

	r, err := regexp.Compile("^(" + strings.Replace(strings.Join(k.paths, "|"), "/", "\\/", -1) + ")")
	if err != nil {
		return err
	}
	k.pathRegexp = r
	return nil
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
		switch strings.ToLower(v) {
		case "true":
			return true, nil
		case "false":
			return false, nil
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

func (k *Key) Int(field string) (int, error) {
	if v, ok := k.extra[field].(int); ok {
		return v, nil
	}
	if v, ok := k.extra[field].(string); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Error converting %s to integer", v)
			return 0, ErrFieldDoesntExistOrNotType
		}
		return i, nil
	}
	return 0, ErrFieldDoesntExistOrNotType
}

func (k *Key) Strings(field string) ([]string, error) {
	if v, ok := k.extra[field].([]interface{}); ok {
		s := make([]string, len(v))
		for i, j := range v {
			s[i] = j.(string)
		}
		return s, nil
	}
	return nil, ErrFieldDoesntExistOrNotType
}
