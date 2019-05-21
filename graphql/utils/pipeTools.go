package utils

import (
	"io"
	"encoding/json"
)

func PipeInput(from, to interface{}) (err error) {
  jReader, jWriter := io.Pipe()
  defer jReader.Close()
  
  go func() {
    defer jWriter.Close()
	  err = json.NewEncoder(jWriter).Encode(from)
	}()
	
	if err != nil {
	  return
	}
	
	return json.NewDecoder(jReader).Decode(to)
}
