package config_test

import (
	"fmt"
	"ray_box/infrastructure/config"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fmt.Println(config.GetConfig("SECRET_KEY"))
}
