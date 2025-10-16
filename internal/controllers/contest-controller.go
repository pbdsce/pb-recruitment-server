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

func (cc *ContestController) ListContests(ctx echo.Context) error {
	pageStr := ctx.QueryParam("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
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

	userID, isAuthenticated := ctx.Get(common.AUTH_USER_ID).(string)

	contest, err := cc.contestService.GetContest(ctx.Request().Context(), contestID, userID, isAuthenticated)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get contest",
		})
	}

	if contest == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "contest not found",
		})
	}

	return ctx.JSON(http.StatusOK, contest)
}
