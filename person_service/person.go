package personservice

import (
	"context"
	"errors"

	"github.com/jcox250/oto-person/gen/server"
)

// Logger defines the logger required by the person service
type Logger interface {
	Info(...interface{})
	Debug(...interface{})
	Error(...interface{})
}

// Person is an implementation of a person service
type Person struct {
	log  Logger
	data map[string]*server.Person
}

func New(l Logger) *Person {
	return &Person{
		log:  l,
		data: make(map[string]*server.Person),
	}
}

// Add adds a person
func (p *Person) Add(ctx context.Context, payload server.AddRequest) (*server.AddResponse, error) {
	person := &server.Person{
		Name: payload.Name,
		Age:  payload.Age,
	}

	p.data[person.Name] = person
	p.log.Debug("msg", "added person")
	return nil, nil
}

// Show shows a person
func (p *Person) Show(ctx context.Context, payload server.ShowRequest) (*server.ShowResponse, error) {
	person, ok := p.data[payload.Name]
	if !ok {
		p.log.Debug("msg", "person not found", "name", payload.Name)
		return nil, errors.New("not found")
	}
	return &server.ShowResponse{Name: person.Name, Age: person.Age}, nil
}
