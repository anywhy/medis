package modules

import (
	"time"
	"github.com/anywhy/medis/pkg/utils/errors"
	"github.com/anywhy/medis/pkg/modules/zk"
)

type Client interface {
	Create(path string, data []byte) error
	Update(path string, data []byte) error
	Delete(path string) error

	Read(path string, must bool) ([]byte, error)
	List(path string, must bool) ([]string, error)

	Close() error
}

var ErrUnknownCoordinator = errors.New("unknown coordinator")

func NewClient(coordinator string, addrs string, timeout time.Duration) (Client, error) {
	switch coordinator {
	case "zk", "zookeeper":
		return zkclient.New(addrs, timeout)
	}
	return nil, errors.Trace(ErrUnknownCoordinator)
}
