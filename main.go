package main

import (
	"fmt"
	"localhost/rakugaki/quotation"
	"localhost/rakugaki/routes"
	"localhost/rakugaki/utils"

	"net/http"
	// "encoding/json"
	//"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	/**********************************
	*                                 *
	*             Routes              *
	*                                 *
	**********************************/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", "home")
	})

	/***********   DB  ************/
	http.HandleFunc("/db/list/", routes.Handle(utils.TableList))
	http.HandleFunc("/db/delete/", routes.Handle(utils.DeleteTableByParam))
	http.HandleFunc("/db/make/quotation/", routes.Handle(utils.MakeQuotationTable))
	http.HandleFunc("/db/make/counter/", routes.Handle(utils.MakeCounterTable))

	/**********  quatation  *********/
	http.HandleFunc("/quotation/counter/", routes.Handle(quotation.ListCounter))
	http.HandleFunc("/quotation/counter/good", routes.Handle(quotation.DetailCounterRes))

	http.HandleFunc("/quotation/random/", routes.Handle(quotation.GetRandomQuot))
	http.HandleFunc("/quotation/post", routes.Handle(quotation.PostQuot))
	// http.HandleFunc("/quotation/det/", routes.Handle(quotation.DetailQuotation))
	http.HandleFunc("/quotation/", routes.Handle(quotation.GetQuot))

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
