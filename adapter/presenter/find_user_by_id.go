package presenter

import (
	"time"

	"github.com/dungnguyen/clean-architecture/domain/entity"
	"github.com/dungnguyen/clean-architecture/usecase"
)

type findUserByIDPresenter struct{}

// NewFindUserByIDPresenter create findUserByIDPresenter
func NewFindUserByIDPresenter() usecase.FindUserByIDPresenter {
	return findUserByIDPresenter{}
}

// Output return the user fetch response by ID
func (f findUserByIDPresenter) Output(u entity.User) usecase.FindUserByIDOutput {
	return usecase.FindUserByIDOutput{
		ID:       u.ID().Value(),
		FullName: u.FullName().Value(),
		Email:    u.Email().Value(),
		Document: usecase.FindUserByIDDocumentOutput{
			Type:  u.Document().Type().String(),
			Value: u.Document().Value(),
		},
		Wallet: usecase.FindUserByIDWalletOutput{
			Currency: u.Wallet().Money().Currency().String(),
			Amount:   u.Wallet().Money().Amount().Value(),
		},
		Roles: usecase.FindUserByIDRolesOutput{
			CanTransfer: u.Roles().CanTransfer,
		},
		Type:      u.TypeUser().String(),
		CreatedAt: u.CreatedAt().Format(time.RFC3339),
	}
}
