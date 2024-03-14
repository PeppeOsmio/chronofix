package structs

type Config struct {
	TryExtractFromExif     bool   `yaml:"try_extract_from_exif"`
	TryExtractFromFileName bool   `yaml:"try_extract_from_file_name"`
	TrySetExifIfNotPresent bool   `yaml:"try_set_exif_if_not_present"`
	LogLevel               string `yaml:"log_level"`
	MediaDir               string `yaml:"media_dir"`
}
