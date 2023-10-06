package repository

import (
	"context"
	"time"

	"github.com/dungnguyen/clean-architecture/domain/entity"
	"github.com/dungnguyen/clean-architecture/infrastructure/database"
)

type (
	// Bson data
	createUserBSON struct {
		ID        string                 `bson:"id"`
		FullName  string                 `bson:"full_name"`
		Email     string                 `bson:"email"`
		Password  string                 `bson:"password"`
		Document  createUserDocumentBSON `bson:"document"`
		Wallet    createUserWalletBSON   `bson:"wallet"`
		Roles     createUserRolesBSON    `bson:"roles"`
		Type      string                 `bson:"type"`
		CreatedAt time.Time              `bson:"created_at"`
	}

	// Bson data
	createUserDocumentBSON struct {
		Type  string `bson:"type"`
		Value string `bson:"value"`
	}

	// Bson data
	createUserWalletBSON struct {
		Currency string `bson:"currency"`
		Amount   int64  `bson:"amount"`
	}

	// Bson data
	createUserRolesBSON struct {
		CanTransfer bool `bson:"can_transfer"`
	}

	createUserRepository struct {
		handler    *database.MongoHandler
		collection string
	}
)

// NewCreateUserRepository create new createUserRepository with it dependencies
func NewCreateUserRepository(handler *database.MongoHandler) entity.UserRepositoryCreator {
	return createUserRepository{
		handler:    handler,
		collection: "users",
	}
}

// Create perform insertOne into database
func (c createUserRepository) Create(ctx context.Context, u entity.User) (entity.User, error) {
	var bson = createUserBSON{
		ID:       u.ID().Value(),
		FullName: u.FullName().Value(),
		Email:    u.Email().Value(),
		Password: u.Password().Value(),
		Document: createUserDocumentBSON{
			Type:  u.Document().Type().String(),
			Value: u.Document().Value(),
		},
		Wallet: createUserWalletBSON{
			Currency: u.Wallet().Money().Currency().String(),
			Amount:   u.Wallet().Money().Amount().Value(),
		},
		Roles: createUserRolesBSON{
			CanTransfer: u.Roles().CanTransfer,
		},
		Type:      u.TypeUser().String(),
		CreatedAt: u.CreatedAt(),
	}

	if _, err := c.handler.Db().Collection(c.collection).InsertOne(ctx, bson); err != nil {
		return entity.User{}, err
	}

	return u, nil
}
