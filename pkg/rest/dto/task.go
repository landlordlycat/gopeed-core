package dto

import "github.com/monkeyWie/gopeed-core/pkg/base"

type CreateTask struct {
	*base.Request `json:"request"`
	*base.Options `json:"options"`
}
