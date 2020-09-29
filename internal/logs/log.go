package logs

import (
	"encoding/json"
	"service-watch/internal/formatter"
	"service-watch/internal/store"
	"time"
)

type Log struct {
	Logs  map[string][]interface{}
	Store string
	Dir   string
}

func (l *Log) StoreLogs() error {
	store := store.Stores[l.Store](l.Dir)

	timestamp := time.Now()

	timestampBytes := formatter.MarshalDateTime(timestamp)

	responses := make([]interface{}, 0)

	for _, response := range l.Logs {
		responses = append(responses, response...)
	}

	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	responsesBytes, _ := json.Marshal(responses)

	keys = append(keys, timestampBytes)
	values = append(values, responsesBytes)

	for status, responses := range l.Logs {
		keys = append(keys, []byte("status:"+status))
		responsesBytes, _ := json.Marshal(responses)
		values = append(values, responsesBytes)
	}

	return store.PutBatch(keys, values)

}
