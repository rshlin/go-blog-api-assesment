package utils

import (
	"github.com/mitchellh/mapstructure"
)

func MapToStruct(source interface{}, dest interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   dest,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(source); err != nil {
		return err
	}

	return nil
}
