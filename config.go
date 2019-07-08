package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct{}

func Configure(fileName string) (Config, error) {
	var cnf Config

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		return Config{}, err
	}

	return cnf, nil
}
