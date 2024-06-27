package mono

import (
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/literalog/go-wise/wise"
)

type CardsRepository wise.MongoRepository[models.Card]

type OrdersRepository wise.MongoRepository[models.Order]
