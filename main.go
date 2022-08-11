package mawsgo

import (
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
//
func MAWSInit() *MAWS {
	//
	return &MAWS{
		AWS: session.Must(session.NewSession()),
	}
}
