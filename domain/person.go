package domain

import "encoding/json"

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *Person) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Person) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
