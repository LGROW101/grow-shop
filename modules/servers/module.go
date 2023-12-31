package servers

import (
	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresHandlers"
	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresRepositories"
	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresUsecases"
	"github.com/LGROW101/lgrow-shop/modules/monitor/monitorHandlers"
	"github.com/LGROW101/lgrow-shop/modules/users/usersHandlers"
	"github.com/LGROW101/lgrow-shop/modules/users/usersRepositories"
	"github.com/LGROW101/lgrow-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}
func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.IMiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup/admin", handler.SignOut)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	// Insotail admin ขึ้นมา 1 คน ใน Db (Insert ใน SQL)
	// Generate Admin key
	// ทุกครั้งที่ทำการสมัคร admin เพิ่ม ให้สิทธิ admin token มาด้วยทุกครั้ง ผ่าน Middleware

}
