package controllers

import (
	"app/internal/common"
	"app/internal/models/dto"
	"app/internal/services"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type SubmissionController struct {
	submissionService *services.SubmissionService
}

func NewSubmissionController(submissionService *services.SubmissionService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
	}
}

func (sc *SubmissionController) SubmitSubmission(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	var req dto.SubmitSubmissionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	response, err := sc.submissionService.SubmitSubmission(ctx.Request().Context(), userID, &req)

	if err != nil {
		ctx.Logger().Error("Error submitting submission:", err)

		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "submission must contain") {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to process submission"})
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (sc *SubmissionController) GetSubmission(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)
	submissionID := ctx.Param("id")

	submission, err := sc.submissionService.GetSubmission(ctx.Request().Context(), userID, submissionID)

	if err != nil {
		if errors.Is(err, errors.New("submission not found")) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		if strings.Contains(err.Error(), "unauthorized access") {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "submission not found"})
		}
		ctx.Logger().Error("Error retrieving submission:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve submission"})
	}
	return ctx.JSON(http.StatusOK, submission)
}

func (sc *SubmissionController) ListProblemSubmissions(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	contestID := ctx.QueryParam("contest_id")
	problemID := ctx.QueryParam("problem_id")
	limitStr := ctx.QueryParam("limit")

	if contestID == "" || problemID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "contest_id and problem_id are required query parameters"})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	response, err := sc.submissionService.ListProblemSubmissions(
		ctx.Request().Context(),
		userID,
		contestID,
		problemID,
		limit,
	)

	if err != nil {
		ctx.Logger().Error("Error listing submissions:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list submissions"})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (sc *SubmissionController) GetSubmissionStatus(ctx echo.Context) error {
	submissionID := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	status, err := sc.submissionService.GetSubmissionStatus(ctx.Request().Context(), userID, submissionID)
	if err != nil {
		if errors.Is(err, errors.New("submission not found")) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "unauthorized access") {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "submission not found"})
		}
		ctx.Logger().Error("Error retrieving submission status:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve submission status"})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"id": submissionID, "status": string(status)})
}

func (sc *SubmissionController) GetSubmissionDetails(ctx echo.Context) error {
	submissionID := ctx.Param("id")
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	details, err := sc.submissionService.GetSubmissionDetails(ctx.Request().Context(), userID, submissionID)
	if err != nil {
		if errors.Is(err, errors.New("submission not found")) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "unauthorized access") {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "submission not found"})
		}
		ctx.Logger().Error("Error retrieving submission details:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve submission details"})
	}

	return ctx.JSON(http.StatusOK, details)
}
