package model

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

type ContractInfo struct {
	Address string `json:"contract_address"`
	Name    string `json:"contract_name"`
}

func LoadAllContractInfos(d *bolt.DB) (map[string]ContractInfo, error) {
	wl := make(map[string]ContractInfo)

	err := d.View(func(tx *bolt.Tx) error {
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

func AddSCToWhitelist(d *bolt.DB, msg string) error {
	commands := strings.Split(msg, " ")
	if len(commands) != 3 {
		return errors.New("invalid command")
	}

	wl := &ContractInfo{
		Name:    commands[1],
		Address: commands[2],
	}

	// Serialize the Whitelist object to JSON
	wlBytes, err := json.Marshal(wl)
	if err != nil {
		return err
	}

	// Add the serialized Whitelist object to the whitelist bucket
	err = d.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sc_whitelist"))
		if b == nil {
			return errors.New("whitelist bucket not found")
		}
		return b.Put([]byte(commands[2]), wlBytes)
	})
	if err != nil {
		return err
	}

	return nil
}

func GetListSC(d *bolt.DB) (map[string]ContractInfo, error) {
	wl := make(map[string]ContractInfo)

	err := d.View(func(tx *bolt.Tx) error {
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
