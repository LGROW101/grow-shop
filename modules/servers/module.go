package servers

import (
	"github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresUsecases"
	monitorHandler "github.com/LGROW101/lgrow-shop/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
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
	repository := middlewaresrepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresusecases.IMiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandler.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}
