package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-lambda-go/lambda"
)

type Comment struct {
	Id      string `dynamodbav:"Id,hash`
	Content string `dynamodbav:"Content"`
}

func main() {
	sess := session.Must(session.NewSession())

	db := dynamodb.New(sess, &aws.Config{
		Region: aws.String("ap-northeast-1"),
		// Output: aws.String("json"),
		Endpoint: aws.String(os.Getenv("DYNAMO_ENDPOINT")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		),
	})
	// ここまでは動く

	// delete table

	delParams := &dynamodb.DeleteTableInput{
		TableName: aws.String("Comments"),
	}
	resp, err := db.DeleteTable(delParams)
	if err != nil {
		log.Printf("%v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", resp)
	})

	http.ListenAndServe(":8000", nil)
}
