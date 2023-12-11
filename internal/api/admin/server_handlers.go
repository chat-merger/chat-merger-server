package admin

import (
	"chatmerger/internal/domain/model"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func (s *Server) createClientHandler(w http.ResponseWriter, r *http.Request) {
	// парсирнг тела как json структуры
	var input model.CreateClient
	var name = r.PostFormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		input = model.CreateClient{Name: name}
	}
	//var err = json.NewDecoder(r.Body).Decode(&input)
	//if err != nil {
	//	log.Printf("failed parse request: %s", err)
	//	w.WriteHeader(http.StatusBadRequest)
	//	var name = r.PostFormValue("name")
	//	if name == "" {
	//		return
	//	} else {
	//		input = model.CreateClient{Name: name}
	//	}
	//}
	log.Printf("body: %#v", input)
	// создание клиента через юзкейс
	err := s.CreateClient(input)
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
	//log.Printf("%#v", r.)
	var idStr = r.FormValue("id")
	if idStr == "" {
		log.Printf("failed parse request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := s.DeleteClients([]model.ID{model.NewID(idStr)})
	if err != nil {
		log.Printf("failed delete clients: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/index.html"))

	var clients, _ = s.ClientsList()
	var resp = map[string][]model.Client{
		"Clients": clients,
	}
	tmpl.Execute(w, resp)
}
