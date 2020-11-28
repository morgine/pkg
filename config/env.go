// Copyright 2019 morgine.com. All rights reserved.

package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strconv"
)

// OsEnvGetter provided OS environment values
var OsEnvGetter = func(namespace, key string) string {
	if namespace != "" {
		key = namespace + "." + key
	}
	return os.Getenv(key)
}

type LoadOsError struct {
	Namespace string
	Key       string
	err       error
}

func (e *LoadOsError) Error() string {
	if e.Namespace != "" {
		return fmt.Sprintf("load OS envirenment %s.%s error: %s", e.Namespace, e.Key, e.err)
	} else {
		return fmt.Sprintf("load OS envirenment %s error: %s", e.Key, e.err)
	}
}

type Configs map[string]interface{}

// Unmarshal for decoding Configs data to schema
func (cs Configs) Unmarshal(schema interface{}) error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(cs)
	if err != nil {
		return err
	}
	return toml.Unmarshal(buf.Bytes(), schema)
}

// UnmarshalSub is shorthand for GetSub and Unmarshal
func (cs Configs) UnmarshalSub(namespace string, schema interface{}) error {
	envs, err := cs.GetSub(namespace)
	if err != nil {
		return err
	}
	return envs.Unmarshal(schema)
}

// GetSub returns sub environments and load os environment values
func (cs Configs) GetSub(namespace string) (env Configs, err error) {
	subs := cs[namespace]
	if subs == nil {
		return make(Configs), nil
	}

	switch tv := subs.(type) {
	case map[string]interface{}:
		env = tv
		err = env.LoadOSEnv(namespace)
		if err != nil {
			return nil, err
		}
		return env, nil
	default:
		return nil, fmt.Errorf("config.GetSub: need object data, got: %t", tv)
	}
}

func (cs Configs) Len() int {
	return len(cs)
}

func (cs Configs) GetStr(name string) string {
	return cs[name].(string)
}

func (cs Configs) GetInt(name string) int {
	return cs[name].(int)
}

func (cs Configs) GetFloat(name string) float64 {
	return cs[name].(float64)
}

func (cs Configs) GetBool(name string) bool {
	return cs[name].(bool)
}

func (cs Configs) GetSliceStr(name string) []string {
	return cs[name].([]string)
}

func (cs Configs) GetSliceInt(name string) []int {
	return cs[name].([]int)
}

func (cs Configs) GetSliceFloat(name string) []float64 {
	return cs[name].([]float64)
}

func (cs Configs) GetSliceBool(name string) []bool {
	return cs[name].([]bool)
}

// LoadOSEnv for loading the OS environment values, if namespace
func (cs Configs) LoadOSEnv(namespace string) (err error) {
	for key, value := range cs {
		ov := OsEnvGetter(namespace, key)
		if ov != "" {
			switch vt := value.(type) {
			case string:
				cs[key] = ov
			case int:
				cs[key], err = strconv.Atoi(ov)
				if err != nil {
					return &LoadOsError{
						Namespace: namespace,
						Key:       key,
						err:       err,
					}
				}
			case bool:
				switch ov {
				case "y", "Y", "yes", "YES", "Yes", "1", "t", "T", "true", "TRUE", "True":
					cs[key] = true
				case "n", "N", "no", "NO", "No", "0", "f", "F", "false", "FALSE", "False":
					cs[key] = false
				default:
					return &LoadOsError{
						Namespace: namespace,
						Key:       key,
						err:       errors.New("need a boolean argument"),
					}
				}
			case float64:
				cs[key], err = strconv.ParseFloat(ov, 64)
				if err != nil {
					return &LoadOsError{
						Namespace: namespace,
						Key:       key,
						err:       err,
					}
				}
			default:
				return &LoadOsError{
					Namespace: namespace,
					Key:       key,
					err:       fmt.Errorf("not support the %v argument", vt),
				}
			}
		}
	}
	return nil
}
