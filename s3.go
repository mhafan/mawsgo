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
	Handle     *s3.S3
}

// ---------------------------------------------------------------------------
//
func MAWSMakeBucket(maws *MAWS, name string) *MAWSBucket {
	//
	return &MAWSBucket{
		BucketName: name,
		Handle:     s3.New(maws.AWS),
	}
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

// ---------------------------------------------------------------------------
// List content of bucket
func (b *MAWSBucket) ListObjects(maws *MAWS) (*s3.ListObjectsV2Output, error) {
	//
	resp, err := b.Handle.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(b.BucketName)})

	//
	return resp, err
}

// ---------------------------------------------------------------------------
// List content of bucket
func (b *MAWSBucket) ListObjectKeys(maws *MAWS) ([]string, error) {
	//
	resp, err := b.ListObjects(maws)

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
func (b *MAWSBucket) DeleteKey(maws *MAWS, key string) error {
	//
	_, err := b.Handle.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(key),
	})

	//
	return err
}
