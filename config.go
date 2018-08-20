package main

import (
	"github.com/kobolog/gorb/core"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Services map[string]*core.ServiceOptions
type Backends map[string]*core.BackendOptions

type Schema struct {
	Backends Backends `yaml:"backends"`
	Services Services `yaml:"services"`
}

func LoadConfiguration(ctx *core.Context) {
	schema := Schema{}

	yamlFile, err := ioutil.ReadFile(*config); if err != nil {
		log.Printf("Unable to read services.yaml")
		return
	}

	err = yaml.Unmarshal(yamlFile, &schema)
	if err != nil {
		log.Printf("Unable to parse services.yaml")
		return
	}
	ctx.Synchronize(schema.Services, schema.Backends)
}
