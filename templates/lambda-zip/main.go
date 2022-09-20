package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mhafan/mawsgo"
)

// ---------------------------------------------------------------------------
// ENV konfigurace Lambdy
var _MAWS = mawsgo.InitRegion(mawsgo.Env("AWS_DEFAULT_REGION", "eu-central-1"))

// ---------------------------------------------------------------------------
// Vstupni udalost - typicky pro bez-SQS spousteni Lambdy
type EventInput struct {
	Cosi string
}

// ---------------------------------------------------------------------------
// Proste volani Lambdy
// navratova hodnota - (string, error) nebo (error)
func HandleRequest(ctx context.Context, event EventInput) (string, error) {
	//
	_enco, _ := json.Marshal(event)

	//
	return string(_enco), nil
}

// ---------------------------------------------------------------------------
//
func HandleRequestSQSInput(ctx context.Context, sqsEvent events.SQSEvent) error {
	// prochazim pres prvky baliku
	for _, msg := range sqsEvent.Records {
		// dekoduju zpravu v baliku
		if _msgDecoded := mawsgo.DecodeMessage(&msg); _msgDecoded != nil {
			// zpracovani
		}
	}

	//
	return nil
}

// ---------------------------------------------------------------------------
// Zavadec Lambdy
func main() {
	// normalni
	lambda.Start(HandleRequest)

	// pro SQS vstup
	//lambda.Start(HandleRequestSQSInput)
}
