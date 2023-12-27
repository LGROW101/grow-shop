package middlewaresusecases

import (
	middlewaresrepositories "github.com/LGROW101/lgrow-shop/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewaresRepository middlewaresrepositories.IMiddlewaresRepository
}

func MiddlewaresRepository(middlewaresRepository middlewaresrepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}
