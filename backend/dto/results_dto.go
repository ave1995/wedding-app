package dto

import "wedding-app/domain/model"

type UserResult struct {
	Result *model.Result `json:"result"`
	User   *model.User   `json:"user"`
}
