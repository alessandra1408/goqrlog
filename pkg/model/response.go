package model

type ErrorResponse struct {
	Message string `json:"message"`
}

type ValidationResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}
