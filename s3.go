package mawsgo

import (
	"bytes"
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

	//
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
}

// ---------------------------------------------------------------------------
// vytvoreni Handle na S3:bucket
func (maws *MAWS) MakeBucket(name string) *Bucket {
	//
	return &Bucket{
		BucketName: name,
		Handle:     s3.New(maws.AWS),
		AWS:        maws.AWS,
		Downloader: s3manager.NewDownloader(maws.AWS),
		Uploader:   s3manager.NewUploader(maws.AWS),
	}
}

// ---------------------------------------------------------------------------
// Predstavuje jeden S3-objekt (soubor) ulozeny v bucket na nejakem prefixu
type BucketKey struct {
	// nazev objektu
	FileName string
	// prefix
	Prefix string
	// handle na prislusny S3-bucket
	Bucket *Bucket
}

// ---------------------------------------------------------------------------
// Konstrukce handle na S3 objekt nadd zadanym S3-bucket
func (b *Bucket) MakeKey(prefix, name string) *BucketKey {
	//
	return &BucketKey{
		FileName: name,
		Prefix:   prefix,
		Bucket:   b,
	}
}

// ---------------------------------------------------------------------------
// prefix+name
func (bk *BucketKey) Key() string {
	//
	return filepath.Join(bk.Prefix, bk.FileName)
}

// ---------------------------------------------------------------------------
// Download souboru S3
func (bk *BucketKey) Download(locFile *LocFile) error {
	// obsluha lokalniho souboru
	f, err := os.Create(locFile.FilePath)
	if err != nil {
		//
		return err
	}

	// ...
	defer f.Close()

	// ...
	_, err2 := bk.Bucket.Downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bk.Bucket.BucketName),
		Key:    aws.String(bk.Key()),
	})

	//
	return err2
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
	return bk.UploadContent(bytes.NewReader(cont))
}

// ---------------------------------------------------------------------------
// Upload of S3 file
func (bk *BucketKey) UploadContent(readStream *bytes.Reader) error {
	//
	_, errUp := bk.Bucket.Uploader.Upload(&s3manager.UploadInput{
		// identifikace
		Bucket: aws.String(bk.Bucket.BucketName),
		Key:    aws.String(bk.Key()),
		// obsah
		Body: readStream,
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
