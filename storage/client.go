package storage

import (
	"github.com/GVPproj/termsheet/models"
)

func CreateClient(name string, address, email, phone *string) (string, error) {
	return CreateEntity("client", name, address, email, phone)
}

func ListClients() ([]models.Entity, error) {
	return ListEntities("client")
}

func UpdateClient(clientID, name string, address, email, phone *string) error {
	return UpdateEntity("client", clientID, name, address, email, phone)
}

func DeleteClient(clientID string) error {
	return DeleteEntity("client", clientID)
}
