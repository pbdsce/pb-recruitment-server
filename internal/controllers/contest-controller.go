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

func (cc *ContestController) GetContestProblemsList(ctx echo.Context) error {
	contestID := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	err := cc.contestService.GetProblemVisibility(ctx.Request().Context(), contestID, userID)
	if err != nil {
		if err == common.ContestNotFoundError {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": common.ContestNotFoundError.Error(),
			})
		} else if err == common.UserNotRegisteredError ||
			err == common.ContestNotRunningError {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": common.FetchContestFailedError.Error(),
		})
	}

	problems, err := cc.contestService.GetContestProblemsList(ctx.Request().Context(), contestID)
	if err != nil {
		if err == common.ContestNotFoundError {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": common.ContestNotFoundError.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get contest problems",
		})
	}

	return ctx.JSON(http.StatusOK, problems)
}

func (cc *ContestController) GetContestProblem(ctx echo.Context) error {
	contestID := ctx.Param("id")
	problemID := ctx.Param("problem_id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	err := cc.contestService.GetProblemVisibility(ctx.Request().Context(), contestID, userID)
	if err != nil {
		if err == common.ContestNotFoundError {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": common.ContestNotFoundError.Error(),
			})
		} else if err == common.UserNotRegisteredError ||
			err == common.ContestNotRunningError {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": common.FetchContestFailedError.Error(),
		})
	}

	problem, err := cc.contestService.GetContestProblem(ctx.Request().Context(), contestID, problemID)
	if err != nil {
		if err == common.ContestNotFoundError {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": common.ContestNotFoundError.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get problem statement",
		})
	}

	return ctx.JSON(http.StatusOK, problem)
}
