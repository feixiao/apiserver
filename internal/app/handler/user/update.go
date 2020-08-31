package user

import (
	h "apiserver/internal/app/handler"
	"apiserver/internal/app/model/db"
	"apiserver/pkg/errno"
	"apiserver/pkg/util"
	"strconv"

	// "github.com/douyu/jupiter/pkg/xlog/lager"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// @Summary Update a user info by the user identifier
// @Description Update a user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Param user body db.UserModel true "The user info"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /user/{id} [put]
func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	userId, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var u db.UserModel
	if err := c.Bind(&u); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// We update the record based on the user id.
	u.Id = uint64(userId)

	// Validate the data.
	if err := u.Validate(); err != nil {
		h.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		h.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// Save changed fields.
	if err := u.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
