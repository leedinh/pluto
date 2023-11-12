package sc

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type Database struct {
	db *bolt.DB
}

type ContractInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func InitDB() (db *Database, closeFunc func() error, err error) {
	dbBolt, err := bolt.Open("sc.db", 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	closeFunc = dbBolt.Close

	// Create a bucket for the whitelist smart contract
	err = dbBolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("sc_whitelist"))
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	return &Database{db: dbBolt}, closeFunc, nil
}

func (d *Database) AddToWhitelist(name string, address string) error {
	// Create a new Whitelist object
	wl := &ContractInfo{
		Name:    name,
		Address: address,
	}

	// Serialize the Whitelist object to JSON
	wlBytes, err := json.Marshal(wl)
	if err != nil {
		return err
	}

	// Add the serialized Whitelist object to the whitelist bucket
	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sc_whitelist"))
		if b == nil {
			return errors.New("whitelist bucket not found")
		}
		return b.Put([]byte(address), wlBytes)
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) LoadAllFromDB() (map[string]ContractInfo, error) {
	var wl map[string]ContractInfo

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sc_whitelist"))
		if b == nil {
			return errors.New("whitelist bucket not found")
		}

		err := b.ForEach(func(k, v []byte) error {
			var ci ContractInfo
			err := json.Unmarshal(v, &ci)
			if err != nil {
				return err
			}
			wl[string(k)] = ci
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wl, nil
}

func (d *Database) GetContractInfo(address string) (*ContractInfo, error) {
	var ci ContractInfo

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sc_whitelist"))
		if b == nil {
			return errors.New("whitelist bucket not found")
		}

		v := b.Get([]byte(address))
		if v == nil {
			return errors.New("contract info not found")
		}

		err := json.Unmarshal(v, &ci)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ci, nil
}

func (d *Database) DeleteContractInfo(address string) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sc_whitelist"))
		if b == nil {
			return errors.New("whitelist bucket not found")
		}

		err := b.Delete([]byte(address))
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
