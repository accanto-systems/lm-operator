package alm

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Version string `yaml:"version"`
}

type LMRelease struct {
	Configurator Service `yaml:"configurator"`
	Conductor    Service `yaml:"conductor"`
	Apollo       Service `yaml:"apollo"`
	Galileo      Service `yaml:"galileo"`
	Talledega    Service `yaml:"talledega"`
	Daytona      Service `yaml:"daytona"`
	Relay        Service `yaml:"relay"`
	Watchtower   Service `yaml:"watchtower"`
	Brent        Service `yaml:"brent"`
	Doki         Service `yaml:"doki"`
	Ishtar       Service `yaml:"ishtar"`
	Nimrod       Service `yaml:"nimrod"`
}

func getLMRelease(url string, reqLogger logr.Logger) (*LMRelease, error) {
	// url := fmt.Sprintf("http://10.220.217.248:8086/accanto/lm-operator-releases/raw/master/%s.yaml", release)

	reqLogger.Info(fmt.Sprintf("Getting release descriptor for release from %s", url))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	yamlFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqLogger.Info(fmt.Sprintf("Release yaml %s", string(yamlFile)))

	r := LMRelease{}
	err = yaml.Unmarshal(yamlFile, &r)
	if err != nil {
		return nil, err
	} else {
		return &r, nil
	}
}
