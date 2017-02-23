//util.go
//
//Copyright 2017-present Zhang. All Rights Reserved.
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
	"encoding/json"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
	"math/rand"
	"net/smtp"
	"strings"
	"crypto/md5"
	"encoding/hex"
)


func Json2map(r string) (s *map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(r), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Json2array(r string) (s *[]interface{}, err error) {
	var result []interface{}
	if err := json.Unmarshal([]byte(r), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Json2arraymap(r string) (s *[]map[string]interface{}, err error) {
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(r), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Map2json(s map[string]interface{})  string {
	result,err := json.Marshal(s)
	if  err != nil {
		return ""
	}
	return string(result)
}

func Array2json(s []interface{})  string {
	result,err := json.Marshal(s)
	if  err != nil {
		return ""
	}
	return string(result)
}

func ArrayMap2json(s []map[string]interface{})  string {
	result,err := json.Marshal(s)
	if  err != nil {
		return ""
	}
	return string(result)
}

func Json(c *gin.Context) (s *map[string]interface{}, err error) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	fmt.Println(string(buf[0:n]));
	fmt.Println(c.Query("where"))
	if req2map, err := Json2map(string(buf[0:n])); err == nil {
		return req2map,nil
	} else {
		return nil,err
	}
}

func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func Bool2string(bool bool)  string{
	if(bool){
		return "true"
	}else {
		return "false"
	}
}

func CheckErr(err error)  {
	if(err!=nil){
		panic(err)
	}
}

func CheckErrWithStr(err error,str string)  {
	if(err!=nil){
		panic(err)
	}else {
		fmt.Println(str)
	}
}

const (
	KC_RAND_KIND_NUM   = 0  // 纯数字
	KC_RAND_KIND_LOWER = 1  // 小写字母
	KC_RAND_KIND_UPPER = 2  // 大写字母
	KC_RAND_KIND_ALL   = 3  // 数字、大小写字母
)

func GetRand() *rand.Rand {
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	return r
}
// 随机字符串
func Krand(size int, kind int,r *rand.Rand) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	r=rand.New(rand.NewSource(time.Now().UnixNano()))
	for i :=0; i < size; i++ {
		if is_all { // random ikind
			ikind = r.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base+r.Intn(scope))
	}
	res := string(result[:])
	return res
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string){
	a_len := len(a)
	for i:=0; i < a_len; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue;
		}
		ret = append(ret, a[i])
	}
	return
}

func MergeStringArray(s ...[]string) (slice []string) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := MergeStringArray(s[1:]...)//...将数组元素打散
		slice = make([]string, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}
	return
}


func is_array(list []string,str string) bool {
	for _,va:=range list{
		if va==str{
			return true
		}
	}
	return false
}

func is_array_interface(list []interface{},str string) bool {
	for _,va:=range list{
		if va.(string)==str{
			return true
		}
	}
	return false
}

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func Md5Encrypt(password string,salt string) string {
	h := md5.New()
	h.Write([]byte(password+salt)) // 需要加密的字符串为 sharejs.com
	tmp1:=hex.EncodeToString(h.Sum(nil))
	h.Reset()
	h.Write([]byte(tmp1+salt))
	tmp2:=hex.EncodeToString(h.Sum(nil))
	h.Reset()
	h.Write([]byte(salt+tmp2))
	return hex.EncodeToString(h.Sum(nil))
}

func SqlStrFilter(str string) string{
	return strings.Replace(str,"'","''",-1)
}