package context

import (
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/conf"
	"github.com/zjbztianya/poppy/models"
)

func WithUser(c *gin.Context, user *models.User) {
	c.Set(conf.Conf.App.UserKey, user)
}

func User(c *gin.Context) *models.User {
	if value := c.Value(conf.Conf.App.UserKey); value != nil {
		if user, ok := value.(*models.User); ok {
			return user
		}
	}
	return nil
}
