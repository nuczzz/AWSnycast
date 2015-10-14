package healthcheck

import (
	"errors"
	"fmt"
	"net"
)

var healthCheckTypes map[string]func(Healthcheck) (HealthChecker, error)

func registerHealthcheck(name string, f func(Healthcheck) (HealthChecker, error)) {
	if healthCheckTypes == nil {
		healthCheckTypes = make(map[string]func(Healthcheck) (HealthChecker, error))
	}
	healthCheckTypes[name] = f
}

type HealthChecker interface {
	Healthcheck() bool
}

type Healthcheck struct {
	Type        string `yaml:"type"`
	Destination string `yaml:"destination"`
	Rise        uint   `yaml:"rise"`
	Fall        uint   `yaml:"fall"`
	Every       uint   `yaml:"every"`
}

func (h Healthcheck) GetHealthChecker() (HealthChecker, error) {
	if constructor, found := healthCheckTypes[h.Type]; found {
		return constructor(h)
	}
	return nil, errors.New(fmt.Sprintf("Healthcheck type '%s' not found in the healthcheck registry", h.Type))
}

func (h *Healthcheck) Default() {
	if h.Rise == 0 {
		h.Rise = 2
	}
	if h.Fall == 0 {
		h.Fall = 3
	}
}

func (h Healthcheck) Validate(name string) error {
	if h.Destination == "" {
		return errors.New(fmt.Sprintf("Healthcheck %s has no destination set", name))
	}
	if net.ParseIP(h.Destination) == nil {
		return errors.New(fmt.Sprintf("Healthcheck %s destination '%s' does not parse as an IP address", name, h.Destination))
	}
	if h.Type != "ping" {
		return errors.New(fmt.Sprintf("Unknown healthcheck type '%s' in %s", h.Type, name))
	}
	if h.Rise == 0 {
		return errors.New(fmt.Sprintf("rise must be > 0 in %s", name))
	}
	if h.Fall == 0 {
		return errors.New(fmt.Sprintf("fall must be > 0 in %s", name))
	}
	return nil
}