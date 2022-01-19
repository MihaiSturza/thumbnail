package internal

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type File struct {
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
}

func IsPdf(name string) bool {
	return strings.HasSuffix(name, ".pdf")
}

// func GeneratePdfThumbnail(bucket, key string, file File) {
// 	info := strings.Split(key, ".")
// 	fileName, extension := info[0], info[1]
// 	localOriginalPdf := fmt.Sprintf("/tmp/original.%s", extension)
// 	// local := fmt.Sprintf("/tmp/%s/%s", bucket, key)

// 	// ensure path is available
// 	dir := filepath.Dir(localOriginalPdf)
// 	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
// 		log.WithError(err).WithField("path", dir).Error("failed to create tmp directory")
// 	}

// 	// create a file locally for original image in S3
// 	f, err := os.Create(localOriginalPdf)
// 	if err != nil {
// 		log.WithError(err).WithField("filename", localOriginalPdf).Error("failed to create file")
// 		return
// 	}

// 	n, err := file.Downloader.Download(f, &s3.GetObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		log.WithError(err).WithFields(log.Fields{
// 			"bucket":   bucket,
// 			"key":      key,
// 			"filename": localOriginalPdf,
// 		}).Error("failed to download file")
// 		return
// 	}

// 	log.WithFields(log.Fields{
// 		"filename": localOriginalPdf,
// 		"bytes":    n,
// 	}).Info("file downloaded")

// 	thumbName := fmt.Sprintf("thumbnail/%s_thumbnail.jpg", fileName)

// 	if err := ConvertPdfToJpg(localOriginalPdf, "/tmp/original.jpg"); err != nil {
// 		log.Fatal(err)
// 	}

// 	// upload thumbnail to S3
// 	up, err := os.Open("/tmp/original.jpg")
// 	if err != nil {
// 		log.WithError(err).WithField("thumbnail", "/tmp/original.jpg").Error("failed to open file")
// 		return
// 	}

// 	result, err := file.Uploader.Upload(&s3manager.UploadInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(thumbName),
// 		Body:   up,
// 	})

// 	if err != nil {
// 		log.WithError(err).WithFields(log.Fields{
// 			"bucket":    bucket,
// 			"thumbnail": thumbName,
// 		}).Error("failed to upload file")
// 	}

// 	log.WithField("location", result.Location).Info("file uploaded")

// }

// func ConvertPdfToJpg(pdfName string, imageName string) error {

// 	// Setup
// 	imagick.Initialize()
// 	defer imagick.Terminate()

// 	mw := imagick.NewMagickWand()
// 	defer mw.Destroy()

// 	// Must be *before* ReadImageFile
// 	// Make sure our image is high quality
// 	if err := mw.SetResolution(300, 300); err != nil {
// 		return err
// 	}

// 	// Load the image file into imagick
// 	if err := mw.ReadImage(pdfName); err != nil {
// 		return err
// 	}

// 	// Must be *after* ReadImageFile
// 	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
// 	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_DISCRETE); err != nil {
// 		return err
// 	}

// 	// Set any compression (100 = max quality)
// 	if err := mw.SetCompressionQuality(95); err != nil {
// 		return err
// 	}

// 	// Select only first page of pdf
// 	mw.SetIteratorIndex(0)

// 	// Convert into JPG
// 	if err := mw.SetFormat("jpg"); err != nil {
// 		return err
// 	}

// 	width := mw.GetImageWidth()
// 	height := mw.GetImageHeight()

// 	fmt.Printf("Width: %d, Height: %d", width, height)
// 	// get the higher integer
// 	//

// 	hWidth := uint(width / 10)
// 	hHeight := uint(height / 10)

// 	if err := mw.ThumbnailImage(hWidth, hHeight); err != nil {
// 		panic(err)
// 	}

// 	// Save File
// 	return mw.WriteImage(imageName)
// }
