package store

func NewMemoryStore() Store {
	return memoryStore{}
}

type memoryStore map[string]int

func (s memoryStore) Save(key string, value int) error {
	s[key] = value
	return nil
}

func (s memoryStore) Delete(key string) error {
	delete(s, key)
	return nil
}

func (s memoryStore) Find(key string) (*Result, error) {
	value, found := s[key]
	if !found {
		return nil, nil
	}
	return &Result{value: value}, nil
}

func (s memoryStore) Clone() (Store, error) {
	cpy := memoryStore{}
	for key, value := range s {
		cpy[key] = value
	}
	return cpy, nil
}

func (s memoryStore) Flush() error {
	return nil
}
