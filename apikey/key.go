package apikey

// Key is the representation of an individual API Key, it supports having extra data
// stored in it which can be collected out using the various methods
type Key struct {
	Key        string
	pathRegexp *regexp.Regexp
	paths      []string
	extra      map[string]interface{}
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
	return 0, ErrFieldDoesntExistOrNotType
}
