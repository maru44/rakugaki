package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SamplePost(w http.ResponseWriter, r *http.Request) error {
	var quatation TQuot
	json.NewDecoder(r.Body).Decode(&quatation)

	fmt.Fprintf(w, "%s %s", quatation.Category, quatation.Content)

	return nil
}
