//user.go
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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
	"net/url"
)

type User struct {
	Id       string
	UserName string
	Role     []string
}

type RolesMap map[string]*[]string

type IBUser struct {

}

var ibuser  IBUser

//获取用户Id和Role，此处可判断是否有此人
func (i *IBUser)GetUserBySession(sessionToken string) (*User, bool) {
	user := User{}
	tmp := QueryWithCache(`select * from public."_User" where "sessionToken"='`+sessionToken+`'`, 0)
	res := *(tmp.(*[]map[string]interface{}))
	if len(res) == 0 {
		return &user, false
	} else {
		user.Id = fmt.Sprint(res[0][`id`])
		user.UserName = fmt.Sprint(res[0][`username`])
		user.Role = ibuser.GetRoles(user.Id)
		return &user, true
	}
}

//获取用户角色列表
func (i *IBUser)GetRoles(id string) []string {
	arr := []string{}
	rolesMap := ibuser.GetRolesMap()
	res := EasyQuery(`select * from public."_Role" where "users"@>array[` + id + `]`)
	if len(*res) != 0 {
		for _, va := range *res {
			tmp := rolesMap[va[`name`].(string)]
			arr = (MergeStringArray(arr, *tmp))
		}
	}
	arr = RemoveDuplicatesAndEmpty(arr)
	return arr
}

//获取角色链映射
func (i *IBUser)GetRolesMap() RolesMap {

	tmp, found := App.Cache.Get(`app:RolesMap`)
	if found {
		return tmp.(RolesMap)
	} else {
		tx, err := App.Db.Begin()
		CheckErr(err)
		defer tx.Commit()

		rolesMap := RolesMap{}

		r, _ := tx.Query(`select * from public."_Role"`)

		res := GetJSON(r)

		for _, value := range *res {
			rolesMap[value[`name`].(string)] = &[]string{}
			var tmp1 **[]string
			var tmp2 *[]string
			tmp2 = rolesMap[value[`name`].(string)]
			*tmp1 = tmp2
			ibuser.GenRoleList(value[`id`].(string), tmp1)
			if !is_array(*rolesMap[value[`name`].(string)], value[`name`].(string)) {
				tmp := append(*rolesMap[value[`name`].(string)], value[`name`].(string))
				rolesMap[value[`name`].(string)] = &tmp
			}
		}

		App.Cache.Set(`app:RolesMap`, rolesMap, -1) //

		return rolesMap
	}
}

//获取角色链
func (i *IBUser)GenRoleList(roleId string, rolesList **[]string) {
	res := QueryWithCache(`select * from public."_Role" where "roles"@>array[`+roleId+`]`, 0) //此处可修改时间，但最好主动清除缓存
	if len(*(res.(*[]map[string]interface{}))) != 0 {
		for _, va := range *(res.(*[]map[string]interface{})) {
			if is_array(**rolesList, va[`name`].(string)) {
				//
			} else {
				rolesList2 := append(**rolesList, va[`name`].(string))
				*rolesList = &rolesList2
				ibuser.GenRoleList(va[`id`].(string), rolesList)
			}
		}
	}
}

//用户登陆
func (i *IBUser)UserLogin(c *Context) (gin.H, error) {

	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()

	data, err := Json(c.GinContext)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
	}
	username, ok := (*data)[`username`]
	if !ok {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少username"))
	}
	password, ok := (*data)[`password`]
	if !ok {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少password"))
	}

	appInfo := GetAppInfo()

	if appInfo[`UserLoginRequireEmailVerified`] == `true` { //未验证邮箱的用户禁止登陆
		tmp, err := tx.Query(`select "sessionToken","id","username","password","emailVerified" from public."_User" where "username"='` + username.(string) + `'`)
		res := GetJSON(tmp)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
		}
		if len(*res) == 0 {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("无此用户"))
		}
		if (((*res)[0])[`password`]).(string) != Md5Encrypt(password.(string), username.(string)) {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("密码错误"))
		}
		if (((*res)[0])[`emailVerified`]).(bool) != true {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("未验证邮箱的用户禁止登陆"))
		}
		return gin.H{"id": (((*res)[0])[`id`]).(int), "sessionToken": (((*res)[0])[`sessionToken`]).(string), "username": (((*res)[0])[`username`]).(string)}, nil
	} else {
		tmp, err := tx.Query(`select "sessionToken","id","username","password" from public."_User" where "username"='` + username.(string) + `'`)
		res := GetJSON(tmp)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
		}
		if len(*res) == 0 {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("无此用户"))
		}
		if (((*res)[0])[`password`]).(string) != Md5Encrypt(password.(string), username.(string)) {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("密码错误"))
		}
		return gin.H{"id": (((*res)[0])[`id`]).(int), "sessionToken": (((*res)[0])[`sessionToken`]).(string), "username": (((*res)[0])[`username`]).(string)}, nil
	}
}

//用户注册
func (i *IBUser)UserRegister(c *Context) (gin.H, error) {

	appInfo := GetAppInfo()
	if appInfo[`UserRegister`] == `true` || c.IsMaster {

		var err error
		tx, err := App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		data, err := Json(c.GinContext)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		username, ok := (*data)[`username`]
		if !ok {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少username"))
		}
		password, ok := (*data)[`password`]
		if !ok {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少password"))
		}

		if appInfo[`UserRegisterNeedEmail`] == `true` {
			r := GetRand()
			email, ok := (*data)[`email`]
			if !ok {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少email"))
			}
			_, err = tx.Query(`INSERT INTO public."_User" ("username","password","email","emailVerified","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + email.(string) + `',false,'` + Krand(24, KC_RAND_KIND_ALL, r) + `') `)
			if err != nil && strings.Contains(err.Error(), `session`) {
				_, err = tx.Query(`INSERT INTO public."_User" ("username","password","email","emailVerified","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + email.(string) + `',false,'` + Krand(24, KC_RAND_KIND_ALL, r) + `') `)
				if err != nil && strings.Contains(err.Error(), `session`) {
					_, err = tx.Query(`INSERT INTO public."_User" ("username","password","email","emailVerified","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + email.(string) + `',false,'` + Krand(24, KC_RAND_KIND_ALL, r) + `') `)
					if err != nil && strings.Contains(err.Error(), `session`) {
						_, err = tx.Query(`INSERT INTO public."_User" ("username","password","email","emailVerified","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + email.(string) + `',false,'` + Krand(24, KC_RAND_KIND_ALL, r) + `') `)
						if err != nil {
							return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
						}
					} // 尝试三次，
				}
			}

			if err != nil && strings.Contains(err.Error(), `username`) {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("用户名已存在"))
			}
			if err != nil && strings.Contains(err.Error(), `email`) {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("邮箱已被使用"))
			}

			tmp := Krand(24, KC_RAND_KIND_ALL, r)
			App.Cache.Set(`emailVerified:`+email.(string), tmp, 0)
			SendToMail(appInfo[`SmtpUser`], appInfo[`SmtpPassword`], appInfo[`SmtpHost`], email.(string), ``, `<a href="`+If(App.IsSSL, "https://", "http://").(string)+appInfo[`AppHost`]+`/v1/emailVerify/`+email.(string)+`/`+tmp+`">请在5分钟以内，点击验证<a/>`, `html`)
			return gin.H{"success": true, "msg": "验证邮件已发送，请在5分钟内登陆邮箱并验证"}, nil

		} else {
			r := GetRand()
			res, err := tx.Query(`INSERT INTO public."_User" ("username","password","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + Krand(24, KC_RAND_KIND_ALL, r) + `') Returning "id","sessionToken" `)
			if err != nil && strings.Contains(err.Error(), `session`) {
				res, err = tx.Query(`INSERT INTO public."_User" ("username","password","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + Krand(24, KC_RAND_KIND_ALL, r) + `') Returning "id","sessionToken" `)
				if err != nil && strings.Contains(err.Error(), `session`) {
					res, err = tx.Query(`INSERT INTO public."_User" ("username","password","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + Krand(24, KC_RAND_KIND_ALL, r) + `') Returning "id","sessionToken" `)
					if err != nil && strings.Contains(err.Error(), `session`) {
						res, err = tx.Query(`INSERT INTO public."_User" ("username","password","sessionToken") values ('` + username.(string) + `','` + Md5Encrypt(password.(string), username.(string)) + `','` + Krand(24, KC_RAND_KIND_ALL, r) + `') Returning "id","sessionToken" `)
						if err != nil {
							return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
						}
					}
				}
			}
			tmp := GetJSON(res)
			if err != nil && strings.Contains(err.Error(), `username`) {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("用户名已存在"))
			}
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
			}
			return gin.H{"success": true, "object": (*tmp)[0]}, nil
		}

	} else {
		return nil, NewError(http.StatusForbidden, 104, errors.New(`禁止访问`))
	}
}

// 邮箱验证地址
func (i *IBUser)EmailVerify(c *Context) (gin.H, error) {
	email := c.GinContext.Param("email")
	yzm := c.GinContext.Param("yzm")
	tmp, found := App.Cache.Get(`emailVerified:` + email)
	if found {
		if yzm == tmp.(string) {
			var err error
			tx, err := App.Db.Begin()
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
			defer tx.Commit()

			_, err = tx.Query(`UPDATE public."_User" SET "emailVerified"=true WHERE "email"='` + email + `'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
			}
			App.Cache.Delete(`emailVerified:` + email)
			return gin.H{"success": true}, nil
		} else {
			return gin.H{"success": false, "msg": "未知验证邮件"}, nil
		}
	} else {
		return gin.H{"success": false, "msg": "验证邮件已过期，请重新请求发送验证邮件"}, nil
	}
}

//请求发送验证邮件
func (i *IBUser)RequestEmailVerify(c *Context) (gin.H, error) {
	data, err := Json(c.GinContext)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
	}
	email, ok := (*data)[`email`]
	if !ok {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少email"))
	}

	res := EasyQuery(`select * from public."_User" where "email"='` + email.(string) + `'`)

	if len(*res) == 0 {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("无此邮箱"))
	}
	if (*res)[0][`emailVerified`].(bool) == true {
		return gin.H{"success": false, "msg": "该邮箱已验证"}, nil
	}
	_, found := App.Cache.Get(`emailVerified:` + email.(string))
	if found {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("无需重复验证"))
	} else {
		r := GetRand()
		tmp := Krand(24, KC_RAND_KIND_ALL, r)
		App.Cache.Set(`emailVerified:`+email.(string), tmp, 0)
		appInfo := GetAppInfo()
		SendToMail(appInfo[`SmtpUser`], appInfo[`SmtpPassword`], appInfo[`SmtpHost`], email.(string), ``, `<a href="`+If(App.IsSSL, "https://", "http://").(string)+appInfo[`AppHost`]+`/v1/emailVerify/`+email.(string)+`/`+tmp+`">请在5分钟以内，点击验证<a/>`, `html`)
		return gin.H{"success": true, "msg": "验证邮件已发送，请在5分钟内登陆邮箱并验证"}, nil
	}

}

//重置sessionToken
func (i *IBUser)ResetSessionToken(c *Context) (gin.H, error) {
	id := c.GinContext.Param("objectId")
	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()

	if (c.IsLogin && c.User.Id == id) || c.IsMaster {
		r := GetRand()
		tmp, err := tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
		if err != nil && strings.Contains(err.Error(), `session`) {
			tmp, err = tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
			if err != nil && strings.Contains(err.Error(), `session`) {
				tmp, err = tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
				if err != nil && strings.Contains(err.Error(), `session`) {
					tmp, err = tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
					if err != nil {
						return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
					}
				} // 尝试三次，
			}
		}
		res := GetJSON(tmp)
		if len(*res) == 0 {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
		}
		return (*res)[0], nil
	} else {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("禁止访问"))
	}

}

//修改密码
func (i *IBUser)ResetPassword(c *Context) (gin.H, error) {
	id := c.GinContext.Param("objectId")
	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()

	if (c.IsLogin && c.User.Id == id) || c.IsMaster {

		appInfo := GetAppInfo()

		data, err := Json(c.GinContext)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}

		if appInfo[`UserResetPasswordRequireOld`] == `true` { //重置密码需要旧密码

			old_password, ok := (*data)[`old_password`]
			if !ok {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少old_password"))
			}
			new_password, ok := (*data)[`new_password`]
			if !ok {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少new_password"))
			}

			tmp, err := tx.Query(`UPDATE public."_User" SET "password"='` + Md5Encrypt(new_password.(string), c.User.UserName) + `' WHERE "id"=` + id + ` and "password"='` + old_password.(string) + `' Returning "id","sessionToken"`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
			}
			res := GetJSON(tmp)
			if len(*res) == 0 {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("密码错误"))
			}

		} else {
			new_password, ok := (*data)[`new_password`]
			if !ok {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少new_password"))
			}

			_, err := tx.Query(`UPDATE public."_User" SET "password"='` + Md5Encrypt(new_password.(string), c.User.UserName) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
			}
		}
		if appInfo[`UserResetPasswordRequireRelogin`] == `true` { //重置密码后强制重新登陆
			r := GetRand()
			tmp, err := tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
			if err != nil && strings.Contains(err.Error(), `session`) {
				tmp, err = tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
				if err != nil && strings.Contains(err.Error(), `session`) {
					tmp, err = tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
					if err != nil && strings.Contains(err.Error(), `session`) {
						tmp, err = tx.Query(`UPDATE public."_User" SET "sessionToken"='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "id"=` + id + ` Returning "id","sessionToken"`)
						if err != nil {
							return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
						}
					} // 尝试三次，
				}
			}
			res := GetJSON(tmp)
			if len(*res) == 0 {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
			}

			return gin.H{"success": true, "msg": "密码修改成功，请重新登陆"}, nil
		} else {
			return gin.H{"success": true, "msg": "密码修改成功"}, nil
		}

	} else {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("禁止访问"))
	}
}

func (i *IBUser)RequestPasswordReset(c *Context) (gin.H, error)  {
	data, err := Json(c.GinContext)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
	}
	email, ok := (*data)[`email`]
	if !ok {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少email"))
	}

	res := EasyQuery(`select * from public."_User" where "email"='` + email.(string) + `'`)

	if len(*res) == 0 {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("无此邮箱"))
	}

	_, found := App.Cache.Get(`passwordReset:` + email.(string))
	if found {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("无需重复发送验证邮件"))
	} else {
		r := GetRand()
		tmp := Krand(24, KC_RAND_KIND_ALL, r)
		App.Cache.Set(`passwordReset:`+email.(string), tmp, 0)
		appInfo := GetAppInfo()
		SendToMail(appInfo[`SmtpUser`], appInfo[`SmtpPassword`], appInfo[`SmtpHost`], email.(string), ``, `<a href="`+appInfo[`UserResetPasswordPageUrl`]+url.QueryEscape(If(App.IsSSL, "https://", "http://").(string)+appInfo[`AppHost`]+`/v1/passwordReset/`+email.(string)+`/`+tmp)+`">请在5分钟以内，点击验证<a/>`, `html`)
		return gin.H{"success": true, "msg": "验证邮件已发送，请在5分钟内登陆邮箱并验证"}, nil
	}
}

// 邮箱验证地址
func (i *IBUser)EmailPasswordReset(c *Context) (gin.H, error) {
	email := c.GinContext.Param("email")
	yzm := c.GinContext.Param("yzm")

	data, err := Json(c.GinContext)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
	}
	new_password, ok := (*data)[`new_password`]
	if !ok {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("缺少new_password"))
	}

	tmp, found := App.Cache.Get(`passwordReset:` + email)
	if found {
		if yzm == tmp.(string) {
			var err error
			tx, err := App.Db.Begin()
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
			defer tx.Commit()

			_, err = tx.Query(`UPDATE public."_User" SET "password"='`+ SqlStrFilter(fmt.Sprint(new_password)) +`' WHERE "email"='` + email + `'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 100, errors.New("未知错误"))
			}
			App.Cache.Delete(`passwordReset:` + email)
			return gin.H{"success": true}, nil
		} else {
			return gin.H{"success": false, "msg": "未知验证邮件"}, nil
		}
	} else {
		return gin.H{"success": false, "msg": "验证邮件已过期，请重新请求发送验证邮件"}, nil
	}
}