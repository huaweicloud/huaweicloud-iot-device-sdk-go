// Copyright (c) 2023-2024 Huawei Cloud Computing Technology Co., Ltd. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of
//    conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used
//    to endorse or promote products derived from this software without specific prior written
//    permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
// EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
// OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
// ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	config2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	device2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/device"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/samples/test_model"
	"os"
	"os/signal"
	"time"
)

func ruleManage() {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Password:     "your password",
		ServerCaPath: "iotda server ca path",
	}
	authConfig.RuleEnable = true
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	device.Client.CommandHandler = func(command model.Command) (bool, interface{}) {
		glog.Infof("command device id is %s", command.ObjectDeviceId)
		glog.Infof("command name is %s", command.CommandName)
		glog.Infof("command serviceId is %s", command.ServiceId)
		glog.Infof("command params is %v", command.Paras)
		return true, map[string]interface{}{
			"cost_time": 12,
		}
	}

	connect := device.Connect()
	glog.Infof("connect result : %v", connect)
	// 上报SDK版本
	device.ReportDeviceInfo("", "")
	time.Sleep(3 * time.Second)

	// 上报属性
	props := model.DevicePropertyEntry{
		ServiceId: "smokeDetector",
		EventTime: iot.GetEventTimeStamp(),
		Properties: test_model.DemoProperties{
			Temperature: 20,
		},
	}

	var content []model.DevicePropertyEntry
	content = append(content, props)
	services := model.DeviceProperties{
		Services: content,
	}
	device.ReportProperties(services)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		<-interrupt
		break
	}
}

func customRuleManage() {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Password:     "your password",
		ServerCaPath: "iotda server ca path",
	}
	authConfig.RuleEnable = true
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	device.Client.RuleActionHandler = func(actions []model.Action) bool {
		for _, action := range actions {
			glog.Infof("action deviceId: %s:", action.DeviceId)
			command := action.Command
			glog.Infof("action command name : %s", command.CommandName)
			glog.Infof("action command body : %v", command.CommandBody)
		}
		return true
	}

	connect := device.Connect()
	glog.Infof("connect result : %v", connect)
	// 上报SDK版本
	device.ReportDeviceInfo("", "")
	time.Sleep(3 * time.Second)

	// 上报属性
	props := model.DevicePropertyEntry{
		ServiceId: "smokeDetector",
		EventTime: iot.GetEventTimeStamp(),
		Properties: test_model.DemoProperties{
			Temperature: 20,
		},
	}

	var content []model.DevicePropertyEntry
	content = append(content, props)
	services := model.DeviceProperties{
		Services: content,
	}
	device.ReportProperties(services)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		<-interrupt
		break
	}
}

func main() {
	ruleManage()
}
