package quotation

import (
	"encoding/json"
	"fmt"
	"localhost/rakugaki/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type TCountersInputResponse struct {
	Status int        `json:"Status"`
	Data   []TCounter `json:"Data"`
}

type TCounterInputResponse struct {
	Status int `json:"Status"`
	Data   int `json:"Count"`
}

type TCounter struct {
	Category string `json:"Category"`
	Count    int    `json:"Count"`
}

func (count TCountersInputResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(count)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	utils.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (count TCounterInputResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(count)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	utils.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func ListCounter(w http.ResponseWriter, r *http.Request) error {
	result := TCountersInputResponse{Status: 200}

	db := utils.AccessDB()
	tableName := "Counters"
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	scanOut, err := db.Scan(input)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	var counts []TCounter
	for _, n := range scanOut.Items {
		var countTemp TCounter
		_ = dynamodbattribute.UnmarshalMap(n, &countTemp)
		counts = append(counts, countTemp)
	}

	result.Data = counts

	result.ResponseWrite(w)
	return nil
}

// 後で削除
func DetailCounterRes(w http.ResponseWriter, r *http.Request) error {
	result := TCounterInputResponse{Status: 200}

	db := utils.AccessDB()
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Counters"),
		Key: map[string]*dynamodb.AttributeValue{
			"Category": {
				S: aws.String("good"),
			},
		},
	}

	res, err := db.GetItem(input)
	if err != nil {
		fmt.Fprintf(w, "Hello, %v", err)
		return nil
	}
	count := TCounter{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &count)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	result.Data = count.Count

	result.ResponseWrite(w)
	return nil
}

func DetailCounter(cat string) int {
	db := utils.AccessDB()
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Counters"),
		Key: map[string]*dynamodb.AttributeValue{
			"Category": {
				S: aws.String(cat),
			},
		},
	}

	res, err := db.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}
	count := TCounter{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &count)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}

	return count.Count
}

func UpdateCounter(cat string, count int) bool {

	db := utils.AccessDB()
	// countは予約語 expressionattributenames で回避
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			/*
				":ncount": {
					N: aws.String(strconv.Itoa(count)),
				},
			*/
			":increment": {
				N: aws.String(strconv.Itoa(1)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#ct": aws.String("Count"),
		},
		TableName: aws.String("Counters"),
		Key: map[string]*dynamodb.AttributeValue{
			"Category": {
				S: aws.String(cat),
			},
		},
		//ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("add #ct :increment"),
	}

	_, err := db.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
