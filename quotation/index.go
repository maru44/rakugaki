package quotation

import (
	"encoding/json"
	"fmt"
	"localhost/rakugaki/utils"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
	utils.SetDefaultResponseHeader(w)
	result := TQuotInputResponse{Status: 200}

	var posted TQuot
	json.NewDecoder(r.Body).Decode(&posted)
	fmt.Print(posted)

	// generate slug
	rand.Seed(time.Now().UnixNano())
	slug := utils.GenRandSlug(12)
	posted.Slug = slug

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

// list or retreive
func GetQuot(w http.ResponseWriter, r *http.Request) error {
	result := TQuotInputResponse{Status: 200}

	query := r.URL.Query()
	/*
		cat := query.Get("c")
		num := query.Get("n")
	*/
	slug := query.Get("s")

	var quots []TQuot
	/*
		if cat != "" && num != "" {
			quots = getDetailQuot(cat, num)
		} else {
			quots = getListQuot()
		}
	*/
	quots = getListQuot()
	if slug != "" {
		quots = getDetailQuot(slug)
	}

	result.Data = quots
	result.ResponseWrite(w)
	return nil
}

func GetRandomQuot(w http.ResponseWriter, r *http.Request) error {
	result := TQuotInputResponse{Status: 200}
	query := r.URL.Query()
	cat := query.Get("c")

	quots := randomQuotation(cat)
	result.Data = quots
	result.ResponseWrite(w)
	return nil
}

func DeleteQuotController(w http.ResponseWriter, r *http.Request) error {
	utils.SetDefaultResponseHeader(w)
	query := r.URL.Query()
	slug := query.Get("s")

	deleteQuot(slug)

	result := utils.TypeJsonResponse{Status: 200}

	var datas []utils.BaseJson
	var data utils.BaseJson
	data.Key = "slug"
	data.Value = slug
	datas = append(datas, data)

	result.Data = datas
	result.ResponseJsonWrite(w)
	return nil
}

/*************  以下不要  **************/

func (quot TQuotInputDetResponse) ResponseWrite(w http.ResponseWriter) bool {
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

// only detail
/*
func DetailQuotation(w http.ResponseWriter, r *http.Request) error {
	result := TQuotInputDetResponse{Status: 200}

	query := r.URL.Query()
	cat := query.Get("c")
	num := query.Get("n")

	quots := getDetailQuot(cat, num)
	quot := quots[0]

	result.Data = quot
	result.ResponseWrite(w)

	return nil
}
*/
