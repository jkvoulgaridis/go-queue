package network

import (
	"encoding/json"
	"fmt"
)

type Box[T any] struct {
	content T
}

func ToCanonicalJson(payload any) string {
	fmt.Println(payload)
	bytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(bytes) + "\n"
}

func FromCanonicalJson[T any](jsonStr string) (T, error) {
	var target T
	err := json.Unmarshal([]byte(jsonStr), &target)
	if err != nil {
		return target, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return target, nil
}
