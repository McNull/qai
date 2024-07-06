package main

import (
	"fmt"

	"github.com/mcnull/qai/api"
)

// getResponseRenderer streams markdown content and applies color based on markdown tags.
func getResponseRenderer() (api.GenerateResponseFunc, error) {
	return func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}, nil
}
