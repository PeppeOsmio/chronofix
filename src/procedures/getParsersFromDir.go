package procedures

import (
	"io/fs"
	"metadata_restorer/structs"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ParsersYaml struct {
	Parsers map[string]*structs.MediaDateTimeParser `yaml:"Parsers"`
}

func GetParsersFromDir(dirPath string) (parsers []*structs.MediaDateTimeParser, err error) {
	err = filepath.WalkDir(dirPath, func(path string, dirEntry fs.DirEntry, _ error) error {
		if dirEntry.IsDir() {
			return nil
		}
		logrus.Info("Found parsers file " + path)
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		parsersYaml := ParsersYaml{}
		err = yaml.Unmarshal(fileContent, &parsersYaml)
		if err != nil {
			return err
		}
		for key, value := range parsersYaml.Parsers {
			parsers = append(parsers, &structs.MediaDateTimeParser{
				Name:              key,
				Regex:             value.Regex,
				YearIndex:         value.YearIndex,
				MonthIndex:        value.MonthIndex,
				DayIndex:          value.DayIndex,
				HoursIndex:        value.HoursIndex,
				MinutesIndex:      value.MinutesIndex,
				SecondsIndex:      value.SecondsIndex,
				MillisecondsIndex: -1,
			})
		}
		logrus.Info("Parsed parsers file " + path)
		return nil
	})
	return parsers, err
}
