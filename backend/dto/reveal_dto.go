package dto

import "wedding-app/domain/model"

type RevealResponse struct {
	Question    *model.Question `json:"question"`
	GoNext      bool            `json:"goNext"`
	NextIndex   int             `json:"nextIndex"`
	TotalQCount int             `json:"totalQCount"`
}
