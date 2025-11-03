package controllers

import (
	"app/internal/common"
	"app/internal/models"
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

// Admin Handlers

func (cc *ContestController) HandleCreateContest(ctx echo.Context) error {
	var newContest models.Contest

	if err := ctx.Bind(&newContest); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: please check your JSON format",
		})
	}

	if newContest.Name == "" || newContest.StartTime == 0 || newContest.EndTime == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "name, start_time, and end_time are required fields",
		})
	}

	createdContest, err := cc.contestService.CreateContest(ctx.Request().Context(), &newContest)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create contest",
		})
	}

	return ctx.JSON(http.StatusCreated, createdContest)
}

func (cc *ContestController) HandleUpdateContest(ctx echo.Context) error {

	contestID := ctx.Param("id")

	var contestToUpdate models.Contest
	if err := ctx.Bind(&contestToUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if contestToUpdate.Name == "" || contestToUpdate.StartTime == 0 || contestToUpdate.EndTime == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "name, start_time, and end_time are required fields",
		})
	}

	contestToUpdate.ID = contestID

	updatedContest, err := cc.contestService.UpdateContest(ctx.Request().Context(), &contestToUpdate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update contest",
		})
	}

	return ctx.JSON(http.StatusOK, updatedContest)
}

func (cc *ContestController) HandleDeleteContest(ctx echo.Context) error {

	contestID := ctx.Param("id")
	if contestID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "contest ID is required",
		})
	}

	err := cc.contestService.DeleteContest(ctx.Request().Context(), contestID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete contest",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message":   "contest deleted successfully",
		"contestID": contestID,
	})
}

func (cc *ContestController) HandleCreateProblem(ctx echo.Context) error {

	contestID := ctx.Param("contestid")
	if contestID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "contest ID is required",
		})
	}

	var newProblem models.Problem
	if err := ctx.Bind(&newProblem); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if newProblem.Name == "" || newProblem.Score <= 0 || newProblem.Type == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "name, score, and type are required fields",
		})
	}

	newProblem.ContestID = contestID

	createdProblem, err := cc.contestService.CreateProblem(ctx.Request().Context(), &newProblem)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create problem",
		})
	}

	return ctx.JSON(http.StatusCreated, createdProblem)
}

func (cc *ContestController) HandleUpdateProblem(ctx echo.Context) error {

	contestID := ctx.Param("contestid")
	problemID := ctx.Param("problemid")
	if contestID == "" || problemID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "contest ID and problem ID are required",
		})
	}

	var problemToUpdate models.Problem
	if err := ctx.Bind(&problemToUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if problemToUpdate.Name == "" || problemToUpdate.Score <= 0 || problemToUpdate.Type == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "name, score, and type are required fields",
		})
	}

	problemToUpdate.ContestID = contestID
	problemToUpdate.ID = problemID

	updatedProblem, err := cc.contestService.UpdateProblem(ctx.Request().Context(), &problemToUpdate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update problem",
		})
	}

	return ctx.JSON(http.StatusOK, updatedProblem)
}

func (cc *ContestController) HandleDeleteProblem(ctx echo.Context) error {

	contestID := ctx.Param("contestid")
	problemID := ctx.Param("problemid")
	if contestID == "" || problemID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "contest ID and problem ID are required",
		})
	}

	err := cc.contestService.DeleteProblem(ctx.Request().Context(), contestID, problemID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete problem",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message":   "problem deleted successfully",
		"contestID": contestID,
		"problemID": problemID,
	})
}

func (cc *ContestController) HandleUpdateLeaderboardUser(ctx echo.Context) error {

	contestID := ctx.Param("contestid")
	userID := ctx.Param("userid")
	if contestID == "" || userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "contest ID and user ID are required",
		})
	}

	var req dto.UpdateLeaderboardUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Hidden == nil && req.Disqualified == nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "at least one field (hidden or disqualified) must be provided",
		})
	}

	err := cc.contestService.UpdateLeaderboardUser(ctx.Request().Context(), contestID, userID, &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update leaderboard",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message":   "leaderboard user updated successfully",
		"contestID": contestID,
		"userID":    userID,
	})
}
