package store

type Store interface {
	Save(key string, value int) error
	Delete(key string) error
	Find(key string) (*Result, error)
	Flush() error
	Clone() (Store, error)
}
