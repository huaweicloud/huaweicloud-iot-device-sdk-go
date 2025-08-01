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
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/gateway"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"time"
)

func createGateway() *gateway.MqttGatewayDevice {
	authConfig := &config.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{MQTT_ACCESS_ADDRESS}:8883",
		Secret:       "your Secret",
		ServerCaPath: "iotda server ca path",
	}
	mqttDevice := gateway.NewMqttGatewayDevice(authConfig)
	connect := mqttDevice.Connect()
	if !connect {
		return nil
	}
	return mqttDevice
}

// 平台通知网关子设备新增
func platformNotifySubDeviceAdd() {
	gatewayDevice := createGateway()
	if gatewayDevice == nil {
		glog.Infof("create gateway device failed.")
		return
	}
	// 网关收到添加子设备的回调
	gatewayDevice.Client.SubDevicesAddHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device add version: %d", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("add sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	time.Sleep(5 * time.Minute)
}

// 平台通知网关子设备删除
func platformNotifySubDeviceDelete() {
	gatewayDevice := createGateway()
	if gatewayDevice == nil {
		glog.Infof("create gateway device failed.")
		return
	}
	// 网关收到删除子设备的回调
	gatewayDevice.Client.SubDevicesDeleteHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device delete version: %d", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("delete sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	time.Sleep(5 * time.Minute)
}

// 网关同步子设备列表
func syncSubDevices() {
	gatewayDevice := createGateway()
	if gatewayDevice == nil {
		glog.Infof("create gateway device failed.")
		return
	}
	//  同步新增子设备响应
	gatewayDevice.Client.SubDevicesAddHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device add version: %d", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("add sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	// 同步删除子设备响应
	gatewayDevice.Client.SubDevicesDeleteHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device delete version: %d", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("delete sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	gatewayDevice.SyncAllVersionSubDevices()
	time.Sleep(5 * time.Minute)
}

// 网关更新子设备状态
func updateSubDeviceStats() {
	gatewayDevice := createGateway()
	if gatewayDevice == nil {
		glog.Infof("create gateway device failed.")
		return
	}
	// 更新子设备状态请求响应
	gatewayDevice.Client.SubDeviceStatusRespHandler = func(response model.SubDeviceStatusResp) {
		if len(response.SuccessfulDevices) > 0 {
			glog.Infof("success update device status.")
			for _, sucessDevice := range response.SuccessfulDevices {
				glog.Infof("update device : %s status: %s", sucessDevice.DeviceId, sucessDevice.Status)
			}
		}
		if len(response.FailedDevices) > 0 {
			glog.Infof("failed to update device status")
			for _, failedDevice := range response.FailedDevices {
				glog.Infof("failed to update device status. deviceId : %s, errCode: %s, errMsg: %s",
					failedDevice.DeviceId, failedDevice.ErrorCode, failedDevice.ErrorMsg)
			}
		}
	}
	status := model.DeviceStatus{
		DeviceId: "sub device id",
		Status:   "ONLINE",
	}
	var statusInfos []model.DeviceStatus
	statusInfos = append(statusInfos, status)
	subDeviceStatus := model.SubDevicesStatus{
		DeviceStatuses: statusInfos,
	}
	gatewayDevice.UpdateSubDeviceState(subDeviceStatus)
	time.Sleep(5 * time.Minute)
}

// 网关添加子设备
func gatewayAddSubDevice() {
	gatewayDevice := createGateway()
	if gatewayDevice == nil {
		glog.Infof("create gateway device failed.")
		return
	}
	// 网关添加子设备请求响应
	gatewayDevice.Client.SubDeviceAddResponseHandler = func(response model.SubDeviceAddResponse) {
		glog.Infof("handle sub device add response")
		deviceList := response.SuccessFulDevices
		for _, deviceInfo := range deviceList {
			glog.Infof("success add sub device. deviceId : %s", deviceInfo.DeviceId)
		}
		failedDeviceList := response.FailedDevices
		for _, failDevice := range failedDeviceList {
			glog.Infof("failed add sub device. deviceId: %s, err: %s", failDevice.NodeId, failDevice.ErrorMsg)
		}
	}
	deviceInfo := model.DeviceInfo{
		DeviceId:       "sub device id",
		NodeId:         "subDevice1",
		Name:           "subDevice1",
		ParentDeviceId: "gateway device id",
		Description:    "subDevice1",
		ProductId:      "product id",
	}
	var infos []model.DeviceInfo
	infos = append(infos, deviceInfo)
	gatewayDevice.AddSubDevices(infos)
	time.Sleep(5 * time.Minute)
}

// 网关删除子设备
func gatewayDeleteDevice() {
	gatewayDevice := createGateway()
	if gatewayDevice == nil {
		glog.Infof("create gateway device failed.")
		return
	}
	// 网关删除子设备请求响应
	gatewayDevice.Client.SubDeviceDeleteResponseHandler = func(response model.SubDeviceDeleteResponse) {
		glog.Infof("handle sub device delete response")
		deviceList := response.SuccessFulDevices
		for _, deviceId := range deviceList {
			glog.Infof("success delete sub device. deviceId : %s", deviceId)
		}
		failedDeviceList := response.FailedDevices
		for _, failDevice := range failedDeviceList {
			glog.Infof("failed delete sub device. deviceId: %s, err: %s", failDevice.DeviceId, failDevice.ErrorMsg)
		}
	}
	// 这里演示网关主动删除子设备请求
	var deviceIds []string
	deviceIds = append(deviceIds, "sub device id")
	gatewayDevice.DeleteSubDevices(deviceIds)
	time.Sleep(5 * time.Minute)
}

func main() {
	gatewayDeleteDevice()
}
