package mawsgo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// ---------------------------------------------------------------------------
//
func SayHello() string {
	return "Hello from MAWS-for-go"
}

// ---------------------------------------------------------------------------
//
type MAWS struct {
	//
	AWS *session.Session
}

// ---------------------------------------------------------------------------
// AWS_DEFAULT_REGION -
func MAWSInit() *MAWS {
	//
	defRegion := os.Getenv("AWS_DEFAULT_REGION")

	//
	return &MAWS{
		AWS: session.Must(session.NewSession(&aws.Config{
			Region: aws.String(defRegion)})),
	}
}
