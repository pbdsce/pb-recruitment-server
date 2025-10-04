package controllers

import (
	"app/internal/common"
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

func (cc *ContestController) RegisterParticipant(ctx echo.Context) error {
	contestID := ctx.Param("id") // /contests/:id/register
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	if err := cc.contestService.RegisterParticipant(contestID, userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to register participant",
		})
	}

	return ctx.NoContent(http.StatusOK)
}

func (cc *ContestController) ListContests(ctx echo.Context) error {
	contests := cc.contestService.ListContests()
	return ctx.JSON(http.StatusOK, contests)
}
