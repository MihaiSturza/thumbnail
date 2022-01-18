package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/h2non/bimg"
	log "github.com/sirupsen/logrus"
)

func IsPdf(name string) bool {
	return strings.HasSuffix(name, ".pdf")
}

func GeneratePdfThumbnail(bucket, key string) {
	info := strings.Split(key, ".")
	fileName, extension := info[0], info[1]
	localOriginal := fmt.Sprintf("/tmp/original.%s", extension)
	// local := fmt.Sprintf("/tmp/%s/%s", bucket, key)

	// ensure path is available
	dir := filepath.Dir(localOriginal)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.WithError(err).WithField("path", dir).Error("failed to create tmp directory")
	}

	// create a file locally for original image in S3
	f, err := os.Create(localOriginal)
	if err != nil {
		log.WithError(err).WithField("filename", localOriginal).Error("failed to create file")
		return
	}

	n, err := Downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"bucket":   bucket,
			"key":      key,
			"filename": localOriginal,
		}).Error("failed to download file")
		return
	}

	log.WithFields(log.Fields{
		"filename": localOriginal,
		"bytes":    n,
	}).Info("file downloaded")

	buffer, err := bimg.Read(localOriginal)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newImage, err := bimg.NewImage(buffer).Convert(bimg.JPEG)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if bimg.NewImage(newImage).Type() == "jpeg" {
		fmt.Fprintln(os.Stderr, "The image was converted into jpeg")
	}

	bimg.Write("/tmp/original.jpg", newImage)

	thumbName := fmt.Sprintf("%s_thumbnail.jpg", fileName)

	// upload thumbnail to S3
	up, err := os.Open("/tmp/original.jpg")
	if err != nil {
		log.WithError(err).WithField("thumbnail", newImage).Error("failed to open file")
		return
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(thumbName),
		Body:   up,
	})

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"bucket":    bucket,
			"thumbnail": thumbName,
		}).Error("failed to upload file")
	}

	log.WithField("location", result.Location).Info("file uploaded")

}
