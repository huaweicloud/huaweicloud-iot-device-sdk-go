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

package callback

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
)

// SubDevicesAddHandler 子设备添加回调函数
type SubDevicesAddHandler func(devices model.SubDeviceInfo)

// SubDevicesDeleteHandler 子设备删除糊掉函数
type SubDevicesDeleteHandler func(devices model.SubDeviceInfo)

type SubDeviceStatusRespHandler func(response model.SubDeviceStatusResp)

type SubDeviceAddResponseHandler func(response model.SubDeviceAddResponse)

type SubDeviceDeleteResponseHandler func(response model.SubDeviceDeleteResponse)

// CommandHandler 处理平台下发的命令
type CommandHandler func(model.Command) (bool, interface{})

// MessageHandler 设备消息
type MessageHandler func(message string) bool

// DevicePropertiesSetHandler 平台设置设备属性
type DevicePropertiesSetHandler func(message model.DevicePropertyDownRequest) bool

// DevicePropertyQueryHandler 平台查询设备属性
type DevicePropertyQueryHandler func(query model.DevicePropertyQueryRequest) model.DevicePropertyEntry

// DeviceUpgradeHandler 设备执行软件/固件升级.upgradeType = 0 软件升级，upgradeType = 1 固件升级
type DeviceUpgradeHandler func(upgradeType byte, info model.UpgradeInfo) model.UpgradeProgress

// SwFwVersionReporter 设备上报软固件版本,第一个返回值为软件版本，第二个返回值为固件版本
type SwFwVersionReporter func() (string, string)

// DeviceShadowQueryResponseHandler 设备获取设备影子数据
type DeviceShadowQueryResponseHandler func(response model.DeviceShadowQueryResponse)

// DeviceStatusLogCollector 设备状态日志收集器
type DeviceStatusLogCollector func(endTime string) []model.DeviceLogEntry

// DevicePropertyLogCollector 设备属性日志收集器
type DevicePropertyLogCollector func(endTime string) []model.DeviceLogEntry

// DeviceMessageLogCollector 设备消息日志收集器
type DeviceMessageLogCollector func(endTime string) []model.DeviceLogEntry

// DeviceCommandLogCollector 设备命令日志收集器
type DeviceCommandLogCollector func(endTime string) []model.DeviceLogEntry

// ConnectionLostHandler 设备断链回调函数
type ConnectionLostHandler func(client mqtt.Client, reason error)

// DeviceConnectHandler 设备建链回调函数
type DeviceConnectHandler func(client mqtt.Client)

// SyncTimeResponseHandler 时间同步响应回调
type SyncTimeResponseHandler func(deviceSendTime, serverRecvTime, serverSendTime int64)
