//master.go
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
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"errors"
	"fmt"
)

type IbMaster struct {

}

var ibmaster IbMaster


/*
post
className 表名
mode 模式
1:限制写入
对象创建者可读、可写，其他人可读、不可写

2:限制读取
对象创建者可读、可写，其他人不可读、不可写

3:限制所有
所有人不可写，仅对象创建者可读
可在服务端使用 MasterKey 绕过权限控制

4:无限制
所有人可读、可写
务必自行增加访问权限控制
*/
func (m *IbMaster) CreateClass(c *Context) (gin.H,error) {
	if c.IsMaster{
		data, err := Json(c.GinContext)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		className:=(*data)[`className`].(string)
		err=CreateTable(className,map[string]interface{}{"get": map[string]interface{}{"type": "all"}, "create": map[string]interface{}{"type": "all"}, "update": map[string]interface{}{"type": "all"}, "delete": map[string]interface{}{"type": "all"}, "find": map[string]interface{}{"type": "all"}})
		if err != nil {
			return nil, err
		}
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		switch (*data)[`mode`].(float64) {
		case 1:
			_,err=tx.Query(`ALTER TABLE public."`+className+`" ALTER COLUMN "ACL" SET DEFAULT '{"*":{"write":false,"read":true},"_owner":{"write":true,"read":true}}'::jsonb; `)
			_,err=tx.Query(`UPDATE public.__field SET "default"='{"*":{"write":false,"read":true},"_owner":{"write":true,"read":true}}' where "className"='`+className+ `' and "fieldName"='ACL'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
		case 2:
			_,err=tx.Query(`ALTER TABLE public."`+className+`" ALTER COLUMN "ACL" SET DEFAULT '{"*":{"write":false,"read":false},"_owner":{"write":true,"read":true}}'::jsonb; `)
			_,err=tx.Query(`UPDATE public.__field SET "default"='{"*":{"write":false,"read":false},"_owner":{"write":true,"read":true}}'  where "className"='`+className+ `' and "fieldName"='ACL'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
		case 3:
			_,err=tx.Query(`ALTER TABLE public."`+className+`" ALTER COLUMN "ACL" SET DEFAULT '{"*":{"write":false,"read":false},"_owner":{"write":false,"read":true}}'::jsonb; `)
			_,err=tx.Query(`UPDATE public.__field SET "default"='{"*":{"write":false,"read":false},"_owner":{"write":false,"read":true}}'  where "className"='`+className+ `' and "fieldName"='ACL'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
		case 4:
			_,err=tx.Query(`ALTER TABLE public."`+className+`" ALTER COLUMN "ACL" SET DEFAULT '{"*":{"write":true,"read":true},"_owner":{"write":true,"read":true}}'::jsonb; `)
			_,err=tx.Query(`UPDATE public.__field SET "default"='{"*":{"write":true,"read":true},"_owner":{"write":true,"read":true}}'  where "className"='`+className+ `' and "fieldName"='ACL'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
		}
		tx.Commit()
		RefreshClassList()
		return gin.H{"success":true},nil
	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}


func (m *IbMaster) GetAllClass(c *Context) (gin.H,error)  {
	if c.IsMaster{
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		tmp,err:=tx.Query(`SELECT * from public."__table"`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		res:=GetJSON(tmp)
		return gin.H{"results":*res},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) SetClassAcl(c *Context) (gin.H,error) {
	if c.IsMaster{
		data, err := Json(c.GinContext)
		className := c.GinContext.Param(`className`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}


		_,err =tx.Query(`UPDATE public."__table" SET "ACL"='`+Map2json(*data)+`'::jsonb WHERE "className"='`+SqlStrFilter(className)+`'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		tx.Commit()
		RefreshClassList()
		return gin.H{"success":true},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) DeleteClass(c *Context) (gin.H,error) {
	if c.IsMaster{
		className := c.GinContext.Param(`className`)
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}


		_,err =tx.Query(`DELETE FROM public."__table" WHERE "className"='`+SqlStrFilter(className)+`'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		_,err =tx.Query(`DROP TABLE public."`+SqlStrFilter(className)+`"`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		_,err =tx.Query(`DELETE FROM public."__field" WHERE "className"='`+SqlStrFilter(className)+`'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		tx.Commit()
		RefreshClassList()
		return gin.H{"success":true},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) DeleteClassAllData(c *Context) (gin.H,error) {
	if c.IsMaster{
		className := c.GinContext.Param(`className`)
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		_,err =tx.Query(`DELETE FROM public."`+className+ `" `)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		return gin.H{"success":true},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) GetField(c *Context) (gin.H,error) {
	if c.IsMaster{
		className := c.GinContext.Param(`className`)
		fieldName := c.GinContext.Param(`fieldName`)
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		tmp,err:=tx.Query(`SELECT * from public."__field" WHERE "className"='`+SqlStrFilter(className)+`' and "fieldName"='`+SqlStrFilter(fieldName)+ `'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		res:=GetJSON(tmp)
		if len(*res)==0{
			return nil, NewError(http.StatusBadRequest, 300, errors.New("没有该字段"))
		}
		return (*res)[0],nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) GetFieldInClass(c *Context) (gin.H,error) {
	if c.IsMaster{
		className := c.GinContext.Param(`className`)
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		tmp,err:=tx.Query(`SELECT * from public."__field" WHERE "className"='`+SqlStrFilter(className)+`'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		res:=GetJSON(tmp)

		return gin.H{"results":*res},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) GetAllField(c *Context) (gin.H,error) {
	if c.IsMaster{
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		tmp,err:=tx.Query(`SELECT * from public."__field"`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		res:=GetJSON(tmp)

		return gin.H{"results":*res},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) CreateField(c *Context) (gin.H,error) {
	if c.IsMaster{
		data, err := Json(c.GinContext)
		className := c.GinContext.Param(`className`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		err =CreateField(className,(*data)[`fieldName`].(string),(*data)[`type`].(string),(*data)[`default`].(string),(*data)[`onlyR`].(bool),(*data)[`ownerRW`].(bool),(*data)[`notNull`].(bool),(*data)[`notSee`].(bool),(*data)[`notes`].(string),(*data)[`relationTo`].(string),(*data)[`autoIncrease`].(bool))
		if err != nil {
			return nil, err
		}
		RefreshClassList()
		return gin.H{"success":true},nil
	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) DeleteField(c *Context) (gin.H,error) {
	if c.IsMaster{
		className := c.GinContext.Param(`className`)
		fieldName := c.GinContext.Param(`fieldName`)
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}


		_,err =tx.Query(`DELETE FROM public."__field" WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		_,err =tx.Query(`alter table public."`+className+ `" drop column "`+fieldName+ `"; `)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		tx.Commit()
		RefreshClassList()
		return gin.H{"success":true},nil
	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

func (m *IbMaster) UpdateField(c *Context) (gin.H,error) {
	if c.IsMaster{
		data, err := Json(c.GinContext)
		className := c.GinContext.Param(`className`)
		fieldName := c.GinContext.Param(`fieldName`)
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}


		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		_default, ok := (*data)[`default`]
		if ok {
			_,err =tx.Query(`UPDATE public."__field" SET "default"='`+SqlStrFilter(fmt.Sprint(_default))+ `' WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
			}
			classList := GetClassList()
			class, ok := classList[className]
			if !ok {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("找不到Class"))
			}
			if va, ok := class.FieldList[fieldName]; ok {

				switch va.Type {
				case "number":

						_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT `+fmt.Sprint(_default)+`;`)
						if err != nil {
							return nil,NewError(http.StatusBadRequest, 101, err)
						}

				case "string":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT '`+SqlStrFilter(fmt.Sprint(_default.(string)))+`';`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				case "boolean":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT `+fmt.Sprint(_default)+`;`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				case "date":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT '`+SqlStrFilter(fmt.Sprint(_default))+`';`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				case "object":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT '`+SqlStrFilter(fmt.Sprint(_default))+`'::jsonb ;`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				case "array":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT '`+SqlStrFilter(fmt.Sprint(_default))+`'::jsonb ;`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				case "pointer":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT '`+SqlStrFilter(fmt.Sprint(_default))+`'::jsonb ;`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				case "relation":
					_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET DEFAULT '`+SqlStrFilter(fmt.Sprint(_default))+`'::jsonb ;`)
					if err != nil {
						return nil,NewError(http.StatusBadRequest, 101, err)
					}
				}

			}else {
				return nil, NewError(http.StatusBadRequest, 300, errors.New("找不到Field"))
			}
		}
		notes, ok := (*data)[`notes`]
		if ok {
			_,err =tx.Query(`UPDATE public."__field" SET "notes"='`+SqlStrFilter(notes.(string))+ `' WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, err)
			}
		}
		onlyR, ok := (*data)[`onlyR`]
		if ok {
			_,err =tx.Query(`UPDATE public."__field" SET "onlyR"=`+fmt.Sprint(onlyR.(bool))+ ` WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, err)
			}
		}
		notNull, ok := (*data)[`notNull`]
		if ok {
			_,err =tx.Query(`UPDATE public."__field" SET "notNull"=`+fmt.Sprint(notNull.(bool))+ ` WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, err)
			}

			if notNull.(bool){
				_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" SET NOT NULL ;`)
				if err != nil {
					return nil, NewError(http.StatusBadRequest, 300, err)
				}
			}else {
				_, err = tx.Query(`ALTER TABLE public."` + className + `" ALTER COLUMN "` + fieldName + `" DROP NOT NULL ;`)
				if err != nil {
					return nil, NewError(http.StatusBadRequest, 300, err)
				}
			}


		}
		notSee, ok := (*data)[`notSee`]
		if ok {
			_,err =tx.Query(`UPDATE public."__field" SET "notSee"=`+fmt.Sprint(notSee.(bool))+ ` WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, err)
			}
		}
		ownerRW, ok := (*data)[`ownerRW`]
		if ok {
			_,err =tx.Query(`UPDATE public."__field" SET "ownerRW"=`+fmt.Sprint(ownerRW.(bool))+ ` WHERE "className"='`+SqlStrFilter(className)+ `' and "fieldName"='`+SqlStrFilter(fieldName)+`'`)
			if err != nil {
				return nil, NewError(http.StatusBadRequest, 300, err)
			}
		}

		tx.Commit()
		RefreshClassList()
		return gin.H{"success":true},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}


func (m *IbMaster) GetAppInfo(c *Context) (gin.H,error)  {
	if c.IsMaster{
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		defer tx.Commit()

		tmp,err:=tx.Query(`SELECT * from public."__config"`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		res:=GetJSON(tmp)
		return gin.H{"results":*res},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}


func (m *IbMaster) SetAppInfo(c *Context) (gin.H,error) {
	if c.IsMaster{
		data, err := Json(c.GinContext)
		key := c.GinContext.Param(`key`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		tx,err:=App.Db.Begin()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}


		_,err =tx.Query(`UPDATE public."__config" SET "value"='`+SqlStrFilter(fmt.Sprint((*data)[`value`]))+`' WHERE "key"='`+SqlStrFilter(key)+`'`)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
		}
		tx.Commit()
		RefreshAppInfo()
		return gin.H{"success":true},nil

	}else {
		return nil,NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}