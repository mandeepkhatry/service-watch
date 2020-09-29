package formatter

import "time"

func MarshalDateTime(dt time.Time) []byte {
	byteKeyTimestamp, _ := dt.MarshalBinary()
	return byteKeyTimestamp

}
