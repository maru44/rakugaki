package utils

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	PKeyName   string
	PKeyType   string
	SKeyName   string
	SKeyType   string
	GSIName    string
	GSPKeyName string
	GSPKeyType string
	GSSKeyName string
	GSSKeyType string
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
	// projection type "ALL" "INCLUDE"
	createParams := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(key.PKeyName),
				AttributeType: aws.String(key.PKeyType),
			},
			/*
				{
					AttributeName: aws.String(key.SKeyName),
					AttributeType: aws.String(key.SKeyType),
				},
			*/
			{
				AttributeName: aws.String(key.GSPKeyName),
				AttributeType: aws.String(key.GSPKeyType),
			},
			{
				AttributeName: aws.String(key.GSSKeyName),
				AttributeType: aws.String(key.GSSKeyType),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(key.PKeyName),
				KeyType:       aws.String("HASH"),
			},
			/*
				{
					AttributeName: aws.String(key.SKeyName),
					KeyType:       aws.String("RANGE"),
				},
			*/
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String(key.GSIName),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String(key.GSPKeyName),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String(key.GSSKeyName),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
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

func CreateCountTable(tName string, key KeyDict) {
	db := AccessDB()

	tableName := tName
	createParams := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
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

const slugLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenRandSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = slugLetters[rand.Intn(len(slugLetters))]
	}
	return string(b)
}
