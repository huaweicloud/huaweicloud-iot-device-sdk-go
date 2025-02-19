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
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/file"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"os"
	"os/signal"
)

// ota升级
func otaUpgrade() {
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

	device.Client.SwFwVersionReporter = func() (string, string) {
		return "v1.0", "v1.0"
	}

	// upgradeType 0：软件升级  1：固件升级  2：OBS软件升级  3：OBS固件升级
	device.Client.DeviceUpgradeHandler = func(upgradeType byte, info model.UpgradeInfo) model.UpgradeProgress {
		glog.Infof("begin to handle upgrade process")
		upgradeProcess := model.UpgradeProgress{}
		currentPath, err := os.Getwd()
		if err != nil {
			glog.Warningf("get executable path failed. err: %s", err.Error())
			upgradeProcess.ResultCode = 255
			upgradeProcess.Description = "get executable path failed."
			return upgradeProcess
		}
		currentPath = currentPath + "\\download\\ota.txt"
		downloadFlag := file.CreateHttpClient().OTADownloadFile(upgradeType, currentPath, info.Url, info.AccessToken)
		if !downloadFlag {
			glog.Errorf("down load file { %s } failed", currentPath)
			upgradeProcess.ResultCode = 10
			upgradeProcess.Description = "down load ota package failed."
			return upgradeProcess
		}
		glog.Infof("download file success.")
		// checkPackage()  校验下载的升级包
		// installPackage()  安装升级包
		upgradeProcess.ResultCode = 0
		upgradeProcess.Version = info.Version
		upgradeProcess.Progress = 100
		return upgradeProcess
	}
	connect := device.Connect()
	glog.Infof("connect result : %v", connect)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		<-interrupt
		break
	}

}

func main() {
	otaUpgrade()
}
