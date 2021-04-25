package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
)

func IsProductionEnv() bool {
	// 本番環境IPリスト
	hosts := []string{
		"aaa",
	}
	host, _ := os.Hostname()

	if runtime.GOOS != "linux" {
		return false
	}
	for _, v := range hosts {
		if v == host {
			return true
		}
	}
	return true
}

func _SetDefaultResponseHeader(w http.ResponseWriter) bool {
	protocol := "http"
	host := "localhost:3000"
	if IsProductionEnv() {
		protocol = "https"
		host = "aaa"
	}
	w.Header().Set("Access-Control-Allow-Origin", protocol+host)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	return true
}

func (json_ TypeJsonResponse) ResponseJsonWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(json_)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	_SetDefaultResponseHeader(w)
	w.Write(res)
	return true
}
