package itbit

import "fmt"

type itbitError struct {
	ErrCode int `json:"code"`
	Description string `json:"description"`
	RequestID string `json:"requestId"`
}

func (e *itbitError) Code () int {
	return e.ErrCode
}

func (e *itbitError) Error () string {
	return fmt.Sprintf("[%d] %s (%s)", e.ErrCode, e.Description, e.RequestID)
}

