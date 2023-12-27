package main

import (
	"os"

	"github.com/LGROW101/lgrow-shop/config"
	"github.com/LGROW101/lgrow-shop/modules/servers"
	"github.com/LGROW101/lgrow-shop/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()

}
