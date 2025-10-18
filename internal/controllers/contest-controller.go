package controllers

import (
	"app/internal/common"
	"app/internal/models/dto"
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

func (cc *ContestController) ModifyRegistration(ctx echo.Context) error {
	contestID := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.ModifyRegistrationRequest)

	if err := cc.contestService.ModifyRegistration(ctx.Request().Context(), contestID, userID, reqBody.Action); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to modify registration",
		})
	}

	return ctx.NoContent(http.StatusOK)
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

func (cc *ContestController) GetContestProblemsList(ctx echo.Context) error {
	contestID := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	problems, err := cc.contestService.GetContestProblemsList(ctx.Request().Context(), contestID, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get contest problems",
		})
	}

	return ctx.JSON(http.StatusOK, problems)
}

func (cc *ContestController) GetContestProblemStatement(ctx echo.Context) error {
	contestID := ctx.Param("id")
	problemID := ctx.Param("problem_id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	problem, err := cc.contestService.GetContestProblemStatement(ctx.Request().Context(), contestID, problemID, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get problem statement",
		})
	}

	if problem == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "problem not found",
		})
	}

	return ctx.JSON(http.StatusOK, problem)
}
