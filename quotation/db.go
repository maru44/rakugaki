package quotation

import (
	"fmt"
	"localhost/rakugaki/utils"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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

type TQuotInputDetResponse struct {
	Status int   `json:"Status"`
	Data   TQuot `json:"Data"`
}

type TQuot struct {
	Slug     string `json:"Slug"`
	Category string `json:"Category"`
	Content  string `json:"Content"`
	Number   int    `json:"Number"`
}

/**************  functions  ******************/

func getListQuot() []TQuot {
	db := utils.AccessDB()
	tableName := "Quotations"
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	scanOut, err := db.Scan(input)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return nil
	}

	var quots []TQuot
	for _, n := range scanOut.Items {
		var quotTemp TQuot
		_ = dynamodbattribute.UnmarshalMap(n, &quotTemp)
		quots = append(quots, quotTemp)
	}

	return quots
}

func getDetailQuot(slug string) []TQuot {
	db := utils.AccessDB()
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Quotations"),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(slug),
			},
		},
	}

	res, err := db.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	quot := TQuot{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &quot)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	ret := []TQuot{quot}

	return ret
}

func getSlugByCat(cat string, num string) string {
	return "aaa"
}

func getDetailQuotByCat(cat string, num string) []TQuot {
	db := utils.AccessDB()
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Quotations"),
		Key: map[string]*dynamodb.AttributeValue{
			"Category": {
				S: aws.String(cat),
			},
			"Number": {
				N: aws.String(num),
			},
		},
	}

	res, err := db.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	quot := TQuot{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &quot)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	ret := []TQuot{quot}

	return ret
}

func getQuotByCatNum(cat string, num int) []TQuot {
	db := utils.AccessDB()
	input := &dynamodb.QueryInput{
		TableName: aws.String("Quotations"),
		IndexName: aws.String("Cat-Num-Index"),
		ExpressionAttributeNames: map[string]*string{
			"#Category": aws.String("Category"),
			"#Number":   aws.String("Number"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cat": {
				S: aws.String(cat),
			},
			":num": {
				N: aws.String(strconv.Itoa(num)),
			},
		},
		KeyConditionExpression: aws.String("#Category = :cat AND #Number = :num"),
	}

	res, err := db.Query(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	quot := TQuot{}
	err = dynamodbattribute.UnmarshalMap(res.Items[0], &quot)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	ret := []TQuot{quot}

	return ret
}

func randomQuotation(cat string) []TQuot {
	count := DetailCounter(cat)

	var ret []TQuot
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		var n int
		n = rand.Intn(count) + 1
		quots := getQuotByCatNum(cat, n)
		if quots[0].Number != 0 {
			ret = quots
			break
		}
	}

	return ret
}
