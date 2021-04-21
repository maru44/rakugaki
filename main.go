package main

import (
	"fmt"
	"net/http"
	"os"

	// "encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-lambda-go/lambda"
)

type Comment struct {
	Person  string `dynamodbav:"Person,hash"`
	Year    string `dynamodbav:"Year,range"`
	Content string `dynamodbav:"Content"`
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db := dynamodb.New(sess, &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String(os.Getenv("DYNAMO_ENDPOINT")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		),
	})

	// delete table

	/*
		delParams := &dynamodb.DeleteTableInput{
			TableName: aws.String("Comments"),
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", db)
	})

	http.ListenAndServe(":8080", nil)
}
