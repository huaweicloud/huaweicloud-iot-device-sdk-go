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

package constants

import "crypto/tls"

const (
	// MessageDownTopic 平台下发消息topic
	MessageDownTopic string = "$oc/devices/{device_id}/sys/messages/down"

	// MessageUpTopic 设备上报消息topic
	MessageUpTopic string = "$oc/devices/{device_id}/sys/messages/up"

	// CommandDownTopic 平台下发命令topic
	CommandDownTopic string = "$oc/devices/{device_id}/sys/commands/#"

	// CommandResponseTopic 设备响应平台命令
	CommandResponseTopic string = "$oc/devices/{device_id}/sys/commands/response/request_id="

	// PropertiesUpTopic 设备上报属性
	PropertiesUpTopic string = "$oc/devices/{device_id}/sys/properties/report"

	// PropertiesSetRequestTopic 平台设置属性topic
	PropertiesSetRequestTopic string = "$oc/devices/{device_id}/sys/properties/set/#"

	// PropertiesSetResponseTopic 设备响应平台属性设置topic
	PropertiesSetResponseTopic string = "$oc/devices/{device_id}/sys/properties/set/response/request_id="

	// PropertiesQueryRequestTopic 平台查询设备属性
	PropertiesQueryRequestTopic string = "$oc/devices/{device_id}/sys/properties/get/#"

	// PropertiesQueryResponseTopic 设备响应平台属性查询
	PropertiesQueryResponseTopic string = "$oc/devices/{device_id}/sys/properties/get/response/request_id="

	// DeviceShadowQueryRequestTopic 设备侧获取平台的设备影子数据
	DeviceShadowQueryRequestTopic string = "$oc/devices/{device_id}/sys/shadow/get/request_id="

	// DeviceShadowQueryResponseTopic 设备侧响应获取平台设备影子
	DeviceShadowQueryResponseTopic string = "$oc/devices/{device_id}/sys/shadow/get/response/#"

	// GatewayBatchReportSubDeviceTopic 网关批量上报子设备属性
	GatewayBatchReportSubDeviceTopic string = "$oc/devices/{device_id}/sys/gateway/sub_devices/properties/report"

	// FileActionUpload 平台下发文件上传和下载URL
	FileActionUpload   string = "upload"
	FileActionDownload string = "download"

	// DeviceToPlatformTopic 设备或网关向平台发送请求
	DeviceToPlatformTopic string = "$oc/devices/{device_id}/sys/events/up"

	// PlatformEventToDeviceTopic 平台向设备下发事件topic
	PlatformEventToDeviceTopic string = "$oc/devices/{device_id}/sys/events/down"
)

const (
	AuthTypePassword uint8 = 0
	AuthTypeX509     uint8 = 1
)

var CipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
}
