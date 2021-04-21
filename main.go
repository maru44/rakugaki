package main

import (
	"fmt"
	"localhost/rakugaki/utils"
	"net/http"
	// "encoding/json"
	//"github.com/aws/aws-lambda-go/lambda"
)

type Contents struct {
	Person  string `dynamodbav:"Person,hash"`
	Year    string `dynamodbav:"Year,range"`
	Content string `dynamodbav:"Content"`
}

func main() {

	// delete table
	/*
		delParams := &dynamodb.DeleteTableInput{
			TableName: aws.String("Contents"),
		}
		resp, err := db.DeleteTable(delParams)
		if err != nil {
			log.Printf("%v", err)
		}
	*/

	// create table
	/*
		tableName := "Contents"
		createParams := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("Year"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("Person"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Person"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Year"),
					KeyType:       aws.String("RANGE"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(tableName),
		}

		_, err := db.CreateTable(createParams)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	*/

	// table list
	/*
		tableListInput := &dynamodb.ListTablesInput{}
		tableList, err := db.ListTables(tableListInput)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	*/

	tableKey := map[string]string{
		"pKeyName": "Person",
		"pkeyType": "S",
		"sKeyName": "Year",
		"sKeyType": "S",
	}

	tableKeyL := []string{
		"Person",
		"S",
		"Year",
		"S",
	}

	utils.CreateTable("Contents", tableKeyL...)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", "resp")
	})

	http.ListenAndServe(":8080", nil)
}

func isExist(el string, list []string) bool {
	for _, listEl := range list {
		if listEl == el {
			return true
		}
	}
	return false
}
