package user

import "apiserver/internal/app/model/db"

type CreateRequest struct {
	Username string `json:"username" binding:"required,gte=1,lte=16"`
	Password string `json:"password" binding:"required"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type ListResponse struct {
	TotalCount uint64         `json:"totalCount"`
	UserList   []*db.UserInfo `json:"userList"`
}

type SwaggerListResponse struct {
	TotalCount uint64        `json:"totalCount"`
	UserList   []db.UserInfo `json:"userList"`
}
