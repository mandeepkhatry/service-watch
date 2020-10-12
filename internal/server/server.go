package server

import (
	"encoding/json"
	"log"
	"net/http"
	"service-watch/internal/formatter"
	"service-watch/internal/logs"
	"service-watch/internal/utils"
	"time"

	"github.com/gorilla/mux"
)

var logsDir string
var storeLog *logs.StoreLog

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	status, _ := query["status"]

	var ids interface{}

	resultBytes, _ := storeLog.Store.Get([]byte("status:" + status[0]))

	json.Unmarshal(resultBytes, &ids)

	keys := make([][]byte, 0)

	timeFormat, _ := utils.FormatConstantDate("2020-10-08T14:18:36.665826703+05:45")

	for _, id := range ids.([]interface{}) {

		t, _ := time.Parse(timeFormat, id.(string))

		byteID := formatter.MarshalDateTime(t)

		keys = append(keys, []byte("timestamp:"+string(byteID)))
	}
	byteKeys, _ := storeLog.Store.GetBatch(keys)

	data := make([]map[string]interface{}, 0)

	for _, eachByteKey := range byteKeys {
		var eachResult map[string]interface{}
		json.Unmarshal(eachByteKey, &eachResult)
		data = append(data, eachResult)

	}

	result := make(map[string]interface{})
	result["data"] = data

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

func RunServer(store *logs.StoreLog) {

	storeLog = store

	router := mux.NewRouter()

	router.HandleFunc("/search", SearchHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}
