package storage

import (
	"github.com/GVPproj/termsheet/models"
)

func CreateClient(name string, address, email, phone *string) (string, error) {
	return CreateEntity("client", name, address, email, phone)
}

func ListClients() ([]models.Client, error) {
	entities, err := ListEntities("client")
	if err != nil {
		return nil, err
	}

	clients := make([]models.Client, len(entities))
	for i, e := range entities {
		clients[i] = models.Client{
			ID:      e.ID,
			Name:    e.Name,
			Address: e.Address,
			Email:   e.Email,
			Phone:   e.Phone,
		}
	}

	return clients, nil
}

func UpdateClient(clientID, name string, address, email, phone *string) error {
	return UpdateEntity("client", clientID, name, address, email, phone)
}

func DeleteClient(clientID string) error {
	return DeleteEntity("client", clientID)
}
