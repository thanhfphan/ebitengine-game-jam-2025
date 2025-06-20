package entity

type Player struct {
	Entity
	Hand []*Card
}

func NewPlayer(name string, ttype Type, posX, posY float64) *Player {
	entity := NewEntity(ttype, name, posX, posY)

	return &Player{
		Entity: *entity,
		Hand:   []*Card{},
	}
}

func (p *Player) AddCard(card *Card) {
	p.Hand = append(p.Hand, card)
}

func (p *Player) RemoveCard(card *Card) {
	for i, c := range p.Hand {
		if c.ID == card.ID {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return
		}
	}
}

func (p *Player) RemoveCardAt(index int) *Card {
	if index < 0 || index >= len(p.Hand) {
		return nil
	}
	removeCard := p.Hand[index]
	p.Hand = append(p.Hand[:index], p.Hand[index+1:]...)
	return removeCard
}

func (p *Player) IsBot() bool {
	return p.EntityType == TypeBot
}
