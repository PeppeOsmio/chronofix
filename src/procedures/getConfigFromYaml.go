package procedures

import (
	"metadata_restorer/structs"
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfigFromYaml(path string) (config *structs.Config, err error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
