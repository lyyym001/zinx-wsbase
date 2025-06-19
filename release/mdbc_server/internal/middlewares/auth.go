package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/helper"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		//fmt.Println("auth = ", auth)
		userClaim, err := helper.AnalyseToken(auth)
		//fmt.Println("auth error ,", err)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Authorization",
			})
			return
		}
		if userClaim == nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Admin",
			})
			return
		}
		c.Set("user_claims", userClaim)
		c.Next()
	}
}
