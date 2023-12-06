package admin

import (
	"chatmerger/internal/domain/model"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) createClientHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	_, err := r.Body.Read(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("body: %#v", string(body))
	var input model.CreateClient

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.CreateClient(input)
	if err != nil {
		log.Println(err)
		log.Printf("err = s.CreateClient(input) err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) getClientsHandler(w http.ResponseWriter, r *http.Request) {

	clients, err := s.ClientsList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(clients)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(response)
}

func (s *Server) deleteClientHandler(w http.ResponseWriter, r *http.Request) {

}
