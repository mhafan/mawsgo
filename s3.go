package mawsgo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ---------------------------------------------------------------------------
// Datovy model pristupu do jednoho bucket
// ---------------------------------------------------------------------------
// API:
// - MAWSMakeBucket - vytvoreni struktury pro pristup na bucket
// - Download(key, local) - stazeni souboru
// - Uploadfile
// - ListObjectKeys
// - DeleteKey
type MAWSBucket struct {
	//
	BucketName string
	Handle     *s3.S3

	// handle na spojeni AWS (kopiruje se)
	AWS *session.Session
}

// ---------------------------------------------------------------------------
// vytvoreni Handle na S3:bucket
func (maws *MAWS) MAWSMakeBucket(name string) *MAWSBucket {
	//
	return &MAWSBucket{
		BucketName: name,
		Handle:     s3.New(maws.AWS),
		AWS:        maws.AWS,
	}
}

// ---------------------------------------------------------------------------
// Download souboru S3
func (b *MAWSBucket) Download(key string, saveto string) error {
	// TODO: musi nutne vznikat pro kazdou operaci?
	// je v tom nejaka vyznamna casova rezie?
	downloader := s3manager.NewDownloader(b.AWS)

	// obsluha lokalniho souboru
	f, err := os.Create(saveto)
	if err != nil {
		//
		return err
	}

	// ...
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
func (b *MAWSBucket) Uploadfile(key string, localname string) error {
	// TODO: podobne jako Download
	uploader := s3manager.NewUploader(b.AWS)

	// obsluha lokalniho souboru
	file, err := os.Open(localname)
	if err != nil {
		return err
	}

	// TODO: nedela to AWS???
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

// ---------------------------------------------------------------------------
// List content of bucket
func (b *MAWSBucket) ListObjects() (*s3.ListObjectsV2Output, error) {
	//
	resp, err := b.Handle.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(b.BucketName),
	})

	//
	return resp, err
}

// ---------------------------------------------------------------------------
// List content of bucket
func (b *MAWSBucket) ListObjectKeys() ([]string, error) {
	//
	resp, err := b.ListObjects()

	//
	if err != nil {
		//
		return []string{}, err
	}

	var _out []string

	//
	for _, i := range resp.Contents {
		//
		_out = append(_out, *i.Key)
	}

	//
	return _out, nil
}

// ---------------------------------------------------------------------------
// List content of bucket
func (b *MAWSBucket) DeleteKey(key string) error {
	//
	_, err := b.Handle.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(key),
	})

	//
	return err
}
