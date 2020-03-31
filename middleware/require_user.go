package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/context"
	"github.com/zjbztianya/poppy/models"
	"net/http"
	"strings"
)

func User(userService models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/images/") {
			c.Next()
			return
		}

		cookie, err := c.Cookie("remember_token")
		if err != nil {
			c.Next()
			return
		}
		user, err := userService.ByRemember(cookie)
		if err != nil {
			c.Next()
			return
		}
		context.WithUser(c, user)
		c.Next()
	}
}

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := context.User(c)
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
