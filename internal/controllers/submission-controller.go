package controllers

import (
	"app/internal/common"
	"app/internal/services"
	"net/http"
	"github.com/labstack/echo/v4"
	"errors"
)

type SubmissionController struct {
	submissionService *services.SubmissionService
}

func NewSubmissionController(submissionService *services.SubmissionService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
	}
}

func(sc *SubmissionController) GetSubmissionStatus(ctx echo.Context) error {
	id := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	sub, err := sc.submissionService.GetSubmissionStatusByID(ctx.Request().Context(), id)
	
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
            return ctx.NoContent(http.StatusNotFound)
        }

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get submission status",
		})
	}

	if(sub.UserID != userID) {
		return ctx.NoContent(http.StatusForbidden)
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"status": string(sub.Status),
	})
}