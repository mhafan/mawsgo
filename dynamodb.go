package mawsgo

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// ---------------------------------------------------------------------------
// Abstrakce nad DynamoDB tabulkou. Zjednodusene schema tabulky, kde
// existuje pouze jeden primarni klic (KeyItem) typu string (typicky UUID)
// ---------------------------------------------------------------------------
// Predpokladame jiz existenci tabulky (zrizene v AWS)
type MAWSDynamoTable struct {
	// jmeno tabulky
	TableName string
	// jeden key atribut typu string !!!
	KeyItem string

	// handle na tabulku
	Handle *dynamodb.DynamoDB
	// handle na spojeni AWS (kopiruje se)
	AWS *session.Session
}

// ---------------------------------------------------------------------------
// vytvoreni Handle na Dynamo DB tabulku splnujici schema
func (maws *MAWS) MAWSMakeDynamoTable(name string, keyItem string) *MAWSDynamoTable {
	//
	return &MAWSDynamoTable{
		TableName: name,
		KeyItem:   keyItem,
		Handle:    dynamodb.New(maws.AWS),
		AWS:       maws.AWS,
	}
}

// ---------------------------------------------------------------------------
// vlozeni v zasade libovolneho typu struct do tabulky
func (tbl *MAWSDynamoTable) PutItem(itm interface{}) error {
	// konverze obsahu na mapu klic-hodnota
	_item, _err := dynamodbattribute.MarshalMap(itm)

	//
	if _err != nil {
		//
		return _err
	}

	// put
	_, __err := tbl.Handle.PutItem(&dynamodb.PutItemInput{
		Item:      _item,
		TableName: aws.String(tbl.TableName)})

	// je ok?
	return __err
}

// ---------------------------------------------------------------------------
// getitem, opet velmi dynamicky-volna typova kontrola
func (tbl *MAWSDynamoTable) GetItem(keyValue string, itm interface{}) error {
	//
	_get, _err := tbl.Handle.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tbl.TableName),
		Key:       DynamoDBStringKey(tbl.KeyItem, keyValue),
	})

	// neco neni oukej
	if _get.Item == nil || _err != nil {
		//
		return _err
	}

	// oukeja
	return dynamodbattribute.UnmarshalMap(_get.Item, &itm)
}

// ---------------------------------------------------------------------------
// deleteitem, key
func (tbl *MAWSDynamoTable) DeleteItem(keyValue string) error {
	//
	_, _err := tbl.Handle.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tbl.TableName),
		Key:       DynamoDBStringKey(tbl.KeyItem, keyValue),
	})

	// oukeja
	return _err
}

// ---------------------------------------------------------------------------
// keyValue - primarni klic zaznamu
// updateCmd - syntax odpovida pozadavkum AWS
// -- set COSI = :arg, COSIJ = :argj, ...
// args - key:value
func (tbl *MAWSDynamoTable) UpdateItem(keyValue string, updateCmd string, args map[string]interface{}) error {
	//
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: DynamoDBKeys(args),
		TableName:                 aws.String(tbl.TableName),
		Key:                       DynamoDBStringKey(tbl.KeyItem, keyValue),
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(updateCmd),
	}

	_, err := tbl.Handle.UpdateItem(input)

	//
	return err
}

// ---------------------------------------------------------------------------
// ValRecord - musi byt pole/slice !!!
// ---------------------------------------------------------------------------
// cond - filtracni condition
// LessThan, GreaterThan, Equal, *Equal
// filt := expression.Name("Age").GreaterThan(expression.Value(10))
func (tbl *MAWSDynamoTable) Scan(cond expression.ConditionBuilder, ValRecord interface{}) error {
	// vytvoreni zaznamu o filtru
	expr, _exprerr := expression.NewBuilder().WithFilter(cond).Build()

	//
	if _exprerr != nil {
		//
		return _exprerr
	}

	// parametry scanu
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tbl.TableName),
	}

	//
	res, scanerr := tbl.Handle.Scan(params)

	//
	if res.Items == nil || scanerr != nil {
		//
		return scanerr
	}

	//
	return dynamodbattribute.UnmarshalListOfMaps(res.Items, ValRecord)
}

// ---------------------------------------------------------------------------
//
func DynamoDBStringKey(key, value string) map[string]*dynamodb.AttributeValue {
	//
	return map[string]*dynamodb.AttributeValue{
		key: {
			S: aws.String(value),
		},
	}
}

// ---------------------------------------------------------------------------
//
func DynamoDBNumKey(key string, value int) map[string]*dynamodb.AttributeValue {
	//
	_val := strconv.Itoa(value)

	//
	return map[string]*dynamodb.AttributeValue{
		key: {
			N: aws.String(_val),
		},
	}
}

// ---------------------------------------------------------------------------
//
func DynamoDBKeys(asmap map[string]interface{}) map[string]*dynamodb.AttributeValue {
	//
	_val := make(map[string]*dynamodb.AttributeValue)

	//
	for k, v := range asmap {
		//
		switch tp := v.(type) {
		case int:
			_val[k] = &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(tp))}
		case float64:
			_tp := fmt.Sprintf("%f", tp)
			_val[k] = &dynamodb.AttributeValue{N: aws.String(_tp)}
		case string:
			_val[k] = &dynamodb.AttributeValue{S: aws.String(tp)}
		default:
			panic("neznam typ")
		}
	}

	//
	return _val
}
