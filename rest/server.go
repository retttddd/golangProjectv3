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

type container struct {
	Getter string `json:"getter"`
	Value  string `json:"value"`
}

type SecretRestAPI struct {
	ssService        service.SecretService
	shutdown         chan struct{}
	shutdownComplete chan struct{}
	port             string
}

func NewSecretRestAPI(sr service.SecretService, portnumber string) *SecretRestAPI {
	return &SecretRestAPI{port: portnumber, ssService: sr, shutdown: make(chan struct{}, 1), shutdownComplete: make(chan struct{}, 1)}
}

func (sr *SecretRestAPI) Stop() {
	log.Println("started stop in the rest")
	sr.shutdown <- struct{}{}
	<-sr.shutdownComplete
}

func (sr *SecretRestAPI) Start(ctx context.Context) error {
	h := chi.NewRouter()
	h.MethodFunc("GET", "/", sr.handlerGet)
	h.MethodFunc("POST", "/", sr.handlerPost)

	httpServer := &http.Server{Addr: ":" + sr.port, Handler: h}
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	var err error
	
	go func() {
	     defer wg.Done()
	     srvErr := httpServer.ListenAndServe()
	     if srvErr != nil && srvErr != http.ErrServerClosed {
		    err = srvErr
	     }
	}()

	go func() {
	       defer wg.Done()
		<-ctx.Done()
		log.Println("Try to shutdown gracefully")
		shutdownCtx, shutDownStopCtx := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutDownStopCtx()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Println("Error trying to shutdown server. Error: ", err)
		}
	}()

        wg.Wait()
        log.Println("HTTP Server stopped successfully")
        if err != nil {
            return err
        }
        
	return nil
}

func (sr *SecretRestAPI) handlerGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("getter")
	result, err := sr.ssService.ReadSecret(key, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	time.Sleep(8 * time.Second)
	_, err = w.Write([]byte(result))
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (sr *SecretRestAPI) handlerPost(w http.ResponseWriter, r *http.Request) {
	var p container
	err := json.NewDecoder(r.Body).Decode(&p)
	//todo: in case missing key in body case assumes empty key value. must be fixed later
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sr.ssService.WriteSecret(p.Getter, p.Value, r.Header.Get("X-Cipher"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
