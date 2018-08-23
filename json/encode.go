package json

import (
	"bytes"
	"encoding/json"
)

func Marshal(object interface{}) ([]byte, error) {

	w := bytes.NewBuffer(nil)

	enc := json.NewEncoder(w)

	err := enc.Encode(object)

	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func MarshalIndent(object interface{}, prefix, indent string) ([]byte, error) {

	w := bytes.NewBuffer(nil)

	enc := json.NewEncoder(w)

	enc.SetIndent(prefix, indent)

	err := enc.Encode(object)

	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
