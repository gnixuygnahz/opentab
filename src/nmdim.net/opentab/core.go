//core.go
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
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// ssl证书 文件存储= master_api  owner=  _File= acl过滤=  严格模式  scan  类封装 批处理  time= 字段过滤=   ip黑名单防恶意刷新   日志  性能记录  数据迁移  备份  conut=  更新缓存

func Run() {
	r := gin.Default()
	loadRouter(r)
	fmt.Println(`应用已启动`)

	appInfo:=GetAppInfo()

	fmt.Println(`AppId:`+appInfo[`AppId`])
	fmt.Println(`AppKey:`+appInfo[`AppKey`])
	fmt.Println(`MasterKey:`+appInfo[`MasterKey`])

	if App.IsSSL {
		r.RunTLS(`:`+App.ListenPosrt, If(App.IsLinux, "/usr/bin/", "").(string)+App.CertFile, If(App.IsLinux, "/usr/bin/", "").(string)+App.KeyFile)
	} else {
		r.Run(`:` + App.ListenPosrt)
	}

}

func loadRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		//==============对象=================

		//创建对象
		v1.POST("/classes/:className", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("create", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Create()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, *res)
				}
			}
		})

		//获取对象
		v1.GET("/classes/:className/:objectId", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("get", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Get()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, *res)
				}
			}
		})

		//更新对象
		v1.PUT("/classes/:className/:objectId", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("update", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Update()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, *res)
				}
			}
		})

		//查询对象
		v1.GET("/classes/:className", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("find", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Find()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, gin.H{"results": *res})
				}
			}
		})

		//删除对象
		v1.DELETE("/classes/:className/:objectId", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("delete", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				err = q.Delete()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, gin.H{"success": true})
				}
			}
		})

		v1.POST("/newClass", func(c *gin.Context) {
			//cont:=NewContext(c)
			CreateTable(c.PostForm("className"), map[string]interface{}{})
		})

		v1.POST("/newfield", func(c *gin.Context) {
			//cont:=NewContext(c)
			CreateField(c.PostForm("className"), c.PostForm("fieldName"), c.PostForm("type"), c.DefaultPostForm("default", ``), false, false, false, false, c.DefaultPostForm("notes", ``), c.PostForm("relationTo"), false)
		})

		//==============用户=================

		//用户注册 用户连接=
		v1.POST("/users", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.UserRegister(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//用户登录=
		v1.GET("/login", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.UserLogin(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//重置用户 sessionToken
		v1.PUT("/users/:objectId/refreshSessionToken	", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.ResetSessionToken(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//更新密码，要求输入旧密码
		v1.PUT("/users/:objectId/updatePassword", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.ResetPassword(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//请求密码重设============
		v1.POST("/requestPasswordReset", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.RequestPasswordReset(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//邮箱-密码重设
		v1.GET("/passwordReset/:email/:yzm", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.EmailPasswordReset(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//请求验证用户邮箱
		v1.POST("/requestEmailVerify", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.RequestEmailVerify(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//验证用户邮箱
		v1.GET("/emailVerify/:email/:yzm", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.EmailVerify(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//========================文件==========================

		//获得上传凭证
		v1.GET(`/files/token`, func(c *gin.Context) {
			cont := NewContext(c)
			q := ibfile.GetUploadQuery(cont)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				c.JSON(http.StatusOK, ibfile.UploadFile(cont))
			}
		})

		//七牛回调接口
		v1.POST(`/files/callback`, func(c *gin.Context) {
			cont := NewContext(c)
			res, err := ibfile.UploadCallback(cont)
			if !cont.ReturnError(err) {
				c.JSON(http.StatusOK, res)
			}

		})

		//删除文件
		v1.DELETE(`/files/:objectId`, func(c *gin.Context) {
			cont := NewContext(c)
			q := ibfile.GetDeleteQuery(cont)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := ibfile.DeleteFile(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}

		})

		//========================master==========================

		//测试返回

		v1.GET(`/master/test`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				if cont.IsMaster {
					c.JSON(http.StatusOK, gin.H{"success": true, "msg": "master connect"})
				} else {
					c.JSON(http.StatusOK, gin.H{"success": false})
				}
			}
		})

		//创建Class
		v1.POST(`/master/class`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.CreateClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//查询所有Class
		v1.GET(`/master/allClass`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetAllClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//设置Class权限
		v1.PUT(`/master/class/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.SetClassAcl(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//删除Class
		v1.DELETE(`/master/class/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.DeleteClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//删除Class所有数据
		v1.DELETE(`/master/allData/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.DeleteClassAllData(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取Class字段信息
		v1.GET(`/master/field/:className/:fieldName`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取Class所有字段信息
		v1.GET(`/master/allFieldInClass/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetFieldInClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取所有字段信息
		v1.GET(`/master/allField`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetAllField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//创建字段
		v1.POST(`/master/field/:className/`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.CreateField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//删除字段
		v1.DELETE(`/master/field/:className/:fieldName`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.DeleteField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//修改字段
		v1.PUT(`/master/field/:className/:fieldName`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.UpdateField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取应用属性
		v1.GET(`/master/appInfo`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetAppInfo(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//修改应用属性
		v1.PUT(`/master/appInfo/:key`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.SetAppInfo(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//=================系统监控===================

		//获取系统资源使用情况
		v1.GET(`/monitor/systemInfo`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmonitor.GetSystemInfo(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

	}
}
