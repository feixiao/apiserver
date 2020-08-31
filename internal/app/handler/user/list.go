package user

import (
	h "apiserver/internal/app/handler"
	"apiserver/internal/app/model/db"
	"apiserver/pkg/errno"
	"apiserver/pkg/util"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary List the users in the database
// @Description List users
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.ListRequest true "List users"
// @Success 200 {object} user.SwaggerListResponse "{"code":0,"message":"OK","data":{"totalCount":1,"userList":[{"id":0,"username":"admin","random":"user 'admin' get random string 'EnqntiSig'","password":"$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG","createdAt":"2018-05-28 00:25:33","updatedAt":"2018-05-28 00:25:33"}]}}"
// @Router /user [get]
func List(c *gin.Context) {
	log.Info("List function called.")
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	infos, count, err := ListUser(r.Username, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		UserList:   infos,
	})
}

func ListUser(username string, offset, limit int) ([]*db.UserInfo, uint64, error) {
	infos := make([]*db.UserInfo, 0)
	users, count, err := db.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}

	ids := []uint64{}
	for _, user := range users {
		ids = append(ids, user.Id)
	}

	wg := sync.WaitGroup{}
	userList := db.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*db.UserInfo, len(users)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// Improve query efficiency in parallel
	for _, u := range users {
		wg.Add(1)
		go func(u *db.UserModel) {
			defer wg.Done()

			shortId, err := util.GenShortId()
			if err != nil {
				errChan <- err
				return
			}

			userList.Lock.Lock()
			defer userList.Lock.Unlock()
			userList.IdMap[u.Id] = &db.UserInfo{
				Id:        u.Id,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("Hello %s", shortId),
				Password:  u.Password,
				CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
