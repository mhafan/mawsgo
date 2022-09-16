package mawsgo

import (
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------------------
// Datova struktura lokalniho souboru
// ---------------------------------------------------------------------------
// Struktura muze vzniknout:
// - tmp adresar + jmeno
// - cela cesta
type LocFile struct {
	// jmeno souboru
	Name string

	// kompletni cesta k souboru
	FilePath string

	// napojeni na S3 bucket
	S3Connect *BucketKey
}

// ---------------------------------------------------------------------------
// Lambda TMP file
// - name muze byt prazdne -> nahodne pridelene jmeno souboru
func TmpFile(name string) *LocFile {
	// pokud je prazdnota
	if len(name) <= 0 {
		// ... generuj
		name = MAWSUUID()
	}

	//
	return &LocFile{
		Name:     name,
		FilePath: filepath.Join("/tmp", name),
	}
}

// ---------------------------------------------------------------------------
// celkova cesta k souboru
func PathFile(efilePath string) *LocFile {
	//
	_, fn := filepath.Split(efilePath)

	//
	return &LocFile{
		Name:     fn,
		FilePath: efilePath,
	}
}

// ---------------------------------------------------------------------------
//
func (lf *LocFile) S3(bucket *Bucket, prefixed string) *LocFile {
	//
	lf.S3Connect = bucket.MakeKey(prefixed, lf.Name)

	return lf
}

// ---------------------------------------------------------------------------
// Ulozeni textu do souboru
func (lf *LocFile) SaveText(text string) error {
	//
	return lf.SaveBin([]byte(text))
}

// ---------------------------------------------------------------------------
// Ulozeni textu do souboru
func (lf *LocFile) SaveBin(content []byte) error {
	//
	f, errOpen := os.Create(lf.FilePath)

	//
	if errOpen != nil {
		return errOpen
	}

	defer f.Close()

	_, err2 := f.Write(content)

	//
	return err2
}

// ---------------------------------------------------------------------------
//
func (lf *LocFile) Read() ([]byte, error) {
	//
	return os.ReadFile(lf.FilePath)
}

// ---------------------------------------------------------------------------
//
func (lf *LocFile) ReadString() string {
	//
	cont, err := lf.Read()

	//
	if err != nil {
		//
		return ""
	}

	//
	return string(cont)
}

// ---------------------------------------------------------------------------
//
func (lf *LocFile) Delete() error {
	//
	return os.Remove(lf.FilePath)
}

// ---------------------------------------------------------------------------
// Transformuj LF -> Bucket+key
func (lf *LocFile) BucketKey(bucket *Bucket, prefixed string) *BucketKey {
	//
	return bucket.MakeKey(prefixed, lf.Name)
}

// ---------------------------------------------------------------------------
// Nahrani souboru do prislusneho bucket-key
func (lf *LocFile) UploadToBucket(bucket *Bucket, prefixed string) error {
	//
	return lf.BucketKey(bucket, prefixed).Upload(lf)
}

// ---------------------------------------------------------------------------
// Nahrani souboru do prislusneho bucket-key
func (lf *LocFile) DownloadFromBucket(bucket *Bucket, prefixed string) error {
	//
	return lf.BucketKey(bucket, prefixed).Download(lf)
}
