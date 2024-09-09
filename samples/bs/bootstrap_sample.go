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
	config2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	device2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/device"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/samples/test_util"
	"time"
)

/*
  演示在发放平台创建设备后，通过秘钥认证设备进行发放，通过引导服务获取真实的服务器地址并接入平台
  使用静态策略，关键字来源为属性上报，上报结构体 “baseStrategyKeyword” 包含设置的策略关键字
*/
func bootstrapSecret() {
	// 发放平台注册的设备ID
	deviceId := "your device id"
	// 设备秘钥
	pwd := "your pwd"

	authConfig := config2.ConnectAuthConfig{
		Id:             deviceId,
		Password:       pwd,
		Servers:        "mqtts://{mqtt access ip}:8883",
		UseBootstrap:   true,
		BsServerCaPath: "bs server ca cert path",
		ServerCaPath:   "iotda server ca cert path",
		BootStrapBody: &model.BootStrapProperties{
			BaseStrategyKeyword: "1111",
		},
	}
	device := device2.NewMqttDevice(&authConfig)
	if device == nil {
		glog.Warningf("create device failed.")
		return
	}
	initRes := device.Connect()
	glog.Infof("connect result : %v", initRes)
	time.Sleep(3 * time.Second)
	// 上报消息
	device.SendMessage(test_util.GenerateMessage())
}

/*
  演示在发放平台创建设备后，通过证书认证设备进行发放，通过引导服务获取真实的服务器地址并接入平台
  使用静态策略，关键字来源为属性上报，上报结构体 “baseStrategyKeyword” 包含设置的策略关键字
*/
func bootstrapCert() {
	// 发放平台注册的设备ID
	deviceId := "your device id"
	authConfig := &config2.ConnectAuthConfig{
		Id:              deviceId,
		Servers:         "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		UseBootstrap:    true,
		AuthType:        constants.AuthTypeX509,
		BsServerCaPath:  "bs server ca cert path",
		ServerCaPath:    "iotda server ca cert path",
		CertFilePath:    "device cert path",
		CertKeyFilePath: "device cert key path",
		BootStrapBody: &model.BootStrapProperties{
			BaseStrategyKeyword: "1111",
		},
	}
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	initRes := device.Connect()
	glog.Infof("connect result : %v", initRes)
	time.Sleep(3 * time.Second)
	// 上报消息
	device.SendMessage(test_util.GenerateMessage())
}

/*
  演示使用设备组秘钥认证方式通过静态策略进行发放设备，通过引导服务获取真实的服务器地址并接入平台，设备组发放时无需在发放平台注册设备，静态策略为数据上报。
  若将设备注册到指定产品则设备ID格式为{product_Id}_xxx, 以productId开头加上下划线后拼接设备id，且仅能存在一个下划线。若存在多个下划线或
  没有下划线，则默认生成一个产品
*/
func bootstrapScopeIdSecretStaticPolicy() {
	// 自定义设备id
	deviceId := "your device id"
	// 注册组秘钥
	pwd := "your password"

	authInfo := &config2.ConnectAuthConfig{
		Id:             deviceId,
		Password:       pwd,
		Servers:        "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		UseBootstrap:   true,
		BsServerCaPath: "bs server ca cert path",
		ServerCaPath:   "iotda server ca cert path",
		ScopeId:        "{scopeId}",
		BootStrapBody: &model.BootStrapProperties{
			BaseStrategyKeyword: "1111",
		},
	}
	device := device2.NewMqttDevice(authInfo)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	initRes := device.Connect()
	glog.Infof("connect result : %v", initRes)

	time.Sleep(3 * time.Second)
	// 上报消息
	device.SendMessage(test_util.GenerateMessage())
}

/*
  演示使用设备组证书认证方式通过静态策略进行发放设备，通过引导服务获取真实的服务器地址并接入平台，设备组发放时无需在发放平台注册设备，静态策略为数据上报。
  若将设备注册到指定产品则设备ID格式为{product_Id}_xxx, 以productId开头加上下划线后拼接设备id，且仅能存在一个下划线。若存在多个下划线或
  没有下划线，则默认生成一个产品
*/
func bootstrapScopeIdCertStaticPolicy() {
	// 自定义设备id
	deviceId := "your device id"
	authConfig := &config2.ConnectAuthConfig{
		Id:              deviceId,
		Servers:         "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		UseBootstrap:    true,
		AuthType:        constants.AuthTypeX509,
		BsServerCaPath:  "bs server ca cert path",
		ServerCaPath:    "iotda server ca cert path",
		CertFilePath:    "device cert path",
		CertKeyFilePath: "device cert key path",
		ScopeId:         "{scopeId}",
		BootStrapBody: &model.BootStrapProperties{
			BaseStrategyKeyword: "1111",
		},
	}
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	initRes := device.Connect()
	glog.Infof("connect result : %v", initRes)

	time.Sleep(3 * time.Second)
	// 上报消息
	device.SendMessage(test_util.GenerateMessage())
}

func main() {
	bootstrapSecret()
}
