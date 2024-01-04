package servers

import (
	"github.com/LGROW101/lgrow-shop/modules/appinfo/appinfoHandlers"
	"github.com/LGROW101/lgrow-shop/modules/appinfo/appinfoRepositories"
	"github.com/LGROW101/lgrow-shop/modules/appinfo/appinfoUsecases"
	"github.com/LGROW101/lgrow-shop/modules/orders/ordersHandlers"
	"github.com/LGROW101/lgrow-shop/modules/orders/ordersRepositories"
	"github.com/LGROW101/lgrow-shop/modules/orders/ordersUsecases"

	"github.com/LGROW101/lgrow-shop/modules/files/filesUsecases"

	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresHandlers"
	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresRepositories"
	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresUsecases"

	"github.com/LGROW101/lgrow-shop/modules/monitor/monitorHandlers"

	"github.com/LGROW101/lgrow-shop/modules/products/productsRepositories"

	"github.com/LGROW101/lgrow-shop/modules/users/usersHandlers"
	"github.com/LGROW101/lgrow-shop/modules/users/usersRepositories"
	"github.com/LGROW101/lgrow-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule() IFilesModule
	ProductsModule() IProductsModule
	OrdersModule()
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

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup/admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignOut)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	// Insotail admin ขึ้นมา 1 คน ใน Db (Insert ใน SQL)
	// Generate Admin key
	// ทุกครั้งที่ทำการสมัคร admin เพิ่ม ให้สิทธิ admin token มาด้วยทุกครั้ง ผ่าน Middleware

}
func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)
	router := m.r.Group("/appinfo")

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) OrdersModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)

	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrdersHandler(m.s.cfg, ordersUsecase)

	router := m.r.Group("/orders")
	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.FindOneOrder)
	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.UpdateOrder)

}
