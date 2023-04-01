package rest

import (
	"awesomeProject3/service"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		w.WriteHeader(500)
	}
}

type container struct {
	Getter string `json:"getter"`
	Value  string `json:"value"`
}

type SecretRestAPI struct {
	ssService service.SecretService

	port string
}

func NewSecretRestAPI(sr service.SecretService, portnumber string) *SecretRestAPI {
	return &SecretRestAPI{port: portnumber, ssService: sr}
}

func (sr *SecretRestAPI) Start() error {
	h := chi.NewRouter()
	h.Method("GET", "/", Handler(sr.handlerGet))
	h.Method("POST", "/", Handler(sr.handlerPost))

	httpServer := &http.Server{Addr: ":" + sr.port, Handler: h}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sign
		log.Println("got termination signal. Try to shutdown gracefully")
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Println("Error trying to shutdown server. Error: ", err)

		}
		serverStopCtx()

	}()

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Println("server stopped successfully")
	<-serverCtx.Done()
	log.Println("Done")
	return nil
}

func (sr *SecretRestAPI) handlerGet(w http.ResponseWriter, r *http.Request) error {
	key := r.URL.Query().Get("getter")

	result, err := sr.ssService.ReadSecret(key, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil
	}
	_, err = w.Write([]byte(result))
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return nil
}

func (sr *SecretRestAPI) handlerPost(w http.ResponseWriter, r *http.Request) error {
	var p container
	err := json.NewDecoder(r.Body).Decode(&p)
	//todo: in case missing key in body case assumes empty key value. must be fixed later
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	err = sr.ssService.WriteSecret(p.Getter, p.Value, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}
