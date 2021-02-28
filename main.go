package main

import (
	"os"

	"github.com/akshay110495/scanner/controller"
	"github.com/akshay110495/scanner/router"
	"github.com/joho/godotenv"
)

var (
	_                        = godotenv.Load()
	PORT                     = os.Getenv("PORT")
	httpRouter router.Router = router.NewMuxRouter()
)

func main() {
	httpRouter.GET("/", controller.Home)
	httpRouter.POST("/", controller.Scan)
	httpRouter.SERVE(":" + PORT)
}
