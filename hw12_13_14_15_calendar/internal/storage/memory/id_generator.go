package memorystorage

import "github.com/hashicorp/go-uuid"

type IDProvider interface {
	GenerateID() (string, error)
}

type UUIDProvider struct{}

func (g UUIDProvider) GenerateID() (string, error) {
	return uuid.GenerateUUID()
}
