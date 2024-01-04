package clients_repository

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var _ domain.ClientsRepository = (*ClientsRepositoryBase)(nil)

type ClientsRepositoryBase struct {
	clients []model.Client
	cfg     Config
	mu      *sync.Mutex
}

func (c *ClientsRepositoryBase) GetClients() ([]model.Client, error) {
	return c.clients, nil
}

func (c *ClientsRepositoryBase) SetClients(clients []model.Client) error {
	c.clients = clients
	err := c.writeToConfig(fromDomain(clients))
	if err != nil {
		return fmt.Errorf("failed write clients: %s", err)
	}
	return nil
}

type Config struct {
	FilePath string
}

func NewClientsRepositoryBase(cfg Config) (*ClientsRepositoryBase, error) {

	fb, err := readFileConfig(cfg.FilePath)
	if err != nil {
		return nil, fmt.Errorf("open or create file: %v", err)
	}

	return &ClientsRepositoryBase{
		clients: fb.convertToDomain(),
		cfg:     cfg,
		mu:      new(sync.Mutex),
	}, nil
}

func fromDomain(clients []model.Client) fileBody {
	var out fileBody
	for _, client := range clients {
		out.Clients = append(out.Clients, fileClient{
			Id:     client.Id.Value(),
			Name:   client.Name,
			ApiKey: client.ApiKey.Value(),
		})
	}
	return out
}

func (r fileBody) convertToDomain() []model.Client {
	var out []model.Client
	for _, client := range r.Clients {
		out = append(out, model.Client{
			Id:     model.NewID(client.Id),
			Name:   client.Name,
			ApiKey: model.NewApiKey(client.ApiKey),
		})
	}
	return out
}

func (c *ClientsRepositoryBase) writeToConfig(body fileBody) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var b []byte
	b, err := json.MarshalIndent(body, "", " ")
	if err != nil {
		return fmt.Errorf("configuration: %v", err)
	}
	file, err := os.OpenFile(c.cfg.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("saving configuration: %v", err)
	}

	return nil
}

func readFileConfig(path string) (*fileBody, error) {
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	f.Close()

	// reads
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}
	// parse
	var body fileBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, fmt.Errorf("wrong clients configuration: %v", err)
	}
	return &body, nil
}

type fileBody struct {
	Clients []fileClient `json:"clients,omitempty"`
}

type fileClient struct {
	Id     string `json:"id"`
	Name   string `json:"name,omitempty"`
	ApiKey string `json:"api_key"`
}
