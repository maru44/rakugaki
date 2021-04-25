package utils

import (
	"encoding/json"
	"fmt"
	"time"

	//"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	//"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	//"github.com/aws/aws-lambda-go/lambda"
	//"github.com/guregu/null"
)

// status(int) & data
type TypeJsonResponse struct {
	Status int        `json:"status"`
	Data   []BaseJson `json:"data"`
}

// strKey & string
type BaseJson struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ListJson struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

type KeyDict struct {
	PKeyName string
	PKeyType string
	SKeyName string
	SKeyType string
}

func AccessDB() *dynamodb.DynamoDB {
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

	return db
}

func CreateTable(tName string, key KeyDict) {
	db := AccessDB()

	tableName := tName
	createParams := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(key.SKeyName),
				AttributeType: aws.String(key.SKeyType),
			},
			{
				AttributeName: aws.String(key.PKeyName),
				AttributeType: aws.String(key.PKeyType),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(key.PKeyName),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(key.SKeyName),
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
}

func DeleteTable(tName string) {
	db := AccessDB()

	tableName := tName
	delParams := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	_, err := db.DeleteTable(delParams)
	if err != nil {
		log.Printf("%v", err)
	}
}

func TableList(w http.ResponseWriter, r *http.Request) error {
	result := TypeJsonResponse{Status: 200}

	db := AccessDB()

	tableListInput := &dynamodb.ListTablesInput{}
	tableList, err := db.ListTables(tableListInput)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	i := 0
	var tables []BaseJson
	for _, n := range tableList.TableNames {
		i++
		s := strconv.Itoa(i)
		temp := BaseJson{s, *n}
		tables = append(tables, temp)
	}

	result.Data = tables
	result.ResponseJsonWrite(w)
	return nil
}

/**************  validation  **************/
type TValid struct {
	Status int    `json:"status"`
	Valid  string `json:"valid"`
}

/**************  quotation post *************/
type TQuotInputResponse struct {
	Status int     `json:"Status"`
	Data   []TQuot `json:"Data"`
}

type TQuot struct {
	Category string `json:"Category"`
	Content  string `json:"Content"`
	Year     string `json:"Year"`
}

func (quot TQuotInputResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(quot)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	_SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func PostQuot(w http.ResponseWriter, r *http.Request) error {
	result := TQuotInputResponse{Status: 200}

	var posted TQuot
	json.NewDecoder(r.Body).Decode(&posted)

	now := time.Now().Format("2006-01-02 15:04:05")
	posted.Year = now

	quot, err := dynamodbattribute.MarshalMap(posted)

	db := AccessDB()
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Quotations"),
		Item:      quot,
	}

	_, err = db.PutItem(input)
	if err != nil {
		fmt.Printf(err.Error())
		return nil
	}

	var postedList []TQuot
	postedList = append(postedList, posted)
	result.Data = postedList
	result.ResponseWrite(w)

	return nil
}

func ListQuot(w http.ResponseWriter, r *http.Request) error {
	result := TQuotInputResponse{Status: 200}

	db := AccessDB()
	tableName := "Quotations"
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	scanOut, err := db.Scan(input)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	var quots []TQuot
	for _, n := range scanOut.Items {
		var quotTemp TQuot
		_ = dynamodbattribute.UnmarshalMap(n, &quotTemp)
		quots = append(quots, quotTemp)
	}

	result.Data = quots
	fmt.Fprintf(w, "%v %v", result.Status, result.Data)
	result.ResponseWrite(w)
	return nil
}
