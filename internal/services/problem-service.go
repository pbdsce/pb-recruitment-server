package services

import (
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/s3"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ProblemService struct {
	s3 *s3.S3
}

func NewProblemService(s3Client *s3.S3) *ProblemService {
	return &ProblemService{
		s3: s3Client,
	}
}

func (ps *ProblemService) CreateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error) {

	if problem.ID == "" {
		problem.ID = uuid.NewString()
	}

	key := fmt.Sprintf("problems/%s/%s.json", problem.ContestID, problem.ID)

	data, err := json.Marshal(problem)
	if err != nil {
		return nil, err
	}

	if err := ps.s3.PutObject(ctx, key, string(data)); err != nil {
		return nil, err
	}

	return problem, nil
}

func (ps *ProblemService) UpdateProblem(ctx context.Context, p *models.Problem) (*models.Problem, error) {

	key := fmt.Sprintf("problems/%s/%s.json", p.ContestID, p.ID)
	data, _ := json.Marshal(p)

	if err := ps.s3.PutObjectOverwrite(ctx, key, string(data)); err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *ProblemService) DeleteProblem(ctx context.Context, contestID string, problemID string) error {

	key := fmt.Sprintf("problems/%s/%s.json", contestID, problemID)
	return ps.s3.DeleteObject(ctx, key)
}

func (ps *ProblemService) GetContestProblemsList(ctx context.Context, contestID string) ([]dto.ProblemOverview, error) {

	prefix := fmt.Sprintf("problems/%s/", contestID)

	keys, err := ps.s3.ListObjects(ctx, prefix)
	if err != nil {
		return nil, err
	}

	var list []dto.ProblemOverview

	for _, key := range keys {
		raw, err := ps.s3.GetObject(ctx, key)
		if err != nil {
			continue
		}

		var p models.Problem
		if err := json.Unmarshal([]byte(raw), &p); err != nil {
			continue
		}

		list = append(list, dto.ProblemOverview{
			ID:    p.ID,
			Name:  p.Name,
			Score: p.Score,
			Type:  p.Type,
		})
	}

	return list, nil
}

func (ps *ProblemService) GetContestProblem(ctx context.Context, contestID string, problemID string) (*dto.GetProblemStatementResponse, error) {

	key := fmt.Sprintf("problems/%s/%s.json", contestID, problemID)

	raw, err := ps.s3.GetObject(ctx, key)
	if err != nil {
		return nil, err
	}

	var p models.Problem
	json.Unmarshal([]byte(raw), &p)

	resp := &dto.GetProblemStatementResponse{
		ProblemID:   p.ID,
		ContestID:   p.ContestID,
		Name:        p.Name,
		Description: p.Description,
		Score:       p.Score,
		Type:        p.Type,
	}

	return resp, nil
}
