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
		w.WriteHeader(503)
	}
}

type container struct {
	Getter string `json:"getter"`
	Value  string `json:"value"`
}

type SecretRestAPI struct {
	ssService  service.SecretService
	httpServer *http.Server
}

func NewSecretRestAPI(sr service.SecretService, portnumber string) *SecretRestAPI {
	a := SecretRestAPI{}

	h := chi.NewRouter()
	h.Method("GET", "/", Handler(a.handlerGet))
	h.Method("POST", "/", Handler(a.handlerPost))

	a.ssService = sr
	a.httpServer = &http.Server{Addr: "localhost:" + portnumber, Handler: h}

	return &a
}

func (sr *SecretRestAPI) Start() {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sign
		log.Println("got termination signal. Try to shutdown gracefully")
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		err := sr.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()

	}()

	err := sr.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("server stopped successfully")
	<-serverCtx.Done()
	log.Println("Done")
}

func (sr *SecretRestAPI) handlerGet(w http.ResponseWriter, r *http.Request) error {
	key := r.URL.Query().Get("getter")

	result, err := sr.ssService.ReadSecret(key, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	_, err = w.Write([]byte(result))
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

func (sr *SecretRestAPI) handlerPost(w http.ResponseWriter, r *http.Request) error {
	var p container
	err := json.NewDecoder(r.Body).Decode(&p)
	//todo: in case missing key in body case assumes empty key value. must be fixed later
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err = sr.ssService.WriteSecret(p.Getter, p.Value, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}
