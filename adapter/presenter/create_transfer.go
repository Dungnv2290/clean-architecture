package presenter

import (
	"time"

	"github.com/dungnguyen/clean-architecture/domain/entity"
	"github.com/dungnguyen/clean-architecture/usecase"
)

type createTransferPresenter struct{}

// NewCreateTransferPresenter create new createTransferPresenter
func NewCreateTransferPresenter() usecase.CreateTransferPresenter {
	return createTransferPresenter{}
}

// Output return the transfer creation response
func (c createTransferPresenter) Output(t entity.Transfer) usecase.CreateTransferOutput {
	return usecase.CreateTransferOutput{
		ID:        t.ID().Value(),
		PayerID:   t.Payer().Value(),
		PayeeID:   t.Payee().Value(),
		Value:     t.Value().Amount().Value(),
		CreatedAt: t.CreatedAt().Format(time.RFC3339),
	}
}
