package user

import (
	"api_server/handler"
	"api_server/model"
	"api_server/pkg/errno"
	"api_server/utils"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// 创建用户

// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {object} user.CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user [post]
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": utils.GetReqID(c)})
	// ------------- 解析参数 ----------------
	// 请求体中的 application/json 或 urlencode 参数
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}
	/*
		// path 参数
		admin2 := c.Param("username")
		log.Infof("URL username: %s", admin2)

		// query 参数
		desc := c.Query("desc")
		log.Infof("URL key param desc: %s", desc)

		// header
		contentType := c.GetHeader("Content-Type")
		log.Infof("Header Content-Type: %s", contentType)

		log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
		if r.Username == "" {
			handler.SendResponse(c, errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not found in db: xx.xx.xx.xx")), nil)
			return
		}

		if r.Password == "" {
			handler.SendResponse(c, fmt.Errorf("password is empty"), nil)
			return
		}
	*/

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data.
	if err := u.Validate(); err != nil {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// 也可自己定义检验函数
	if err := r.checkParam(); err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}
	// Insert the user to the database.
	if err := u.Create(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	handler.SendResponse(c, nil, rsp)
}

func (r *CreateRequest) checkParam() error {
	if r.Username == "我是错误的 name" {
		return errno.New(errno.ErrValidation, nil).Add("username is error.")
	}

	return nil
}
