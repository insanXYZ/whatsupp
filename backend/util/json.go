package util

import "encoding/json"

func MarshalIndent(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
