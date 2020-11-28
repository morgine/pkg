// Copyright 2020 morgine.com. All rights reserved.

package config

import (
	"fmt"
	"os"
)

func ExampleConfigs_UnmarshalSub() {
	// config structure
	var mysql = struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	}{}

	// config data
	var data = `
[mysql]
host = "127.0.0.1"
port= "3306"
`

	// os environment will cover config data
	err := os.Setenv("mysql.host", "localhost")
	if err != nil {
		panic(err)
	}

	configs, err := UnmarshalMemory([]byte(data))
	if err != nil {
		panic(err)
	}

	err = configs.UnmarshalSub("mysql", &mysql)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:%s", mysql.Host, mysql.Port)
	// Output:
	// localhost:3306
}
