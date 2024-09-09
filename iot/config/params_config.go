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

package config

import "github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/callback"

type DeviceParamsConfig struct {
	CommandHandler                   callback.CommandHandler
	MessageHandlers                  []callback.MessageHandler
	PropertiesSetHandlers            []callback.DevicePropertiesSetHandler
	PropertyQueryHandler             callback.DevicePropertyQueryHandler
	DeviceShadowQueryResponseHandler callback.DeviceShadowQueryResponseHandler
	SubDevicesAddHandler             callback.SubDevicesAddHandler
	SubDevicesDeleteHandler          callback.SubDevicesDeleteHandler
	SubDeviceStatusRespHandler       callback.SubDeviceStatusRespHandler
	SubDeviceAddResponseHandler      callback.SubDeviceAddResponseHandler
	SubDeviceDeleteResponseHandler   callback.SubDeviceDeleteResponseHandler
	SwFwVersionReporter              callback.SwFwVersionReporter
	DeviceUpgradeHandler             callback.DeviceUpgradeHandler
	DeviceStatusLogCollector         callback.DeviceStatusLogCollector
	DevicePropertyLogCollector       callback.DevicePropertyLogCollector
	DeviceMessageLogCollector        callback.DeviceMessageLogCollector
	DeviceCommandLogCollector        callback.DeviceCommandLogCollector
	FileUrls                         map[string]string
	Lcc                              *LogCollectionConfig
	ConnectionLostHandler            callback.ConnectionLostHandler
	ConnectHandler                   callback.DeviceConnectHandler
	SyncTimeResponseHandler          callback.SyncTimeResponseHandler
	RuleActionHandler                callback.RuleActionHandler
}

func (config *DeviceParamsConfig) AddCommandHandler(handler callback.CommandHandler) {
	if handler == nil {
		return
	}
	config.CommandHandler = handler
}

func (config *DeviceParamsConfig) AddMessageHandler(handler callback.MessageHandler) {
	if handler == nil {
		return
	}
	config.MessageHandlers = append(config.MessageHandlers, handler)
}

func (config *DeviceParamsConfig) AddPropertiesSetHandler(handler callback.DevicePropertiesSetHandler) {
	if handler == nil {
		return
	}
	config.PropertiesSetHandlers = append(config.PropertiesSetHandlers, handler)
}

func (config *DeviceParamsConfig) AddDevicePropertyQueryHandler(handler callback.DevicePropertyQueryHandler) {
	if handler == nil {
		return
	}
	config.PropertyQueryHandler = handler
}

func (config *DeviceParamsConfig) AddDeviceShadowQueryResponseHandler(handler callback.DeviceShadowQueryResponseHandler) {
	if handler == nil {
		return
	}
	config.DeviceShadowQueryResponseHandler = handler
}

func (config *DeviceParamsConfig) AddSubDevicesAddHandler(handler callback.SubDevicesAddHandler) {
	if handler == nil {
		return
	}
	config.SubDevicesAddHandler = handler
}

func (config *DeviceParamsConfig) AddSubDevicesDeleteHandler(handler callback.SubDevicesDeleteHandler) {
	if handler == nil {
		return
	}
	config.SubDevicesDeleteHandler = handler
}

func (config *DeviceParamsConfig) SetSwFwVersionReporter(handler callback.SwFwVersionReporter) {
	config.SwFwVersionReporter = handler
}

func (config *DeviceParamsConfig) SetDeviceUpgradeHandler(handler callback.DeviceUpgradeHandler) {
	config.DeviceUpgradeHandler = handler
}

func (config *DeviceParamsConfig) SetPropertyQueryHandler(handler callback.DevicePropertyQueryHandler) {
	config.PropertyQueryHandler = handler
}

func (config *DeviceParamsConfig) SetDeviceStatusLogCollector(collector callback.DeviceStatusLogCollector) {
	config.DeviceStatusLogCollector = collector
}

func (config *DeviceParamsConfig) SetDevicePropertyLogCollector(collector callback.DevicePropertyLogCollector) {
	config.DevicePropertyLogCollector = collector
}

func (config *DeviceParamsConfig) SetDeviceMessageLogCollector(collector callback.DeviceMessageLogCollector) {
	config.DeviceMessageLogCollector = collector
}

func (config *DeviceParamsConfig) SetDeviceCommandLogCollector(collector callback.DeviceCommandLogCollector) {
	config.DeviceCommandLogCollector = collector
}
