package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

func UnmarshalMemory(data []byte) (Configs, error) {
	cs := make(Configs)
	err := toml.Unmarshal(data, &cs)
	if err != nil {
		return nil, err
	} else {
		return cs, nil
	}
}

func UnmarshalFile(filename string) (Configs, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return UnmarshalMemory(data)
}
