package storage

import (
	"github.com/anywhy/medis/pkg/models"
)

func GetFrameworkId(client models.Client) string {
	if data, err := client.Read(models.FrameworkIdPath(), true); err != nil {
		return nil
	} else {
		return string(data)
	}
}

func SetFrameworkId(client models.Client, fwId string) error {
	return client.Create(models.FrameworkIdPath(), fwId)
}