package mawsgo

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ---------------------------------------------------------------------------
//
type MessageQueue struct {
	//
	QueueName string
	QueueURL  string

	//
	AWS    *session.Session
	Handle *sqs.SQS
}

// ---------------------------------------------------------------------------
//
func (maws *MAWS) MakeMessageQueue(qName string) (*MessageQueue, error) {
	//
	var _sqs = &MessageQueue{
		AWS:       maws.AWS,
		QueueName: qName,
		Handle:    sqs.New(maws.AWS),
	}

	//
	_url, _err := _sqs.Handle.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(_sqs.QueueName),
	})

	//
	if _err != nil {
		//
		return nil, _err
	}

	//
	_sqs.QueueURL = *_url.QueueUrl

	//
	return _sqs, nil
}

// ---------------------------------------------------------------------------
//
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
//
func (msgs *MessageQueue) SendMsgJSON(body, argKey string, arg interface{}) error {
	//
	_, era := msgs.Handle.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			argKey: MAWSEncodeMessage(&arg),
		},
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(msgs.QueueURL),
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

// ---------------------------------------------------------------------------
//
func MAWSEncodeMessage(rec interface{}) *sqs.MessageAttributeValue {
	//
	_enco, _err := json.Marshal(rec)

	//
	if _err != nil {
		//
		panic("nelze zakodovat")
	}

	//
	return &sqs.MessageAttributeValue{
		//
		DataType:    aws.String("String"),
		StringValue: aws.String(string(_enco)),
	}
}

// ---------------------------------------------------------------------------
//
func MAWSDecodeMessage(msg *sqs.Message, argKey string, rec interface{}) error {
	//
	_args, _ok := msg.Attributes[argKey]

	//
	if _ok == false {
		//
		return errors.New("arg not found")
	}

	//
	_err := json.Unmarshal([]byte(*_args), &rec)

	//
	if _err != nil {
		//
		panic("nelze zakodovat")
	}

	//
	return _err
}
