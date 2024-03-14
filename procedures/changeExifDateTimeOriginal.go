package procedures

import (
	"errors"
	"os"
	"time"

	"github.com/dsoprea/go-exif/v3"
	jis "github.com/dsoprea/go-jpeg-image-structure/v2"
)

func ChangeExifDateTimeOriginal(path string, timeToSet time.Time) error {
	// Decode the EXIF data
	mediaContext, err := jis.NewJpegMediaParser().ParseFile(path)
	if err != nil {
		return errors.New("Error extracting media context: " + err.Error())
	}
	segmentList := mediaContext.(*jis.SegmentList)
	ifdBuilder, err := segmentList.ConstructExifBuilder()
	if err != nil {
		return errors.New("Error extracting root IFD: " + err.Error())
	}
	exifIfd, err := exif.GetOrCreateIbFromRootIb(ifdBuilder, "IFD/Exif")
	if err != nil {
		return errors.New("Error extracting IFD IFD/Exif : " + err.Error())
	}
	// DateTimeOriginal EXIF tag
	const dateTimeOriginalExifTag = 0x9003
	err = exifIfd.SetStandard(dateTimeOriginalExifTag, timeToSet)
	if err != nil {
		return errors.New("Error setting EXIF: " + err.Error())
	}
	err = segmentList.SetExif(ifdBuilder)
	if err != nil {
		return errors.New("Error setting EXIF: " + err.Error())
	}
	// Open the image file
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return errors.New("Error opening file for saving EXIF: " + err.Error())
	}
	defer file.Close()
	segmentList.Write(file)
	return err
}
