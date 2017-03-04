//acl.go
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
)

/*
权限优先级，在acl中附加属性
*/

//
func Authenticate(c *Context) error {
	id, ok := c.GinContext.Request.Header[`X-Ic-Id`]
	if !ok || len(id) == 0 {
		return NewError(http.StatusBadRequest, 200, errors.New("未知访问"))
	}
	key, ok := c.GinContext.Request.Header[`X-Ic-Key`]
	if !ok || len(key) == 0 {
		return NewError(http.StatusBadRequest, 200, errors.New("未知访问"))
	}
	appInfo := GetAppInfo()
	aId := appInfo[`AppId`]
	aKey := appInfo[`AppKey`]
	aMasterKey := appInfo[`MasterKey`]

	if id[0] == aId && key[0] == aKey { //普通

		c.IsMaster = false

	} else if id[0] == aId && key[0] == aMasterKey+`,master` { //master

		c.IsMaster = true

	} else {
		fmt.Println(aId + ` ` + aKey)
		return NewError(http.StatusBadRequest, 200, errors.New("未知访问"))
	}

	session, ok := c.GinContext.Request.Header[`X-Ic-Session`]
	if ok && len(session) != 0 {
		user, isLogin := ibuser.GetUserBySession(session[0])
		c.IsLogin = isLogin
		c.User = *user
	}

	return nil
}

func AclFilter(c *Context, q *Query) error {
	//master
	//表操作权限过滤
	//字段权限过滤，字段填充
	//ACL条件过滤

	classList := GetClassList()
	class, ok := classList[q.ClassName]
	if !ok {
		return NewError(http.StatusBadRequest, 200, errors.New("找不到Class"))
	}
	classAcl := class.Acl

	q.IsMaster = c.IsMaster
	//判断表操作权限
	if !c.IsMaster {
		switch q.Method {
		case `scan`:
		default:
			switch (classAcl[q.Method].(map[string]interface{}))[`type`].(string) {
			case `all`:

			case `sessionUser`:
				if !c.IsLogin {
					return NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
				}
			case `special`:
				if is_array_interface((classAcl[q.Method].(map[string]interface{}))[`objects`].([]interface{}), c.User.Id) {

				} else {
					flat := true
					for _, va := range c.User.Role {
						if is_array_interface((classAcl[q.Method].(map[string]interface{}))[`objects`].([]interface{}), `role:`+va) {

							flat = false
							break
						}
					}
					if flat {
						return NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
					}
				}
			default:
			}
		}

	}

	//过滤和补全字段
	switch q.Method {
	case `find`, `get`:
		if q.Keys != nil && len(*q.Keys) != 0 {
			if !c.IsMaster{
				tmp := *q.Keys
				q.Keys = &[]string{}
				err := FieldFilter1(q.ClassName, q, c, ``, tmp)
				if err != nil {
					return NewError(http.StatusBadRequest, 201, err)
				}
			}
		} else {
			q.Keys = &[]string{}
			err := FieldFilter2(q.ClassName, q, c, ``)
			if err != nil {
				return NewError(http.StatusBadRequest, 201, err)
			}
		}
	case `update`:
		if !c.IsMaster {
			err := FieldFilter3(q, c)
			if err != nil {
				return NewError(http.StatusBadRequest, 201, err)
			}
		}
		(*q.Data)[`updatedAt`]=`current_timestamp`
	case `create`:
		if !c.IsMaster {
			err := FieldFilter4(q, c)
			if err != nil {
				return NewError(http.StatusBadRequest, 201, err)
			}
		}

		(*q.Data)[`createdAt`]=`current_timestamp`
		(*q.Data)[`updatedAt`]=`current_timestamp`
	}

	return nil
}

// create
func FieldFilter4(q *Query, c *Context) error {
	isUser := false
	if q.ClassName == `_User` {
		isUser = true
	}
	classList := GetClassList()
	class, ok := classList[q.ClassName]
	if !ok {
		return errors.New("找不到Class")
	}
	for key, _ := range *q.Data {
		if va, ok := class.FieldList[key]; ok {
			if !va.OnlyR {
				if isUser && va.OwnerRW && c.IsLogin { //_User 只限拥有者读写  已登录的用户
					if va, ok := (*q.Data)[`id`]; ok && va == c.User.Id {

					} else {
						delete(*q.Data, key)
					}
				} else if (isUser && !va.OwnerRW) || (!isUser) { // _User 不限只限拥有者读写 或 不是_User

				} else {
					delete(*q.Data, key)
				}

			} else {
				delete(*q.Data, key)
			}
		} else {
			return errors.New("找不到该field")
		}
	}

	if va,ok:=(*q.Data)[`ACL`];ok{
		if  va2,ok2:=(va.(map[string]interface{}))[`_owner`];ok2{
			tmp:=va2
			delete((*q.Data)[`ACL`].(map[string]interface{}),`_owner`)
			if c.IsLogin{
				((*q.Data)[`ACL`].(map[string]interface{}))[c.User.Id]=tmp
			}
		}
	}else {
		tmp1,err:=Json2map(class.FieldList[`ACL`].Default)
		if err!=nil {
			return errors.New("参数解析错误")
		}
		(*q.Data)[`ACL`]= *tmp1
		if  va2,ok2:=(va.(map[string]interface{}))[`_owner`];ok2{
			tmp:=va2
			delete((*q.Data)[`ACL`].(map[string]interface{}),`_owner`)
			if c.IsLogin{
				((*q.Data)[`ACL`].(map[string]interface{}))[c.User.Id]=tmp
			}
		}
	}

	return nil
}

// update

func FieldFilter3(q *Query, c *Context) error {
	isUser := false
	if q.ClassName == `_User` {
		isUser = true
	}
	classList := GetClassList()
	class, ok := classList[q.ClassName]
	if !ok {
		return errors.New("找不到Class")
	}
	for key, _ := range *q.Data {
		if va, ok := class.FieldList[key]; ok {
			if !va.OnlyR {
				if isUser && va.OwnerRW && c.IsLogin { //_User 只限拥有者读写  已登录的用户
					if q.Id != c.User.Id {
						delete(*q.Data, key)
					}
				} else if (isUser && !va.OwnerRW) || (!isUser) { // _User 不限拥有者读写 或 不是_User

				} else {
					delete(*q.Data, key)
				}

			} else {
				delete(*q.Data, key)
			}
		} else {
			return errors.New("找不到该field")
		}
	}

	return nil
}

//建立在 有Keys  get和find请求下
func FieldFilter1(className string, q *Query, c *Context, prefix string, keys []string) error {
	isUser := false
	if className == `_User` {
		isUser = true
	}
	classList := GetClassList()
	class, ok := classList[className]
	if !ok {
		return errors.New("找不到Class")
	}

		for _, va := range class.FieldList {
			if is_array(keys, getFieldName(prefix, va.FieldName)) {
				if !va.NotSee {
					if !c.IsLogin && isUser && va.OwnerRW { //未登录状态操作 只限用户读写字段
						DeleteIncludeValue(q, getFieldName(prefix, va.FieldName))
					} else {
						flat := false
						if isUser && va.OwnerRW {
							if q.DeleteKeys == nil {
								q.DeleteKeys = &map[string]bool{}
							}
							q.UserId = c.User.Id
							flat = true
							(*q.DeleteKeys)[(getFieldName(prefix, va.FieldName))] = false
						}

						if va.Type == `pointer` {

							if q.Include != nil && len(*q.Include) != 0 {
								if is_array(*q.Include, getFieldName(prefix, va.FieldName)) {
									if flat {
										(*q.DeleteKeys)[getFieldName(prefix, va.FieldName)] = true
									}
									FieldFilter1(va.RelationTo, q, c, getFieldName(prefix, va.FieldName), keys)
								} else {
									tmp := append(*q.Keys, getFieldName(prefix, va.FieldName))
									q.Keys = &tmp
								}
							} else {
								tmp := append(*q.Keys, getFieldName(prefix, va.FieldName))
								q.Keys = &tmp
							}

						} else if va.Type == `relation` {

							tmp := append(*q.Keys, getFieldName(prefix, `$`+va.FieldName))
							q.Keys = &tmp

						} else {
							tmp := append(*q.Keys, va.FieldName)
							q.Keys = &tmp
						}
					}

				} else {
					DeleteIncludeValue(q, getFieldName(prefix, va.FieldName))
				}
			} else {
				DeleteIncludeValue(q, getFieldName(prefix, va.FieldName))
			}

		}

	return nil
}

//建立在 无Keys  get和find请求下
func FieldFilter2(className string, q *Query, c *Context, prefix string) error {
	isUser := false
	if className == `_User` {
		isUser = true
	}
	classList := GetClassList()
	class, ok := classList[className]
	if !ok {
		return errors.New("找不到Class")
	}

	if !c.IsMaster {
		for _, va := range class.FieldList {
			if !va.NotSee {
				if !c.IsLogin && isUser && va.OwnerRW { //未登录状态操作 只限用户读写字段
					DeleteIncludeValue(q, getFieldName(prefix, va.FieldName))
				} else {
					flat := false
					if isUser && va.OwnerRW {
						if q.DeleteKeys == nil {
							q.DeleteKeys = &map[string]bool{}
						}
						q.UserId = c.User.Id
						flat = true
						(*q.DeleteKeys)[getFieldName(prefix, va.FieldName)] = false
					}

					if va.Type == `pointer` {

						if q.Include != nil && len(*q.Include) != 0 {
							if is_array(*q.Include, getFieldName(prefix, va.FieldName)) {
								if flat {
									(*q.DeleteKeys)[getFieldName(prefix, va.FieldName)] = true
								}
								FieldFilter2(va.RelationTo, q, c, getFieldName(prefix, va.FieldName))
							} else {
								tmp := append(*q.Keys, getFieldName(prefix, va.FieldName))
								q.Keys = &tmp
							}
						} else {
							tmp := append(*q.Keys, getFieldName(prefix, va.FieldName))
							q.Keys = &tmp
						}

					} else if va.Type == `relation` {

						tmp := append(*q.Keys, getFieldName(prefix, `$`+va.FieldName))
						q.Keys = &tmp

					} else {
						tmp := append(*q.Keys, va.FieldName)
						q.Keys = &tmp
					}
				}

			} else {
				DeleteIncludeValue(q, getFieldName(prefix, va.FieldName))
			}
		}
	} else {
		for _, va := range class.FieldList {
			if va.Type == `pointer` {
				if q.Include != nil && len(*q.Include) != 0 {
					if is_array(*q.Include, getFieldName(prefix, va.FieldName)) {
						FieldFilter2(va.RelationTo, q, c, getFieldName(prefix, va.FieldName))
					} else {
						tmp := append(*q.Keys, getFieldName(prefix, va.FieldName))
						q.Keys = &tmp
					}
				} else {
					tmp := append(*q.Keys, getFieldName(prefix, va.FieldName))
					q.Keys = &tmp
				}
			} else if va.Type == `relation` {
				tmp := append(*q.Keys, getFieldName(prefix, `$`+va.FieldName))
				q.Keys = &tmp
			} else {
				tmp := append(*q.Keys, va.FieldName)
				q.Keys = &tmp
			}

		}

	}

	return nil
}

func getFieldName(prefix string, field string) string {
	if prefix == `` {
		return field
	} else {
		return prefix + `.` + field
	}
}

func DeleteIncludeValue(q *Query, value string) {
	if q.Include != nil {
		for key2, va2 := range *q.Include {
			if va2 == value {
				(*q.Include)[key2] = ``
			}
		}
	}
}
