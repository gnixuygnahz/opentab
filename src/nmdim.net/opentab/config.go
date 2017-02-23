//config.go
//
//Copyright 2017-present Zhang Yuxing. All Rights Reserved.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package opentab

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/goconfig"
	_ "github.com/lib/pq"
	"github.com/pmylund/go-cache"
	"gopkg.in/gin-gonic/gin.v1"
	"os"
)

type AppContext struct {
	Db          *sql.DB
	ListenPosrt string
	Cache       *cache.Cache
	IsLinux     bool
	IsSSL       bool
	CertFile    string
	KeyFile     string
}

type Context struct {
	IsMaster    bool
	IsHTTP      bool
	IsLogin     bool
	User        User
	GinContext  *gin.Context
	IsReturnErr bool
}

func (c *Context) ReturnError(err error) bool {
	if !c.IsReturnErr && err != nil {
		str := err.Error()
		tmp := strings.Split(str, `;`)
		stauts, _ := strconv.Atoi(tmp[0])
		code, _ := strconv.Atoi(tmp[1])
		c.GinContext.JSON(stauts, gin.H{"code": code, "error": tmp[2]})
		c.IsReturnErr = true
	}
	return c.IsReturnErr
}

func (c *Context) ReturnJson(res gin.H) {
	c.GinContext.JSON(http.StatusOK, res)
}

//
func NewContext(c *gin.Context) *Context {
	context := &Context{}
	context.GinContext = c
	context.IsReturnErr = false
	return context
}

var App AppContext

func InitApp() {

	fmt.Println(`正在初始化应用`)
	fmt.Println(`加载配置文件`)
	isLinux:=false;

	isDev:=os.Getenv("GIN_MODE")

	var err error
	var c  *goconfig.ConfigFile

	if isDev==`release` {
		c, err = goconfig.LoadConfigFile("config.ini")
		if err!=nil{
			c, err = goconfig.LoadConfigFile("/usr/bin/config.ini")
			CheckErr(err)
			isLinux=true
		}
	}else {
		fmt.Println(``)
		c, err = goconfig.LoadConfigFile("config-dev.ini")
		if err!=nil{
			c, err = goconfig.LoadConfigFile("/usr/bin/config-dev.ini")
			CheckErr(err)
			isLinux=true
		}
	}





	// GetValue
	host, err := c.GetValue("db", "host")
	CheckErr(err)
	port, err := c.GetValue("db", "port")
	CheckErr(err)
	dbname, err := c.GetValue("db", "dbname")
	CheckErr(err)
	dbuser, err := c.GetValue("db", "dbuser")
	CheckErr(err)
	dbpwd, err := c.GetValue("db", "dbpwd")
	CheckErr(err)

	maxIdleConns, err := c.GetValue("app", "maxIdleConns")
	CheckErr(err)
	maxOpenConns, err := c.GetValue("app", "maxOpenConns")
	CheckErr(err)
	listenPort, err := c.GetValue("app", "listenPort")
	CheckErr(err)

	db, err := sql.Open("postgres", "host="+host+" port="+port+" user="+dbuser+" password="+dbpwd+" dbname="+dbname+" sslmode=disable")
	CheckErr(err)
	maxIdleConns2, err := strconv.Atoi(maxIdleConns)
	CheckErr(err)
	maxOpenConns2, err := strconv.Atoi(maxOpenConns)
	CheckErr(err)
	db.SetMaxIdleConns(maxIdleConns2)
	db.SetMaxOpenConns(maxOpenConns2)

	isSSL, err := c.GetValue("ssl", "isDone")
	CheckErr(err)
	certFile, err := c.GetValue("ssl", "certFile")
	CheckErr(err)
	keyFile, err := c.GetValue("ssl", "keyFile")
	CheckErr(err)

	App.Db = db

	App.ListenPosrt = listenPort

	App.IsLinux=isLinux

	if isSSL == `true` {
		App.IsSSL = true
	} else {
		App.IsSSL = false
	}

	App.CertFile = certFile
	App.KeyFile = keyFile

	App.Cache = cache.New(5*time.Minute, 30*time.Second)

	fmt.Println(`初始化数据库`)
	err = InitDb(*App.Db)
	CheckErrWithStr(err, `初始化数据库完成`)

	//fmt.Println(`保存配置`)
	//goconfig.SaveConfigFile(c, If(App.IsLinux, "/usr/bin/"+If(isDev,"config-dev.ini","config.ini").(string), If(isDev,"config-dev.ini","config.ini").(string)).(string))
	fmt.Println(`初始化完成`)


}
