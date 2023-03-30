package main

import (
	"ginchat/router"
	"ginchat/utils"

	"github.com/spf13/viper"
)

func main() {
	utils.InitRedis()
	utils.InitConfig()
	utils.InitMySQL()
	r := router.Router()

	r.Run(viper.GetString("port.server"))
}
