package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type TOperation struct {
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}

func (operation TOperation) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(operation)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func DeleteTableByParam(w http.ResponseWriter, r *http.Request) error {
	result := TOperation{Status: 200}
	query := r.URL.Query()
	tName := query.Get("table")

	if tName != "" {
		DeleteTable(tName)
		result.Message = "delete success"
	} else {
		result.Status = 400
		result.Message = "failure"
	}

	result.ResponseWrite(w)
	return nil
}

func MakeQuotationTable(w http.ResponseWriter, r *http.Request) error {
	result := TOperation{Status: 200}

	quots := KeyDict{
		PKeyName:   "Slug",
		PKeyType:   "S",
		GSIName:    "Cat-Num-Index",
		GSPKeyName: "Category",
		GSPKeyType: "S",
		GSSKeyName: "Number",
		GSSKeyType: "N",
	}
	CreateTable("Quotations", quots)

	result.Message = "success"

	result.ResponseWrite(w)
	return nil
}

func MakeCounterTable(w http.ResponseWriter, r *http.Request) error {
	result := TOperation{Status: 200}

	counts := KeyDict{
		PKeyName: "Category",
		PKeyType: "S",
	}
	CreateCountTable("Counters", counts)

	type TCounter struct {
		Category string
		Count    int
	}
	db := AccessDB()
	item := TCounter{
		Category: "good",
		Count:    0,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
		result.Status = 400
		result.Message = "failure"
		return nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Counters"),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		fmt.Printf(err.Error())
		result.Status = 400
		result.Message = "failure"
		return nil
	}

	result.Message = "success"
	result.ResponseWrite(w)
	return nil
}
