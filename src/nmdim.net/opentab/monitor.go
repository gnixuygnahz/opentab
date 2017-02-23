//monitor.go
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
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
)

type IbMonitor struct {
}

var ibmonitor IbMaster

func (i *IbMaster) GetSystemInfo(c *Context) (gin.H, error) {
	if c.IsMaster {
		v, _ := mem.VirtualMemory()
		c, _ := cpu.Info()
		d, _ := disk.Usage("/")
		n, _ := host.Info()

		res := gin.H{}
		res[`Men`] = gin.H{"Total": v.Total, "Available": v.Available, "Used": v.Used, "UsedPercent": v.UsedPercent}
		res1, _ := cpu.Percent(0, false)

		res2:=[]gin.H{}
			for _, sub_cpu := range c {
				modelname := sub_cpu.ModelName
				cores := sub_cpu.Cores
				res2=append(res2,gin.H{"ModelName":modelname,"Cores":cores})
			}


		res[`CPU`] = gin.H{"Percent":res1[0],"Info":res2}
		res[`HD`]=gin.H{"Total":d.Total,"Used":d.Used,"Free":d.Free}
		res[`OS`]=gin.H{"OS":n.OS,"Platform":n.Platform,"PlatformVersion":n.PlatformVersion,"Hostname":n.Hostname}
		return res, nil
	} else {
		return nil, NewError(http.StatusBadRequest, 201, errors.New("没有权限"))
	}
}

