package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dungnguyen/clean-architecture/adapter/api/response"
	"github.com/dungnguyen/clean-architecture/adapter/logger"
	"github.com/dungnguyen/clean-architecture/domain/vo"
	"github.com/dungnguyen/clean-architecture/usecase"
	"github.com/google/uuid"
)

type (
	// Request data
	CreateUserRequest struct {
		FullName string
		Email    string
		Password string
		Document CreateUserDocumentRequest
		Wallet   CreateUserWalletRequest
		Type     string
	}

	// Request data
	CreateUserDocumentRequest struct {
		Type  string
		Value string
	}

	// Request data
	CreateUserWalletRequest struct {
		Currency string
		Amount   int64
	}

	// CreateUserHandler define the dependencies of the HTTP handler for the use case
	CreateUserHandler struct {
		uc     usecase.CreateUserUseCase
		log    logger.Logger
		logKey string
	}
)

// NewCreateUserHandler create new CreateUserHandler with its dependencies
func NewCreateUserHandler(uc usecase.CreateUserUseCase, l logger.Logger) CreateUserHandler {
	return CreateUserHandler{
		uc:     uc,
		log:    l,
		logKey: "create_user",
	}
}

// Handle handles http request
func (c CreateUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	c.log = c.log.WithFields(logger.Fields{
		"correlation_id": r.Context().Value("correlation_id"),
	})

	var reqData CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		c.log.WithFields(logger.Fields{
			"key":         c.logKey,
			"error":       err.Error(),
			"http_status": http.StatusBadRequest,
		}).Errorf("Failed to marshal message")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	input, errs := c.validate(reqData)
	if len(errs) > 0 {
		c.log.WithFields(logger.Fields{
			"key":         c.logKey,
			"error":       "invalid input",
			"http_status": http.StatusInternalServerError,
		}).Errorf("failed to valiate data")

		response.NewErrors(errs, http.StatusInternalServerError).Send(w)
		return
	}

	output, err := c.uc.Execute(r.Context(), input)
	if err != nil {
		c.log.WithFields(logger.Fields{
			"key":         c.logKey,
			"error":       err.Error(),
			"http_status": http.StatusInternalServerError,
		}).Errorf("error while creating a new user")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	c.log.WithFields(logger.Fields{
		"key":         c.logKey,
		"http_status": http.StatusCreated,
	}).Infof("success creating user")

	response.NewSuccess(http.StatusCreated, output).Send(w)
}

func (c CreateUserHandler) validate(i CreateUserRequest) (usecase.CreateUserInput, []error) {
	var errs []error
	id, err := vo.NewUuid(uuid.New().String())
	if err != nil {
		errs = append(errs, err)
	}
	doc, err := vo.NewDocument(vo.TypeDocument(i.Document.Type), i.Document.Value)
	if err != nil {
		errs = append(errs, err)
	}
	email, err := vo.NewEmail(i.Email)
	if err != nil {
		errs = append(errs, err)
	}
	currency, err := vo.NewCurrency(i.Wallet.Currency)
	if err != nil {
		errs = append(errs, err)
	}
	amount, err := vo.NewAmount(i.Wallet.Amount)
	if err != nil {
		errs = append(errs, err)
	}
	money := vo.NewMoney(currency, amount)
	wallet := vo.NewWallet(money)
	typeUser, err := vo.NewTypeUser(i.Type)
	if err != nil {
		errs = append(errs, err)
	}

	return usecase.CreateUserInput{
		ID:        id,
		FullName:  vo.NewFullName(i.FullName),
		Document:  doc,
		Email:     email,
		Password:  vo.NewPassword(i.Password),
		Wallet:    wallet,
		Type:      typeUser,
		CreatedAt: time.Now(),
	}, errs
}
