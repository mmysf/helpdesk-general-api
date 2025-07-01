package domain

import "github.com/Yureka-Teknologi-Cipta/yureka/response"

type ResponseList struct {
	response.List
	TotalPage int64 `json:"totalPage"`
	UnreadCount int64 `json:"unreadCount,omitempty"`
}
