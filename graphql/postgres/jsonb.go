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
	default:
		return fmt.Errorf("Unsure how to continue")
	}
	
	return j.UnmarshalJSON(buf.Bytes())
}
