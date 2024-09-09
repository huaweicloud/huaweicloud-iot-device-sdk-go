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

package gateway

import (
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/device"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
)

type MqttGatewayDevice struct {
	// 网关继承MqttDevice通用方法
	device.MqttDevice
}

func NewMqttGatewayDevice(authConfig *config.ConnectAuthConfig) *MqttGatewayDevice {
	mqttDevice := device.NewMqttDevice(authConfig)
	return &MqttGatewayDevice{
		MqttDevice: *mqttDevice,
	}
}

func (gatewayDevice *MqttGatewayDevice) Connect() bool {
	connect := gatewayDevice.MqttDevice.Connect()
	if !connect {
		return false
	}
	return true
}

func (gatewayDevice *MqttGatewayDevice) UpdateSubDeviceState(subDevicesStatus model.SubDevicesStatus) bool {
	glog.Infof("begin to update sub-devices status")
	subDeviceCounts := len(subDevicesStatus.DeviceStatuses)

	batchUpdateSubDeviceState := 0
	if subDeviceCounts%gatewayDevice.ConnectionAuthInfo.BatchSubDeviceSize == 0 {
		batchUpdateSubDeviceState = subDeviceCounts / gatewayDevice.ConnectionAuthInfo.BatchSubDeviceSize
	} else {
		batchUpdateSubDeviceState = subDeviceCounts/gatewayDevice.ConnectionAuthInfo.BatchSubDeviceSize + 1
	}

	for i := 0; i < batchUpdateSubDeviceState; i++ {
		begin := i * gatewayDevice.ConnectionAuthInfo.BatchSubDeviceSize
		end := (i + 1) * gatewayDevice.ConnectionAuthInfo.BatchSubDeviceSize
		if end > subDeviceCounts {
			end = subDeviceCounts
		}

		sds := model.SubDevicesStatus{
			DeviceStatuses: subDevicesStatus.DeviceStatuses[begin:end],
		}

		requestEventService := model.DataEntry{
			ServiceId: "$sub_device_manager",
			EventType: "sub_device_update_status",
			EventTime: iot.GetEventTimeStamp(),
			Paras:     sds,
		}

		request := model.Data{
			ObjectDeviceId: gatewayDevice.ConnectionAuthInfo.Id,
			Services:       []model.DataEntry{requestEventService},
		}
		if !gatewayDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, gatewayDevice.ConnectionAuthInfo.Id), gatewayDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
			glog.Warningf("gateway %s update sub devices status failed", gatewayDevice.ConnectionAuthInfo.Id)
			return false
		}
	}

	glog.Info("gateway  update sub devices status failed", gatewayDevice.ConnectionAuthInfo.Id)
	return true
}

func (gatewayDevice *MqttGatewayDevice) DeleteSubDevices(deviceIds []string) bool {
	glog.Infof("begin to delete sub-devices %s", deviceIds)

	subDevices := struct {
		Devices []string `json:"devices"`
	}{
		Devices: deviceIds,
	}

	requestEventService := model.DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "delete_sub_device_request",
		EventTime: iot.GetEventTimeStamp(),
		Paras:     subDevices,
	}

	request := model.Data{
		ObjectDeviceId: gatewayDevice.ConnectionAuthInfo.Id,
		Services:       []model.DataEntry{requestEventService},
	}

	if !gatewayDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, gatewayDevice.ConnectionAuthInfo.Id), gatewayDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
		glog.Warningf("gateway %s delete sub devices request send failed", gatewayDevice.ConnectionAuthInfo.Id)
		return false
	}

	glog.Warningf("gateway %s delete sub devices request send success", gatewayDevice.ConnectionAuthInfo.Id)
	return true
}

func (gatewayDevice *MqttGatewayDevice) AddSubDevices(deviceInfos []model.DeviceInfo) bool {
	devices := struct {
		Devices []model.DeviceInfo `json:"devices"`
	}{
		Devices: deviceInfos,
	}

	requestEventService := model.DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "add_sub_device_request",
		EventTime: iot.GetEventTimeStamp(),
		Paras:     devices,
	}

	request := model.Data{
		ObjectDeviceId: gatewayDevice.ConnectionAuthInfo.Id,
		Services:       []model.DataEntry{requestEventService},
	}

	if !gatewayDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, gatewayDevice.ConnectionAuthInfo.Id), gatewayDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
		glog.Warningf("gateway %s add sub devices request send failed", gatewayDevice.ConnectionAuthInfo.Id)
		return false
	}

	glog.Warningf("gateway %s add sub devices request send success", gatewayDevice.ConnectionAuthInfo.Id)
	return true
}

func (gatewayDevice *MqttGatewayDevice) SyncAllVersionSubDevices() {
	dataEntry := model.DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "sub_device_sync_request",
		EventTime: iot.GetEventTimeStamp(),
		Paras: struct {
		}{},
	}

	var dataEntries []model.DataEntry
	dataEntries = append(dataEntries, dataEntry)

	data := model.Data{
		Services: dataEntries,
	}

	if !gatewayDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, gatewayDevice.ConnectionAuthInfo.Id), gatewayDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(data)) {
		glog.Errorf("send sub device sync request failed")
	}
}

func (gatewayDevice *MqttGatewayDevice) SyncSubDevices(version int) {
	syncParas := struct {
		Version int `json:"version"`
	}{
		Version: version,
	}

	dataEntry := model.DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "sub_device_sync_request",
		EventTime: iot.GetEventTimeStamp(),
		Paras:     syncParas,
	}

	var dataEntries []model.DataEntry
	dataEntries = append(dataEntries, dataEntry)

	data := model.Data{
		Services: dataEntries,
	}

	if !gatewayDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, gatewayDevice.ConnectionAuthInfo.Id), gatewayDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(data)) {
		glog.Errorf("send sync sub device request failed")
	}
}
