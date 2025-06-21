package card

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"os"
	"sort"

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

	ingData, err := os.ReadFile("assets/configs/decks/" + theme + "/ingredients.json")
	if err != nil {
		return err
	}
	rcpData, err := os.ReadFile("assets/configs/decks/" + theme + "/recipes.json")
	if err != nil {
		return err
	}

	var ingFile IngredientFile
	var rcpFile RecipeFile

	if err := json.Unmarshal(ingData, &ingFile); err != nil {
		return err
	}
	if err := json.Unmarshal(rcpData, &rcpFile); err != nil {
		return err
	}

	mapIng := make(map[string]IngredientConfig)
	for _, ing := range ingFile.Ingredients {
		mapIng[ing.ID] = ing
	}

	for _, r := range rcpFile.Recipes {
		card := &entity.Card{
			Entity:              *entity.NewEntity(entity.TypeCard, r.Name, 0, 0),
			Type:                entity.CardRecipe,
			RequiredIngredients: r.Requires,
		}
		m.Deck = append(m.Deck, card)

		for _, ingID := range r.Requires {
			ing, ok := mapIng[ingID]
			if !ok {
				return fmt.Errorf("recipe %s requires unknown ingredient %s", r.ID, ingID)
			}

			m.Deck = append(m.Deck, &entity.Card{
				Entity:       *entity.NewEntity(entity.TypeCard, ing.Name, 0, 0),
				Type:         entity.CardIngredient,
				IngredientID: ingID,
			})
		}
	}

	m.shuffle(m.Deck)
	return nil
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

func (m *Manager) PlayCard(player *entity.Player, handIndex int) error {
	removeCard := player.RemoveCardAt(handIndex)
	if removeCard == nil {
		return fmt.Errorf("invalid card index: %d", handIndex)
	}

	m.TableStack.AddCard(removeCard, player.ID)

	if m.OnPlayCard != nil {
		m.OnPlayCard(player, removeCard)
	}

	return nil
}

func (m *Manager) TryMakeDish() bool {
	if len(m.TableStack.Receipes) == 0 {
		return false
	}

	// LIFO traversal
	for rIdx := len(m.TableStack.Receipes) - 1; rIdx >= 0; rIdx-- {
		recipe := m.TableStack.Receipes[rIdx]

		need := make(map[string]int)
		for _, ing := range recipe.RequiredIngredients {
			need[ing]++
		}

		var usedIngIdx []int
		var usedCardIDs []string
		for iIdx, card := range m.TableStack.Ingredients {
			if cnt, ok := need[card.IngredientID]; ok && cnt > 0 {
				need[card.IngredientID]--
				usedIngIdx = append(usedIngIdx, iIdx)
				usedCardIDs = append(usedCardIDs, card.ID)
			}
		}

		missing := false
		for _, c := range need {
			if c > 0 {
				missing = true
				break
			}
		}
		if missing {
			continue
		}

		m.TableStack.RemoveReceipeAt(rIdx)
		sort.Sort(sort.Reverse(sort.IntSlice(usedIngIdx)))
		for _, idx := range usedIngIdx {
			m.TableStack.RemoveIngredientAt(idx)
		}
		for _, cardID := range usedCardIDs {
			m.TableStack.RemoveCardFromPlayerTracking(cardID)
		}

		if m.OnDishMade != nil {
			m.OnDishMade(recipe)
		}

		return true
	}

	return false
}
