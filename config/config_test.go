package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := DefaultCommon()
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))
}
