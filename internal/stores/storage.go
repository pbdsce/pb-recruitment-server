package stores

type Storage struct {
	Contests interface {
		//todo: add contest store
	}
	Users interface {
		//todo: add user store
	}
	Submissions interface {
		//todo: add submission store
	}
	Rankings interface {
		//todo: add ranking store
	}
	Problems interface {
		//todo: add problem store
	}
}

func NewStorage() *Storage {
	return &Storage{
		Contests:    NewContestStore(),
		Users:       NewUserStore(),
		Submissions: NewSubmissionStore(),
		Rankings:    NewRankingStore(),
		Problems:    NewProblemStore(),
	}
}
