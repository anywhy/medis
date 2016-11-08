package build

import (
	"github.com/anywhy/medis/pkg/core/instance"
	"github.com/anywhy/medis/pkg/core/state"
)

type Service interface {
	BuildApp(app state.AppDefinition) *instance.Instance
}
