package json

import (
	"bytes"
	"encoding/json"

	"github.com/hailongz/kk-lib/dynamic"
)

func Unmarshal(data []byte, object interface{}) error {
	rd := bytes.NewReader(data)
	dec := json.NewDecoder(rd)
	dec.UseNumber()
	var v interface{} = nil
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	dynamic.SetValue(object, v)
	return nil
}
