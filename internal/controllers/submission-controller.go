package controllers

import (
	"app/internal/services"
)

type SubmissionController struct {
	submissionService *services.SubmissionService
}

func NewSubmissionController(submissionService *services.SubmissionService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
	}
}
