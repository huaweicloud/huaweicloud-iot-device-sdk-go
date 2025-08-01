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
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/device"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/samples/test_util"
	"os"
	"os/signal"
)

// 使用秘钥接入
func connectWithSecret() {
	authConfig := &config.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Secret:       "your Secret",
		ServerCaPath: "iotda server ca path",
	}
	mqttDevice := device.NewMqttDevice(authConfig)
	if mqttDevice == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	initResult := mqttDevice.Connect()
	glog.Info("connect result is : ", initResult)
	// 上报消息
	mqttDevice.SendMessage(test_util.GenerateMessage())
}

// 使用证书接入
func connectWithX509Cert() {
	authConfig := &config.ConnectAuthConfig{
		Id:              "your device id",
		Servers:         "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		AuthType:        constants.AuthTypeX509,
		ServerCaPath:    "iotda server ca path",
		CertFilePath:    "device cert path",
		CertKeyFilePath: "device key path",
	}
	mqttDevice := device.NewMqttDevice(authConfig)
	if mqttDevice == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	initResult := mqttDevice.Connect()
	glog.Info("connect result is : ", initResult)
	// 上报消息
	mqttDevice.SendMessage(test_util.GenerateMessage())
}

// 关闭自动重连，通过回调函数实现自定义重连。
func connectWithRetry() {
	var autoReconnect = true
	authConfig := &config.ConnectAuthConfig{
		Id:            "your device id",
		Servers:       "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Secret:        "your Secret",
		AutoReconnect: &autoReconnect,
		ServerCaPath:  "iotda server ca path",
	}
	authConfig.MaxBackOffTime = 30 * 1000
	authConfig.MinBackOffTime = 1000
	authConfig.MaxBufferMessage = 100
	mqttDevice := device.NewMqttDevice(authConfig)
	if mqttDevice == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	// 关闭自动重连后， 可以在此处回调函数内实现自定义断链重连功能
	mqttDevice.Client.ConnectionLostHandler = func(client mqtt.Client, reason error) {
		glog.Warningf("connect lost from server. you can customize auto reconnect logic here")
	}
	// 可以在此处实现与平台建链后的一些自定义逻辑
	mqttDevice.Client.ConnectHandler = func(client mqtt.Client) {
		glog.Infof("connect to server success.")
	}
	initResult := mqttDevice.Connect()
	glog.Info("connect result is : ", initResult)
	// 上报消息
	mqttDevice.SendMessage(test_util.GenerateMessage())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		<-interrupt
		break
	}
}

func main() {
	connectWithRetry()
}
