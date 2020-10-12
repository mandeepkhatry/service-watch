package server

import (
	"encoding/json"
	"log"
	"net/http"
	"service-watch/internal/formatter"
	"service-watch/internal/logs"
	encode "service-watch/internal/server/response"
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
	data := make([]map[string]interface{}, 0)

	resultBytes, err := storeLog.Store.Get([]byte("status:" + status[0]))

	if err != nil {
		encode.JsonEncode(w, data, http.StatusInternalServerError)
	}

	json.Unmarshal(resultBytes, &ids)

	if ids == nil {
		encode.JsonEncode(w, data, http.StatusNotFound)
		return
	}

	keys := make([][]byte, 0)

	timeFormat, _ := utils.FormatConstantDate("2020-10-08T14:18:36.665826703+05:45")

	for _, id := range ids.([]interface{}) {

		t, _ := time.Parse(timeFormat, id.(string))

		byteID := formatter.MarshalDateTime(t)

		keys = append(keys, []byte("timestamp:"+string(byteID)))
	}
	byteKeys, _ := storeLog.Store.GetBatch(keys)

	for _, eachByteKey := range byteKeys {
		var eachResult map[string]interface{}
		json.Unmarshal(eachByteKey, &eachResult)
		data = append(data, eachResult)

	}

	encode.JsonEncode(w, data, http.StatusOK)
	return
}

func RunServer(store *logs.StoreLog) {

	storeLog = store

	router := mux.NewRouter()

	router.HandleFunc("/search", SearchHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}
