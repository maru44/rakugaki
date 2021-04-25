package main

import (
	"fmt"
	"localhost/rakugaki/routes"
	"localhost/rakugaki/utils"

	"net/http"
	// "encoding/json"
	//"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	/*
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
	*/

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

	http.HandleFunc("/create/table/quot", func(w http.ResponseWriter, r *http.Request) {
		quotes := utils.KeyDict{
			PKeyName: "Category",
			PKeyType: "S",
			SKeyName: "Year",
			SKeyType: "S",
		}
		utils.CreateTable("Quotations", quotes)

		fmt.Fprintf(w, "aaa %v", "bbb")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", "home")
	})

	/***********   DB  ************/
	http.HandleFunc("/db/list", routes.Handle(utils.TableList))

	/**********  quatation  *********/
	http.HandleFunc("/quotation/post/sample", routes.Handle(utils.SamplePost))
	http.HandleFunc("/quotation/post", routes.Handle(utils.PostQuot))
	http.HandleFunc("/quotation/", routes.Handle(utils.ListQuot))

	/***********  Serve  **********/
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
