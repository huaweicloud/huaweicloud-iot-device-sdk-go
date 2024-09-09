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

package model

import "net"

// UpgradeInfo 平台下发的升级信息
type UpgradeInfo struct {
	Version     string `json:"version"`      // 软固件包版本号
	Url         string `json:"url"`          // 软固件包下载地址
	FileSize    int    `json:"file_size"`    // 软固件包文件大小
	AccessToken string `json:"access_token"` // 软固件包url下载地址的临时token
	Expires     int64  `json:"expires"`      // access_token的超期时间
	Sign        string `json:"sign"`         // 软固件包MD5值
}

// UpgradeProgress 设备升级状态响应，用于设备向平台反馈进度，错误信息等
// ResultCode： 设备的升级状态，结果码定义如下：
// 0：处理成功
// 1：设备使用中
// 2：信号质量差
// 3：已经是最新版本
// 4：电量不足
// 5：剩余空间不足
// 6：下载超时
// 7：升级包校验失败
// 8：升级包类型不支持
// 9：内存不足
// 10：安装升级包失败
// 255： 内部异常
type UpgradeProgress struct {
	ResultCode  int    `json:"result_code"`
	Progress    int    `json:"progress"`    // 设备的升级进度，范围：0到100
	Version     string `json:"version"`     // 设备当前版本号
	Description string `json:"description"` // 升级状态描述信息，可以返回具体升级失败原因。
}

// Command 设备命令
type Command struct {
	ObjectDeviceId string      `json:"object_device_id"`
	ServiceId      string      `json:"service_id"`
	CommandName    string      `json:"command_name"`
	Paras          interface{} `json:"paras"`
}

type CommandResponse struct {
	ResultCode   byte        `json:"result_code"`
	ResponseName string      `json:"response_name"`
	Paras        interface{} `json:"paras"`
}

// Message 设备消息
type Message struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

// 定义平台和设备之间的数据交换结构体

type Data struct {
	ObjectDeviceId string      `json:"object_device_id,omitempty"`
	Services       []DataEntry `json:"services"`
}

type DataEntry struct {
	ServiceId string      `json:"service_id"`
	EventType string      `json:"event_type"`
	EventTime string      `json:"event_time"`
	Paras     interface{} `json:"paras"` // 不同类型的请求paras使用的结构体不同
}

// SubDevicesStatus 网关更新子设备状态
type SubDevicesStatus struct {
	DeviceStatuses []DeviceStatus `json:"device_statuses"`
}

type DeviceStatus struct {
	DeviceId string `json:"device_id"`
	Status   string `json:"status"` // 子设备状态。 OFFLINE：设备离线 ONLINE：设备上线
}

// SubDeviceInfo 添加子设备
type SubDeviceInfo struct {
	Devices []DeviceInfo `json:"devices"`
	Version int          `json:"version"`
}

type DeviceInfo struct {
	ParentDeviceId string      `json:"parent_device_id,omitempty"`
	NodeId         string      `json:"node_id,omitempty"`
	DeviceId       string      `json:"device_id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	ManufacturerId string      `json:"manufacturer_id,omitempty"`
	Model          string      `json:"model,omitempty"`
	ProductId      string      `json:"product_id"`
	FwVersion      string      `json:"fw_version,omitempty"`
	SwVersion      string      `json:"sw_version,omitempty"`
	Status         string      `json:"status,omitempty"`
	ExtensionInfo  interface{} `json:"extension_info,omitempty"`
}

type SubDeviceDeleteResponse struct {
	SuccessFulDevices []string       `json:"successful_devices"`
	FailedDevices     []FailedDevice `json:"failed_devices"`
}

// SubDeviceAddResponse 网关添加删除子设备响应
type SubDeviceAddResponse struct {
	SuccessFulDevices []DeviceInfo   `json:"successful_devices"`
	FailedDevices     []FailedDevice `json:"failed_devices"`
}

// SubDeviceStatusResp 网关更新子设备状态响应
type SubDeviceStatusResp struct {
	SuccessfulDevices []SuccessDevice `json:"successful_devices"`
	FailedDevices     []FailedDevice  `json:"failed_devices"`
}

type SuccessDevice struct {
	DeviceId string `json:"device_id"`
	Status   string `json:"status"`
}

type FailedDevice struct {
	DeviceId  string `json:"device_id"`
	ProductId string `json:"product_id"`
	NodeId    string `json:"node_id"`
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// 设备属性相关

// DeviceProperties 设备属性
type DeviceProperties struct {
	Services []DevicePropertyEntry `json:"services"`
}

type BootStrapProperties struct {
	BaseStrategyKeyword string `json:"baseStrategyKeyword"`
}

// DevicePropertyEntry 设备的一个属性
type DevicePropertyEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// DevicePropertyDownRequest 平台设置设备属性
type DevicePropertyDownRequest struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []DevicePropertyDownRequestEntry `json:"services"`
}

type DevicePropertyDownRequestEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
}

// DevicePropertyQueryRequest 平台设置设备属性
type DevicePropertyQueryRequest struct {
	ObjectDeviceId string `json:"object_device_id"`
	ServiceId      string `json:"service_id"`
}

type DeviceShadowQueryResponse struct {
	ObjectDeviceId string             `json:"object_device_id"`
	Shadow         []DeviceShadowData `json:"shadow"`
}

type DeviceShadowData struct {
	ServiceId string                     `json:"service_id"`
	Desired   DeviceShadowPropertiesData `json:"desired"`
	Reported  DeviceShadowPropertiesData `json:"reported"`
	Version   int                        `json:"version"`
}
type DeviceShadowPropertiesData struct {
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// 网关批量上报子设备属性

type DevicesService struct {
	Devices []DeviceService `json:"devices"`
}

type DeviceService struct {
	DeviceId string                `json:"device_id"`
	Services []DevicePropertyEntry `json:"services"`
}

// DeviceMessage 消息下发格式
type DeviceMessage struct {
	ObjectDeviceId string `json:"object_device_id"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Content        string `json:"content"`
}

// FileRequest 设备获取文件上传下载请求体
type FileRequest struct {
	ObjectDeviceId string                    `json:"object_device_id"`
	Services       []FileRequestServiceEvent `json:"services"`
}

// FileResponse 平台下发文件上传和下载URL响应
type FileResponse struct {
	ObjectDeviceId string                     `json:"object_device_id"`
	Services       []FileResponseServiceEvent `json:"services"`
}

type FileResultResponse struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []FileResultResponseServiceEvent `json:"services"`
}

type BaseServiceEvent struct {
	ServiceId string `json:"service_id"`
	EventType string `json:"event_type"`
	EventTime string `json:"event_time,omitempty"`
}

type FileRequestServiceEvent struct {
	BaseServiceEvent
	Paras FileRequestServiceEventParas `json:"paras"`
}

// TimeSyncRequest 时间同步相关
type TimeSyncRequest struct {
	Services []TimeSyncRequestServiceEvent `json:"services"`
}

type TimeSyncRequestServiceEvent struct {
	BaseServiceEvent
	Paras TimeSyncRequestServiceEventParas `json:"paras"`
}

type FileResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResponseServiceEventParas `json:"paras"`
}

type FileResultResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResultServiceEventParas `json:"paras"`
}

// FileRequestServiceEventParas 设备获取文件上传下载URL参数
type FileRequestServiceEventParas struct {
	FileName       string      `json:"file_name"`
	FileAttributes interface{} `json:"file_attributes"`
}

type TimeSyncRequestServiceEventParas struct {
	DeviceSendTime int64 `json:"device_send_time"`
}

type TimeSyncResponse struct {
	DeviceSendTime int64 `json:"device_send_time"`
	ServerRecvTime int64 `json:"server_recv_time"`
	ServerSendTime int64 `json:"server_send_time"`
}

// FileResponseServiceEventParas 平台下发响应参数
type FileResponseServiceEventParas struct {
	Url            string      `json:"url"`
	BucketName     string      `json:"bucket_name"`
	ObjectName     string      `json:"object_name"`
	Expire         int         `json:"expire"`
	FileAttributes interface{} `json:"file_attributes"`
}

// FileResultServiceEventParas 上报文件上传下载结果参数
type FileResultServiceEventParas struct {
	ObjectName        string `json:"object_name"`
	ResultCode        int    `json:"result_code"`
	StatusCode        int    `json:"status_code"`
	StatusDescription string `json:"status_description"`
}

// ReportDeviceInfoRequest 上报设备信息请求
type ReportDeviceInfoRequest struct {
	ObjectDeviceId string                         `json:"object_device_id,omitempty"`
	Services       []ReportDeviceInfoServiceEvent `json:"services,omitempty"`
}

type ReportDeviceInfoServiceEvent struct {
	BaseServiceEvent
	Paras ReportDeviceInfoEventParas `json:"paras,omitempty"`
}

// ReportDeviceInfoEventParas 设备信息上报请求参数
type ReportDeviceInfoEventParas struct {
	DeviceSdkVersion string `json:"device_sdk_version,omitempty"`
	SwVersion        string `json:"sw_version,omitempty"`
	FwVersion        string `json:"fw_version,omitempty"`
}

// ReportDeviceLogRequest 上报设备日志请求
type ReportDeviceLogRequest struct {
	Services []ReportDeviceLogServiceEvent `json:"services,omitempty"`
}

type ReportDeviceLogServiceEvent struct {
	BaseServiceEvent
	Paras DeviceLogEntry `json:"paras,omitempty"`
}

type DeviceLogEntry struct {
	Timestamp string `json:"timestamp"` // 日志产生时间
	Type      string `json:"type"`      // 日志类型：DEVICE_STATUS，DEVICE_PROPERTY ，DEVICE_MESSAGE ，DEVICE_COMMAND
	Content   string `json:"content"`   // 日志内容
}

// Session 网关相关
type Session struct {
	DeviceId string
	Conn     net.Conn
}

// OTAQueryVersion OTA 软固件升级相关
type OTAQueryVersion struct {
	TaskId         string      `json:"task_id"`
	TaskExtInfo    interface{} `json:"task_ext_info"`
	SubDeviceCount int         `json:"sub_device_count"`
}

type OTAPackageInfo struct {
	Url            string `json:"url"`
	Version        string `json:"version"`
	AccessToken    string `json:"access_token"`
	TaskId         string `json:"task_id"`
	FileSize       int    `json:"file_size"`
	FileName       string `json:"file_name"`
	Expires        int    `json:"expires"`
	Sign           string `json:"sign"`
	CustomInfo     string `json:"custom_info"`
	SubDeviceCount string `json:"sub_device_count"`
	TaskExtInfo    string `json:"task_ext_info"`
}

// RuleInfo 端侧规则

type RuleParas struct {
	RulesInfos []RuleInfo `json:"rulesInfos"`
}

type RuleInfo struct {
	RuleId              string      `json:"ruleId"`
	RuleName            string      `json:"ruleName"`
	Logic               string      `json:"logic"`
	TimeRange           TimeRange   `json:"timeRange"`
	Status              string      `json:"status"`
	Conditions          []Condition `json:"conditions"`
	Actions             []Action    `json:"actions"`
	RuleVersionInShadow int         `json:"ruleVersionInShadow"`
}

type TimeRange struct {
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	DaysOfWeek string `json:"daysOfWeek"`
}

type Condition struct {
	Type           string         `json:"type"`
	StartTime      string         `json:"startTime"`
	RepeatInterval int            `json:"repeatInterval"`
	RepeatCount    int            `json:"repeatCount"`
	Time           string         `json:"time"`
	DaysOfWeek     string         `json:"daysOfWeek"`
	Operator       string         `json:"operator"`
	DeviceInfo     RuleDeviceInfo `json:"deviceInfo"`
	Value          string         `json:"value"`
	InValue        []string       `json:"inValues"`
}

type RuleDeviceInfo struct {
	DeviceId string `json:"deviceId"`
	Path     string `json:"path"`
}

type Action struct {
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	DeviceId string      `json:"deviceId"`
	Command  RuleCommand `json:"command"`
}

type RuleVersion struct {
	Version int `json:"version"`
}

//  端侧规则 end

// DeviceRuleEvent 设备规则事件

type DeviceEvents struct {
	DeviceId string            `json:"object_device_id"`
	Services []DeviceRuleEvent `json:"services"`
}

type DeviceRuleEvent struct {
	BaseServiceEvent
	Paras DeviceRuleRequestEventParams `json:"paras"`
}

type DeviceRuleRequestEventParams struct {
	RuleIds []string `json:"ruleIds"`
	DelIds  []string `json:"delIds"`
}

type RuleCommand struct {
	ServiceId   string      `json:"serviceId"`
	CommandName string      `json:"commandName"`
	CommandBody interface{} `json:"commandBody"`
}

type BufferMessage struct {
	Topic   string
	Qos     byte
	Message string
}

type ServerInfo struct {
	ServerUri string `json:"server_uri"`
	Secret    string `json:"secret"`
	Port      int    `json:"port"`
}
