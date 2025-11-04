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
		if err == common.ContestRegistrationClosedError ||
			err == common.InvalidYearError {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		} else if err == common.ContestNotFoundError ||
			err == common.UserNotFoundError {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		} else if err == common.UserAlreadyExistsError {
			return ctx.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to modify registration",
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
