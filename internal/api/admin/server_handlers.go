package admin

import (
	"chatmerger/internal/domain/model"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
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
	log.Printf("body: %#v", input)
	// создание клиента через юзкейс
	err := s.CreateClient(input)
	if err != nil {
		log.Println(err)
		log.Printf("err = s.CreateClient(input) err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tmpl = template.Must(template.ParseFiles("web/clients_table.html"))

	var clients, _ = s.ClientsList()
	var resp = map[string][]model.Client{
		"Clients": clients,
	}
	tmpl.Execute(w, resp)
}

func (s *Server) getClientsHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/clients_table.html"))

	clients, err := s.ClientsList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resp = map[string][]model.Client{
		"Clients": clients,
	}
	tmpl.Execute(w, resp)
}

func (s *Server) deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
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
	file, err := os.ReadFile("web/index.html")
	if err != nil {
		log.Printf("failed read index.html file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(file)
}
