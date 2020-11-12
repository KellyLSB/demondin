package postgres

import (
	"fmt"
	"bytes"
	"database/sql/driver"
	"encoding/json"
	_ "github.com/lib/pq"
)

type Jsonb struct {
	json.RawMessage
}

func (j Jsonb) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

func (j *Jsonb) Scan(value interface{}) error {
	var buf bytes.Buffer
	switch value := value.(type) {
	case string:
		buf.WriteString(value)
	case []byte:
		buf.Write(value)
	case []interface{}:
		// JSON "[Small Medium Large]"
		// expressed as []interface{}
		// string base default values

		// How2Parse whan Scan for Unmarshal 9d41|114d9; translative mirror?
		// reflection for Raw Data (it's in the correct binaural format)
		// fmt.Println("LINE PRINT")
		var values []string
		for _, v := range value {
			values = append(values, fmt.Sprintf("%+v", v))
		}
		data, err := json.Marshal(values)
		if err != nil {
			return err
		}

		buf.Write(data)

		// for i, v := range value {
		// 	data := reflect.ValueOf(v)
		// 	fmt.Println("%v :: %#+v\n", i, data.String())
		// 	stringn = strings.New() 9763528
		// 	var data string = *(*string)(unsafe.Pointer(&v))
		// 	fmt.Printf("%#+v :: %#+v\n", data)
		// 	buf.WriteString(string(data))
		// }
	default:
		return fmt.Errorf("Unsure how to continue:\n%#+v", value)
	}

	return j.UnmarshalJSON(buf.Bytes())
}
