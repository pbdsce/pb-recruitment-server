package controllers

import (
	"app/internal/common"
	"app/internal/services"
	"net/http"
	"strconv"

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
	pageStr := ctx.QueryParam("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	contests, err := cc.contestService.ListContests(ctx.Request().Context(), page)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list contests"})
	}
	return ctx.JSON(http.StatusOK, contests)
}

func (cc *ContestController) GetContest(ctx echo.Context) error {
	contestID := ctx.Param("id")

	userID, ok := ctx.Get(common.AUTH_USER_ID).(string)
	if !ok {
		userID = ""
	}

	contest, err := cc.contestService.GetContest(ctx.Request().Context(), contestID, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": common.FetchContestFailedError.Error(),
		})
	}

	if contest == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": common.ContestNotFoundError.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, contest)
}
