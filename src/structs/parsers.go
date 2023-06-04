package structs

type MediaRegex string

type MediaDateTimeParser struct {
	Name              string `yaml:"Name"`
	Regex             string `yaml:"Regex"`
	YearIndex         int    `yaml:"YearIndex"`
	MonthIndex        int    `yaml:"MonthIndex"`
	DayIndex          int    `yaml:"DayIndex"`
	HoursIndex        int    `yaml:"HoursIndex"`
	MinutesIndex      int    `yaml:"MinutesIndex"`
	SecondsIndex      int    `yaml:"SecondsIndex"`
	MillisecondsIndex int    `yaml:"MillisecondsIndex"`
}
