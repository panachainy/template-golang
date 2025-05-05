//go:build wireinject
// +build wireinject

//go:generate wire
package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	Provide,
)

func Wire() (*Config, error) {
	wire.Build(ProviderSet)
	return &Config{}, nil
}
