package storage

import (
	"github.com/anywhy/medis/pkg/modules"
)

func GetFrameworkId(client modules.Client) string {
	if data, err := client.Read(modules.FrameworkIdPath(), true); err != nil {
		return nil
	} else {
		return string(data)
	}
}

func SetFrameworkId(client modules.Client, fwId string) error {
	return client.Create(modules.FrameworkIdPath(), fwId)
}