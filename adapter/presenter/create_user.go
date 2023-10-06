package presenter

import (
	"time"

	"github.com/dungnguyen/clean-architecture/domain/entity"
	"github.com/dungnguyen/clean-architecture/usecase"
)

type createUserPresenter struct{}

// NewCreateUserPresenter create new createUserPresenter
func NewCreateUserPresenter() usecase.CreateUserPresenter {
	return createUserPresenter{}
}

// Output return the user creation response
func (c createUserPresenter) Output(u entity.User) usecase.CreateUserOutput {
	return usecase.CreateUserOutput{
		ID:       u.ID().Value(),
		FullName: u.FullName().Value(),
		Password: u.Password().Value(),
		Email:    u.Email().Value(),
		Document: usecase.CreateUserDocumentOutput{
			Type:  u.Document().Type().String(),
			Value: u.Document().Value(),
		},
		Wallet: usecase.CreateUserWalletOutput{
			Currency: u.Wallet().Money().Currency().String(),
			Amount:   u.Wallet().Money().Amount().Value(),
		},
		Roles: usecase.CreateUserRolesOutput{
			CanTransfer: u.Roles().CanTransfer,
		},
		Type:      u.TypeUser().String(),
		CreatedAt: u.CreatedAt().Format(time.RFC3339),
	}
}
