package mawsgo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// ---------------------------------------------------------------------------
// Zakladni ridici struktura pro MAWS
type MAWS struct {
	// handle na spojeni
	AWS *session.Session
}

// ---------------------------------------------------------------------------
// AWS_DEFAULT_REGION - vyzaduje se pritomnost v env
func MAWSInit() *MAWS {
	//
	defRegion := os.Getenv("AWS_DEFAULT_REGION")

	//
	return &MAWS{
		AWS: session.Must(session.NewSession(&aws.Config{
			Region: aws.String(defRegion)})),
	}
}

// ---------------------------------------------------------------------------
// AWS_DEFAULT_REGION - vyzaduje se pritomnost v env
func MAWSInitRegion(regionCode string) *MAWS {
	//
	return &MAWS{
		AWS: session.Must(session.NewSession(&aws.Config{
			Region: aws.String(regionCode)})),
	}
}

// ---------------------------------------------------------------------------
//
func MAWSSetRegion(regionCode string) {
	//
	os.Setenv("AWS_DEFAULT_REGION", regionCode)
}
