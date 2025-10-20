package controllers

import (
	"app/internal/common"
	"app/internal/services"
	"net/http"
	"github.com/labstack/echo/v4"
	"errors"
	"app/internal/models/dto"
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

	if sub.UserID != userID {
		return ctx.NoContent(http.StatusForbidden)
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"status": string(sub.Status),
	})
}

func(sc *SubmissionController) GetSubmissionDetails(ctx echo.Context) error {
	id := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	sub, err := sc.submissionService.GetSubmissionDetailsByID(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get submission details",
		})
	}

	if sub.UserID != userID {
		return ctx.NoContent(http.StatusForbidden)
	}

	return ctx.JSON(http.StatusOK, sub)
}

func(sc *SubmissionController) ListUserSubmissions(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	req, ok := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.ListProblemSubmissionsRequest)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal error: Request DTO not found in context",
		})
	}

	submissions, err := sc.submissionService.ListUserSubmissionsByProblemID(ctx.Request().Context(), userID, req.ProblemID, req.Page)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list user submissions",
		})
	}

	return ctx.JSON(http.StatusOK, dto.ListProblemSubmissionsResponse{
		Submissions: submissions,
	})
}	