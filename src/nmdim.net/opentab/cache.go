//cache.go
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

import "time"

func QueryWithCache(query string,t time.Duration) interface{} {
	tmp, found := App.Cache.Get(`db:`+query)
	if found {
		return tmp
	}else {
		tx,err:=App.Db.Begin()
		CheckErr(err)
		defer tx.Commit()

		r,_:=tx.Query(query)

		res :=GetJSON(r)

		App.Cache.Set(`db:`+query,res,t)

		return res
	}
}





func RefreshCacheValue(key string)  {
	App.Cache.Delete(key)
}

