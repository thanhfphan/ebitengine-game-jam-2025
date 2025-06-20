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
	X, Y       float64
}

func NewEntity(entityType Type, name string, posX, posY float64) *Entity {
	return &Entity{
		ID:         uuid.NewString(),
		EntityType: entityType,
		Name:       name,
		X:          posX,
		Y:          posY,
	}
}
