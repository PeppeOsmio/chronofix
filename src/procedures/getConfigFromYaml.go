package procedures

import (
	"metadata_restorer/structs"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigYaml struct {
	Config *structs.Config `yaml:"Config"`
}

func GetConfigFromYaml(path string) (config *structs.Config, err error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	configYaml := ConfigYaml{}
	err = yaml.Unmarshal(fileContent, &configYaml)
	if err != nil {
		return nil, err
	}
	config = configYaml.Config
	return config, nil
}
