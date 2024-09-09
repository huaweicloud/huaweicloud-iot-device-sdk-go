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
	"time"
)

func createMqttDevice() *device2.MqttDevice {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Password:     "your password",
		ServerCaPath: "iotda server ca path",
	}
	device := device2.NewMqttDevice(authConfig)
	// 注册平台设置属性callback,当应用通过API设置设备属性时，会调用此callback，支持注册多个callback
	device.Client.AddPropertiesSetHandler(func(propertiesSetRequest model.DevicePropertyDownRequest) bool {
		glog.Infof("I get property set command")
		glog.Infof("request is %s", iot.Interface2JsonString(propertiesSetRequest))
		return true
	})

	// 注册平台查询设备属性callback，当平台查询设备属性时此callback被调用，仅支持设置一个callback
	device.Client.SetPropertyQueryHandler(func(query model.DevicePropertyQueryRequest) model.DevicePropertyEntry {
		return model.DevicePropertyEntry{
			ServiceId: "smokeDetector",
			Properties: test_model.DemoProperties{
				Temperature: 27,
			},
			EventTime: "2024-05-28 14:23:24",
		}
	})

	device.Client.AddDeviceShadowQueryResponseHandler(func(response model.DeviceShadowQueryResponse) {
		shadow := response.Shadow
		glog.Infof("receive shadow msg from device.")
		glog.Infof("on_shadow_get device_id:  %s", response.ObjectDeviceId)
		glog.Infof("shadow service_id: %s", shadow[0].ServiceId)
		glog.Infof("shadow desired properties: %v", shadow[0].Desired)
		glog.Infof("shadow reported: %v", shadow[0].Reported)
	})
	return device
}

func createDeviceProperty() []model.DevicePropertyEntry {
	// 设备上报属性
	props := model.DevicePropertyEntry{
		ServiceId: "smokeDetector",
		EventTime: iot.GetEventTimeStamp(),
		Properties: test_model.DemoProperties{
			Temperature: 28,
		},
	}

	var content []model.DevicePropertyEntry
	return append(content, props)
}

func main() {
	device := createMqttDevice()
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	device.Connect()
	glog.Infof("device connected: %v\n", device.IsConnected())

	content := createDeviceProperty()

	services := model.DeviceProperties{
		Services: content,
	}
	device.ReportProperties(services)

	// 设备查询设备影子数据
	device.QueryDeviceShadow(model.DevicePropertyQueryRequest{
		ServiceId: "smokeDetector",
	})

	// 批量上报子设备属性
	subDevice1 := model.DeviceService{
		DeviceId: "668619463030681a7de586d6_sub-device-1",
		Services: content,
	}
	subDevice2 := model.DeviceService{
		DeviceId: "668619463030681a7de586d6_sub-device-2",
		Services: content,
	}

	subDevice3 := model.DeviceService{
		DeviceId: "668619463030681a7de586d6_sub-device-3",
		Services: content,
	}

	var devices []model.DeviceService
	devices = append(devices, subDevice1, subDevice2, subDevice3)

	device.BatchReportSubDevicesProperties(model.DevicesService{
		Devices: devices,
	})
	time.Sleep(1 * time.Minute)
}
