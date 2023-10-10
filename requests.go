package main

type CreateSurveyRequest struct {
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	Rka       string `json:"rka"`
	RcName    string `json:"rc_name"`
	Adress    string `json:"adress"`
	Questions []struct {
		Descripton string `json:"descripton"`
		Answers    []struct {
			Name string `json:"name"`
		} `json:"answers"`
	} `json:"questions"`
	CreatedAt string `json:"created_at"`
	UserID    int    `json:"user_id"`
}



