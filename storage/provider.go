package storage

import (
	"github.com/GVPproj/termsheet/models"
)

func CreateProvider(name string, address, email, phone *string) (string, error) {
	return CreateEntity("provider", name, address, email, phone)
}

func ListProviders() ([]models.Entity, error) {
	return ListEntities("provider")
}

func UpdateProvider(providerID, name string, address, email, phone *string) error {
	return UpdateEntity("provider", providerID, name, address, email, phone)
}

func DeleteProvider(providerID string) error {
	return DeleteEntity("provider", providerID)
}
