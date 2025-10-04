package controllers

import (
	"app/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ContestController struct {
	contestService *services.ContestService
}

func NewContestController(contestService *services.ContestService) *ContestController {
	return &ContestController{
		contestService: contestService,
	}
}

func (cc *ContestController) ListContests(ctx echo.Context) error {
	contests := cc.contestService.ListContests()
	return ctx.JSON(http.StatusOK, contests)
}
