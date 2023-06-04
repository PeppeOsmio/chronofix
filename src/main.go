package main

import (
	"flag"
	"io/fs"
	"metadata_restorer/procedures"
	"metadata_restorer/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// TODO check if EXIF is present first. If it is, set modified time like EXIF. If EXIF is not present,
	// set modified time from the file name. Optionally, set the EXIF from the file name.
	// Add support for more file formats
	config_path := flag.String("config", "./config.yml", "Path of the config file for metadata_parser")
	config, err := procedures.GetConfigFromYaml(*config_path)
	if err != nil {
		logrus.Error("Can't open config file " + *config_path + ": " + err.Error())
		return
	}
	flag.Parse()
	logsDir := "./logs"
	os.Mkdir(logsDir, 0755)
	logFileName := "metadata_restorer_" + time.Now().Format("20060102-150405") + ".log"
	logFilePath := filepath.Join(logsDir, logFileName)
	logFile, err := utils.OpenOrCreate(logFilePath, os.O_RDWR)
	logrus.SetLevel(logrus.DebugLevel)
	if err != nil {
		logrus.Error("Can't open log file: " + err.Error())
		return
	}
	logrus.SetOutput(logFile)
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Error("Invalid log level")
		return
	}
	logrus.SetLevel(logLevel)
	ignored_files_dir := "./ignored_files"
	os.Mkdir(ignored_files_dir, 0755)
	ignored_file_name := filepath.Join(ignored_files_dir, time.Now().Format("20060102-150405")+".txt")
	ignoredFile, err := utils.OpenOrCreate(ignored_file_name, os.O_APPEND|os.O_WRONLY)
	if err != nil {
		logrus.Error("Can't open file: " + err.Error())
	}
	defer ignoredFile.Close()
	fileNamesParsed := 0
	exifRead := 0
	exifSet := 0
	succeeded := 0
	failed := 0
	srcDir := "./src_files"
	defaultParsersDir := "./default_parsers"
	customParsersDir := "./custom_parsers"
	parsers, err := procedures.GetParsersFromDir(defaultParsersDir)
	if err != nil {
		logrus.Error("Can't read default parsers from directory " + customParsersDir + ": " + err.Error())
		return
	}
	customParsers, err := procedures.GetParsersFromDir(customParsersDir)
	if err != nil {
		logrus.Error("Can't read custom parsers from directory " + customParsersDir + ": " + err.Error())
		return
	}
	parsers = append(parsers, customParsers...)
	waitGroup := sync.WaitGroup{}
	filepath.WalkDir(srcDir, func(path string, dirEntry fs.DirEntry, err error) error {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err != nil {
				failed++
				logrus.Error(err.Error())
				return
			}
			// check if we are dealing with a dir or a file
			if dirEntry.IsDir() {
				logrus.Debug("Found directory " + path)
				return
			} else {
				logrus.Debug("Found file " + path)
			}
			// try to get the EXIF DateTimeOriginal from the file
			exifDateTimeOriginal := time.Time{}
			timeToSet := time.Time{}
			success := false
			// extract EXIF only if it's written in the config
			if config.TryExtractFromExif {
				exifDateTimeOriginal, err = procedures.GetExifDateTimeOriginal(path)
				if err != nil {
					logrus.Debug(path + ": " + err.Error())
				} else {
					// if the EXIF DateTimeOriginal is present, set the modified time of the file to it
					logrus.Debug(path + ": " + "EXIF DateTimeOriginal: " + utils.IsoDateTime(exifDateTimeOriginal))
					timeToSet = exifDateTimeOriginal
					success = true
					exifRead++
				}
			}
			// if we can try to extract the time from the file name and we didn't already extract it from the EXIF, try from the
			// file name
			if !success && config.TryExtractFromFileName {
				// otherwise, try to get the created time of the file from its name
				matchedIndex, t, err := procedures.GetDateTimeFromFileName(dirEntry.Name(), parsers)
				timeToSet = t
				if err != nil {
					logrus.Error(path + ": " + err.Error())
					return
				} else if matchedIndex != -1 {
					logrus.Debug(path + ": " + " matched pattern " + parsers[matchedIndex].Name)
					fileNamesParsed++
					success = true
				} else {
					failed++
					logrus.Debug(path + ": " + " did not match any patterns")
					ignoredFile.WriteString(path + "\n")
				}
				if success && config.TrySetExifIfNotPresent {
					err = procedures.ChangeExifDateTimeOriginal(path, timeToSet)
					if err != nil {
						logrus.Debug(path + ": Could not set EXIF: " + err.Error())
					} else {
						exifSet++
						logrus.Debug(path + ": " + "Set DateTimeOriginal EXIF to " + utils.IsoDateTime(timeToSet))
					}
				}
			}
			// if we still haven't succeeded in extracting the time, just return...
			if !success {
				failed++
				return
			}
			//Set the modified time of the file to the one extracted from EXIF or from the file name
			err = os.Chtimes(path, time.Now(), timeToSet)
			if err != nil {
				failed++
				logrus.Error("Error setting modified datetime: " + err.Error())
				return
			}
			succeeded++
			logrus.Debug(path + ": " + "set modified time to " + utils.IsoDateTime(timeToSet))
		}()
		return nil
	})
	waitGroup.Wait()
	logrus.Info("Total files: " + strconv.Itoa(succeeded+failed))
	logrus.Info("File names parsed successfully: " + strconv.Itoa(fileNamesParsed))
	logrus.Info("EXIF read successfully: " + strconv.Itoa(exifRead))
	logrus.Info("EXIF updated: " + strconv.Itoa(exifSet))
	logrus.Info("Files not updated: " + strconv.Itoa(failed))
}
