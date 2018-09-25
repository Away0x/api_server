package user

import (
	service "github.com/Away0x/api_server/service/user"

	"github.com/Away0x/api_server/handler"
	"github.com/Away0x/api_server/pkg/errno"
	"github.com/gin-gonic/gin"
)

// List list the users in the database.
/*
一般在 handler 中主要做解析参数、返回数据操作，简单的逻辑也可以在 handler 中做
像新增用户、删除用户、更新用户，代码量不大，所以也可以放在 handler 中。
有些代码量很大的逻辑就不适合放在 handler 中，因为这样会导致 handler 逻辑不是很清晰，
这时候实际处理的部分通常放在 service 包中
*/

/*
// @Summary List the users in the database
// @Description List users
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.ListRequest true "List users"
// @Success 200 {object} user.SwaggerListResponse "{"code":0,"message":"OK","data":{"totalCount":1,"userList":[{"id":0,"username":"admin","random":"user 'admin' get random string 'EnqntiSig'","password":"$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG","createdAt":"2018-05-28 00:25:33","updatedAt":"2018-05-28 00:25:33"}]}}"
// @Router /user [get]
*/
func List(c *gin.Context) {
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// service.ListUser() 函数用来做具体的查询处理
	infos, count, err := service.ListUser(r.Username, r.Offset, r.Limit)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	handler.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		UserList:   infos,
	})
}
