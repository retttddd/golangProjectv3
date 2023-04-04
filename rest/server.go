package rest

import (
	"awesomeProject3/service"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
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
	shutdown  chan struct{}
	port      string
}

func NewSecretRestAPI(sr service.SecretService, portnumber string) *SecretRestAPI {
	return &SecretRestAPI{port: portnumber, ssService: sr, shutdown: make(chan struct{}, 1)}
}

func (sr *SecretRestAPI) Stop() {
	log.Println("started stop in the rest")
	sr.shutdown <- struct{}{}
}

func (sr *SecretRestAPI) Start() error {
	h := chi.NewRouter()
	h.MethodFunc("GET", "/", sr.handlerGet)
	h.MethodFunc("POST", "/", sr.handlerPost)

	httpServer := &http.Server{Addr: ":" + sr.port, Handler: h}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	go func() {
		<-sr.shutdown
		log.Println("got termination signal. Try to shutdown gracefully")
		shutdownCtx, shutDownStopCtx := context.WithTimeout(serverCtx, 30*time.Second)

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Println("Error trying to shutdown server. Error: ", err)

		}
		serverStopCtx()
		shutDownStopCtx()

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

func (sr *SecretRestAPI) handlerGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("getter")
	result, err := sr.ssService.ReadSecret(key, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	time.Sleep(11 * time.Second)
	_, err = w.Write([]byte(result))
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func (sr *SecretRestAPI) handlerPost(w http.ResponseWriter, r *http.Request) {
	var p container
	err := json.NewDecoder(r.Body).Decode(&p)
	//todo: in case missing key in body case assumes empty key value. must be fixed later
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	err = sr.ssService.WriteSecret(p.Getter, p.Value, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	w.WriteHeader(http.StatusCreated)

}
