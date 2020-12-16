package postgres

import (
	"fmt"
	"bytes"
	"database/sql/driver"
	"encoding/json"
	_ "github.com/lib/pq"
)

// Jsonb provides an interface for interaction with
// json structures and properties when interfacing
// remote Databases or API
type Jsonb struct {
	json.RawMessage
}

// Value returns the data as intended for fs storage
func (j Jsonb) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

// Scan loads the data for jsonb storage in a database
func (j *Jsonb) Scan(value interface{}) error {
	var buf bytes.Buffer

	switch value := value.(type) {
	case string:
		buf.WriteString(value)
	case []byte:
		buf.Write(value)
	case map[string]interface{}:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}

		buf.Write(data)
	case interface{}:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}

		buf.Write(data)
	default:
		return fmt.Errorf("Unsure how to continue:\n%#+v", value)
	}

	//<<<///?
	//values := string(value).([]byte)
	//return j.UnmarshalJSON(values)
  ///>>>???

	return j.UnmarshalJSON(buf.Bytes())
}
