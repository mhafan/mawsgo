package mawsgo

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ---------------------------------------------------------------------------
// SQS Fronta
// ---------------------------------------------------------------------------
// Zprava ve fronte se sklada z:
// - body - string
// - messageAttributes - map[string] value
// - metadat
// ---------------------------------------------------------------------------
// Je veci koncepce, zda-li zpravu pojmout jako
// 1) Body=identifikator + MessageAttributes jako obsah/telo zpravy
// 2) Body=JSON s celkovym obsahem
type MessageQueue struct {
	// zakladni pojmenovani fronty (bezny identifikator)
	QueueName string
	// intern URL/ARN zdroje v AWS
	QueueURL string

	// handle na AWS + na frontu
	AWS    *session.Session
	Handle *sqs.SQS
}

// ---------------------------------------------------------------------------
// Struktura vyjadrujici zakladni atributy SQSMessage
type PlainMessage struct {
	//
	Body  string
	Attrs map[string]string
}

// ---------------------------------------------------------------------------
// Dekodovani prichozi zpravy do jeji zakladni podstaty
func DecodeMessage(inm *events.SQSMessage) *PlainMessage {
	//
	_out := &PlainMessage{
		Body:  inm.Body,
		Attrs: map[string]string{},
	}

	//
	for k, v := range inm.MessageAttributes {
		//
		_out.Attrs[k] = *v.StringValue
	}

	//
	return _out
}

// ---------------------------------------------------------------------------
// Vytvoreni handle na SQS frontu
func (maws *MAWS) MakeMessageQueue(qName string) (*MessageQueue, error) {
	// handle
	var _sqs = &MessageQueue{
		AWS:       maws.AWS,
		QueueName: qName,
		Handle:    sqs.New(maws.AWS),
	}

	// zjisteni meta informaci o zdroji
	_url, _err := _sqs.Handle.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(_sqs.QueueName),
	})

	//
	if _err != nil {
		//
		return nil, _err
	}

	// ...
	_sqs.QueueURL = *_url.QueueUrl

	//
	return _sqs, nil
}

// ---------------------------------------------------------------------------
// Vytvoreni handle na SQS frontu
func (maws *MAWS) MakeMessageQueue_(qName string) *MessageQueue {
	//
	_sqs, _err := maws.MakeMessageQueue(qName)

	//
	if _err != nil {
		//
		panic(_err)
	}

	//
	return _sqs
}

// ---------------------------------------------------------------------------
// Proste poslani textoveho obsahu do fronty
// -> body
func (msgs *MessageQueue) SendMsg(body string) error {
	//
	_, era := msgs.Handle.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(msgs.QueueURL),
	})

	//
	return era
}

// ---------------------------------------------------------------------------
// Zprava ve formatu Body + MessageAttributes
func (msgs *MessageQueue) SendMsgAttrs(plain *PlainMessage) error {
	//
	var _attrs = make(map[string]*sqs.MessageAttributeValue)

	//
	for key, value := range plain.Attrs {
		//
		_attrs[key] = &sqs.MessageAttributeValue{
			//
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	//
	return msgs.SendMsgAttrs_(plain.Body, _attrs)
}

// ---------------------------------------------------------------------------
// Zprava ve formatu Body + MessageAttributes
func (msgs *MessageQueue) SendMsgAttrs_(body string, attrs map[string]*sqs.MessageAttributeValue) error {
	//
	_, era := msgs.Handle.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: attrs,
		MessageBody:       aws.String(body),
		QueueUrl:          aws.String(msgs.QueueURL),
	})

	//
	return era
}

// ---------------------------------------------------------------------------
//
func (msgs *MessageQueue) ReceiveMsgs(limit, timeOut int) ([]*sqs.Message, error) {
	//
	msgResult, err := msgs.Handle.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &msgs.QueueURL,
		MaxNumberOfMessages: aws.Int64(int64(limit)),
		VisibilityTimeout:   aws.Int64(int64(timeOut)),
	})

	//
	return msgResult.Messages, err
}

// ---------------------------------------------------------------------------
//
func (msgs *MessageQueue) DeleteMsg(msg *sqs.Message) error {
	//
	_, __err := msgs.Handle.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &msgs.QueueURL,
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	})

	//
	return __err
}

// ---------------------------------------------------------------------------
//
func (msgs *MessageQueue) DeleteMsgHandle(msg string) error {
	//
	_, __err := msgs.Handle.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &msgs.QueueURL,
		ReceiptHandle: aws.String(msg),
	})

	//
	return __err
}
