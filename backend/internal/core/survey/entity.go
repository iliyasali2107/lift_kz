package survey

import (
	"errors"
	"time"
)

var (
	ErrSurvey         = errors.New("surveyRequirements is nil")
	ErrSurveyName     = errors.New("name field is empty")
	ErrSurveyQuestion = errors.New("questions field is empty")
	ErrAlreadyExist   = errors.New("survey with given email or nickname already exist")
	ErrNotFound       = errors.New("survey not found")
)

type SurveyRequirements struct {
	ID         string     `json:"id,omitempty"`
	UserID     int        `json:"user_id,omitempty"`
	Name       string     `json:"name"`
	Rka        string     `json:"rka,omitempty"`
	RcName     string     `json:"rc_name,omitempty"`
	Adress     string     `json:"address,omitempty"`
	Questions  []Question `json:"questions"`
	CreateDate string     `json:"create_date,omitempty"`
}

type Question struct {
	// Name string `json:"name"`
	Description string `json:"description"`
} //return []question id

type SurveyResponse struct {
	ID          int       `json:"id"`
	Name        string    `'json:"name"`
	Status      bool      `'json:"status"`
	Rka         string    `'json:"rka"`
	Rc_name     string    `'json:"rc_name"`
	Adress      string    `'json:"adress"`
	Question_id []int     `'json:"question_id"`
	CreatedAt   time.Time `'json:"created_at"`
	User_id     int       `'json:"user_id"`
}
