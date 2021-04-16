package redisclient

import "encoding/json"

type encodeJson struct {
	value interface{}
}

func (this_ *encodeJson) MarshalBinary() (data []byte, err error) {
	return json.Marshal(this_.value)
}

type decodeJson struct {
	value interface{}
}

func (this_ *decodeJson) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, this_.value)
}

func NewDecoder(v interface{}) *decodeJson {
	return &decodeJson{value: v}
}

func NewEncoder(v interface{}) *encodeJson {
	return &encodeJson{value: v}
}
