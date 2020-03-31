package context

import (
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/models"
)

const (
	userKey string = "user"
)

func WithUser(c *gin.Context, user *models.User) {
	c.Set(userKey, user)
}

func User(c *gin.Context) *models.User {
	if value := c.Value(userKey); value != nil {
		if user, ok := value.(*models.User); ok {
			return user
		}
	}
	return nil
}
