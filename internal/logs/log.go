package logs

import (
	"encoding/json"
	"service-watch/internal/formatter"
	"service-watch/internal/store"
	"time"
)

type StoreLog struct {
	Store store.Store
}

func NewLog(storeSpace string, dir string) *StoreLog {
	return &StoreLog{
		Store: store.Stores[storeSpace](dir),
	}
}

func (l *StoreLog) StoreLogs(logs map[string][]interface{}) error {

	store := l.Store

	responses := make([]interface{}, 0)
	recordsByStatus := make(map[string][]time.Time)

	for status, response := range logs {
		if _, present := recordsByStatus[status]; !present {
			recordsByStatus[status] = make([]time.Time, 0)
		}
		for _, eachResponse := range response {
			recordsByStatus[status] = append(recordsByStatus[status], eachResponse.(map[string]interface{})["timestamp"].(time.Time))
		}
		responses = append(responses, response...)
	}

	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	for _, response := range responses {

		timestampByte := formatter.MarshalDateTime(response.(map[string]interface{})["timestamp"].(time.Time))
		keys = append(keys, []byte("timestamp:"+string(timestampByte)))

		responseByte, _ := json.Marshal(response)
		values = append(values, responseByte)

	}

	for status, recordTimstamps := range recordsByStatus {
		keys = append(keys, []byte("status:"+status))
		timestampsBytes, _ := json.Marshal(recordTimstamps)
		values = append(values, timestampsBytes)
	}

	return store.PutBatch(keys, values)

}
