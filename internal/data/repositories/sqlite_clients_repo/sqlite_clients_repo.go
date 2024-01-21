package sqlite_clients_repo

import (
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"database/sql"
)

var _ repository.ClientsRepository = (*ClientsRepository)(nil)

type ClientsRepository struct {
	db *sql.DB
}

func NewClientsRepository(db *sql.DB) *ClientsRepository {
	return &ClientsRepository{db: db}
}

func (c *ClientsRepository) GetClients(filter model.ClientsFilter) ([]model.Client, error) {
	stmt, err := c.db.Prepare(`
		with arg(f_id, f_name, f_api_key, f_status) as (
		    select ?,?,?,?
		)
		select id, name, api_key, status
		from client, arg
		where 	(f_id is null or id = f_id) and
				(f_name is null or name = f_name) and
				(f_api_key is null or api_key = f_api_key) and
				(f_status is null or status = f_status)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(filter.Id, filter.Name, filter.ApiKey, filter.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clients []model.Client
	for rows.Next() {
		var client model.Client
		err = rows.Scan(&client.Id, &client.Name, &client.ApiKey, &client.Status)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (c *ClientsRepository) Create(client model.Client) error {
	stmt, err := c.db.Prepare(`
		insert into client (id, name, api_key, status)
		values (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(client.Id, client.Name, client.ApiKey, client.Status)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientsRepository) Delete(id model.ID) error {
	stmt, err := c.db.Prepare(`
		delete from client where id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientsRepository) Update(id model.ID, new model.Client) error {
	stmt, err := c.db.Prepare(`
		update client 
		set id = ?,
		    name = ?,
		    api_key = ?,
		    status = ?
		where id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(new.Id, new.Name, new.ApiKey, new.Status, id)
	if err != nil {
		return err
	}
	return nil
}
