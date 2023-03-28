package rest

import (
	"awesomeProject3/service"
	"encoding/json"
	"log"
	"net/http"
)

type container struct {
	Getter string `json:"getter"`
	Value  string `json:"value"`
}

type httpServer struct {
	ssService service.SecretService
}

func NewHttpServer(sr service.SecretService) httpServer {
	return httpServer{
		ssService: sr,
	}
}

func (sr *httpServer) Start() {
	http.HandleFunc("/", sr.homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func (sr *httpServer) homePage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key := r.URL.Query().Get("getter")

		result, err := sr.ssService.ReadSecret(key, r.Header.Get("X-Cipher"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_, err = w.Write([]byte(result))
		if nil != err {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "POST":
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

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
