package opentab

import (
	"gopkg.in/gin-gonic/gin.v1"

	"strconv"

	"net/http"
	"time"
	"fmt"
)

type Middleware struct {

}

var middleware Middleware

func (m *Middleware)IpLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		va, found := App.Cache.Get(`ip:`+c.ClientIP())
		if found {
			fmt.Println("found"+` ip:`+c.ClientIP())
			appinfo:=GetAppInfo()
			*va.(*int)++
			num,err:=strconv.Atoi(appinfo[`IpLimit`])
			if err!=nil{
				panic(err)
			}
			if *va.(*int)>num{
				c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "error":"禁止访问"})
				c.Abort()
			}
		}else {
			tmp:=1
			fmt.Println("create "+` ip:`+c.ClientIP())
			App.Cache.Set(`ip:`+c.ClientIP(),&tmp,1*time.Minute)
		}
	}
}