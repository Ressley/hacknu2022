package apiserver

import (
	"io"
	"net/http"

	"github.com/Ressley/hacknu/internal/app/apiserver/controllers"
	"github.com/Ressley/hacknu/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//APIserver ...
type APIserver struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

//New ...
func New(config *Config) *APIserver {
	return &APIserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

//Start ...
func (s *APIserver) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIserver) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIserver) configureRouter() {

	s.router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	s.router.HandleFunc("/login", controllers.Login).Methods("POST")
	s.router.HandleFunc("/download", controllers.DownloadFile).Methods("GET")
	s.router.HandleFunc("/upload", controllers.UploadFile).Methods("POST")
	s.router.HandleFunc("/community", controllers.CreateCommunity).Methods("POST")
	s.router.HandleFunc("/follow", controllers.Follow).Methods("POST")
	s.router.HandleFunc("/unfollow", controllers.Unfollow).Methods("POST")
	s.router.HandleFunc("/event/create", controllers.CreateEvent).Methods("POST")
	s.router.HandleFunc("/event", controllers.DeleteEvent).Methods("DELETE")
	s.router.HandleFunc("/event", controllers.GetEvent).Methods("POST")
	s.router.HandleFunc("/mypage", controllers.GetUser).Methods("GET")

}

func (s *APIserver) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *APIserver) handleSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "henlo")
	}
}
