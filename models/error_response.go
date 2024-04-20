package models

type ErrorInfo struct {
	Code    *int   `json:"code,omitempty"`
	Message string `json:"message"`
}
