package provider

import (
	"fmt"
	"github.com/mcnull/qai/shared/jsonmap"
)

type IConfig interface {
	Merge(other IConfig) error
}

func InitConfig(providerConfig IConfig, factory ConfigFactory, profileSettings jsonmap.JsonMap) (IConfig, error) {
	a := factory()      // coded defaults
	b := providerConfig // from config
	c := factory()      // from profile

	err := profileSettings.ToStruct(c)

	if err != nil {
		err = fmt.Errorf("error converting profile settings to struct: %w", err)
		return nil, err
	}

	err = a.Merge(b)

	if err != nil {
		err = fmt.Errorf("error merging provider config into base: %w", err)
		return nil, err
	}

	err = a.Merge(c)

	if err != nil {
		err = fmt.Errorf("error merging profile config into base: %w", err)
		return nil, err
	}

	return a, nil
}
