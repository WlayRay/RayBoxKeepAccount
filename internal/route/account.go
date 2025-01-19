package route

import (
	"ray_box/internal/controller"

	"github.com/kataras/iris/v12"
)

func AccountRoute(app *iris.Application) {
	app.Post("/account/login", controller.AccountController.Login)
}
