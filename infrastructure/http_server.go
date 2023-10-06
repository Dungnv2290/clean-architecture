package infrastructure

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dungnguyen/clean-architecture/adapter/api/handler"
	adapterhttp "github.com/dungnguyen/clean-architecture/adapter/http"
	adapterlogger "github.com/dungnguyen/clean-architecture/adapter/logger"
	"github.com/dungnguyen/clean-architecture/adapter/presenter"
	adapterqueue "github.com/dungnguyen/clean-architecture/adapter/queue"
	"github.com/dungnguyen/clean-architecture/adapter/repository"
	"github.com/dungnguyen/clean-architecture/infrastructure/database"
	infrahttp "github.com/dungnguyen/clean-architecture/infrastructure/http"
	"github.com/dungnguyen/clean-architecture/infrastructure/logger"
	"github.com/dungnguyen/clean-architecture/infrastructure/queue"
	"github.com/dungnguyen/clean-architecture/infrastructure/router"
	"github.com/dungnguyen/clean-architecture/usecase"
)

// HTTPServer define an application structure
type HTTPServer struct {
	database *database.MongoHandler
	logger   adapterlogger.Logger
	router   router.Router
	queue    *queue.RabbitMQHandler
}

// NewHTTPServer create new HTTPServer with its dependencies
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		database: database.NewMongoHandler(),
		logger:   logger.NewLogrus(),
		router:   router.NewMux(),
		queue:    queue.NewRabbitMQHandler(),
	}
}

// Start run the application
func (a HTTPServer) Start() {
	a.router.GET("health", healthCheck)

	a.router.POST("/users", a.createUserHandler())
	a.router.GET("/users/{user_id}", a.findUserByIDHandler())

	a.router.POST("/transfers", a.createTransferHandler())

	a.logger.WithFields(adapterlogger.Fields{"port": os.Getenv("APP_PORT")}).Infof("Starting HTTP Server")
	a.router.SERVE(os.Getenv("APP_PORT"))
}

func (a HTTPServer) createTransferHandler() http.HandlerFunc {
	authorizer := adapterhttp.NewAuthorizer(
		infrahttp.NewClient(
			infrahttp.NewRequest(
				infrahttp.WithRetry(infrahttp.NewRetry(3, []int{http.StatusInternalServerError}, 400*time.Millisecond)),
				infrahttp.WithTimeout(5*time.Second),
			),
		),
		a.logger,
	)

	notifier := adapterhttp.NewNotifier(
		infrahttp.NewClient(
			infrahttp.NewRequest(
				infrahttp.WithRetry(infrahttp.NewRetry(3, []int{http.StatusInternalServerError}, 400*time.Millisecond)),
				infrahttp.WithTimeout(5*time.Second),
			),
		),
		adapterqueue.NewProducer(a.queue.Channel(), a.queue.Queue().Name, a.logger),
		a.logger,
	)

	uc := usecase.NewCreateTransferInteractor(
		repository.NewCreateTransferRepository(a.database),
		repository.NewUpdateUserWalletRepository(a.database),
		repository.NewFindUserByIDRepository(a.database),
		presenter.NewCreateTransferPresenter(),
		authorizer,
		notifier,
	)

	return handler.NewCreateTransferHandler(uc, a.logger).Handle
}

func (a HTTPServer) createUserHandler() http.HandlerFunc {
	uc := usecase.NewCreateUserInteractor(
		repository.NewCreateUserRepository(a.database),
		presenter.NewCreateUserPresenter())

	return handler.NewCreateUserHandler(uc, a.logger).Handle
}

func (a HTTPServer) findUserByIDHandler() http.HandlerFunc {
	uc := usecase.NewFindUserByIDInteractor(
		repository.NewFindUserByIDRepository(a.database),
		presenter.NewFindUserByIDPresenter())

	return handler.NewFindUserByIDHandler(uc, a.logger).Handle
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: http.StatusText(http.StatusOK)})
}
