package config

import (
	"errors"
	"github.com/spf13/viper"
)

func SaveWith(kvs ...any) error {
	if len(kvs)%2 != 0 {
		return errors.New("kvs MUST be even! Tey must be in the format 'k1, v1, k2, v2, etc")
	}

	for i := 0; i < len(kvs); i += 2 {
		viper.Set(kvs[i].(string), kvs[i+1])
	}

	return viper.WriteConfig()
}
