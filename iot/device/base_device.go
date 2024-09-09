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

package device

import (
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/callback"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
)

// BaseDevice 设备基类， 所有设备均继承该类
type BaseDevice interface {
	Connect() bool
	DisConnect(timeout uint)
	IsConnected() bool

	SendMessage(message model.Message) bool
	ReportProperties(properties model.DeviceProperties) bool
	BatchReportSubDevicesProperties(service model.DevicesService) bool
	QueryDeviceShadow(query model.DevicePropertyQueryRequest, handler callback.DeviceShadowQueryResponseHandler)
	UploadFile(filename, filePath string) bool
	DownloadFile(filename, filePath string) bool
	ReportDeviceInfo(swVersion, fwVersion string)
	ReportLogs(logs []model.DeviceLogEntry) bool
	RequestTimeSync() bool
}
