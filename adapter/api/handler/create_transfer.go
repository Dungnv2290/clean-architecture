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
	CreateTransferRequest struct {
		PayerID string `json:"payser_id"`
		PayeeID string `json:"payee_id"`
		Value   int64  `json:"value"`
	}

	// CreateTransferHandler define the dependencies of the HTTP handler for the use case
	CreateTransferHandler struct {
		uc     usecase.CreateTransferUseCase
		log    logger.Logger
		logKey string
	}
)

// NewCreateTransferHandler create new CreateTransferHandler with its dependencies
func NewCreateTransferHandler(uc usecase.CreateTransferUseCase, l logger.Logger) CreateTransferHandler {
	return CreateTransferHandler{
		uc:     uc,
		log:    l,
		logKey: "create_transfer",
	}
}

// Handle handle http request
func (c CreateTransferHandler) Handle(w http.ResponseWriter, r *http.Request) {
	c.log = c.log.WithFields(logger.Fields{
		"correlation_id": r.Context().Value("correlation_id"),
	})

	var reqData CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		c.log.WithFields(logger.Fields{
			"key":         c.logKey,
			"error":       err.Error(),
			"http_status": http.StatusBadRequest,
		}).Errorf("failed to marshal message")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	input, errs := c.validate(reqData)
	if len(errs) > 0 {
		c.log.WithFields(logger.Fields{
			"key":         c.logKey,
			"error":       "invalid input",
			"http_status": http.StatusBadRequest,
		}).Errorf("failed to data")

		response.NewErrors(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := c.uc.Execute(r.Context(), input)
	if err != nil {
		c.log.WithFields(logger.Fields{
			"key":         c.logKey,
			"error":       err.Error(),
			"http_status": http.StatusInternalServerError,
		}).Errorf("error when creating a new transfer")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	c.log.WithFields(logger.Fields{
		"key":         c.logKey,
		"http_status": http.StatusCreated,
	}).Infof("success creating transfer")

	response.NewSuccess(http.StatusCreated, output).Send(w)
}

func (c CreateTransferHandler) validate(i CreateTransferRequest) (usecase.CreateTransferInput, []error) {
	var errs []error
	id, err := vo.NewUuid(uuid.New().String())
	if err != nil {
		errs = append(errs, err)
	}
	payerID, err := vo.NewUuid(i.PayerID)
	if err != nil {
		errs = append(errs, err)
	}
	payeeID, err := vo.NewUuid(i.PayeeID)
	if err != nil {
		errs = append(errs, err)
	}
	amount, err := vo.NewAmount(i.Value)
	if err != nil {
		errs = append(errs, err)
	}

	return usecase.CreateTransferInput{
		ID:       id,
		PayerID:  payerID,
		PayeeID:  payeeID,
		Value:    vo.NewMoneyBRL(amount),
		CreateAt: time.Now(),
	}, errs
}
