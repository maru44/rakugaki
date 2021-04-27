package quotation

import (
	"fmt"
	"localhost/rakugaki/utils"
	"log"

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

func getDetailQuot(cat string, num string) []TQuot {
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
