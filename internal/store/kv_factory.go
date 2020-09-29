package store

import (
	baderdb "service-watch/internal/store/badgerdb"
)

var Stores = map[string]func(dir string) Store{
	"badgerdb": func(dir string) Store {
		return NewBadgerFactory(dir)
	},
}

func NewBadgerFactory(dbDIR string) Store {
	badger := &baderdb.StoreClient{}
	err := badger.NewClient(dbDIR)
	if err != nil {
		panic(err)
	}
	return badger
}
