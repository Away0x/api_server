package user

import "api_server/model"

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

// 用于查询列表的参数和响应
type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type ListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	UserList   []*model.UserInfo `json:"userList"`
}
