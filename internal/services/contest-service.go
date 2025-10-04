package services

type ContestService struct {
}

func NewContestService() *ContestService {
	return &ContestService{}
}

func (cs *ContestService) RegisterParticipant(contestID string, userID string) error {
	// Registration logic would go here
	return nil
}

func (cs *ContestService) ListContests() []string {
	// Realistically you would fetch this data from a database
	return []string{"Contest1", "Contest2"}
}
