package main

import (
	"fmt"
	"localhost/rakugaki/quotation"
	"localhost/rakugaki/routes"
	"localhost/rakugaki/utils"

	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	// "encoding/json"
	//"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	// delete table
	//utils.DeleteTable("Quotations")

	// create table
	/*
		quots := utils.KeyDict{
			PKeyName: "Category",
			PKeyType: "S",
			SKeyName: "Number",
			SKeyType: "N",
		}
		utils.CreateTable("Quotations", quots)
	*/
	/*
		counts := utils.KeyDict{
			PKeyName: "Category",
			PKeyType: "S",
			SKeyName: "Count",
			SKeyType: "N",
		}
		utils.CreateTable("Counter", counts)
	*/
	db := utils.AccessDB()
	item := {
		Categoty: "GOOD",
		Count: 0,
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Quotations"),
		Item: ,
	}

	_, err := db.PutItem(input)
	if err != nil {
		fmt.Printf(err.Error())
	}

	/**********************************
	*                                 *
	*             Routes              *
	*                                 *
	**********************************/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", "home")
	})

	/***********   DB  ************/
	http.HandleFunc("/db/list", routes.Handle(utils.TableList))

	/**********  quatation  *********/
	http.HandleFunc("/quotation/post/sample", routes.Handle(utils.SamplePost))
	http.HandleFunc("/quotation/post", routes.Handle(quotation.PostQuot))
	http.HandleFunc("/quotation/", routes.Handle(quotation.ListQuot))

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
