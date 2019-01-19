package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"

	"api_server/handler"
	"api_server/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/willf/pad"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs the each request.
/*
该中间件由于会清空 body 然后再重赋值，比较吃性能，所以可以在debug时用，线上环境就不要用了
*/
/*
1. return后不会走当前middleware后面的代码，但会走除当前middleware外的其它代码middleware和路由函数，
   至于为什么还会走下面的代码，这个是gin框架自身的特性，可以跟下代码

2. return唯一影响的就是当前middleware后面的代码无法执行了，其它无影响

3. return代表返回，c.Next代表继续执行下个中间件 -> 路由函数，之后再执行Next之后的代码
*/
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		// 该中间件只记录业务请求，比如 /v1/user 和 /login 路径
		reg := regexp.MustCompile("(/v1/user|/login)")
		if !reg.MatchString(path) {
			return
		}

		// Skip for the health check requests.
		if path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" || path == "/sd/disk" {
			return
		}

		// Read the Body content
		/*
			该中间件需要截获 HTTP 的请求信息，然后打印请求信息，
			因为 HTTP 的请求 Body，在读取过后会被置空，所以这里读取完后会重新赋值
		*/
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// The basic informations.
		method := c.Request.Method
		ip := c.ClientIP()

		//log.Debugf("New request come in, path: %s, Method: %s, body `%s`", path, method, string(bodyBytes))
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// Continue.
		c.Next() // 执行下一个中间件，然后执行路由函数，之后再执行下面的代码

		/*
			截获 HTTP 的 Response 更麻烦些，原理是重定向 HTTP 的 Response 到指定的 IO 流
		*/
		// Calculates the latency.
		end := time.Now().UTC()
		latency := end.Sub(start)

		code, message := -1, ""

		// get code and message
		var response handler.Response
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			log.Errorf(err, "response body can not unmarshal to model.Response struct, body: `%s`", blw.body.Bytes())
			code = errno.InternalServerError.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		// 截获 HTTP 的 Request 和 Response 后，就可以获取需要的信息，最终程序通过 log.Infof() 记录 HTTP 的请求信息
		log.Infof("%-13s | %-12s | %s %s | {code: %d, message: %s}", latency, ip, pad.Right(method, 5, ""), path, code, message)
	}
}
