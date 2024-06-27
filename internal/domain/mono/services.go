package mono

import (
	"context"
	"time"

	"github.com/KBingsoo/entities/pkg/models"
	event "github.com/KBingsoo/mono.git/pkg/models/card_event"
	"github.com/google/uuid"
)

type Manager interface {
	CreateCard(ctx context.Context, card *models.Card) error
	GetAllCards(ctx context.Context) ([]models.Card, error)
	GetCardByID(ctx context.Context, id string) (models.Card, error)
	UpdateCard(ctx context.Context, card *models.Card) error
	DeleteCard(ctx context.Context, id string) (models.Card, error)
	GetOrderByID(ctx context.Context, orderID string) (models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) error
	FulfillOrder(ctx context.Context, orderID string) error
}

type itemsMap struct {
	internalMap map[string]bool
	order       models.Order
	succeed     int
}

type manager struct {
	cardsRepository  CardsRepository
	ordersRepository OrdersRepository
	orderMap         map[string]itemsMap
	participants     []Participant //lista de participantes q manager conhece
}

func NewManager(cardsRepository CardsRepository, ordersRepository OrdersRepository, participants []Participant) *manager {
	return &manager{
		cardsRepository:  cardsRepository,
		ordersRepository: ordersRepository,
		participants:     participants,
	}
}

func (m *manager) CreateCard(ctx context.Context, card *models.Card) error {
	if card.ID == "" {
		card.ID = uuid.NewString()
	}

	transactionID := uuid.NewString()

	for _, participant := range m.participants {
		if err := participant.Prepare(ctx, transactionID); err != nil { // se a ação de "prepare" de algum participante falhar, faz rollback, retorna erro
			for _, p := range m.participants {
				p.Rollback(ctx, transactionID)
			}
			return err
		}
	}

	if err := m.cardsRepository.Upsert(ctx, card.ID, *card); err != nil { // se a ação de insert/update falhar, faz rollback, retorna erro
		for _, participant := range m.participants {
			participant.Rollback(ctx, transactionID)
		}
		return err
	}

	for _, participant := range m.participants {
		if err := participant.Commit(ctx, transactionID); err != nil { // se a ação de commit falhar, retorna erro
			return err
		}
	}

	return nil
}

func (m *manager) GetAllCards(ctx context.Context) ([]models.Card, error) {
	return m.cardsRepository.FindAll(ctx)
}

func (m *manager) GetCardByID(ctx context.Context, id string) (models.Card, error) {
	return m.cardsRepository.Find(ctx, id)
}

func (m *manager) UpdateCard(ctx context.Context, card *models.Card) error {
	transactionID := uuid.NewString()

	for _, participant := range m.participants {
		if err := participant.Prepare(ctx, transactionID); err != nil { // se a ação de "prepare" de algum participante falhar, faz rollback, retorna erro
			for _, p := range m.participants {
				p.Rollback(ctx, transactionID)
			}
			return err
		}
	}

	if err := m.cardsRepository.Upsert(ctx, card.ID, *card); err != nil { // se a ação de insert/update falhar, faz rollback, retorna erro
		for _, participant := range m.participants {
			participant.Rollback(ctx, transactionID)
		}
		return err
	}

	for _, participant := range m.participants {
		if err := participant.Commit(ctx, transactionID); err != nil { // se a ação de commit falhar, retorna erro
			return err
		}
	}

	return nil
}

func (m *manager) DeleteCard(ctx context.Context, id string) (models.Card, error) {
	transactionID := uuid.NewString()

	for _, participant := range m.participants {
		if err := participant.Prepare(ctx, transactionID); err != nil { // se a ação de "prepare" de algum participante falhar, faz rollback, retorna erro
			for _, p := range m.participants {
				p.Rollback(ctx, transactionID)
			}
			return models.Card{}, err
		}
	}

	card, err := m.cardsRepository.Delete(ctx, id)
	if err != nil { // se a ação de deletar falhar, faz rollback, retorna erro
		for _, participant := range m.participants {
			participant.Rollback(ctx, transactionID)
		}
		return models.Card{}, err
	}

	for _, participant := range m.participants {
		if err := participant.Commit(ctx, transactionID); err != nil {
			return models.Card{}, err
		}
	}

	return card, nil
}

func (m *manager) GetOrderByID(ctx context.Context, orderID string) (models.Order, error) {
	return m.ordersRepository.Find(ctx, orderID)
}

func (m *manager) CreateOrder(ctx context.Context, order *models.Order) error {
	if order.ID == "" {
		order.ID = uuid.NewString()
	}

	return m.ordersRepository.Upsert(ctx, order.ID, *order)
}

func (m *manager) FulfillOrder(ctx context.Context, orderID string) error {
	order, err := m.ordersRepository.Find(ctx, orderID)
	if err != nil {
		return err
	}

	m.orderMap[orderID] = itemsMap{
		internalMap: make(map[string]bool),
		order:       order,
	}

	for _, card := range order.Cards {
		event := event.Event{
			Type: event.OrderFulfill,
			Time: time.Now(),
			Card: models.Card{
				ID: card,
			},
			OrderID: orderID,
			Context: ctx,
		}

		m.orderMap[orderID].internalMap[card] = false
	}

	return nil
}
