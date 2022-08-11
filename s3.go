package mawsgo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ---------------------------------------------------------------------------
// Datovy model pristupu do jednoho bucket
type MAWSBucket struct {
	//
	BucketName string
}

// ---------------------------------------------------------------------------
// Download souboru S3
func (b *MAWSBucket) Download(maws *MAWS, key string, saveto string) error {
	//
	downloader := s3manager.NewDownloader(maws.AWS)

	//
	f, err := os.Create(saveto)
	if err != nil {
		//
		return err
	}
	defer f.Close()

	// ...
	_, err2 := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(key),
	})

	//
	return err2
}

// ---------------------------------------------------------------------------
// Upload of S3 file
func (b *MAWSBucket) Uploadfile(maws *MAWS, key string, localname string) error {
	//
	uploader := s3manager.NewUploader(maws.AWS)

	//
	file, err := os.Open(localname)
	//
	if err != nil {
		return err
	}
	defer file.Close()

	//
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	//
	return err
}
