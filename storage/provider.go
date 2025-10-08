package storage

import (
	"github.com/GVPproj/termsheet/models"
)

func CreateProvider(name string, address, email, phone *string) (string, error) {
	return CreateEntity("provider", name, address, email, phone)
}

func ListProviders() ([]models.Provider, error) {
	entities, err := ListEntities("provider")
	if err != nil {
		return nil, err
	}

	providers := make([]models.Provider, len(entities))
	for i, e := range entities {
		providers[i] = models.Provider{
			ID:      e.ID,
			Name:    e.Name,
			Address: e.Address,
			Email:   e.Email,
			Phone:   e.Phone,
		}
	}

	return providers, nil
}

func UpdateProvider(providerID, name string, address, email, phone *string) error {
	return UpdateEntity("provider", providerID, name, address, email, phone)
}

func DeleteProvider(providerID string) error {
	return DeleteEntity("provider", providerID)
}
