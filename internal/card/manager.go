package card

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	mrand "math/rand"

	"github.com/thanhfphan/ebitengj2025/assets/configs"
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

type Manager struct {
	Deck []*entity.Card
	rand *mrand.Rand

	TableStack *entity.TableStack
	OnDishMade func(recipe *entity.Card)
	OnPlayCard func(player *entity.Player, card *entity.Card)
}

func NewManager() *Manager {
	var seed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &seed)

	mgr := &Manager{
		rand:       mrand.New(mrand.NewSource(seed)),
		TableStack: entity.NewTableStack(),
	}

	if err := mgr.LoadDeck("default"); err != nil {
		panic(err)
	}
	return mgr
}

func (m *Manager) LoadDeck(theme string) error {
	m.Deck = []*entity.Card{}
	m.TableStack = entity.NewTableStack()

	var ingFile IngredientFile
	var rcpFile RecipeFile

	if err := json.Unmarshal(configs.DefaultIngredientsJSON, &ingFile); err != nil {
		return err
	}
	if err := json.Unmarshal(configs.DefaultRecipesJSON, &rcpFile); err != nil {
		return err
	}

	mapIng := make(map[string]IngredientConfig)
	for _, ing := range ingFile.Ingredients {
		mapIng[ing.ID] = ing
	}

	for _, r := range rcpFile.Recipes {
		card := &entity.Card{
			Entity:              *entity.NewEntity(entity.TypeCard, r.Name),
			Type:                entity.CardTypeRecipe,
			RequiredIngredients: r.Requires,
		}
		m.Deck = append(m.Deck, card)

		for _, ingID := range r.Requires {
			ing, ok := mapIng[ingID]
			if !ok {
				return fmt.Errorf("recipe %s requires unknown ingredient %s", r.ID, ingID)
			}

			m.Deck = append(m.Deck, &entity.Card{
				Entity:       *entity.NewEntity(entity.TypeCard, ing.Name),
				Type:         entity.CardTypeIngredient,
				IngredientID: ingID,
			})
		}
	}

	m.shuffle(m.Deck)
	return nil
}

func (m *Manager) GetMapIngredientNames() map[string]string {
	result := make(map[string]string)
	for _, card := range m.Deck {
		if card.Type == entity.CardTypeIngredient {
			result[card.IngredientID] = card.Name
		}
	}
	return result
}

// shuffle performs Fisher‑Yates in‑place on the given slice.
func (m *Manager) shuffle(cards []*entity.Card) {
	for i := len(cards) - 1; i > 0; i-- {
		j := m.rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func (m *Manager) DealHands(players []*entity.Player) {
	j := 0
	numPlayers := len(players)
	for _, card := range m.Deck {
		players[j%numPlayers].AddCard(card)
		j++
	}
}

func (m *Manager) PlayCard(player *entity.Player, cardID string) error {
	removeCard := player.GetCard(cardID)
	if removeCard == nil {
		return fmt.Errorf("invalid card id: %s", cardID)
	}
	player.RemoveCard(removeCard.ID)

	m.TableStack.AddCard(removeCard, player.ID)

	if m.OnPlayCard != nil {
		m.OnPlayCard(player, removeCard)
	}

	return nil
}

func (m *Manager) TryMakeDish() bool {
	var recipes []*entity.Card
	var ingredients []*entity.Card

	for _, card := range m.TableStack.GetAllCardsInReverseOrder() {
		if card.Type == entity.CardTypeRecipe {
			recipes = append(recipes, card)
		} else if card.Type == entity.CardTypeIngredient {
			ingredients = append(ingredients, card)
		}
	}

	for _, r := range recipes {
		need := make(map[string]int)
		for _, ing := range r.RequiredIngredients {
			need[ing]++
		}

		var usedIDs []string
		usedCount := 0

		for _, ing := range ingredients {
			if cnt, ok := need[ing.IngredientID]; ok && cnt > 0 {
				need[ing.IngredientID]--
				usedIDs = append(usedIDs, ing.ID)
				usedCount++
				if usedCount == len(r.RequiredIngredients) {
					break
				}
			}
		}

		if usedCount == len(r.RequiredIngredients) {
			m.TableStack.RemoveCard(r.ID)
			for _, id := range usedIDs {
				m.TableStack.RemoveCard(id)
			}
			if m.OnDishMade != nil {
				m.OnDishMade(r)
			}
			return true
		}
	}

	return false
}
