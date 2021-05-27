package personservice

import (
	"context"

	"github.com/jcox250/oto-person/domain"
	"github.com/jcox250/oto-person/gen/server"
)

// Logger defines the logger required by the person service
type Logger interface {
	Info(...interface{})
	Debug(...interface{})
	Error(...interface{})
}

type PersonStore interface {
	AddPerson(ctx context.Context, p *domain.Person) error
	GetPerson(ctx context.Context, key string) (*domain.Person, error)
}

// Person is an implementation of a person service
type Person struct {
	log   Logger
	store PersonStore
}

func New(l Logger, ps PersonStore) *Person {
	return &Person{
		log:   l,
		store: ps,
	}
}

// Add adds a person
func (p *Person) Add(ctx context.Context, payload server.AddRequest) (*server.AddResponse, error) {
	person := &domain.Person{
		Name: payload.Name,
		Age:  payload.Age,
	}

	if err := p.store.AddPerson(ctx, person); err != nil {
		return nil, err
	}
	return nil, nil
}

// Show shows a person
func (p *Person) Show(ctx context.Context, payload server.ShowRequest) (*server.ShowResponse, error) {
	person, err := p.store.GetPerson(ctx, payload.Name)
	if err != nil {
		return nil, err
	}
	return &server.ShowResponse{Name: person.Name, Age: person.Age}, nil
}
