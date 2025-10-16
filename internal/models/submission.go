package models

type SubmissionType string

const (
	MCQ  SubmissionType = "MCQ"
	Code SubmissionType = "Code"
)

type SubmissionStatus string

const (
	Pending           SubmissionStatus = "Pending"
	Accepted          SubmissionStatus = "Accepted"
	WrongAnswer       SubmissionStatus = "Wrong Answer"
	TimeLimitExceed   SubmissionStatus = "Time Limit Exceeded"
	MemoryLimitExceed SubmissionStatus = "Memory Limit Exceeded"
	RuntimeError      SubmissionStatus = "Runtime Error"
	CompilationError  SubmissionStatus = "Compilation Error"
)

type TestCaseResult struct {
	TestCaseID string `json:"test_case_id"`
	Status	   string `json:"status"`   
	Runtime    int64  `json:"runtime"` 
	Memory     int64  `json:"memory"`  
}

type Submission struct {
	ID        		string           `json:"id"`
	UserID    		string           `json:"user_id"`
	ContestID 		string           `json:"contest_id"`
	ProblemID 		string           `json:"problem_id"`
	Type      		SubmissionType   `json:"type"`
	Language  		string           `json:"language,omitempty"` // For code submissions
	Code      		string           `json:"code,omitempty"`     // Base64 encoded code for code submissions
	Option    		[]int            `json:"option,omitempty"`   // Selected option(s) for MCQ submissions
	Status    		SubmissionStatus `json:"status"`             // e.g., "Pending", "Accepted", "Wrong Answer", etc.
	CreatedAt 		int64            `json:"created_at"`         // Unix timestamp
	Runtime   		int64            `json:"runtime,omitempty"` 
	Memory    		int64            `json:"memory,omitempty"`
	TestCaseResults []TestCaseResult `json:"test_case_results,omitempty"`
}