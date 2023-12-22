package client

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
)

type Config struct {
	Username string
}

func SaveConfig(c Config) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	p := path.Join(dir, ".super-insecure")
	finfo, err := os.Stat(p)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		} else {
			err = os.Mkdir(p, 0700)
			if err != nil {
				return err
			}
		}
	} else {
		if !finfo.IsDir() {
			return errors.New("file in the way")
		}
	}

	f, err := os.Create(path.Join(p, "config.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(c)
}

func LoadConfig() (*Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	p := path.Join(dir, ".super-insecure/config.json")
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
