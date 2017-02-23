//file.go
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
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"qiniupkg.com/api.v7/kodo"
)

type IbFile struct {
}

var ibfile = IbFile{}

func (f *IbFile) Init(ak string, sk string) {
	kodo.SetMac(ak, sk)
}

func (f *IbFile) GetUploadQuery(cont *Context) *Query {
	q := &Query{}
	q.ClassName = "_File"
	q.Method = "create"
	q.Data = &map[string]interface{}{}
	q.Context = cont
	return q
}

func (f *IbFile) UploadFile(cont *Context) gin.H {

	appInfo := GetAppInfo()

	callbackurl := If(App.IsSSL, "https://", "http://").(string) + appInfo[`AppHost`] + "/callback"
	// 设置CallbackBody字段
	callbackbody := `{"key":$(key), "hash":$(etag),"filesize":$(fsize),"bucket":$(bucket),"mimeType":$(mimeType),"fname";$(fname),"masterKey":"` + appInfo[`MasterKey`] + `","owner":` + If(cont.IsLogin, cont.User.Id, `-1`).(string) + `}`

	// 创建一个Client
	c := kodo.New(0, nil)
	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: appInfo[`QiniuBucket`],
		//设置Token过期时间
		Expires:          3600,
		CallbackUrl:      callbackurl,
		CallbackBody:     callbackbody,
		CallbackBodyType: `application/json`,
	}
	// 生成一个上传token
	token := c.MakeUptoken(policy)

	return gin.H{"token": token}
}

func (f *IbFile) UploadCallback(cont *Context) (gin.H, error) {
	data, err := Json(cont.GinContext)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
	}
	appInfo := GetAppInfo()
	if appInfo[`MasterKey`] == (*data)[`masterKey`] {
		q := Query{}
		q.ClassName = "_File"
		q.Data = &map[string]interface{}{}
		(*q.Data)[`mime_type`] = (*data)[`mimeType`].(string)
		(*q.Data)[`key`] = (*data)[`key`].(string)
		(*q.Data)[`name`] = (*data)[`fname`].(string)
		(*q.Data)[`url`] = appInfo[`QiniuUrl`] + (*data)[`key`].(string)
		(*q.Data)[`bucket`] = (*data)[`bucket`].(string)
		(*q.Data)[`metaData`] = gin.H{"owner": (*data)[`owner`].(float64), "size": (*data)[`filesize`].(float64), "mime_type": (*data)[`mimeType`].(string)}
		res, err := q.Create()
		if err != nil {
			return nil, NewError(http.StatusBadRequest, 100, errors.New("参数解析错误"))
		}
		return gin.H{"success": true, "id": (*res)[`id`], "url": appInfo[`QiniuUrl`] + (*data)[`key`].(string)}, nil
	} else {
		return nil, NewError(http.StatusBadRequest, 100, errors.New("非法访问"))
	}
}

func (f *IbFile) GetDeleteQuery(cont *Context) *Query {
	q := &Query{}
	q.ClassName = "_File"
	q.Id = cont.GinContext.Param("id")
	q.Method = "delete"
	q.Context = cont
	return q
}

func (f *IbFile) DeleteFile(cont *Context) (gin.H, error) {
	id := cont.GinContext.Param("id")
	var err error
	tx, err := App.Db.Begin()
	if err != nil {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("服务器内部错误"))
	}
	defer tx.Commit()

	tmp, err := tx.Query(`DELETE FROM public."_File" WHERE "id"=` + id + ` Returning "key";`)
	res := GetJSON(tmp)
	if len(*res) == 0 {
		return nil, NewError(http.StatusBadRequest, 300, errors.New("文件不存在"))
	}
	// new一个Bucket管理对象
	c := kodo.New(0, nil)
	appInfo := GetAppInfo()
	p := c.Bucket(appInfo[`QiniuBucket`])
	// 调用Delete方法删除文件
	err = p.Delete(nil, (*res)[0][`key`].(string))
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, 300, errors.New("删除失败"))
	}
	return gin.H{"success": true}, nil
}
