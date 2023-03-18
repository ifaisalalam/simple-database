package db

import (
	"sync"

	"github.com/ifaisalalam/simple-database/internal/store"
)

var defaultDatabase = NewDatabaseWithMemoryStore

func NewDatabase() Database {
	return defaultDatabase()
}

func NewDatabaseWithMemoryStore() Database {
	memoryStore := store.NewMemoryStore()
	return NewDatabaseWithStore(memoryStore)
}

func NewDatabaseWithStore(store store.Store) Database {
	return &db{
		store: store,
	}
}

type Database interface {
	// Set the key to given value
	Set(key string, value int)

	// Get the value for the given key, set 'ok' to true if key exists
	Get(key string) (value int, ok bool)

	// Unset the key, making it just like that key was never set
	Unset(key string)

	// Begin opens a new transaction
	Begin()

	// Commit closes all open transaction blocks, permanently apply the
	// changes made in them.
	Commit() error

	// Rollback undoes all the commands issued in the most recent
	// transaction block, and closes the block.
	Rollback() error
}

type db struct {
	sync.Mutex

	store       store.Store
	transaction *db
}

// Set the key to given value
func (d *db) Set(key string, value int) {
	d.Lock()
	defer d.Unlock()

	if d.transaction != nil {
		d.transaction.Set(key, value)
		return
	}

	_ = d.store.Save(key, value)
}

// Get the value for the given key, set 'ok' to true if key exists
func (d *db) Get(key string) (value int, ok bool) {
	d.Lock()
	defer d.Unlock()

	if d.transaction != nil {
		value, ok = d.transaction.Get(key)
		return
	}

	if result, err := d.store.Find(key); nil == err && nil != result {
		value, ok = result.GetValue().(int)
	}

	return
}

// Unset the key, making it just like that key was never set
func (d *db) Unset(key string) {
	d.Lock()
	defer d.Unlock()

	if d.transaction != nil {
		d.transaction.Unset(key)
		return
	}

	_ = d.store.Delete(key)
}

// Begin opens a new transaction
func (d *db) Begin() {
	d.Lock()
	defer d.Unlock()

	if d.transaction != nil {
		d.transaction.Begin()
		return
	}

	clonedStore, _ := d.store.Clone()
	d.transaction = NewDatabaseWithStore(clonedStore).(*db)
}

// Commit closes all open transaction blocks, permanently apply the
// changes made in them.
func (d *db) Commit() error {
	d.Lock()
	defer d.Unlock()

	if d.transaction == nil {
		return ErrNoActiveTransaction
	}

	transaction := d.transaction
	for transaction.transaction != nil {
		transaction = transaction.transaction
	}

	d.store = transaction.GetStore()
	d.transaction = nil

	return d.store.Flush()
}

// Rollback undoes all the commands issued in the most recent
// transaction block, and closes the block.
func (d *db) Rollback() error {
	d.Lock()
	defer d.Unlock()

	if d.transaction == nil {
		return ErrNoActiveTransaction
	}

	parent := d
	transaction := d.transaction
	for transaction.transaction != nil {
		parent = transaction
		transaction = transaction.transaction
	}

	transaction = nil
	parent.transaction = nil

	return nil
}
