package userM

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var TOKEN_MAP = make(map[string]AuthItem)

// 白名单
var whiteList = []string{"/login"}

type AuthItem struct {
	Token     string
	TimeStamp int64
}

// 判断接口鉴权的切面
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.RequestURI)
		for _, v := range whiteList {
			if v == c.Request.RequestURI {
				c.Next()
				return
			}
		}
		// 通过http header中的token解析来认证
		token := c.Request.Header.Get("studyhouse")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 10014,
				"msg":  "请求未携带token，无权限访问",
				"data": nil,
			})
			c.Abort()
			return
		}
		log.Print("get token: ", token)
		// 通过token拿到auth实例
		auth := NewAuth(token)
		// 解析token中包含的相关信息(有效载荷)
		err := auth.Validate()

		if err != nil {
			// token过期
			if err.Error() == "TokenExpired" {
				c.JSON(http.StatusOK, gin.H{
					"code": 10015,
					"msg":  "token授权已过期，请重新申请授权",
					"data": nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"code": 10000,
				"msg":  err.Error(),
				"data": nil,
			})
			c.Abort()
			return
		}
		auth.Update()
	}
}

func NewAuth(token string) AuthItem {
	v := TOKEN_MAP[token]
	log.Println("TOKEN_MAP--", TOKEN_MAP, v)
	return v
}

// 判断是否过期
func (auth AuthItem) Validate() error {
	now := time.Now()
	tokenTime := time.Unix(auth.TimeStamp, 0)
	duration := int64(now.Sub(tokenTime).Minutes())
	if duration > 30 {
		return errors.New("TokenExpired")
	}
	return nil
}

// 更新token
func (auth AuthItem) Update() {
	auth.TimeStamp = time.Now().Unix()
	TOKEN_MAP[auth.Token] = auth
}

// 创建token
func GenerateToken() AuthItem {
	var auth AuthItem
	auth.Token = uuid.NewV4().String()
	auth.TimeStamp = time.Now().Unix()
	TOKEN_MAP[auth.Token] = auth
	return auth
}
