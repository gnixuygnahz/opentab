//db.go
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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"qiniupkg.com/api.v7/kodo"
)

/*
错误处理
事务
参数过滤
relation
time
字段首字母小写  表名首字母大写

*/

type Query struct {
	Method     string
	ClassName  string
	Id         string
	Include    *[]string
	Limit      int
	Skip       int
	Order      *[]string
	Keys       *[]string
	RelationTo string
	Where      *map[string]interface{}
	Data       *map[string]interface{}
	NameMap    *map[string]string
	Action     string
	IsCount    bool

	DeleteKeys *map[string]bool //_User 权限 需要判断删除的 字段
	UserId     string           //判断所需的用户Id
	IsMaster   bool

	FetchWhenSave bool

	Context *Context
}

func GenQuery(method string, c *gin.Context, context *Context) (*Query, error) {
	var err error
	r := Query{}
	r.Method = strings.ToLower(method)
	r.Context = context
	//r.Where=&map[string]interface{}{}
	nameMap := map[string]string{}
	r.NameMap = &nameMap
	r.IsCount = false
	r.FetchWhenSave = false
	//r.Include=&[]string{}
	//r.Order=&[]string{}
	//r.Keys=&[]string{}

	switch r.Method {
	case "create":
		r.ClassName = c.Param("className")
		r.Data, err = Json(c)
		if err != nil {
			return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
	case "get":
		r.ClassName = c.Param("className")
		r.Id = c.Param("objectId")
		keys := c.DefaultQuery("keys", "")
		include := c.DefaultQuery("include", "")

		fetchWhenSave := c.DefaultQuery("fetchWhenSave", "")

		if fetchWhenSave == `true` {
			r.FetchWhenSave = true
		}

		if keys != "" {
			item := strings.Split(keys, `,`)
			r.Keys = &item
		}
		if include != "" {
			item := strings.Split(include, `,`)
			r.Include = &item
		}

	case "update":
		r.ClassName = c.Param("className")
		r.Id = c.Param("objectId")
		where := c.DefaultQuery("where", "")
		fetchWhenSave := c.DefaultQuery("fetchWhenSave", "")

		if fetchWhenSave == `true` {
			r.FetchWhenSave = true
		}

		if where != "" {
			r.Where, _ = Json2map(where)
		}

		r.Data, err = Json(c)
		if err != nil {
			return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
	case "scan":
		r.ClassName = c.Param("className")
		r.Action = c.Param("action")

		where := c.DefaultQuery("where", "")

		if where != "" {
			r.Where, err = Json2map(where)
			if err != nil {
				return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
			}
		}

		r.Data, err = Json(c)
		if err != nil {
			return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
	case "delete":
		r.ClassName = c.Param("className")
		r.Id = c.Param("objectId")
		where := c.DefaultQuery("where", "")

		if where != "" {
			r.Where, err = Json2map(where)
			if err != nil {
				return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
			}
		}
	case "find":
		r.ClassName = c.Param("className")

		limit := c.DefaultQuery("limit", "-1")
		skip := c.DefaultQuery("skip", "-1")
		order := c.DefaultQuery("order", "")
		include := c.DefaultQuery("include", "")
		keys := c.DefaultQuery("keys", "")
		where := c.DefaultQuery("where", "")
		count := c.DefaultQuery("count", "0")

		if count != `0` {
			r.IsCount = true
		}

		if limit != "" {
			r.Limit, err = strconv.Atoi(limit)
			if err != nil {
				return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
			}
		} else {
			r.Limit = -1
		}
		if skip != "" {
			r.Skip, err = strconv.Atoi(skip)
			if err != nil {
				return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
			}
		} else {
			r.Skip = -1
		}
		if order != "" {
			item := strings.Split(order, ",")
			r.Order = &item
		}
		if include != "" {
			item := strings.Split(include, ",")
			r.Include = &item
		}
		if keys != "" {
			item := strings.Split(keys, ",")
			r.Keys = &item
		}
		if where != "" {
			r.Where, err = Json2map(where)
			if err != nil {
				return &r, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
			}
		}
	}

	return &r, nil
}

func CreateTable(className string, acl map[string]interface{}) error {
	var err error

	tx, err := App.Db.Begin()
	if err != nil {
		return NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()
	_, err = tx.Query(`CREATE TABLE public."` + className + `"
	(
    	id bigserial NOT NULL,
   	"ACL" jsonb DEFAULT '{"*":{"read":true,"write":true}}'::jsonb,
    	"createdAt" timestamp with time zone,
    	"updatedAt" timestamp with time zone,
    	PRIMARY KEY (id)
	)
	WITH (
	    OIDS = FALSE
	)
	TABLESPACE pg_default;

	ALTER TABLE public."` + className + `"
    	OWNER to postgres;`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}
	a := Map2json(acl)
	_, err = tx.Query(`INSERT INTO public."__table" ("className", "ACL") VALUES ('` + className + `', '` + a + `'::jsonb)`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}
	_, err = tx.Query(`INSERT INTO public.__field(
	 "className", "fieldName", "_type", "relationTo", "default", "notes", "onlyR", "ownerRW", "notNull", "notSee","autoIncrease")
	VALUES ( '` + className + `', 'id', 'id', '', '', '', false, false,false,false,false)`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}
	_, err = tx.Query(`INSERT INTO public.__field(
	 "className", "fieldName", "_type", "relationTo", "default", "notes", "onlyR", "ownerRW", "notNull", "notSee","autoIncrease")
	VALUES ( '` + className + `', 'ACL', 'object', '', '` + Map2json(map[string]interface{}{"*": map[string]interface{}{"write": true, "read": true}}) + `', '', false, false,false,false,false);`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}
	_, err = tx.Query(`INSERT INTO public.__field(
	 "className", "fieldName", "_type", "relationTo", "default", "notes", "onlyR", "ownerRW", "notNull", "notSee","autoIncrease")
	VALUES ( '` + className + `', 'createdAt', 'date', '', '', '', false, false,false,false,false);`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}
	_, err = tx.Query(`INSERT INTO public.__field(
	 "className", "fieldName", "_type", "relationTo", "default", "notes", "onlyR", "ownerRW", "notNull", "notSee","autoIncrease")
	VALUES ( '` + className + `', 'updatedAt', 'date', '', '', '', false, false,false,false,false);`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}

	return err
}

func CreateField(className string, fieldName string, _type string, _default string, onlyR bool, ownerRW bool, notNull bool, notSee bool, notes string, relationTo string, autoIncrease bool) error {

	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	switch strings.ToLower(_type) { //实际插入表属性
	case "number":
		if autoIncrease {
			_, err = tx.Query(`ALTER TABLE public."` + className + `"
    			ADD COLUMN "` + fieldName + `" bigserial ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT `+_default, ``).(string) + `;`)
			if err != nil {
				return NewError(http.StatusBadRequest, 101, err)
			}
		} else {
			_, err = tx.Query(`ALTER TABLE public."` + className + `"
    			ADD COLUMN "` + fieldName + `" double precision ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT `+_default, ``).(string) + `;`)
			if err != nil {
				return NewError(http.StatusBadRequest, 101, err)
			}
		}
	case "string":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "` + fieldName + `" text ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT '`+_default+`'`, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	case "boolean":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "` + fieldName + `" boolean ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT `+_default, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	case "date":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
   		ADD COLUMN "` + fieldName + `" timestamp  with time zone ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT '`+_default+`'`, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	case "object":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "` + fieldName + `" jsonb ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT '`+_default+`'::jsonb`, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	case "array":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "` + fieldName + `" jsonb ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT '`+_default+`'::jsonb`, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	case "pointer":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "` + fieldName + `" jsonb ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT '`+_default+`'::jsonb`, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	case "relation":
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "` + fieldName + `" jsonb ` + If(notNull, `NOT NULL`, ``).(string) + `  ` + If(_default != ``, `DEFAULT '`+_default+`'::jsonb`, ``).(string) + `;`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
		_, err = tx.Query(`ALTER TABLE public."` + className + `"
    		ADD COLUMN "$` + fieldName + `" jsonb ` + `DEFAULT '{ "__type": "Relation","className":"` + className + `" }'::jsonb`)
		if err != nil {
			return NewError(http.StatusBadRequest, 101, err)
		}
	default:
		return NewError(http.StatusBadRequest, 300, errors.New("无此类型"))
	}
	//记录属性
	_, err = tx.Query(`INSERT INTO public."__field"(
	 "className", "fieldName", "_type", "relationTo", "default", "notes", "onlyR", "ownerRW", "notNull", "notSee","autoIncrease")
	VALUES ( '` + className + `', '` + fieldName + `', '` + strings.ToLower(_type) + `', ` + If(relationTo == ``, `NULL`, `'`+relationTo+`'`).(string) + `, '` + _default + `', '` + notes + `', ` + Bool2string(onlyR) + `, ` + Bool2string(ownerRW) + `,` + Bool2string(notNull) + `,` + Bool2string(notSee) + `,` + Bool2string(autoIncrease) + `);`)
	if err != nil {
		return NewError(http.StatusBadRequest, 101, err)
	}

	tx.Commit()
	return nil
}

//
func (q *Query) Create() (*map[string]interface{}, error) {
	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()

	if q.ClassName == `_File` {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("禁止操作_File"))
	}
	if _, ok := (*q.Data)["id"]; ok {

		id := strconv.Itoa(int((*q.Data)["id"].(float64)))
		delete(*q.Data, "id")

		keys := []string{}
		values := []interface{}{}
		i := 1
		for key, value := range *q.Data {
			if key == `createdAt` || key == `updatedAt` {
				keys = append(keys, ` "`+key+`"= `+value.(string)+` `)
			} else {
				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Array:
					keys = append(keys, ` "`+key+`"='`+SqlStrFilter(Array2json(value.([]interface{})))+`'::jsonb`)
				case reflect.Map:
					keys = append(keys, ` "`+key+`"='`+SqlStrFilter(Map2json(value.(map[string]interface{})))+`'::jsonb`)
				default:
					keys = append(keys, ` "`+key+`"=$`+strconv.Itoa(i))
					if rv.Kind() == reflect.String {
						values = append(values, SqlStrFilter(value.(string)))
						i++
					} else {
						values = append(values, value)
						i++
					}

				}
			}

		}
		row, err := tx.Query(`UPDATE public."`+q.ClassName+`" SET `+strings.Join(keys, `,`)+` WHERE "id"=`+id+` and `+getFieldAclStr(q.Context, `"ACL"`, IB_CAN_WRITE)+` Returning "id","createdAt"`+If(q.FetchWhenSave, `,"`+strings.Join(keys, `","`)+`";`, ``).(string), values...)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 101, err)
		}
		res := GetJSON(row)
		if err != nil {
			return &((*res)[0]), NewError(http.StatusBadRequest, 101, err)
		}
		return &((*res)[0]), nil
	} else {
		keys := []string{}
		s := []string{}

		values := []interface{}{}
		i := 1
		for key, value := range *q.Data {

			keys = append(keys, key)

			if key == `createdAt` || key == `updatedAt` {
				s = append(s, ` `+value.(string)+` `)
			} else {

				rv := reflect.ValueOf(value)
				switch rv.Kind() {
				case reflect.Array:
					s = append(s, `'`+SqlStrFilter(Array2json(value.([]interface{})))+`'::jsonb`)
				case reflect.Map:
					s = append(s, `'`+SqlStrFilter(Map2json(value.(map[string]interface{})))+`'::jsonb`)
				default:
					s = append(s, "$"+strconv.Itoa(i))
					if rv.Kind() == reflect.String {
						values = append(values, SqlStrFilter(value.(string)))
						i++
					} else {
						values = append(values, value)
						i++
					}
				}
			}
		}
		row, err := tx.Query(`INSERT INTO public."`+q.ClassName+`" ("`+strings.Join(keys, `","`)+`")
		VALUES ( `+strings.Join(s, `,`)+`) Returning "id","createdAt"`+If(q.FetchWhenSave, `,"`+strings.Join(keys, `","`)+`";`, ``).(string), values...)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 101, err)
		}
		res := GetJSON(row)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 101, err)
		}
		return &((*res)[0]), nil
	}

}

//
func (q *Query) Update() (*map[string]interface{}, error) {
	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()

	if q.ClassName == `_File` {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("禁止操作_File"))
	}

	delete(*q.Data, "id")
	keys := []string{}
	values := []interface{}{}
	i := 1
	for key, value := range *q.Data {
		if key == `createdAt` || key == `updatedAt` {
			keys = append(keys, ` "`+key+`"= `+value.(string)+` `)
		} else {
			switch value.(type) {
			case map[string]interface{}:
				if _, ok := (value.(map[string]interface{}))["__op"]; ok {
					switch strings.ToLower((value.(map[string]interface{}))["__op"].(string)) {
					case "Add":
						if len((value.(map[string]interface{}))["objects"].([]interface{})) != 0 {
							keys = append(keys, ` "`+key+`"="`+key+`"||'`+SqlStrFilter(Array2json((value.(map[string]interface{}))["objects"].([]interface{})))+`'::jsonb `)
						}
					case "Remove":
						if len((value.(map[string]interface{}))["objects"].([]int)) != 0 {
							query := ` "` + key + `"="` + key + `"`
							for i := 0; i < len((value.(map[string]interface{}))["objects"].([]int)); i++ {
								query += `-` + strconv.Itoa(((value.(map[string]interface{}))["objects"].([]int))[i])
							}
							keys = append(keys, query)
						}
					case "AddRelation":
						if len((value.(map[string]interface{}))["objects"].([]map[string]interface{})) != 0 {
							query := []string{}
							for _, value1 := range (value.(map[string]interface{}))["objects"].([]map[string]interface{}) {
								query = append(query, strconv.Itoa(value1["id"].(int)))
							}
							keys = append(keys, ` "`+key+`"="`+key+`"||'["`+SqlStrFilter(strings.Join(query, `","`))+`"]'::jsonb`)
						}
					case "RemoveRelation":
						if len((value.(map[string]interface{}))["objects"].([]map[string]interface{})) != 0 {
							query := ` "` + key + `"="` + key + `" `
							for _, value1 := range (value.(map[string]interface{}))["objects"].([]map[string]interface{}) {
								query += ` -'` + strconv.Itoa(value1["id"].(int)) + `'`
							}
							keys = append(keys, query)
						}
					case "Increment":
						keys = append(keys, ` "`+key+`"="`+key+`"+`+strconv.Itoa((value.(map[string]interface{}))["amount"].(int)))
					case "Decrement":
						keys = append(keys, ` "`+key+`"="`+key+`"-`+strconv.Itoa((value.(map[string]interface{}))["amount"].(int)))
					}
				} else {
					keys = append(keys, ` "`+key+`"='`+SqlStrFilter(Map2json(value.(map[string]interface{})))+`';:jsonb`)
				}
			case []interface{}:
				keys = append(keys, ` "`+key+`"='`+SqlStrFilter(Array2json(value.([]interface{})))+`';:jsonb`)
			default:
				keys = append(keys, ` "`+key+`"=$`+strconv.Itoa(i))

				rv := reflect.ValueOf(value)
				if rv.Kind() == reflect.String {
					values = append(values, SqlStrFilter(value.(string)))
					i++
				} else {
					values = append(values, value)
					i++
				}
			}
		}
	}
	if len(keys) > 0 {
		tmp, err := tx.Query(`UPDATE public."`+q.ClassName+`" SET `+strings.Join(keys, `,`)+` WHERE "id"=`+q.Id+` and `+getFieldAclStr(q.Context, `"ACL"`, IB_CAN_WRITE)+If(q.Where != nil, ` and `+RecWhere(*q.Where), ``).(string)+` Returning "updatedAt"`+If(q.FetchWhenSave, `,"`+strings.Join(keys, `","`)+`";`, ``).(string), values...)
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 101, err)
		}
		res := GetJSON(tmp)
		return &((*res)[0]), nil
	}
	return nil, NewError(http.StatusBadRequest, 100, errors.New(`参数解析错误`))

}

//
func (q *Query) Get() (*map[string]interface{}, error) {
	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()
	//初始化查询字符串
	query := "select "
	//建立表名映射集
	//$pointMap=[];

	//重建字段集
	for n := 0; n < len(*(q.Keys)); n++ {

		tmp := strings.Split((*q.Keys)[n], ".")
		last := tmp[len(tmp)-1]
		tmp = tmp[0 : len(tmp)-1]

		if len(tmp) != 0 {
			(*(q.Keys))[n] = ` "` + strings.Join(tmp, ".") + `"."` + last + `" as "` + (*(q.Keys))[n] + `"`
		} else {
			(*(q.Keys))[n] = ` "` + q.ClassName + `"."` + (*(q.Keys))[n] + `" as "` + (*(q.Keys))[n] + `" `
		}
	}

	//重建条件集

	query = query + strings.Join((*(q.Keys)), `,`)

	query = query + ` from "` + q.ClassName + `" as "` + q.ClassName + `" `
	if q.Include != nil {
		for n := 0; n < len(*(q.Include)); n++ {
			if (*(q.Include))[n] != `` {
				tmp := strings.Split((*(q.Include))[n], ".")
				last := tmp[len(tmp)-1]
				tmp = tmp[0 : len(tmp)-1]
				tmp2, err := q.QueryTableName((*(q.Include))[n])
				if err != nil {
					return nil, NewError(http.StatusBadRequest, 102, err)
				}
				query = query + ` left join "` + tmp2 + `" as "` + (*(q.Include))[n] + `" on "` + (*(q.Include))[n] + `".id = ` + If(len(strings.Split((*(q.Include))[n], ".")) == 1, `("`+q.ClassName+`".`+(*(q.Include))[n]+`->>'id')::int`, `("`+strings.Join(tmp, ".")+`".`+last+"->>'id')::int ").(string) + ` and ` + getFieldAclStr(q.Context, `"`+(*(q.Include))[n]+`"`+`."ACL"`, IB_CAN_READ) + ` `
			}
		}
	}
	query += ` where "id"=` + q.Id + ` and ` + getFieldAclStr(q.Context, `"ACL"`, IB_CAN_READ)
	re, err := tx.Query(query)
	if err != nil {
		fmt.Println(query)
		return nil, NewError(http.StatusBadRequest, 101, err)
	}
	res := GetJSON(re)

	if q.DeleteKeys != nil && len(*res) != 0 {
		for key, value := range *q.DeleteKeys {
			if q.UserId != ((*res)[0])[q.GetId(key)].(string) {
				if value {
					tmp := []string{}
					for key1, _ := range (*res)[0] {
						if strings.IndexAny(key1, key) == 0 {
							tmp = append(tmp, key1)
						}
					}
					for _, value1 := range tmp {
						delete((*res)[0], value1)
					}

				} else {
					delete((*res)[0], key)
				}
			}
		}
	}

	return &((*res)[0]), nil
}

//
func (q *Query) Find() (*[]map[string]interface{}, error) {

	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()
	//初始化查询字符串
	query := "select "
	//建立表名映射集
	//$pointMap=[];
	if !q.IsCount {
		//重建字段集 可拓展至where
		for n := 0; n < len(*(q.Keys)); n++ {
			tmp := strings.Split((*q.Keys)[n], ".")
			last := tmp[len(tmp)-1]
			tmp = tmp[0 : len(tmp)-1]

			if len(tmp) != 0 {
				(*(q.Keys))[n] = ` "` + strings.Join(tmp, ".") + `"."` + last + `" as "` + (*(q.Keys))[n] + `"`
			} else {
				(*(q.Keys))[n] = ` "` + q.ClassName + `"."` + (*(q.Keys))[n] + `" as "` + (*(q.Keys))[n] + `" `
			}
		}

		//重建条件集

		query = query + strings.Join((*(q.Keys)), `,`)
	} else {
		query += `COUNT(*)`
	}

	query = query + ` from "` + q.ClassName + `" as "` + q.ClassName + `" `

	if q.Include != nil {
		for n := 0; n < len(*(q.Include)); n++ {
			if (*(q.Include))[n] != `` {
				tmp := strings.Split((*(q.Include))[n], ".")
				last := tmp[len(tmp)-1]
				tmp = tmp[0 : len(tmp)-1]
				tmp2, err := q.QueryTableName((*(q.Include))[n])
				if err != nil {
					return nil, NewError(http.StatusBadRequest, 102, err)
				}
				query = query + ` left join "` + tmp2 + `" as "` + (*(q.Include))[n] + `" on "` + (*(q.Include))[n] + `".id = ` + If(len(strings.Split((*(q.Include))[n], ".")) == 1, `("`+q.ClassName+`".`+(*(q.Include))[n]+`->>'id')::int`, `("`+strings.Join(tmp, ".")+`".`+last+"->>'id')::int ").(string) + ` and ` + getFieldAclStr(q.Context, `"`+(*(q.Include))[n]+`"`+`."ACL"`, IB_CAN_READ) + ` `

			}
		}
	}

	query += ` where ` + getFieldAclStr(q.Context, `"ACL"`, IB_CAN_READ)

	if q.Where != nil {
		query += ` and ` + RecWhere(*q.Where)
	}
	if q.Limit != -1 && q.Skip != -1 {
		query += ` limit ` + fmt.Sprint(q.Limit) + ` offset ` + fmt.Sprint(q.Skip)
	} else if q.Limit != -1 && q.Skip == -1 {
		query += ` limit ` + fmt.Sprint(q.Limit) + ` `
	} else if q.Limit == -1 && q.Skip != -1 {
		query += ` limit 100 offset ` + fmt.Sprint(q.Skip)
	} else {
		query += ` limit 100 `
	}
	if q.Order != nil {
		query += ` order by `
		tmp := []string{}
		for _, value := range *q.Order {
			if value[0] == '-' {
				tmp = append(tmp, `"`+value[1:]+`" asc`)
			} else {
				tmp = append(tmp, `"`+value+`" desc`)
			}
		}
		query += strings.Join(tmp, `,`)
	}

	re, err := tx.Query(query)
	if err != nil {
		fmt.Println(query)
		return nil, NewError(http.StatusBadRequest, 101, err)
	}
	res := GetJSON(re)

	if !q.IsCount {
		if q.DeleteKeys != nil && len(*res) != 0 {
			for _, va1 := range *res {
				for key, value := range *q.DeleteKeys {
					if q.UserId != va1[q.GetId(key)].(string) {
						if value {
							tmp := []string{}
							for key1, _ := range va1 {
								if strings.IndexAny(key1, key) == 0 {
									tmp = append(tmp, key1)
								}
							}
							for _, value1 := range tmp {
								delete(va1, value1)
							}

						} else {
							delete(va1, key)
						}
					}
				}
			}

		}
	}

	return res, nil
}

func (q *Query) GetId(str string) string {
	tmp := strings.Split(str, ".")
	tmp = tmp[0 : len(tmp)-1]

	if len(tmp) != 0 {
		return strings.Join(tmp, ".") + `.id`
	} else {
		return `id`
	}
}

//
func (q *Query) Delete() error {
	tx, err := App.Db.Begin()
	defer tx.Commit()
	_, err = tx.Query(`DELETE FROM public."` + q.ClassName + `" WHERE "id"=` + q.Id + If(q.Where != nil, ` and `+RecWhere(*q.Where), ``).(string))
	if err != nil {
		return NewError(http.StatusBadRequest, 103, errors.New(`删除失败`))
	}
	return nil
}

func InitDb(db sql.DB) error {
	var err error
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Query(`
	CREATE TABLE public.__table
	(
    	"id" serial NOT NULL,
    	"className" character varying(30),
    	"ACL" jsonb,
    	PRIMARY KEY ("id")
	)
	WITH (
    	OIDS = FALSE
	)
	TABLESPACE pg_default;

	ALTER TABLE public.__table
    	OWNER to postgres;`)
	if err != nil {
		if strings.Contains(err.Error(), `exists`) {
			return nil
		}
		return err
	}
	_, err = tx.Query(`
	CREATE TABLE public.__field
	(
    	id bigserial NOT NULL,
    	"className" character varying(30) NOT NULL,
    	"fieldName" character varying(30) NOT NULL,
    	"_type" character varying(10) NOT NULL,
    	"relationTo" character varying(30),
    	"default" text,
    	"notes" text,
    	"onlyR" boolean,
    	"ownerRW" boolean,
    	"notNull" boolean,
    	"notSee" boolean,
    	"autoIncrease" boolean,
    	PRIMARY KEY (id)
	)
	WITH (
    	OIDS = FALSE
	)
	TABLESPACE pg_default;

	ALTER TABLE public.__field
    	OWNER to postgres;`)
	if err != nil {
		return err
	}
	_, err = tx.Query(`
	CREATE TABLE public.__config
(
    key text,
    value text,
    notes text
)
	`)
	if err != nil {
		return err
	}
	r := GetRand()
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('AppName','未定义','应用名称')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('AppHost','未定义','应用绑定域名（或IP）')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('AppId','` + Krand(24, KC_RAND_KIND_ALL, r) + `','')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('AppKey','` + Krand(24, KC_RAND_KIND_ALL, r) + `','')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('MasterKey','` + Krand(24, KC_RAND_KIND_ALL, r) + `','')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('IpLimit','1000','每分钟IP访问次数限制')`)

	//SMTP 设置
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('SmtpUser','-','')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('SmtpPassword','-','')`) //
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('SmtpHost','-','')`)

	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('QiniuAccessKey','-','七牛AccessKey')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('QiniuSecretKey','-','七牛SecretKey')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('QiniuBucket','-','七牛Bucket')`)
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('QiniuUrl','-','七牛Url访问地址')`)

	//
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('UserRegister','true','允许客户端注册')`)                        //允许客户端注册
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('UserRegisterNeedEmail','true','注册需要验证邮箱')`)              //注册需要验证邮箱
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('UserLoginRequireEmailVerified','true','未验证邮箱的用户禁止登陆')`)  //未验证邮箱的用户禁止登陆
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('UserResetPasswordRequireRelogin','true','重置密码后强制重新登陆')`) //重置密码后强制重新登陆
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('UserResetPasswordRequireOld','true','重置密码需要旧密码')`)       //重置密码需要旧密码
	_, err = tx.Query(`INSERT INTO public."__config"("key","value","notes") VALUES('UserResetPasswordPageUrl','-','密码重设跳转地址')`)

	tx.Commit()

	CreateTable(`_User`, map[string]interface{}{"get": map[string]interface{}{"type": "all"}, "create": map[string]interface{}{"type": "all"}, "update": map[string]interface{}{"type": "all"}, "delete": map[string]interface{}{"type": "special", "objects": []string{}}, "find": map[string]interface{}{"type": "special", "objects": []string{}}})
	CreateField(`_User`, `salt`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `email`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `sessionToken`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `password`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `username`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `emailVerified`, `boolean`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `mobilePhoneNumber`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_User`, `mobilePhoneVerified`, `boolean`, ``, false, false, false, false, ``, ``, false)

	EasyQuery(`ALTER TABLE public."_User"
    ADD CONSTRAINT username_uk UNIQUE ("username");

ALTER TABLE public."_User"
    ADD CONSTRAINT email_uk UNIQUE ("email");

ALTER TABLE public."_User"
    ADD CONSTRAINT sessionToken_uk UNIQUE ("sessionToken"); `)

	CreateTable(`_Role`, map[string]interface{}{"get": map[string]interface{}{"type": "all"}, "create": map[string]interface{}{"type": "special", "objects": []string{}}, "update": map[string]interface{}{"type": "special", "objects": []string{}}, "delete": map[string]interface{}{"type": "special", "objects": []string{}}, "find": map[string]interface{}{"type": "all"}})
	CreateField(`_Role`, `roles`, `relation`, ``, false, false, false, false, ``, `_Role`, false)
	CreateField(`_Role`, `users`, `relation`, ``, false, false, false, false, ``, `_User`, false)
	CreateField(`_Role`, `name`, `string`, ``, false, false, false, false, ``, ``, false)

	CreateTable(`_File`, map[string]interface{}{"get": map[string]interface{}{"type": "all"}, "create": map[string]interface{}{"type": "special", "objects": []string{}}, "update": map[string]interface{}{"type": "special", "objects": []string{}}, "delete": map[string]interface{}{"type": "special", "objects": []string{}}, "find": map[string]interface{}{"type": "all"}})
	CreateField(`_File`, `mime_type`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_File`, `key`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_File`, `name`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_File`, `url`, `string`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_File`, `metaData`, `object`, ``, false, false, false, false, ``, ``, false)
	CreateField(`_File`, `bucket`, `string`, ``, false, false, false, false, ``, ``, false)

	return nil
}

//func (q *Query) GenWhereString() string {
//	where:=""
//
//}

func RecWhere(where map[string]interface{}) string {
	arr := []string{}
	for key, value := range where {
		switch key {
		case "$relatedTo": //可拓展 条件查询
			if reflect.TypeOf(value) == reflect.TypeOf(map[string]interface{}{}) {
				if _, ok := value.(map[string]interface{})[`key`]; !ok {
					panic(`key not exist`)
				}
				if _, ok := value.(map[string]interface{})[`object`]; !ok {
					panic(`object not exist`)
				}
				object := value.(map[string]interface{})[`object`].(map[string]interface{})
				tmp := `( select public."` + value.(map[string]interface{})[`key`].(string) + `" from public."` + object[`className`].(string) + `" where "id"=` + fmt.Sprint(object[`id`].(float64)) + `)@>array[id]`
				arr = append(arr, tmp)
			} else {
				panic(`$relatedTo value is not a object`)
			}
		case "$or":
			if reflect.TypeOf(value) == reflect.TypeOf([]interface{}{}) && len(value.([]interface{})) != 0 {
				arrtmp := []string{}
				for _, value2 := range value.([]interface{}) {
					arrtmp = append(arrtmp, RecWhere(value2.(map[string]interface{})))
				}
				arr = append(arr, `( `+strings.Join(arrtmp, ` or `)+` )`)
			} else {
				panic(`$or value is not a array`)
			}
		default:
			switch value.(type) {
			case string:
				arr = append(arr, `( `+FieldName2SqlStr(key)+`='`+SqlStrFilter(fmt.Sprint(value))+`' )`)
			case float64:
				arr = append(arr, `( `+FieldName2SqlStr(key)+`=`+fmt.Sprint(value)+` )`)
			case map[string]interface{}:
				if va, ok := (value.(map[string]interface{}))[`__type`]; ok {
					switch va {
					case `Pointer`:
						arr = append(arr, ` `+FieldName2SqlStr(key)+`->>'id'=`+SqlStrFilter(fmt.Sprint((value.(map[string]interface{}))[`id`]))+`::text `)
					case `Relation`:
						arr = append(arr, ` `+FieldName2SqlStr(key)+`@>'[`+SqlStrFilter(fmt.Sprint((value.(map[string]interface{}))[`id`]))+`]'::jsonb `)
					}
				} else {
					for key2, value2 := range value.(map[string]interface{}) {
						switch key2 {
						case `$ne`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`<>'`+SqlStrFilter(fmt.Sprint(value2))+`' `)
							case float64:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`<>`+fmt.Sprint(value2)+` `)
							case map[string]interface{}:
								if va, ok := (value2.(map[string]interface{}))[`__type`]; ok {
									switch va.(string) {
									case `Date`:
										arr = append(arr, ` `+FieldName2SqlStr(key)+`<>'`+SqlStrFilter(fmt.Sprint((value2.(map[string]interface{}))[`iso`]))+`'::timestamp with time zone `)
									}
								} else {
									panic(`__type is null`)
								}
							}
						case `$lt`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`<'`+SqlStrFilter(fmt.Sprint(value2))+`' `)
							case float64:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`<`+fmt.Sprint(value2)+` `)
							case map[string]interface{}:
								if va, ok := (value2.(map[string]interface{}))[`__type`]; ok {
									switch va.(string) {
									case `Date`:
										arr = append(arr, ` `+FieldName2SqlStr(key)+`<'`+SqlStrFilter(fmt.Sprint((value2.(map[string]interface{}))[`iso`]))+`'::timestamp with time zone `)
									}
								} else {
									panic(`__type is null`)
								}
							}
						case `$lte`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`<='`+SqlStrFilter(fmt.Sprint(value2))+`' `)
							case float64:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`<=`+fmt.Sprint(value2)+` `)
							case map[string]interface{}:
								if va, ok := (value2.(map[string]interface{}))[`__type`]; ok {
									switch va.(string) {
									case `Date`:
										arr = append(arr, ` `+FieldName2SqlStr(key)+`<='`+SqlStrFilter(fmt.Sprint((value2.(map[string]interface{}))[`iso`]))+`'::timestamp with time zone `)
									}
								} else {
									panic(`__type is null`)
								}
							}
						case `$gt`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`>'`+SqlStrFilter(fmt.Sprint(value2))+`' `)
							case float64:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`>`+fmt.Sprint(value2)+` `)
							case map[string]interface{}:
								if va, ok := (value2.(map[string]interface{}))[`__type`]; ok {
									switch va.(string) {
									case `Date`:
										arr = append(arr, ` `+FieldName2SqlStr(key)+`>'`+SqlStrFilter(fmt.Sprint((value2.(map[string]interface{}))[`iso`]))+`'::timestamp with time zone `)
									}
								} else {
									panic(`__type is null`)
								}
							}
						case `$gte`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`>='`+SqlStrFilter(fmt.Sprint(value2))+`' `)
							case float64:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`>=`+fmt.Sprint(value2)+` `)
							case map[string]interface{}:
								if va, ok := (value2.(map[string]interface{}))[`__type`]; ok {
									switch va.(string) {
									case `Date`:
										arr = append(arr, ` `+FieldName2SqlStr(key)+`>='`+SqlStrFilter(fmt.Sprint((value2.(map[string]interface{}))[`iso`]))+`'::timestamp with time zone `)
									}
								} else {
									panic(`__type is null`)
								}
							}
						case `$regex`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+` ~ '`+SqlStrFilter(fmt.Sprint(value2))+`' `)
							}
						case `$in`:
							switch value2.(type) {
							case string:
								arr = append(arr, `  `+FieldName2SqlStr(key)+`<@array['`+SqlStrFilter(fmt.Sprint(value2))+`'] `)
							case float64:
								arr = append(arr, `  `+FieldName2SqlStr(key)+`<@array[`+fmt.Sprint(value2)+`] `)
							case []interface{}:
								if len(value2.([]interface{})) != 0 && reflect.TypeOf(value2.([]interface{})[0]) == reflect.TypeOf(``) {
									tmp := []string{}
									for _, value3 := range value2.([]interface{}) {
										tmp = append(tmp, SqlStrFilter(fmt.Sprint(value3)))
									}
									arr = append(arr, `  `+FieldName2SqlStr(key)+`<@array['`+fmt.Sprint(strings.Join(tmp, `','`))+`'] `)
								}
								if len(value2.([]interface{})) != 0 && reflect.TypeOf(value2.([]interface{})[0]) == reflect.TypeOf(float64(1)) {
									tmp := []string{}
									for _, value3 := range value2.([]interface{}) {
										tmp = append(tmp, fmt.Sprint(value3))
									}
									arr = append(arr, `  `+FieldName2SqlStr(key)+`<@array[`+fmt.Sprint(strings.Join(tmp, `,`))+`] `)
								}
							}
						case `$nin`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` not `+FieldName2SqlStr(key)+`<@array['`+SqlStrFilter(fmt.Sprint(value2))+`'] `)
							case float64:
								arr = append(arr, ` not `+FieldName2SqlStr(key)+`<@array[`+fmt.Sprint(value2)+`] `)
							case []interface{}:
								if len(value2.([]interface{})) != 0 && reflect.TypeOf(value2.([]interface{})[0]) == reflect.TypeOf(``) {
									tmp := []string{}
									for _, value3 := range value2.([]interface{}) {
										tmp = append(tmp, SqlStrFilter(fmt.Sprint(value3)))
									}
									arr = append(arr, ` not `+FieldName2SqlStr(key)+`<@array['`+fmt.Sprint(strings.Join(tmp, `','`))+`'] `)
								}
								if len(value2.([]interface{})) != 0 && reflect.TypeOf(value2.([]interface{})[0]) == reflect.TypeOf(float64(1)) {
									tmp := []string{}
									for _, value3 := range value2.([]interface{}) {
										tmp = append(tmp, fmt.Sprint(value3))
									}
									arr = append(arr, ` not `+FieldName2SqlStr(key)+`<@array[`+fmt.Sprint(strings.Join(tmp, `,`))+`] `)
								}
							}
						case `$all`:
							switch value2.(type) {
							case string:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`@>array['`+SqlStrFilter(fmt.Sprint(value2))+`'] `)
							case float64:
								arr = append(arr, ` `+FieldName2SqlStr(key)+`@>array[`+fmt.Sprint(value2)+`] `)

							case []interface{}:
								if len(value2.([]interface{})) != 0 && reflect.TypeOf(value2.([]interface{})[0]) == reflect.TypeOf(``) {
									tmp := []string{}
									for _, value3 := range value2.([]interface{}) {
										tmp = append(tmp, SqlStrFilter(fmt.Sprint(value3)))
									}
									arr = append(arr, `  `+FieldName2SqlStr(key)+`@>array['`+fmt.Sprint(strings.Join(tmp, `','`))+`'] `)
								}
								if len(value2.([]interface{})) != 0 && reflect.TypeOf(value2.([]interface{})[0]) == reflect.TypeOf(float64(1)) {
									tmp := []string{}
									for _, value3 := range value2.([]interface{}) {
										tmp = append(tmp, fmt.Sprint(value3))
									}
									arr = append(arr, `  `+FieldName2SqlStr(key)+`@>array[`+fmt.Sprint(strings.Join(tmp, `,`))+`] `)
								}
							}
						}
					}
				}

			}
		}
	}
	query := `( ` + strings.Join(arr, ` and `) + ` )`
	return query
}

func GetJSON(rows *sql.Rows) *[]map[string]interface{} {
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return &[]map[string]interface{}{}
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if reflect.TypeOf(val) == reflect.TypeOf([]uint8{}) {
				tmp, err := Json2map(string(b))
				if err == nil {
					v = *tmp
				} else {
					tmp, err := Json2array(string(b))
					if err == nil {
						v = *tmp
					} else {
						v = string(b)
					}
				}
			} else if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return &tableData
}

//func (qu *Query) GenWhereQuery() (string,error) {
//
//}

func (qu *Query) QueryTableName(queryName string) (string, error) {
	nameMap := qu.NameMap
	if _, ok := (*nameMap)[queryName]; ok {
		return (*qu.NameMap)[queryName], nil
	}

	query := strings.Split(queryName, ",")
	le := -1 //查询级
	for n := (len(query) - 2); n >= 0; n-- {
		q := ""
		for nn := 0; nn <= n; nn++ {
			q = q + query[nn] + If(nn == n, ``, `.`).(string)
		}
		if _, ok := (*nameMap)[q]; ok {
			le = n
			break
		}
	}
	//
	if le == -1 {

		class, ok := GetClassList()[qu.ClassName]
		if !ok {
			return ``, errors.New(`没有此Class`)
		}
		field, ok := ((class).FieldList)[query[0]]
		if !ok {
			return ``, errors.New(`没有此field`)
		}

		(*nameMap)[query[0]] = field.RelationTo
		le = 0
	}

	q := ""
	for nn := 0; nn <= le; nn++ {
		q = q + query[nn] + If(nn == le, ``, `.`).(string)
	}

	for le = le + 1; le < len(query); le++ {

		class, ok := GetClassList()[(*nameMap)[q]]
		if !ok {
			return ``, errors.New(`没有此Class`)
		}
		field, ok := (class.FieldList)[query[le]]
		if !ok {
			return ``, errors.New(`没有此field`)
		}

		(*nameMap)[q+`.`+query[le]] = field.RelationTo
		q = q + `.` + query[le]
	}
	//
	return (*nameMap)[queryName], nil
}

func EasyQuery(query string) *[]map[string]interface{} {
	tx, err := App.Db.Begin()
	CheckErr(err)
	defer tx.Commit()

	r, _ := tx.Query(query)

	res := GetJSON(r)

	return res

}

/*

 */

type ClassList map[string]*Class

type Class struct {
	ClassName string
	Acl       map[string]interface{}
	FieldList map[string]*Field
}

type Field struct {
	FieldName    string
	OnlyR        bool
	OwnerRW      bool
	NotNull      bool
	NotSee       bool
	AutoIncrease bool
	RelationTo   string
	Type         string
	Default      string
}

//
func GetClassList() (ret ClassList) {
	ret = ClassList{}
	tmp, found := App.Cache.Get(`app:ClassList`)
	if found {
		ret = tmp.(ClassList)
	} else {

		classRes := EasyQuery(`select * from public."__table"`)
		fieldRes := EasyQuery(`select * from public."__field"`)

		for _, class := range *classRes {
			fmt.Println(class[`ACL`].(map[string]interface{}))
			ret[class[`className`].(string)] = &Class{Acl: class[`ACL`].(map[string]interface{}), FieldList: map[string]*Field{}, ClassName: class[`className`].(string)}
			for _, field := range *fieldRes {
				if field[`className`].(string) == class[`className`].(string) {
					((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)] = &Field{}
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).FieldName = field[`fieldName`].(string)
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).OnlyR = field[`onlyR`].(bool)
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).OwnerRW = field[`ownerRW`].(bool)
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).NotNull = field[`notNull`].(bool)
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).NotSee = field[`notSee`].(bool)
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).AutoIncrease = field[`autoIncrease`].(bool)
					if tmp, ok := field[`relationTo`]; ok && tmp != nil {
						(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).RelationTo = field[`relationTo`].(string)
					} else {
						(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).RelationTo = ``
					}
					if tmp, ok := field[`default`]; ok && tmp != nil {
						(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).Default = field[`default`].(string)
					} else {
						(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).Default = ``
					}
					(((ret[class[`className`].(string)]).FieldList)[field[`fieldName`].(string)]).Type = field[`_type`].(string)
				}
			}
		}
		App.Cache.Set(`app:ClassList`, ret, -1)
	}
	return
}

func RefreshClassList() {
	App.Cache.Delete(`app:ClassList`)
}

/**/

type AppInfo map[string]string

//获取应用信息
func GetAppInfo() (ret AppInfo) {
	ret = AppInfo{}
	tmp, found := App.Cache.Get(`app:AppInfo`)
	if found {
		ret = tmp.(AppInfo)
	} else {
		tmp := EasyQuery(`select * from public."__config"`)
		for _, va := range *tmp {
			ret[va[`key`].(string)] = va[`value`].(string)
		}
	}

	return
}

const (
	IB_CAN_READ  = `read`
	IB_CAN_WRITE = `write`
)

func getFieldAclStr(cont *Context, fieldName string, mode string) string {
	str := ``
	if cont.IsMaster {
		return ` true `
	} else {
		if cont.IsLogin {
			str += ` (` + fieldName + `?'` + cont.User.Id + `' and ` + fieldName + `->'` + cont.User.Id + `'->>'` + mode + `'=true::text) or (not ` + fieldName + `?'` + cont.User.Id + `' and (`
			if len(cont.User.Role) != 0 {
				str += ` (not ` + fieldName + `?'` + strings.Join(cont.User.Role, `' and not `+fieldName+`?'`) + `' and ` + fieldName + `->'*'->>'` + mode + `'=true::text) `
				for _, va := range cont.User.Role {
					str += ` or ( ` + fieldName + `?'role:` + va + `' and ` + fieldName + `->'role;` + va + `'->>'` + mode + `'=true::text) `
				}
				str += `))`
			} else {
				str += fieldName + `->'*'->>'` + mode + `'=true::text)) `
			}
		} else {
			str += fieldName + `->'*'->>'` + mode + `'=true::text `
		}
		return str
	}
}

//刷新 应用信息缓存
func RefreshAppInfo() {
	App.Cache.Delete(`app:AppInfo`)
	appInfo := GetAppInfo()
	kodo.SetMac(appInfo[`QiniuAccessKey`], appInfo[`QiniuSecretKey`])
}

//重置MasterKey
func ResetMasterKey() {
	tx, err := App.Db.Begin()
	CheckErr(err)
	defer tx.Commit()
	r := GetRand()
	tx.Query(`UPDATE public.__config SET value='` + Krand(24, KC_RAND_KIND_ALL, r) + `' WHERE "key"='MasterKey'`)
	RefreshAppInfo()
}

func setConfigValue(key string, value string) {
	tx, err := App.Db.Begin()
	CheckErr(err)
	defer tx.Commit()
	tx.Query(`UPDATE public.__config SET value='` + SqlStrFilter(value) + `' WHERE "key"='` + SqlStrFilter(key) + `'`)
	RefreshAppInfo()
}

// ("b"."ACL"?'b' and "ACL"->>'b'='1') or (not "ACL"?'b' and "ACL"?'role' and ("ACL"->'stu'->>'read'='true' or "ACL"->'stu2'->>'read'='true')) or (not "ACL"?'b' and not "ACL"?'role' and "ACL"->'*'->>'read'='true');

func FieldName2SqlStr(str string) string {
	tmp := strings.Split(str, ".")
	last := tmp[len(tmp)-1]
	tmp = tmp[0 : len(tmp)-1]
	if len(tmp) != 0 {
		return `"` + strings.Join(tmp, ".") + `"."` + last + `"`
	} else {
		return `"` + last + `"`
	}
}
