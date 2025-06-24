package entity

type Player struct {
	Entity
	Hand      map[string]*Card
	OrderHand []string
}

func NewPlayer(name string, ttype Type, posX, posY float64) *Player {
	entity := NewEntity(ttype, name, posX, posY)

	return &Player{
		Entity:    *entity,
		Hand:      make(map[string]*Card),
		OrderHand: []string{},
	}
}

func (p *Player) AddCard(card *Card) {
	p.Hand[card.ID] = card
	p.OrderHand = append(p.OrderHand, card.ID)
}

func (p *Player) GetCard(id string) *Card {
	return p.Hand[id]
}

func (p *Player) GetCards() []*Card {
	cards := make([]*Card, 0, len(p.Hand))
	for _, card := range p.Hand {
		cards = append(cards, card)
	}
	return cards
}

func (p *Player) RemoveCard(id string) {
	for i, cardID := range p.OrderHand {
		if cardID == id {
			p.OrderHand = append(p.OrderHand[:i], p.OrderHand[i+1:]...)
			delete(p.Hand, id)
			break
		}
	}
}

func (p *Player) IsBot() bool {
	return p.EntityType == TypeBot
}
