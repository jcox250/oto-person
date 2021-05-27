package domain

import "encoding/json"

type Person struct {
	ID   string
	Name string
	Age  int
}

func (p *Person) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
