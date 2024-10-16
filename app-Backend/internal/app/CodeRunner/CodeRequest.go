package coderunner

type CodeRequest struct {
	Lang   string `json:"lang"`
	Code   string `json:"code"`
	Action string `json:"action,omitempty"`
}
