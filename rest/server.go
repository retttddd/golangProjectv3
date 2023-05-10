package rest

import (
	"awesomeProject3/service"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"moul.io/chizap"
	"net/http"
	"sync"
	"time"
)

type responseBodyErr struct {
	Error string `json:"error"`
}
type responseBodyValue struct {
	Value string `json:"value"`
}
type container struct {
	Getter string `json:"getter"`
	Value  string `json:"value"`
}

type SecretRestAPI struct {
	ssService service.SecretService
	port      string
}

func NewSecretRestAPI(sr service.SecretService, portnumber string) *SecretRestAPI {
	return &SecretRestAPI{port: portnumber, ssService: sr}
}

func (sr *SecretRestAPI) Start(ctx context.Context) error {
	h := chi.NewRouter()
	logger := zap.NewExample()
	h.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	h.MethodFunc("GET", "/health", sr.handlerHealth)
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

func jsonWriter(err error, value string, errorFunc func(error)) string {
	if err != nil {
		errorMsgJson, err := json.Marshal(responseBodyErr{Error: err.Error()})
		if err != nil {
			errorFunc(err)
		}
		return string(errorMsgJson)
	} else {
		valueJson, err := json.Marshal(responseBodyValue{Value: value})
		if err != nil {
			errorFunc(err)
		}
		return string(valueJson)
	}
}

func (sr *SecretRestAPI) handlerGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("getter")
	result, err := sr.ssService.ReadSecret(key, r.Header.Get("X-Cipher"))
	resultMsgJson := jsonWriter(err, *result.Value, func(err2 error) {
		http.Error(w, "error:"+err2.Error(), http.StatusInternalServerError)
	})
	if err != nil {
		errorMsgJson := jsonWriter(err, "", func(err2 error) {
			http.Error(w, "error:"+err2.Error(), http.StatusInternalServerError)
		})
		http.Error(w, errorMsgJson, http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(resultMsgJson))
	if nil != err {
		errorMsgJson := jsonWriter(err, "", func(err2 error) {
			http.Error(w, "error:"+err2.Error(), http.StatusInternalServerError)
		})
		http.Error(w, errorMsgJson, http.StatusBadRequest)
		return
	}

}

func (sr *SecretRestAPI) handlerPost(w http.ResponseWriter, r *http.Request) {
	var p container
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		errorMsgJson := jsonWriter(err, "", func(err2 error) {
			http.Error(w, "error:"+err2.Error(), http.StatusInternalServerError)
		})
		http.Error(w, errorMsgJson, http.StatusBadRequest)
		return
	}

	err = sr.ssService.WriteSecret(p.Getter, &service.SecretServiceModel{Value: &p.Value}, r.Header.Get("X-Cipher"))
	if err != nil {
		errorMsgJson := jsonWriter(err, "", func(err2 error) {
			http.Error(w, "error:"+err2.Error(), http.StatusInternalServerError)
		})
		http.Error(w, errorMsgJson, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
func (sr *SecretRestAPI) handlerHealth(w http.ResponseWriter, _ *http.Request) {
	resultMsgJson := "OK"
	_, err := w.Write([]byte(resultMsgJson))
	if nil != err {
		errorMsgJson := jsonWriter(err, "", func(err2 error) {
			http.Error(w, "error:"+err2.Error(), http.StatusInternalServerError)
		})
		http.Error(w, errorMsgJson, http.StatusBadRequest)
		return
	}
}
