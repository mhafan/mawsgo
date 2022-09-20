package mawsgo

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

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
type Bucket struct {
	//
	BucketName string
	Handle     *s3.S3

	// handle na spojeni AWS (kopiruje se)
	AWS *session.Session
}

// ---------------------------------------------------------------------------
//
type BucketKey struct {
	//
	FileName string
	Prefix   string

	//
	Bucket *Bucket
}

// ---------------------------------------------------------------------------
//
func (b *Bucket) MakeKey(prefix, name string) *BucketKey {
	//
	return &BucketKey{
		FileName: name,
		Prefix:   prefix,
		Bucket:   b,
	}
}

// ---------------------------------------------------------------------------
//
func (bk *BucketKey) Key() string {
	//
	return filepath.Join(bk.Prefix, bk.FileName)
}

// ---------------------------------------------------------------------------
// vytvoreni Handle na S3:bucket
func (maws *MAWS) MakeBucket(name string) *Bucket {
	//
	return &Bucket{
		BucketName: name,
		Handle:     s3.New(maws.AWS),
		AWS:        maws.AWS,
	}
}

// ---------------------------------------------------------------------------
// Download souboru S3
func (bk *BucketKey) Download(locFile *LocFile) error {
	// TODO: musi nutne vznikat pro kazdou operaci?
	// je v tom nejaka vyznamna casova rezie?
	downloader := s3manager.NewDownloader(bk.Bucket.AWS)

	// obsluha lokalniho souboru
	f, err := os.Create(locFile.FilePath)
	if err != nil {
		//
		return err
	}

	// ...
	defer f.Close()

	// ...
	_, err2 := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bk.Bucket.BucketName),
		Key:    aws.String(bk.Key()),
	})

	//
	return err2
}

// ---------------------------------------------------------------------------
//
func S3Uploads(inarr []*LocFile) error {
	//
	for _, lf := range inarr {
		//
		if lf.S3Connect == nil {
			//
			return errors.New("Chybi S3Connect na LF")
		}

		//
		if _err := lf.S3Connect.Upload(lf); _err != nil {
			//
			return _err
		}
	}

	//
	return nil
}

// ---------------------------------------------------------------------------
// Upload of S3 file
func (bk *BucketKey) Upload(locFile *LocFile) error {
	// obsluha lokalniho souboru
	cont, err := locFile.Read()

	//
	if err != nil {
		return err
	}

	//
	return bk.UploadContent(cont)
}

// ---------------------------------------------------------------------------
// Upload of S3 file
func (bk *BucketKey) UploadContent(cont []byte) error {
	// ...
	uploader := s3manager.NewUploader(bk.Bucket.AWS)

	//
	_, errUp := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bk.Bucket.BucketName),
		Key:    aws.String(bk.Key()),
		Body:   bytes.NewReader(cont),
	})

	//
	return errUp
}

// ---------------------------------------------------------------------------
// List content of bucket
func (b *Bucket) ListObjects() (*s3.ListObjectsV2Output, error) {
	//
	resp, err := b.Handle.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(b.BucketName),
	})

	//
	return resp, err
}

// ---------------------------------------------------------------------------
// List content of bucket
func (b *Bucket) ListObjectKeys() ([]string, error) {
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
func (b *Bucket) DeleteKey(key string) error {
	//
	_, err := b.Handle.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(key),
	})

	//
	return err
}
