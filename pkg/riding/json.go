package riding

import "encoding/json"

func EncodeMessage[V any](v V) ([]byte, error) {
	data, err := json.Marshal(v)
	return data, err
}

func DecodeMessage[V any](data []byte, v V) error {
	return json.Unmarshal(data, &v)
}
