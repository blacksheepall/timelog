package datasource

import (
	"fmt"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/adapter/datasource/maimemo"
	"github.com/blacksheepaul/timelog/internal/ports"
)

// Registry holds enabled DataSources by name.
type Registry struct {
	sources map[string]ports.DataSource
}

// NewRegistry builds a registry from config. Disabled or unknown types are skipped.
func NewRegistry(cfgs []config.DatasourceConfig) (*Registry, error) {
	sources := make(map[string]ports.DataSource)
	for _, cfg := range cfgs {
		if !cfg.Enabled {
			continue
		}
		if cfg.Name == "" {
			return nil, fmt.Errorf("datasource config is missing name")
		}

		var source ports.DataSource
		var err error

		switch cfg.Type {
		case "maimemo":
			source, err = maimemo.NewDataSource(cfg.Name, cfg)
		default:
			return nil, fmt.Errorf("unknown datasource type %q for %q", cfg.Type, cfg.Name)
		}

		if err != nil {
			return nil, fmt.Errorf("create datasource %q: %w", cfg.Name, err)
		}
		sources[cfg.Name] = source
	}
	return &Registry{sources: sources}, nil
}

// Get returns the named DataSource or an error if it is not registered.
func (r *Registry) Get(name string) (ports.DataSource, error) {
	s, ok := r.sources[name]
	if !ok {
		return nil, fmt.Errorf("datasource %q not found", name)
	}
	return s, nil
}

// Names returns the names of all registered datasources.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.sources))
	for name := range r.sources {
		names = append(names, name)
	}
	return names
}
