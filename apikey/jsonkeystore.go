package api

// JSONKeyStore implements the KeyStore interface and loads its data from the file
// provided when you instantate the type
type JSONKeyStore struct {
	filename string
	keys     map[string]*Key
	lock     sync.RWMutex
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
	j.lock.Lock()
	defer j.lock.Unlock()

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
			if field == "paths" {
				continue
			}
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
