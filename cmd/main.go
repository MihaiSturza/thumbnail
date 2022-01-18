package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/MihaiSturza/thumbnail/internal"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

// S3 Session to use
var sess = session.Must(session.NewSession())

// Create an uploader with session and default option
var uploader = s3manager.NewUploader(sess)

// Create a downloader with session and default option
var downloader = s3manager.NewDownloader(sess)

func main() {
	fmt.Println("lambda running....")
	lambda.Start(ReceiveEvent)
}

func ReceiveEvent(ctx context.Context, request Data) (string, error) {
	file := internal.File{Uploader: uploader, Downloader: downloader}
	if isImage(request.Key) && !strings.Contains(request.Key, "_thumbnail") {
		generateThumbnail(request.Bucket, request.Key)
	}
	if internal.IsPdf(request.Key) {
		fmt.Print(file)
		internal.GeneratePdfThumbnail(request.Bucket, request.Key, file)

	}
	return fmt.Sprintf("Data received and processed: %v", request), nil
}

func generateThumbnail(bucket, key string) {
	info := strings.Split(key, ".")
	fileName, extension := info[0], info[1]
	localOriginal := fmt.Sprintf("/tmp/original_image.%s", extension)
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

	n, err := downloader.Download(f, &s3.GetObjectInput{
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

	img, err := imaging.Open(localOriginal)
	if err != nil {
		panic(err)
	}
	thumb := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)

	// create a new blank image
	dst := imaging.New(100, 100, color.NRGBA{0, 0, 0, 0})

	// paste thumbnails into the new image
	dst = imaging.Paste(dst, thumb, image.Pt(0, 0))

	// save the combined image to file
	fname := strings.Join(strings.Split(fileName, "/")[1:], "/")
	thumbName := fmt.Sprintf("%s_thumbnail.%s", fname, extension)
	localThumbnail := fmt.Sprintf("/tmp/thumbnail_image.%s", extension)

	err = imaging.Save(dst, localThumbnail)
	if err != nil {
		log.WithError(err).WithField("thumbnail", localThumbnail).Error("failed to generate thumbnail")
		return
	}

	// upload thumbnail to S3
	up, err := os.Open(localThumbnail)
	if err != nil {
		log.WithError(err).WithField("thumbnail", localThumbnail).Error("failed to open file")
		return
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("thumbnail/%s", thumbName)),
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

func isImage(name string) bool {
	if strings.HasSuffix(name, ".jpg") {
		return true
	}
	if strings.HasSuffix(name, ".jpeg") {
		return true
	}

	if strings.HasSuffix(name, ".png") {
		return true
	}
	return false
}
