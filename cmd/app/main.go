package main

import (
	"github.com/bubalync/uni-auth/internal/app"
	"github.com/bubalync/uni-auth/internal/config"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
