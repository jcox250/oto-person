package cache

import (
	"context"

	"github.com/jcox250/oto-person/domain"
)

type Logger interface {
	Info(keyvals ...interface{})
	Debug(keyvals ...interface{})
	Error(keyvals ...interface{})
}

type KVAddGetter interface {
	Add(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type PersonCache struct {
	log Logger
	kv  KVAddGetter
}

func NewPersonCache(l Logger, kv KVAddGetter) *PersonCache {
	return &PersonCache{
		log: l,
		kv:  kv,
	}
}

func (p *PersonCache) AddPerson(ctx context.Context, person *domain.Person) error {
	if err := p.kv.Add(ctx, person.ID, person); err != nil {
		p.log.Error("msg", "failed adding person from key value store", "err", err)
		return ErrInternal
	}
	return nil
}

func (p *PersonCache) GetPerson(ctx context.Context, key string) (*domain.Person, error) {
	b, err := p.kv.Get(ctx, key)
	if err != nil {
		p.log.Error("msg", "failed getting person from key value store", "err", err)
		return nil, ErrInternal
	}

	if b == nil {
		return nil, ErrNotFound
	}

	person := &domain.Person{}
	if err := person.UnmarshalBinary(b); err != nil {
		p.log.Error("msg", "failed to unmarshal bytes from key value store to person", "err", err)
		return nil, ErrInternal
	}

	return person, nil
}
