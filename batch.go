package mawsgo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

// ---------------------------------------------------------------------------
//
type MAWSBatchQueue struct {
	//
	QueueName string
	Handle    *batch.Batch

	//
	AWS *session.Session
}

// ---------------------------------------------------------------------------
//
type MAWSBatchJob struct {
	//
	JobName       string
	JobDefinition string
	ArraySize     int
	Envs          map[string]string
}

// ---------------------------------------------------------------------------
// vytvoreni Handle na Batch Frontu
func (maws *MAWS) MAWSMakeBatchQueue(name string) *MAWSBatchQueue {
	//
	return &MAWSBatchQueue{
		//
		QueueName: name,
		Handle:    batch.New(maws.AWS),
		AWS:       maws.AWS,
	}
}

// ---------------------------------------------------------------------------
//
func (bq *MAWSBatchQueue) SubmitJobRec(rec MAWSBatchJob) (*batch.SubmitJobOutput, error) {
	//
	var _env = []*batch.KeyValuePair{}
	var _array *batch.ArrayProperties = nil

	//
	if len(rec.JobName) == 0 {
		//
		rec.JobName = MAWSUUID()
	}

	//
	if rec.ArraySize > 0 {
		//
		_array = &batch.ArrayProperties{Size: aws.Int64(int64(rec.ArraySize))}
	}

	//
	for k, v := range rec.Envs {
		//
		_env = append(_env, &batch.KeyValuePair{Name: aws.String(k), Value: aws.String(v)})
	}

	//
	input := &batch.SubmitJobInput{
		JobDefinition:   aws.String(rec.JobDefinition),
		JobName:         aws.String(rec.JobName),
		JobQueue:        aws.String(bq.QueueName),
		ArrayProperties: _array,
		ContainerOverrides: &batch.ContainerOverrides{
			Environment: _env,
		},
	}

	//
	result, err := bq.Handle.SubmitJob(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				fmt.Println(batch.ErrCodeClientException, aerr.Error())
			case batch.ErrCodeServerException:
				fmt.Println(batch.ErrCodeServerException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		//
		return result, err
	}

	//
	return result, err
}
