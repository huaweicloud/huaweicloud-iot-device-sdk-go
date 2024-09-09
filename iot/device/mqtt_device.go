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
	"flag"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/client"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/file"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"github.com/panjf2000/ants/v2"
	uuid "github.com/satori/go.uuid"
	"time"
)

type MqttDevice struct {
	Client             client.MqttDeviceClient
	ConnectionAuthInfo *config.ConnectAuthConfig
}

func NewMqttDevice(authConfig *config.ConnectAuthConfig) *MqttDevice {
	// 初始化日志， 设置日志级别
	flag.Parse()
	err := flag.Set("stderrthreshold", "info")
	if err != nil {
		glog.Warningf("init log level failed.")
	}
	if authConfig == nil {
		glog.Warningf("connection auth config is nil.")
		return nil
	}
	if !checkAuthConfig(authConfig) {
		glog.Infof("auth config params was invalid.")
		return nil
	}
	pool, err := ants.NewPool(authConfig.ThreadNum)
	if err != nil {
		glog.Warningf("init go routing pool failed. err: %s", err.Error())
		return nil
	}
	device := &MqttDevice{
		ConnectionAuthInfo: authConfig,
		Client: client.MqttDeviceClient{
			ConnectAuthConfig: authConfig,
			Pool:              pool,
		},
	}
	if authConfig.MaxBufferMessage > 0 {
		device.Client.Queue = iot.NewCircularQueue(authConfig.MaxBufferMessage)
	}
	device.Client.FileUrls = make(map[string]string)
	return device
}

func checkAuthConfig(authConfig *config.ConnectAuthConfig) bool {
	if authConfig == nil {
		return false
	}
	if len(authConfig.Id) == 0 {
		glog.Warning("device id is empty.")
		return false
	}
	if authConfig.AuthType == 0 && len(authConfig.Password) == 0 {
		glog.Warning("password is empty when auth type is password.")
		return false
	}
	if len(authConfig.Servers) == 0 {
		glog.Warning("params server is empty.")
		return false
	}
	if authConfig.ConnectTimeOut == 0 {
		authConfig.ConnectTimeOut = 30 * time.Second
	}
	if authConfig.AutoReconnect == nil {
		var b = true
		authConfig.AutoReconnect = &b
	}
	if authConfig.BackOffTime <= 0 {
		authConfig.BackOffTime = 1000
	}
	if authConfig.MinBackOffTime <= 0 {
		authConfig.MinBackOffTime = 1000
	}
	if authConfig.MaxBackOffTime <= 0 {
		authConfig.MaxBackOffTime = 30 * 1000
	}
	if authConfig.ThreadNum <= 0 {
		authConfig.ThreadNum = 10
	}
	if authConfig.InflightMessages <= 0 {
		authConfig.InflightMessages = 20
	}
	if authConfig.BatchSubDeviceSize <= 0 {
		authConfig.BatchSubDeviceSize = 10
	}
	return true
}

func (mqttDevice *MqttDevice) Connect() bool {
	return mqttDevice.Client.Connect()
}

func (mqttDevice *MqttDevice) SendMessage(message model.Message) bool {
	topic := message.Topic
	if topic == "" {
		topic = constants.MessageUpTopic
	}
	return mqttDevice.Client.PublishMessage(iot.FormatTopic(topic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, message.Payload)
}

func (mqttDevice *MqttDevice) ReportProperties(properties model.DeviceProperties) bool {
	propertiesData := iot.Interface2JsonString(properties)
	result := mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.PropertiesUpTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, propertiesData)
	if result && mqttDevice.ConnectionAuthInfo.RuleEnable {
		mqttDevice.Client.RuleManageService.HandleRule(properties.Services, mqttDevice.Client.CreateRuleActionHandler())
	}
	return result
}

func (mqttDevice *MqttDevice) BatchReportSubDevicesProperties(service model.DevicesService) bool {
	subDeviceCounts := len(service.Devices)

	batchReportSubDeviceProperties := 0
	if subDeviceCounts%mqttDevice.ConnectionAuthInfo.BatchSubDeviceSize == 0 {
		batchReportSubDeviceProperties = subDeviceCounts / mqttDevice.ConnectionAuthInfo.BatchSubDeviceSize
	} else {
		batchReportSubDeviceProperties = subDeviceCounts/mqttDevice.ConnectionAuthInfo.BatchSubDeviceSize + 1
	}

	for i := 0; i < batchReportSubDeviceProperties; i++ {
		begin := i * mqttDevice.ConnectionAuthInfo.BatchSubDeviceSize
		end := (i + 1) * mqttDevice.ConnectionAuthInfo.BatchSubDeviceSize
		if end > subDeviceCounts {
			end = subDeviceCounts
		}

		sds := model.DevicesService{
			Devices: service.Devices[begin:end],
		}
		return mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.GatewayBatchReportSubDeviceTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(sds))
	}

	return true
}

func (mqttDevice *MqttDevice) QueryDeviceShadow(query model.DevicePropertyQueryRequest) {
	requestId := uuid.NewV4()
	message := mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceShadowQueryRequestTopic, mqttDevice.ConnectionAuthInfo.Id)+requestId.String(), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(query))
	if message {
		glog.Warningf("device %s query device shadow data failed,request id = %s", mqttDevice.ConnectionAuthInfo.Id, requestId)
	}
}

func (mqttDevice *MqttDevice) UploadFile(filename, filePath string) bool {
	request := mqttDevice.generateUploadFileRequest(filename, filePath)
	if request == nil {
		return false
	}
	if !mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(*request)) {
		glog.Warningf("publish file upload request url failed")
		return false
	}
	glog.Info("publish file upload request url success")

	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ticker:
			_, ok := mqttDevice.Client.FileUrls[filename+constants.FileActionUpload]
			if ok {
				glog.Infof("platform send file upload url success")
				goto BreakPoint
			}

		}
	}
BreakPoint:

	if len(mqttDevice.Client.FileUrls[filename+constants.FileActionUpload]) == 0 {
		glog.Errorf("get file upload url failed")
		return false
	}
	glog.Infof("file upload url is %s", mqttDevice.Client.FileUrls[filename+constants.FileActionUpload])

	uploadFlag := file.CreateHttpClient().UploadFile(filePath, mqttDevice.Client.FileUrls[filename+constants.FileActionUpload])
	if !uploadFlag {
		glog.Errorf("upload file failed")
		return false
	}

	response := file.CreateFileUploadDownLoadResultResponse(filename, constants.FileActionUpload, uploadFlag)

	if !mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.PlatformEventToDeviceTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(response)) {
		glog.Error("report file upload file result failed")
		return false
	}
	return true
}

func (mqttDevice *MqttDevice) generateUploadFileRequest(filename, filePath string) *model.FileRequest {
	length, fromFile, err := iot.Sha256FromFile(filePath)
	if err != nil {
		glog.Warningf("sha256 file err. %s", err.Error())
		return nil
	}
	fileAttributes := make(map[string]interface{})
	fileAttributes["hash_code"] = fromFile
	fileAttributes["size"] = length
	// 构造获取文件上传URL的请求
	requestParas := model.FileRequestServiceEventParas{
		FileName:       filename,
		FileAttributes: fileAttributes,
	}

	serviceEvent := model.FileRequestServiceEvent{
		Paras: requestParas,
	}
	serviceEvent.ServiceId = "$file_manager"
	serviceEvent.EventTime = iot.GetEventTimeStamp()
	serviceEvent.EventType = "get_upload_url"

	var services []model.FileRequestServiceEvent
	services = append(services, serviceEvent)
	request := &model.FileRequest{
		Services: services,
	}
	return request
}

func (mqttDevice *MqttDevice) DownloadFile(filename, filePath string) bool {
	// 构造获取文件上传URL的请求
	requestParas := model.FileRequestServiceEventParas{
		FileName: filename,
	}

	serviceEvent := model.FileRequestServiceEvent{
		Paras: requestParas,
	}
	serviceEvent.ServiceId = "$file_manager"
	serviceEvent.EventTime = iot.GetEventTimeStamp()
	serviceEvent.EventType = "get_download_url"

	var services []model.FileRequestServiceEvent
	services = append(services, serviceEvent)
	request := model.FileRequest{
		Services: services,
	}
	if !mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
		glog.Warningf("publish file download request url failed")
		return false
	}

	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ticker:
			_, ok := mqttDevice.Client.FileUrls[filename+constants.FileActionDownload]
			if ok {
				glog.Infof("platform send file upload url success")
				goto BreakPoint
			}

		}
	}
BreakPoint:

	if len(mqttDevice.Client.FileUrls[filename+constants.FileActionDownload]) == 0 {
		glog.Errorf("get file download url failed")
		return false
	}

	downloadFlag := file.CreateHttpClient().DownloadFile(filePath, mqttDevice.Client.FileUrls[filename+constants.FileActionDownload], "")
	if !downloadFlag {
		glog.Errorf("down load file { %s } failed", filename)
		return false
	}

	response := file.CreateFileUploadDownLoadResultResponse(filename, constants.FileActionDownload, downloadFlag)

	if !mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.PlatformEventToDeviceTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(response)) {
		glog.Error("report file download file result failed")
		return false
	}

	return true
}

func (mqttDevice *MqttDevice) ReportDeviceInfo(swVersion, fwVersion string) {
	event := model.ReportDeviceInfoServiceEvent{
		BaseServiceEvent: model.BaseServiceEvent{
			ServiceId: "$sdk_info",
			EventType: "sdk_info_report",
			EventTime: iot.GetEventTimeStamp(),
		},
		Paras: model.ReportDeviceInfoEventParas{
			DeviceSdkVersion: iot.SdkInfo()["sdk-version"],
			SwVersion:        swVersion,
			FwVersion:        fwVersion,
		},
	}

	request := model.ReportDeviceInfoRequest{
		ObjectDeviceId: mqttDevice.ConnectionAuthInfo.Id,
		Services:       []model.ReportDeviceInfoServiceEvent{event},
	}

	if !mqttDevice.Client.PublishMessage(iot.FormatTopic(constants.DeviceToPlatformTopic, mqttDevice.ConnectionAuthInfo.Id), mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
		glog.Error("report device info failed.")
	}
}

func (mqttDevice *MqttDevice) ReportLogs(logs []model.DeviceLogEntry) bool {
	var services []model.ReportDeviceLogServiceEvent

	for _, logEntry := range logs {
		service := model.ReportDeviceLogServiceEvent{
			BaseServiceEvent: model.BaseServiceEvent{
				ServiceId: "$log",
				EventType: "log_report",
				EventTime: iot.GetEventTimeStamp(),
			},
			Paras: logEntry,
		}

		services = append(services, service)
	}

	request := model.ReportDeviceLogRequest{
		Services: services,
	}

	glog.Infof("report log request is : %s", iot.Interface2JsonString(request))

	topic := iot.FormatTopic(constants.DeviceToPlatformTopic, mqttDevice.ConnectionAuthInfo.Id)

	mqttDevice.Client.PublishMessage(topic, mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request))

	if !mqttDevice.Client.PublishMessage(topic, mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
		glog.Errorf("device %s report log failed", mqttDevice.ConnectionAuthInfo.Id)
		return false
	}
	return true

}

func (mqttDevice *MqttDevice) RequestTimeSync() bool {
	var services []model.TimeSyncRequestServiceEvent

	timestamp := time.Now().Unix()
	para := model.TimeSyncRequestServiceEventParas{
		DeviceSendTime: timestamp,
	}
	timeEvent := model.TimeSyncRequestServiceEvent{
		Paras: para,
	}
	timeEvent.ServiceId = "$time_sync"
	timeEvent.EventType = "time_sync_request"
	timeEvent.EventTime = iot.GetEventTimeStamp()

	services = append(services, timeEvent)
	request := model.TimeSyncRequest{
		Services: services,
	}

	topic := iot.FormatTopic(constants.DeviceToPlatformTopic, mqttDevice.ConnectionAuthInfo.Id)

	if !mqttDevice.Client.PublishMessage(topic, mqttDevice.ConnectionAuthInfo.Qos, iot.Interface2JsonString(request)) {
		glog.Errorf("device %s report log failed", mqttDevice.ConnectionAuthInfo.Id)
		return false
	}
	return true
}

func (mqttDevice *MqttDevice) DisConnect(timeout uint) {
	mqttDevice.Client.Close(timeout)
}

func (mqttDevice *MqttDevice) IsConnected() bool {
	return mqttDevice.Client.IsConnect()
}
