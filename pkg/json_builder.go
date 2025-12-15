package pkg

import "encoding/json"

type JSONBBuilder map[string]any

func NewJSONB() JSONBBuilder {
	return make(JSONBBuilder)
}

func (j JSONBBuilder) Add(key string, value any) {
	if value == nil {
		return
	}
	j[key] = value
}

func (j JSONBBuilder) Bytes() ([]byte, error) {
	if len(j) == 0 {
		return []byte(`{}`), nil
	}
	return json.Marshal(j)
}
