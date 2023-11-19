package db

import (
	"github.com/boltdb/bolt"
)

type Database struct {
	Db *bolt.DB
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

	// Create a bucket for the tracker check points
	err = dbBolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tracker_check_points"))
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	return &Database{Db: dbBolt}, closeFunc, nil
}

// func (d *Database) LoadAllFromDB() (map[string]model.ContractInfo, error) {
// 	var wl map[string]model.ContractInfo

// 	err := d.Db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("sc_whitelist"))
// 		if b == nil {
// 			return errors.New("whitelist bucket not found")
// 		}

// 		err := b.ForEach(func(k, v []byte) error {
// 			var ci model.ContractInfo
// 			err := json.Unmarshal(v, &ci)
// 			if err != nil {
// 				return err
// 			}
// 			wl[string(k)] = ci
// 			return nil
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return wl, nil
// }

// func (d *Database) GetContractInfo(address string) (*model.ContractInfo, error) {
// 	var ci model.ContractInfo

// 	err := d.Db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("sc_whitelist"))
// 		if b == nil {
// 			return errors.New("whitelist bucket not found")
// 		}

// 		v := b.Get([]byte(address))
// 		if v == nil {
// 			return errors.New("contract info not found")
// 		}

// 		err := json.Unmarshal(v, &ci)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ci, nil
// }

// func (d *Database) DeleteContractInfo(address string) error {
// 	err := d.Db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("sc_whitelist"))
// 		if b == nil {
// 			return errors.New("whitelist bucket not found")
// 		}

// 		err := b.Delete([]byte(address))
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
