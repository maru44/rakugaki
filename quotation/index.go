package quotation

import (
	"encoding/json"
	"fmt"
	"localhost/rakugaki/utils"
	"log"
	"net/http"

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

type TQuot struct {
	Category string `json:"Category"`
	Content  string `json:"Content"`
	Number   int    `json:"Number"`
}

func (quot TQuotInputResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(quot)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	utils.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func PostQuot(w http.ResponseWriter, r *http.Request) error {
	result := TQuotInputResponse{Status: 200}

	var posted TQuot
	json.NewDecoder(r.Body).Decode(&posted)

	//now := time.Now().Format("2006-01-02 15:04:05")
	cat := posted.Category
	nowCount := DetailCounter(cat)
	counterUpdated := UpdateCounter(cat, nowCount)
	nowCount++
	posted.Number = nowCount

	if counterUpdated == false {
		return nil
	}

	quot, err := dynamodbattribute.MarshalMap(posted)

	db := utils.AccessDB()
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

	db := utils.AccessDB()
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
	result.ResponseWrite(w)
	return nil
}
