package controllers

import (
	"app/internal/common"
	"app/internal/services"
	"net/http"
	"github.com/labstack/echo/v4"
	"errors"
	"context"
	"app/internal/models"
	"app/internal/models/dto"
)

type ContestService interface {
	IsUserRegistered(ctx context.Context, contestID, userID string) (bool, error)
}

type SubmissionController struct {
	submissionService *services.SubmissionService
	contestService     ContestService
}

func NewSubmissionController(submissionService *services.SubmissionService, contestService ContestService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
		contestService:     contestService,
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

func(sc *SubmissionController) SubmitSolution(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	req, ok := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.SubmitSubmissionRequest)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal error: SubmitSubmissionRequest DTO not found in context",
		})
	}

	isRegistered, err := sc.contestService.IsUserRegistered(reqCtx, userID, req.ContestID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to check contest registration",
		})
	}
	if !isRegistered {
		return ctx.NoContent(http.StatusForbidden)
	}

	submissionType := models.MCQ
	if req.Code != "" || req.Language != "" {
		submissionType = models.Code
	}
	
	submissionID, err := sc.submissionService.CreateSubmission(reqCtx, userID, submissionType, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, dto.SubmitSubmissionResponse{
		SubmissionID: submissionID,
	})
}