package entity

import "github.com/google/uuid"

type Type int

const (
	TypePlayer Type = iota
	TypeBot
	TypeCard
)

type Entity struct {
	ID         string
	Name       string
	EntityType Type
}

func NewEntity(entityType Type, name string) *Entity {
	return &Entity{
		ID:         uuid.NewString(),
		EntityType: entityType,
		Name:       name,
	}
}
