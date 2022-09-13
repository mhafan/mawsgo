package mawsgo

import (
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
// ...
func InitRegion(regionCode string) *MAWS {
	//
	return &MAWS{
		AWS: session.Must(session.NewSession(&aws.Config{
			Region: aws.String(regionCode)})),
	}
}
