package model

type CreateClient struct {
	Name string `json:"name"`
}

type ClientsFilter struct {
	Id     *ID
	Name   *string
	ApiKey *ApiKey
	Status ConnStatus
}

func (f ClientsFilter) ExceptStatus() ClientsFilterExceptStatus {
	return ClientsFilterExceptStatus{
		Id:     f.Id,
		Name:   f.Name,
		ApiKey: f.ApiKey,
	}
}

type ClientsFilterExceptStatus struct {
	Id     *ID
	Name   *string
	ApiKey *ApiKey
}
