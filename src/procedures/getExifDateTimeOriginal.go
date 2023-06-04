package procedures

import (
	"errors"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	jis "github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/sirupsen/logrus"
)

func GetExifDateTimeOriginal(path string) (dateTimeOriginal time.Time, err error) {
	if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		// Decode the EXIF data
		mediaContext, err := jis.NewJpegMediaParser().ParseFile(path)
		if err != nil {
			return dateTimeOriginal, errors.New("Error extracting media context: " + err.Error())
		}
		segmentList := mediaContext.(*jis.SegmentList)
		ifdBuilder, err := segmentList.ConstructExifBuilder()
		if err != nil {
			return dateTimeOriginal, errors.New("Error extracting root IFD: " + err.Error())
		}
		exifIfd, err := exif.GetOrCreateIbFromRootIb(ifdBuilder, "IFD/Exif")
		if err != nil {
			return dateTimeOriginal, errors.New("Error extracting IFD IFD/Exif : " + err.Error())
		}
		builderTag, err := exifIfd.FindTagWithName("DateTimeOriginal")
		if err != nil {
			return dateTimeOriginal, errors.New("Error extracting DateTimeOriginal : " + err.Error())
		}
		value := builderTag.Value()
		if value == nil {
			return dateTimeOriginal, errors.New("DateTimeOriginal is not present")
		}
		dateTimeOriginalString, _ := strings.CutSuffix(string(value.Bytes()), "\x00")
		// try to get the TimeZoneOffset from EXIF. Otherwise use the system's timezone
		location := time.Local
		builderTag, err = exifIfd.FindTagWithName("OffsetTimeOriginal")
		if err == nil {
			value = builderTag.Value()
			locationString, _ := strings.CutSuffix(string(value.Bytes()), "\x00")
			logrus.Debug(path + ": found TimeZoneOffset " + locationString)
		} else {
			logrus.Debug(path + ": no TimeZoneOffset found")
		}
		if value == nil {
			return dateTimeOriginal, errors.New("DateTimeOriginal is not present")
		}
		dateTimeOriginal, err = time.Parse("2006:01:02 15:04:05", dateTimeOriginalString)
		if err != nil {
			return dateTimeOriginal, errors.New("Can't parse DateTimeOriginal : " + err.Error())
		}
		// set the timezone
		dateTimeOriginal = time.Date(dateTimeOriginal.Year(), dateTimeOriginal.Month(), dateTimeOriginal.Day(),
			dateTimeOriginal.Hour(), dateTimeOriginal.Minute(), dateTimeOriginal.Second(), dateTimeOriginal.Nanosecond(), location)
		return dateTimeOriginal, nil
	}
	return dateTimeOriginal, errors.New("format not supported")
}
