package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TQuot struct {
	Category string `json:"Category"`
	Content  string `json:"Content"`
	Year     string `json:"Year"`
}

func SamplePost(w http.ResponseWriter, r *http.Request) error {
	var quatation TQuot
	json.NewDecoder(r.Body).Decode(&quatation)

	fmt.Fprintf(w, "%s %s", quatation.Category, quatation.Content)

	return nil
}
