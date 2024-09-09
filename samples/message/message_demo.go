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
	config2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	device2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/device"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"os"
	"os/signal"
	"time"
)

// 使用平台默认topic接收消息下发的消息
func messageDeliveryDefault() {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Password:     "your password",
		ServerCaPath: "iotda server ca path",
	}
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	device.Connect()

	// 注册平台下发消息的callback，当收到平台下发的消息时，调用此callback.
	// 支持注册多个callback，并且按照注册顺序调用
	device.Client.AddMessageHandler(func(message string) bool {
		glog.Infof("first callback called : %s", message)
		return true
	})
	device.Client.AddMessageHandler(func(message string) bool {
		glog.Infof("second callback called : %s", message)
		return true
	})

	// 向平台发送消息
	message := model.Message{
		Payload: "first message",
	}

	sendMsgResult := device.SendMessage(message)
	glog.Infof("send message %v", sendMsgResult)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		<-interrupt
		break
	}
}

// 使用平台自定义topic接收消息下发的消息
func messageDeliveryCustomize() {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Password:     "your password",
		ServerCaPath: "iotda server ca path",
	}
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	// 可以在此处实现与平台建链后的一些自定义逻辑，比如建链后订阅一些自定义的topic
	device.Client.ConnectHandler = func(client mqtt.Client) {
		glog.Infof("connect to server success.")
		// 使用自定义topic接收平台下发的消息
		device.Client.SubscribeCustomizeTopic("$oc/devices/668619463030681a7de586d6_dsdasaczxz/user/testabsda", func(message string) bool {
			glog.Infof("first callback called: %s", message)
			return true
		})
	}
	device.Connect()

	// 向平台发送消息
	message := model.Message{
		Topic:   "$oc/devices/668619463030681a7de586d6_dsdasaczxz/user/testabsda",
		Payload: "first message",
	}

	for i := 0; i < 100; i++ {
		sendMsgResult := device.SendMessage(message)
		glog.Infof("send message %v", sendMsgResult)
		time.Sleep(5 * time.Second)
	}
}

// 使用平台策略topic接收消息下发的消息
func messageDeliveryPolicy() {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Password:     "your password",
		ServerCaPath: "iotda server ca path",
	}
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	// 可以在此处实现与平台建链后的一些自定义逻辑，比如建链后订阅一些自定义的topic
	device.Client.ConnectHandler = func(client mqtt.Client) {
		// 使用自定义topic接收平台下发的消息
		device.Client.SubscribeCustomizeTopic("testdevicetopic", func(message string) bool {
			glog.Infof("first callback called %s", message)
			return true
		})
	}

	device.Connect()
	// 向平台发送消息
	message := model.Message{
		Topic:   "testdevicetopic",
		Payload: "first message",
	}

	for i := 0; i < 100; i++ {
		sendMsgResult := device.SendMessage(message)
		glog.Infof("send message %v", sendMsgResult)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	messageDeliveryDefault()
}
