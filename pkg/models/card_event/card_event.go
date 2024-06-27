package event

import (
	"context"
	"time"

	"github.com/KBingsoo/entities/pkg/models"
)

type EventType string

const (
	Create       EventType = "card_create"
	Update       EventType = "card_update"
	OrderFulfill EventType = "card_order_fulfill"
	OrderRevert  EventType = "card_order_revert"
	Delete       EventType = "card_delete"

	Succeed EventType = "card_succeed"
	Error   EventType = "card_error"
)

type Event struct {
	Type    EventType       `json:"type"`
	Time    time.Time       `json:"time"`
	OrderID string          `json:"order_id"`
	Card    models.Card     `json:"card"`
	Context context.Context `json:"context"`
}
