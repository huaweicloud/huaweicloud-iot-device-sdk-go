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
	device2 "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/device"
	"os"
	"time"
)

func main() {
	// 创建一个设备并初始化
	authConfig := &config2.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Secret:       "your Secret",
		ServerCaPath: "iotda server ca path",
	}
	device := device2.NewMqttDevice(authConfig)
	if device == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	connect := device.Connect()
	glog.Infof("connect result : %v", connect)
	currentPath, err := os.Getwd()
	if err != nil {
		glog.Warningf("get executable path failed. err: %s", err.Error())
		return
	}

	fileName := "test_upload.txt"
	uploadFilePath := currentPath + "\\download\\test_upload.txt"
	device.UploadFile(fileName, uploadFilePath)
	time.Sleep(10 * time.Second)
	downloadFilePath := currentPath + "\\download\\test_download.txt"
	device.DownloadFile(fileName, downloadFilePath)
}
