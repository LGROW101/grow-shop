package myTests

import (
	"encoding/json"

	"github.com/LGROW101/lgrow-shop/config"
	"github.com/LGROW101/lgrow-shop/modules/servers"
	"github.com/LGROW101/lgrow-shop/pkg/databases"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnect(cfg.Db())

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
