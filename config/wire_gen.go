// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package config

import (
	"github.com/google/wire"
)

// Injectors from wire.go:

func Wire() (*Config, error) {
	config := Provide()
	return config, nil
}

// wire.go:

var ProviderSet = wire.NewSet(
	Provide,
)
