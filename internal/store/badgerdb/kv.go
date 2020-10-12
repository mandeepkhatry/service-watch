package baderdb

import (
	"bytes"
	"os"
	"path/filepath"
	"service-watch/internal/def"

	"github.com/dgraph-io/badger"
)

type StoreClient struct {
	DB *badger.DB
}

func (s *StoreClient) NewClient(dir string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir = filepath.Join(pwd, dir)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	defaultOptions := badger.DefaultOptions(dir)
	db, err := badger.Open(defaultOptions)
	if err != nil {
		return err
	}

	s.DB = db
	return nil
}

//CloseClient closes badgerDB
func (s *StoreClient) CloseClient() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

//Put inserts key,val to badgerDB
func (s *StoreClient) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return def.KeyEmpty
	}

	txn := s.DB.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set(key, value)
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

//PutBatch inserts key,val pairs in batch
func (s *StoreClient) PutBatch(keys [][]byte, values [][]byte) error {
	//create a new write batch
	wb := s.DB.NewWriteBatch()
	defer wb.Cancel()

	for i := 0; i < len(keys); i++ {
		err := wb.Set(keys[i], values[i])
		if err != nil {
			return err
		}
	}
	err := wb.Flush()
	if err != nil {
		return err
	}
	return nil
}

//Get reads value for given key
func (s *StoreClient) Get(key []byte) ([]byte, error) {

	if len(key) == 0 {
		return []byte{}, def.KeyEmpty
	}
	value := make([]byte, 0)
	err := s.DB.View(func(txn *badger.Txn) error {

		item, err := txn.Get(key)

		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		value = val

		return nil

	})

	if err != nil {
		return []byte{}, err
	}

	return value, nil
}

//GetBatch retrieves values for given pair of keys in batch
func (s *StoreClient) GetBatch(keys [][]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return [][]byte{}, def.ResultsNotFound
	}

	values := make([][]byte, 0)
	for i := 0; i < len(keys); i++ {
		err := s.DB.View(func(txn *badger.Txn) error {
			item, err := txn.Get(keys[i])

			//if some key is not found, proceed to next one w/o throwing err
			if err == badger.ErrKeyNotFound {
				return nil
			}

			if err != nil {
				return err
			}

			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			values = append(values, val)

			return nil
		})
		if err != nil {
			return [][]byte{}, err
		}
	}
	return values, nil
}

//DeleteKey deletes given key from badgerDB
func (s *StoreClient) DeleteKey(key []byte) error {
	if len(key) == 0 {
		return def.EmptyKeyCannotBeDeleted
	}
	//delete given key
	txn := s.DB.NewTransaction(true)
	defer txn.Discard()

	err := txn.Delete(key)
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

//DeleteKeyRange deletes key,val pairs from startKey to endKey from badgerDB
func (s *StoreClient) DeleteKeyRange(startKey []byte, endKey []byte) error {
	if len(startKey) == 0 && len(endKey) == 0 {
		return def.StartOrEndKeyEmpty
	}
	err := s.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		wb := s.DB.NewWriteBatch()
		defer wb.Cancel()

		for it.Seek(startKey); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			err := wb.Delete(k)
			if err != nil {
				return err
			}

			if bytes.Compare(k, endKey) == 0 {
				break
			}

		}

		//delete in batch
		err := wb.Flush()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil

}

//Scan iterates from startKey to endKey for closed set [startKey,endKey] upto within limit
func (s *StoreClient) Scan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	counter := 0
	err := s.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(startKey); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			keys = append(keys, k)
			values = append(values, val)
			//include [startKey,endKey]
			if bytes.Compare(k, endKey) == 0 || counter > limit {
				break
			}
			counter += 1
		}
		return nil
	})

	if err != nil {
		return [][]byte{}, [][]byte{}, err
	}
	return keys, values, nil
}

//ReverseScan takes startKey, endKey and limit to scan in reverse direction in closed set [endKey,startKey]
//returns key,value,error
func (s *StoreClient) ReverseScan(startKey []byte, endKey []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	if len(startKey) == 0 {
		return [][]byte{}, [][]byte{}, def.StartKeyUnknown
	}

	opts := badger.DefaultIteratorOptions
	opts.Reverse = true
	counter := 0
	err := s.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(endKey); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			keys = append(keys, k)
			values = append(values, val)
			//include [endKey,startKey]
			if bytes.Compare(k, startKey) == 0 || counter > limit {
				break
			}
			counter += 1

		}
		return nil
	})

	if err != nil {
		return [][]byte{}, [][]byte{}, err
	}
	return keys, values, nil

}

//PrefixScan scans over [StartKey,endKey] for valid prefix upto limit
//if limit is zero, it returns whole set of [startKey,endKey]
func (s *StoreClient) PrefixScan(startKey []byte, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	opts := badger.DefaultIteratorOptions
	opts.Prefix = prefix

	if limit == 0 {
		err := s.DB.View(func(txn *badger.Txn) error {

			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Seek(startKey); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()

				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				keys = append(keys, k)
				values = append(values, val)
			}

			return nil

		})
		if err != nil {
			return [][]byte{}, [][]byte{}, err
		}
		return keys, values, nil

	} else {

		counter := 0
		err := s.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Seek(startKey); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				if counter > limit {
					break
				}
				keys = append(keys, k)
				values = append(values, val)
				counter += 1
			}
			return nil
		})

		if err != nil {
			return [][]byte{}, [][]byte{}, err
		}
		return keys, values, nil
	}
}

//ReversePrefixScan scans over [endKey,startKey] from reverse for valid prefix upto limit
//if limit is zero, it returns full result set
func (s *StoreClient) ReversePrefixScan(endKey []byte, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	opts := badger.DefaultIteratorOptions
	opts.Prefix = prefix
	opts.Reverse = true

	if limit == 0 {
		err := s.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Seek(endKey); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()

				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				keys = append(keys, k)
				values = append(values, val)

			}
			return nil
		})

		if err != nil {
			return [][]byte{}, [][]byte{}, err
		}
		return keys, values, nil
	} else {
		//if limit is not set to zero, scan in reverse for limit x
		counter := 0
		err := s.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Seek(endKey); it.Valid() && counter <= limit; it.Next() {
				item := it.Item()
				k := item.Key()
				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				keys = append(keys, k)
				values = append(values, val)
				counter += 1

			}
			return nil
		})

		if err != nil {
			return [][]byte{}, [][]byte{}, err
		}
		return keys, values, nil
	}

}

func (s *StoreClient) GetAllKeys() ([][]byte, error) {

	keys := make([][]byte, 0)

	err := s.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			keys = append(keys, k)
		}
		return nil
	})

	return keys, err
}

func (s *StoreClient) DeleteBatch(keys [][]byte) error {

	wb := s.DB.NewWriteBatch()
	defer wb.Cancel()

	for _, key := range keys {
		err := wb.Delete(key)
		if err != nil {
			return err
		}
	}
	err := wb.Flush()
	return err
}
