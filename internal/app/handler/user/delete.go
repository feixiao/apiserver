package user

import (
	"strconv"

	h "apiserver/internal/app/handler"
	"apiserver/internal/app/model/db"
	"apiserver/pkg/errno"

	"github.com/gin-gonic/gin"
)

// @Summary Delete an user by the user identifier
// @Description Delete user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /v1/user/{id} [delete]
func Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := db.DeleteUser(uint64(userId)); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
