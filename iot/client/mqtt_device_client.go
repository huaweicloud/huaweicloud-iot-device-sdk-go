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

package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/callback"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/rule"
	"github.com/panjf2000/ants/v2"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

type MqttDeviceClient struct {
	config.DeviceParamsConfig
	client            mqtt.Client
	ConnectAuthConfig *config.ConnectAuthConfig
	RuleManageService rule.RuleManageService
	Pool              *ants.Pool
	Queue             *iot.CircularQueue
	retryTimes        int64
	calculate         int64
}

func (mqttClient *MqttDeviceClient) Connect() bool {
	// 退避重试。默认最大退避时间30s
	minBackoffTime := mqttClient.ConnectAuthConfig.MinBackOffTime
	maxBackoffTime := mqttClient.ConnectAuthConfig.MaxBackOffTime
	backOffTime := mqttClient.ConnectAuthConfig.BackOffTime
	// 建链失败进行重试
	internal := mqttClient.connectMqttBroker()
	// 只有开启了自动重连，sdk还会进行自动重连
	for !internal && *mqttClient.ConnectAuthConfig.AutoReconnect {
		lowBound := int64(float64(backOffTime) * 0.8)
		highBound := int64(float64(backOffTime) * 1.0)
		rand.Seed(time.Now().Unix())
		randomBackoff := rand.Int63n(highBound - lowBound)
		// 防止幂次方计算出现超大值
		if mqttClient.calculate > 20 {
			mqttClient.calculate = 0
		}
		pow := math.Pow(2, float64(mqttClient.calculate))
		backoffWithJitter := int64(pow) * (randomBackoff + lowBound)

		var waitTimeMs int64
		if minBackoffTime+backoffWithJitter > maxBackoffTime {
			waitTimeMs = maxBackoffTime
		} else {
			waitTimeMs = minBackoffTime + backoffWithJitter
		}
		glog.Warningf("client will retry to reconnect after %d ms", waitTimeMs)
		time.Sleep(time.Duration(waitTimeMs) * time.Millisecond)
		mqttClient.retryTimes++
		mqttClient.calculate++
		internal = mqttClient.connectMqttBroker()
		if !internal {
			glog.Warningf("connect mqtt go broker retry. times: %d", mqttClient.retryTimes)
		}
	}
	mqttClient.retryTimes = 0
	return internal
}

func (mqttClient *MqttDeviceClient) configureTLS(options *mqtt.ClientOptions) error {
	if !strings.ContainsAny(mqttClient.ConnectAuthConfig.Servers, "tls|ssl|mqtts") {
		return nil
	}

	ca, err := ioutil.ReadFile(mqttClient.ConnectAuthConfig.ServerCaPath)
	if err != nil {
		glog.Error("load server ca failed\n")
		return err
	}
	serverCaPool := x509.NewCertPool()
	serverCaPool.AppendCertsFromPEM(ca)
	tlsConfig := &tls.Config{
		RootCAs:            serverCaPool,
		InsecureSkipVerify: true,
		MaxVersion:         tls.VersionTLS13,
		MinVersion:         tls.VersionTLS12,
		CipherSuites:       constants.CipherSuites,
		VerifyConnection:   iot.VerifyConnection(serverCaPool),
	}

	// 设备使用x.509证书认证
	if mqttClient.ConnectAuthConfig.AuthType == constants.AuthTypeX509 {
		if len(mqttClient.ConnectAuthConfig.ServerCaPath) == 0 || len(mqttClient.ConnectAuthConfig.CertFilePath) == 0 || len(mqttClient.ConnectAuthConfig.CertKeyFilePath) == 0 {
			glog.Error("device use x.509 auth but not set cert")
			return err
		}
		deviceCert, err := tls.LoadX509KeyPair(mqttClient.ConnectAuthConfig.CertFilePath, mqttClient.ConnectAuthConfig.CertKeyFilePath)
		if err != nil {
			glog.Error("load device cert failed")
			return err
		}
		var clientCerts []tls.Certificate
		clientCerts = append(clientCerts, deviceCert)
		tlsConfig.Certificates = clientCerts
	}
	options.SetTLSConfig(tlsConfig)
	return nil
}

func (mqttClient *MqttDeviceClient) connectMqttBroker() bool {
	options := mqtt.NewClientOptions()
	newPwd, err := iot.HmacSha256(mqttClient.ConnectAuthConfig.Secret, iot.TimeStamp())
	if err != nil {
		return false
	}
	options.SetPassword(newPwd)
	if mqttClient.ConnectAuthConfig.UseBootstrap {
		info := mqttClient.getServerInfo()
		if info == nil {
			return false
		}
		if len(info.Secret) != 0 {
			newPwd, err = iot.HmacSha256(info.Secret, iot.TimeStamp())
			if err != nil {
				return false
			}
			options.SetPassword(newPwd)
		}
		options.AddBroker(info.ServerUri)
	} else {
		options.AddBroker(mqttClient.ConnectAuthConfig.Servers)
	}

	options.SetClientID(assembleClientId(*mqttClient.ConnectAuthConfig))
	options.SetUsername(mqttClient.ConnectAuthConfig.Id)
	options.SetConnectTimeout(20 * time.Second)
	// 关闭sdk内部重连，使用自定义重连刷新时间戳
	options.SetAutoReconnect(false)
	options.SetConnectRetry(false)
	options.SetMaxResumePubInFlight(mqttClient.ConnectAuthConfig.InflightMessages)

	if mqttClient.ConnectAuthConfig.ConnectTimeout <= 12000 && mqttClient.ConnectAuthConfig.ConnectTimeout >= 30 {
		options.SetKeepAlive(time.Duration(mqttClient.ConnectAuthConfig.ConnectTimeout) * time.Second)
	} else {
		options.SetKeepAlive(120 * time.Second)
	}
	options.OnConnectionLost = mqttClient.createConnectionLostHandler()
	options.OnConnect = mqttClient.createConnectHandler()
	options.DefaultPublishHandler = mqttClient.createDefaultMessageHandler()

	if err := mqttClient.configureTLS(options); err != nil {
		return false
	}
	mqttClient.client = mqtt.NewClient(options)
	if token := mqttClient.client.Connect(); token.WaitTimeout(mqttClient.ConnectAuthConfig.ConnectTimeOut) && token.Error() != nil {
		glog.Warningf("device %s init failed,error = %v", mqttClient.ConnectAuthConfig.Id, token.Error())
		return false
	}
	mqttClient.RuleManageService = rule.RuleManageService{
		RuleIdList:       make(map[string]bool),
		RuleInfoMap:      make(map[string]model.RuleInfo),
		TimerRuleMap:     make(map[string]rule.TimerRuleInstance),
		ConditionExecute: rule.ConditionExecute{},
	}
	go logFlush()

	return true
}

func (mqttClient *MqttDeviceClient) getServerInfo() *model.ServerInfo {
	serverInfo := iot.GetServer()
	if serverInfo != nil {
		return serverInfo
	}
	bootstracpClient, err := NewBootstrapClient(*mqttClient.ConnectAuthConfig)
	if err != nil {
		glog.Warningf("create bootstrap client failed,err %s\n", err)
		return nil
	}
	serverAddress, deviceSecret := bootstracpClient.Boot(mqttClient.ConnectAuthConfig.BootStrapBody)
	if len(serverAddress) == 0 {
		glog.Warningf("get server address from bootstrap server failed")
		bootstracpClient.Close()
		return nil
	}
	serverInfo = &model.ServerInfo{
		ServerUri: serverAddress,
		Secret:    deviceSecret,
	}
	err = iot.SavePwd(*serverInfo)
	if err != nil {
		return nil
	}
	return serverInfo
}

func (mqttClient *MqttDeviceClient) Close(timeout uint) {
	mqttClient.client.Disconnect(timeout)
	mqttClient.Pool.Release()
}

func (mqttClient *MqttDeviceClient) IsConnect() bool {
	return mqttClient.client.IsConnected()
}

func (mqttClient *MqttDeviceClient) PublishMessage(topic string, qos byte, message string) bool {
	if !mqttClient.client.IsConnected() {
		if mqttClient.Queue != nil {
			bufferMessage := model.BufferMessage{
				Topic:   topic,
				Qos:     qos,
				Message: message,
			}
			mqttClient.Queue.Push(bufferMessage)
		}
		return false
	}
	if token := mqttClient.client.Publish(topic, qos, false, message); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s send message failed", mqttClient.ConnectAuthConfig.Id)
		return false
	}
	glog.Infof("public message success. topic: %s,  message: %s", topic, message)
	return true
}

func (mqttClient *MqttDeviceClient) PublishBufferMessage() {
	if mqttClient.Queue == nil {
		return
	}
	for mqttClient.Queue.Len() > 0 {
		pop := mqttClient.Queue.Pop()
		message, ok := pop.(model.BufferMessage)
		if ok {
			mqttClient.PublishMessage(message.Topic, message.Qos, message.Message)
			continue
		}
		glog.Warningf("message type was not match.")
	}
}

func (mqttClient *MqttDeviceClient) SubscribeCustomizeTopic(topic string, handler callback.MessageHandler) {
	if token := mqttClient.client.Subscribe(topic, mqttClient.ConnectAuthConfig.Qos, func(client mqtt.Client, message mqtt.Message) {
		payload := message.Payload()
		handler(string(payload))
	}); token.Wait() && token.Error() != nil {
		glog.Warningf("device subscribe customize message topic failed. deviceId: %s, topic: %s", mqttClient.ConnectAuthConfig.Id, topic)
		return
	}
	glog.Infof("subscribe customize topic success. topic: %s", topic)
}

func (mqttClient *MqttDeviceClient) CreateRuleActionHandler() func(actionList []model.Action) bool {
	return func(actionList []model.Action) bool {
		if mqttClient.RuleActionHandler != nil {
			return mqttClient.RuleActionHandler(actionList)
		}
		for _, action := range actionList {
			if !strings.EqualFold(action.DeviceId, mqttClient.ConnectAuthConfig.Id) {
				glog.Warningf("action device is not match. target: %s, action: %s", mqttClient.ConnectAuthConfig.Id, action.DeviceId)
				continue
			}
			if mqttClient.CommandHandler == nil {
				glog.Warningf("command handler is not define.")
				continue
			}
			command := action.Command
			deviceCommand := model.Command{
				CommandName: command.CommandName,
				ServiceId:   command.ServiceId,
				Paras:       command.CommandBody,
			}
			success, _ := mqttClient.CommandHandler(deviceCommand)
			if !success {
				glog.Warningf("handle command failed.")
			}
		}
		return true
	}
}

func (mqttClient *MqttDeviceClient) createDefaultMessageHandler() func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		err := mqttClient.Pool.Submit(func() {
			topic := message.Topic()
			glog.Infof("receive message from device. topic: %s, message: %s", topic, string(message.Payload()))
			if strings.Contains(topic, "/messages/down") {
				mqttClient.handleDeviceMessageDown(client, message)
				return
			}
			if strings.Contains(topic, "sys/commands/request_id") {
				mqttClient.handleDeviceCommand(client, message)
				return
			}
			if strings.Contains(topic, "/sys/properties/set/request_id") {
				mqttClient.handleDevicePropertiesSet(client, message)
				return
			}
			if strings.Contains(topic, "/sys/properties/get/request_id") {
				mqttClient.handleDevicePropertiesQuery(client, message)
				return
			}
			if strings.Contains(topic, "/sys/shadow/get/response") {
				mqttClient.handleDevicePropertiesQueryResponse(client, message)
				return
			}
			if strings.Contains(topic, "/sys/events/down") {
				mqttClient.handleDeviceEvent(client, message)
				return
			}
		})
		if err != nil {
			glog.Warningf("submit message failed. topic: %s, err: %s", message.Topic(), err.Error())
		}
	}
}

func (mqttClient *MqttDeviceClient) handleDeviceCommand(client mqtt.Client, message mqtt.Message) {
	command := &model.Command{}
	if json.Unmarshal(message.Payload(), command) != nil {
		glog.Warningf("unmarshal platform command failed,device id = %s，message = %s", mqttClient.ConnectAuthConfig.Id, message)
	}

	flag, response := mqttClient.CommandHandler(*command)
	var res string
	if flag {
		glog.Infof("device %s handle command success", mqttClient.ConnectAuthConfig.Id)
		res = iot.Interface2JsonString(model.CommandResponse{
			ResultCode: 0,
			Paras:      response,
		})
	} else {
		glog.Warningf("device %s handle command failed", mqttClient.ConnectAuthConfig.Id)
		res = iot.Interface2JsonString(model.CommandResponse{
			ResultCode: 1,
			Paras:      response,
		})
	}
	if token := mqttClient.client.Publish(iot.FormatTopic(constants.CommandResponseTopic, mqttClient.ConnectAuthConfig.Id)+iot.GetTopicRequestId(message.Topic()),
		1, false, res); token.Wait() && token.Error() != nil {
		glog.Infof("device %s send command response failed", mqttClient.ConnectAuthConfig.Id)
	}
}

func (mqttClient *MqttDeviceClient) handleDevicePropertiesSet(client mqtt.Client, message mqtt.Message) {
	propertiesSetRequest := &model.DevicePropertyDownRequest{}
	if json.Unmarshal(message.Payload(), propertiesSetRequest) != nil {
		glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", mqttClient.ConnectAuthConfig.Id, message)
	}

	handleFlag := true
	for _, handler := range mqttClient.PropertiesSetHandlers {
		handleFlag = handleFlag && handler(*propertiesSetRequest)
	}
	// 端侧规则的信息则需要单独处理
	if mqttClient.ConnectAuthConfig.RuleEnable {
		for _, service := range propertiesSetRequest.Services {
			if strings.Contains(service.ServiceId, "$device_rule") {
				mqttClient.RuleManageService.ModifyRule(service, func(event model.DeviceEvents) bool {
					return mqttClient.reportEvent(event)
				})
			}
		}
	}
	var res string
	response := struct {
		ResultCode byte   `json:"result_code"`
		ResultDesc string `json:"result_desc"`
	}{}
	if handleFlag {
		response.ResultCode = 0
		response.ResultDesc = "Set property success."
		res = iot.Interface2JsonString(response)
	} else {
		response.ResultCode = 1
		response.ResultDesc = "Set properties failed."
		res = iot.Interface2JsonString(response)
	}
	if token := mqttClient.client.Publish(iot.FormatTopic(constants.PropertiesSetResponseTopic, mqttClient.ConnectAuthConfig.Id)+iot.GetTopicRequestId(message.Topic()),
		mqttClient.ConnectAuthConfig.Qos, false, res); token.Wait() && token.Error() != nil {
		glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", mqttClient.ConnectAuthConfig.Id, message)
	}

}

func (mqttClient *MqttDeviceClient) handleDeviceMessageDown(client mqtt.Client, message mqtt.Message) {
	for _, handler := range mqttClient.MessageHandlers {
		handler(string(message.Payload()))
	}
}

func (mqttClient *MqttDeviceClient) handleDevicePropertiesQuery(client mqtt.Client, message mqtt.Message) {
	propertiesQueryRequest := &model.DevicePropertyQueryRequest{}
	if json.Unmarshal(message.Payload(), propertiesQueryRequest) != nil {
		glog.Warningf("device %s unmarshal properties query request failed %s", mqttClient.ConnectAuthConfig.Id, message)
	}

	queryResult := mqttClient.PropertyQueryHandler(*propertiesQueryRequest)
	responseToPlatform := iot.Interface2JsonString(queryResult)
	if token := mqttClient.client.Publish(iot.FormatTopic(constants.PropertiesQueryResponseTopic, mqttClient.ConnectAuthConfig.Id)+iot.GetTopicRequestId(message.Topic()),
		mqttClient.ConnectAuthConfig.Qos, false, responseToPlatform); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s send properties query response failed.", mqttClient.ConnectAuthConfig.Id)
	}
}
func (mqttClient *MqttDeviceClient) handleDevicePropertiesQueryResponse(client mqtt.Client, message mqtt.Message) {
	propertiesQueryResponse := &model.DeviceShadowQueryResponse{}
	if json.Unmarshal(message.Payload(), propertiesQueryResponse) != nil {
		glog.Warningf("device %s unmarshal property response failed,message %s", mqttClient.ConnectAuthConfig.Id, iot.Interface2JsonString(message))
	}
	mqttClient.DeviceShadowQueryResponseHandler(*propertiesQueryResponse)
}

func (mqttClient *MqttDeviceClient) createConnectionLostHandler() func(client mqtt.Client, reason error) {
	// 断链后进行自定义重连
	connectionLostHandler := func(client mqtt.Client, reason error) {
		if mqttClient.ConnectionLostHandler != nil {
			glog.Warningf("connection lost from server. reason: %s\n", reason.Error())
			mqttClient.ConnectionLostHandler(client, reason)
		}
		if *mqttClient.ConnectAuthConfig.AutoReconnect {
			glog.Warningf("connection lost from server. begin to reconnect broker. reason: %s\n", reason.Error())
			connected := mqttClient.Connect()
			if connected {
				glog.Infof("reconnect mqtt go broker success.")
			}
		}
	}
	return connectionLostHandler
}

func (mqttClient *MqttDeviceClient) createConnectHandler() func(client mqtt.Client) {
	// 断链后进行自定义重连
	onConnectHandler := func(client mqtt.Client) {
		glog.Infof("connect from server.")
		err := mqttClient.Pool.Submit(func() {
			mqttClient.PublishBufferMessage()
		})
		if err != nil {
			glog.Warningf("submit buffer ")
		}
		if mqttClient.ConnectHandler != nil {
			mqttClient.ConnectHandler(client)
			return
		}
	}
	return onConnectHandler
}

// 平台向设备下发的事件callback
func (mqttClient *MqttDeviceClient) handleDeviceEvent(client mqtt.Client, message mqtt.Message) {
	data := &model.Data{}
	err := json.Unmarshal(message.Payload(), data)
	if err != nil {
		glog.Warningf("handle platform to device data failed. err: %s", err.Error())
		return
	}
	for _, entry := range data.Services {
		serviceId := entry.ServiceId
		switch serviceId {
		case "$sub_device_manager":
			mqttClient.handleSubDeviceService(entry)
		case "$file_manager":
			mqttClient.handleFileService(entry)
		case "$ota":
			mqttClient.handleOtaService(entry)
		case "$log":
			mqttClient.handleDeviceLogService(entry)
		case "$time_sync":
			mqttClient.handleTimeSyncService(entry)
		case "$device_rule":
			mqttClient.handleDeviceRuleService(entry)
		}
	}
}

func (mqttClient *MqttDeviceClient) handleDeviceRuleService(entry model.DataEntry) {
	if strings.EqualFold("device_rule_config_response", entry.EventType) {
		if !mqttClient.ConnectAuthConfig.RuleEnable {
			return
		}
		paras := model.RuleParas{}
		err := json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), &paras)
		if err != nil {
			glog.Warningf("convert event msg to rule info failed. err: %s", err)
			return
		}
		mqttClient.RuleManageService.QueryRuleResponse(paras.RulesInfos, mqttClient.CreateRuleActionHandler())
	}
}

func (mqttClient *MqttDeviceClient) handleDeviceLogService(entry model.DataEntry) {
	if strings.EqualFold("log_config", entry.EventType) {
		// 平台下发日志收集通知
		glog.Infof("platform send log collect command")
		logConfig := &config.LogCollectionConfig{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), logConfig) != nil {
			return
		}

		lcc := &config.LogCollectionConfig{
			LogCollectSwitch: logConfig.LogCollectSwitch,
			EndTime:          logConfig.EndTime,
		}
		mqttClient.Lcc = lcc
		mqttClient.reportLogsWorker()
	}
}

func (mqttClient *MqttDeviceClient) handleTimeSyncService(entry model.DataEntry) {
	if strings.EqualFold("time_sync_response", entry.EventType) {
		timeSyncResponse := &model.TimeSyncResponse{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), timeSyncResponse) != nil {
			return
		}
		mqttClient.SyncTimeResponseHandler(timeSyncResponse.DeviceSendTime, timeSyncResponse.ServerRecvTime, timeSyncResponse.ServerSendTime)
	}
}

func (mqttClient *MqttDeviceClient) handleOtaService(entry model.DataEntry) {
	eventType := entry.EventType
	switch eventType {
	case "version_query":
		// 查询软固件版本
		mqttClient.reportVersion()
	case "firmware_upgrade", "firmware_upgrade_v2":
		upgradeInfo := &model.UpgradeInfo{}
		jsonString := iot.Interface2JsonString(entry.Paras)
		err := json.Unmarshal([]byte(jsonString), upgradeInfo)
		if err != nil {
			glog.Warningf("unmarshal firware upgrade failed. err: %s", err.Error())
			return
		}
		if "firmware_upgrade" == eventType {
			mqttClient.upgradeDevice(1, upgradeInfo)
		} else {
			mqttClient.upgradeDevice(3, upgradeInfo)
		}

	case "software_upgrade", "software_upgrade_v2":
		upgradeInfo := &model.UpgradeInfo{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), upgradeInfo) != nil {
			return
		}
		if "software_upgrade" == eventType {
			mqttClient.upgradeDevice(0, upgradeInfo)
		} else {
			mqttClient.upgradeDevice(2, upgradeInfo)
		}
	}
}

func (mqttClient *MqttDeviceClient) handleFileService(entry model.DataEntry) {
	eventType := entry.EventType
	switch eventType {
	case "get_upload_url_response":
		// 获取文件上传URL
		fileResponse := &model.FileResponseServiceEventParas{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), fileResponse) != nil {
			return
		}
		mqttClient.FileUrls[fileResponse.ObjectName+constants.FileActionUpload] = fileResponse.Url
	case "get_download_url_response":
		fileResponse := &model.FileResponseServiceEventParas{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), fileResponse) != nil {
			return
		}
		mqttClient.FileUrls[fileResponse.ObjectName+constants.FileActionDownload] = fileResponse.Url
	}
}

func (mqttClient *MqttDeviceClient) handleSubDeviceService(entry model.DataEntry) {
	eventType := entry.EventType
	switch eventType {
	case "add_sub_device_notify":
		// 子设备添加
		subDeviceInfo := &model.SubDeviceInfo{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceInfo) != nil {
			return
		}
		if mqttClient.SubDevicesAddHandler == nil {
			return
		}
		mqttClient.SubDevicesAddHandler(*subDeviceInfo)
	case "delete_sub_device_notify":
		subDeviceInfo := &model.SubDeviceInfo{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceInfo) != nil {
			return
		}
		if mqttClient.SubDevicesDeleteHandler == nil {
			return
		}
		mqttClient.SubDevicesDeleteHandler(*subDeviceInfo)
	case "sub_device_update_status_response":
		subDeviceStatusResp := &model.SubDeviceStatusResp{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceStatusResp) != nil {
			return
		}
		if mqttClient.SubDeviceStatusRespHandler == nil {
			return
		}
		mqttClient.SubDeviceStatusRespHandler(*subDeviceStatusResp)
	case "add_sub_device_response":
		subDeviceResponse := &model.SubDeviceAddResponse{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceResponse) != nil {
			return
		}
		if mqttClient.SubDeviceAddResponseHandler == nil {
			return
		}
		mqttClient.SubDeviceAddResponseHandler(*subDeviceResponse)
	case "delete_sub_device_response":
		subDeviceResponse := &model.SubDeviceDeleteResponse{}
		if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceResponse) != nil {
			return
		}
		if mqttClient.SubDeviceDeleteResponseHandler == nil {
			return
		}
		mqttClient.SubDeviceDeleteResponseHandler(*subDeviceResponse)
	}
}

func (mqttClient *MqttDeviceClient) reportEvent(event model.DeviceEvents) bool {
	eventStr := iot.Interface2JsonString(event)
	topic := iot.FormatTopic(constants.DeviceToPlatformTopic, mqttClient.ConnectAuthConfig.Id)
	if token := mqttClient.client.Publish(topic, mqttClient.ConnectAuthConfig.Qos, false, eventStr); token.Wait() && token.Error() != nil {
		glog.Errorf("device %s report version query failed,type %d", mqttClient.ConnectAuthConfig.Id)
		return false
	}
	glog.Infof("report event success. topic: %s, msg: %s", topic, eventStr)
	return true
}

func (mqttClient *MqttDeviceClient) reportVersion() {
	sw, fw := mqttClient.SwFwVersionReporter()
	dataEntry := model.DataEntry{
		ServiceId: "$ota",
		EventType: "version_report",
		EventTime: iot.GetEventTimeStamp(),
		Paras: struct {
			SwVersion string `json:"sw_version"`
			FwVersion string `json:"fw_version"`
		}{
			SwVersion: sw,
			FwVersion: fw,
		},
	}
	data := model.Data{
		ObjectDeviceId: mqttClient.ConnectAuthConfig.Id,
		Services:       []model.DataEntry{dataEntry},
	}
	glog.Infof("report version query. body: %s", iot.Interface2JsonString(data))
	if token := mqttClient.client.Publish(iot.FormatTopic(constants.DeviceToPlatformTopic, mqttClient.ConnectAuthConfig.Id), mqttClient.ConnectAuthConfig.Qos,
		false, iot.Interface2JsonString(data)); token.Wait() && token.Error() != nil {
		glog.Errorf("device %s report version query failed,type %d", mqttClient.ConnectAuthConfig.Id)
	}
}

func (mqttClient *MqttDeviceClient) upgradeDevice(upgradeType byte, upgradeInfo *model.UpgradeInfo) {
	progress := mqttClient.DeviceUpgradeHandler(upgradeType, *upgradeInfo)
	dataEntry := model.DataEntry{
		ServiceId: "$ota",
		EventType: "upgrade_progress_report",
		EventTime: iot.GetEventTimeStamp(),
		Paras:     progress,
	}
	data := model.Data{
		ObjectDeviceId: mqttClient.ConnectAuthConfig.Id,
		Services:       []model.DataEntry{dataEntry},
	}

	if token := mqttClient.client.Publish(iot.FormatTopic(constants.DeviceToPlatformTopic, mqttClient.ConnectAuthConfig.Id), mqttClient.ConnectAuthConfig.Qos,
		false, iot.Interface2JsonString(data)); token.Wait() && token.Error() != nil {
		glog.Errorf("device %s upgrade failed,type %d", mqttClient.ConnectAuthConfig.Id, upgradeType)
	}
}

func (mqttClient *MqttDeviceClient) reportLogsWorker() {
	mqttClient.deviceStatusLogCollect()
	mqttClient.devicePropertyReportLogCollect()
	mqttClient.deviceMessageLogCollect()
	mqttClient.deviceCommandLogCollect()
}

func (mqttClient *MqttDeviceClient) deviceCommandLogCollect() {
	go func() {
		for {
			if !mqttClient.Lcc.GetLogCollectSwitch() {
				break
			}
			logs := mqttClient.DeviceCommandLogCollector(mqttClient.Lcc.GetEndTime())
			if len(logs) == 0 {
				glog.Warningf("no log about device command")
				break
			}
			mqttClient.reportLogs(logs)
		}

	}()
}

func (mqttClient *MqttDeviceClient) deviceMessageLogCollect() {
	go func() {
		for {
			if !mqttClient.Lcc.GetLogCollectSwitch() {
				break
			}
			logs := mqttClient.DeviceMessageLogCollector(mqttClient.Lcc.GetEndTime())
			if len(logs) == 0 {
				glog.Warningf("no log about device message")
				break
			}
			mqttClient.reportLogs(logs)
		}

	}()
}

func (mqttClient *MqttDeviceClient) devicePropertyReportLogCollect() {
	go func() {
		for {
			if !mqttClient.Lcc.GetLogCollectSwitch() {
				break
			}
			logs := mqttClient.DevicePropertyLogCollector(mqttClient.Lcc.GetEndTime())
			if len(logs) == 0 {
				glog.Warningf("no log about device property")
				break
			}
			mqttClient.reportLogs(logs)
		}

	}()
}

func (mqttClient *MqttDeviceClient) deviceStatusLogCollect() {
	go func() {
		for {
			if !mqttClient.Lcc.GetLogCollectSwitch() {
				break
			}
			logs := mqttClient.DeviceStatusLogCollector(mqttClient.Lcc.GetEndTime())
			if len(logs) == 0 {
				glog.Warningf("no log about device status")
				break
			}
			mqttClient.reportLogs(logs)
		}

	}()
}

func (mqttClient *MqttDeviceClient) reportLogs(logs []model.DeviceLogEntry) {
	var dataEntries []model.DataEntry
	for _, log := range logs {
		dataEntry := model.DataEntry{
			ServiceId: "$log",
			EventType: "log_report",
			EventTime: iot.GetEventTimeStamp(),
			Paras:     log,
		}
		dataEntries = append(dataEntries, dataEntry)
	}
	data := model.Data{
		Services: dataEntries,
	}

	reportedLog := iot.Interface2JsonString(data)
	mqttClient.client.Publish(iot.FormatTopic(constants.DeviceToPlatformTopic, mqttClient.ConnectAuthConfig.Id), 0, false, reportedLog)
}

func assembleClientId(device config.ConnectAuthConfig) string {
	segments := make([]string, 4)
	segments[0] = device.Id
	segments[1] = "0"
	if device.VerifyTimestamp {
		segments[2] = "1"
	} else {
		segments[2] = "0"
	}
	segments[3] = iot.TimeStamp()

	return strings.Join(segments, "_")
}

func logFlush() {
	ticker := time.Tick(5 * time.Second)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-ticker:
			glog.Flush()
			break
		case <-interrupt:
			return
		}
	}
}
