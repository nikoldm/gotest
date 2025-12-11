package utils

import (
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

var XToken = "Authorization"

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie(XToken, "", -1, "/", "", false, false)
	} else {
		c.SetCookie(XToken, "", -1, "/", host, false, false)
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie(XToken, token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie(XToken, token, maxAge, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get(XToken)
	// token 不存在则再Cooke里面找
	if token == "" {
		token, _ = c.Cookie(XToken)
		claims, err := ParseToken(token)
		if err != nil {
			fmt.Println("重新写入cookie token失败,未能成功解析token,请检查请求头是否存在x-token且claims是否为规定结构")
			return token
		}
		SetToken(c, token, int(claims.ExpiresAt.Unix()-time.Now().Unix()))
	}
	return token
}
