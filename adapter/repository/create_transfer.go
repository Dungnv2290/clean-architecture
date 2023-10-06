package repository

import (
	"context"

	"github.com/dungnguyen/clean-architecture/domain/entity"
	"github.com/dungnguyen/clean-architecture/infrastructure/database"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Bson data
	createTransferBSON struct {
		ID        string `json:"id"`
		PayerID   string `json:"payer"`
		PayeeID   string `json:"payee"`
		Value     int64  `json:"value"`
		CreatedAt string `json:"created_at"`
	}

	createTransferRepository struct {
		handler    *database.MongoHandler
		collection string
	}
)

// NewCreateTransferRepository creates new createTransferRepository with its dependencies
func NewCreateTransferRepository(handler *database.MongoHandler) entity.TransferRepositoryCreator {
	return createTransferRepository{
		handler:    handler,
		collection: "transfers",
	}
}

// Create perform insertOne into database
func (c createTransferRepository) Create(ctx context.Context, t entity.Transfer) (entity.Transfer, error) {
	var bson = createTransferBSON{
		ID:        t.ID().Value(),
		PayerID:   t.Payer().Value(),
		PayeeID:   t.Payee().Value(),
		Value:     t.Value().Amount().Value(),
		CreatedAt: t.CreatedAt().String(),
	}

	if _, err := c.handler.Db().Collection(c.collection).InsertOne(ctx, bson); err != nil {
		return entity.Transfer{}, errors.Wrap(err, entity.ErrCreateTransfer.Error())
	}

	return t, nil
}

func (c createTransferRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := fn(sessCtx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	session, err := c.handler.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return nil
}
