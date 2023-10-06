package handler

import (
	"errors"
	"net/http"

	"github.com/dungnguyen/clean-architecture/adapter/api/response"
	"github.com/dungnguyen/clean-architecture/adapter/logger"
	"github.com/dungnguyen/clean-architecture/domain/entity"
	"github.com/dungnguyen/clean-architecture/domain/vo"
	"github.com/dungnguyen/clean-architecture/usecase"
	"github.com/gorilla/mux"
)

// FindUserByIDHandler define the dependenies of the HTTP handler for the use case
type FindUserByIDHandler struct {
	uc     usecase.FindUserByIDUseCase
	log    logger.Logger
	logKey string
}

// NewFindUserByIDHandler create new FindUserByIDHandler with its dependencies
func NewFindUserByIDHandler(uc usecase.FindUserByIDUseCase, l logger.Logger) FindUserByIDHandler {
	return FindUserByIDHandler{
		uc:     uc,
		log:    l,
		logKey: "find_user_by_id",
	}
}

// Handler handle http request
func (f FindUserByIDHandler) Handle(w http.ResponseWriter, r *http.Request) {
	f.log = f.log.WithFields(logger.Fields{
		"correlation_id": r.Context().Value("correlation_id"),
	})

	reqID := mux.Vars(r)["user_id"]
	if reqID == "" {
		err := errors.New("invalid parameter")
		f.log.WithFields(logger.Fields{
			"key":         f.logKey,
			"error":       err.Error(),
			"http_status": http.StatusBadRequest,
		}).Errorf("invalid parameter")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	ID, err := vo.NewUuid(reqID)
	if err != nil {
		err := errors.New("invalid uuid")
		f.log.WithFields(logger.Fields{
			"key":         f.logKey,
			"error":       err.Error(),
			"http_status": http.StatusBadRequest,
		}).Errorf("invalid uuid")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	output, err := f.uc.Execute(r.Context(), usecase.FindUserByIDInput{ID: ID})
	if err != nil {
		switch err {
		case entity.ErrNotFoundUser:
			f.log.WithFields(logger.Fields{
				"key":         f.logKey,
				"error":       err.Error(),
				"http_status": http.StatusNotFound,
			}).Errorf("error fetching user by ID")

			response.NewError(err, http.StatusNotFound).Send(w)
		default:
			f.log.WithFields(logger.Fields{
				"key":         f.logKey,
				"error":       err.Error(),
				"http_status": http.StatusInternalServerError,
			}).Errorf("error fetching user by id")

			response.NewError(err, http.StatusInternalServerError).Send(w)
		}

		return
	}

	f.log.WithFields(logger.Fields{
		"key":         f.logKey,
		"http_status": http.StatusOK,
	}).Infof("success when returning user by id")

	response.NewSuccess(http.StatusOK, output).Send(w)
}
