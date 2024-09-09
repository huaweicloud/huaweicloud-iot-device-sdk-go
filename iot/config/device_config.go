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

import (
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"sync"
	"time"
)

// ConnectAuthConfig 用于创建设备的参数
type ConnectAuthConfig struct {
	Id                 string
	Password           string
	VerifyTimestamp    bool
	Servers            string
	Qos                byte  // qos default 0
	BatchSubDeviceSize int   // 一次上报数据的子设备数量 默认10， 若超过该值， 则默认会分多次上报
	AuthType           uint8 // 认证类型， 密码认证或证书认证
	BsServerCaPath     string
	ServerCaPath       string
	CertFilePath       string
	CertKeyFilePath    string
	UseBootstrap       bool // 使用设备引导功能开关，true-使用，false-不使用
	ScopeId            string
	BootStrapBody      *model.BootStrapProperties // 使用设备引导功能时，静态策略为数据上报时的结构体
	ConnectTimeOut     time.Duration              // 与平台建链超时时间
	AutoReconnect      *bool                      // 是否支持断链重连，默认为True
	BackOffTime        int64                      // 退避系数，默认1000ms
	MinBackOffTime     int64                      // 最小退避时间, 默认1000ms
	MaxBackOffTime     int64                      // 最大退避时间, 默认30000ms
	ThreadNum          int                        // 协程数量，用于处理平台的消息,默认10
	RuleEnable         bool                       // 是否开启端侧规则
	MaxBufferMessage   int                        // max buffer max
	InflightMessages   int                        // qos1时最多可以同时发布多条消息，默认20条
}

type ScopeConfig struct {
	ScopeId      string
	ScopeType    uint8  // 认证类型， 密码认证或证书认证
	CertFilePath string // 注册组证书认证方式的设备侧证书路径
	CertKeyPath  string // 注册组证书认证方式的设备侧证书私钥路径
}

type LogCollectionConfig struct {
	rw               sync.RWMutex
	LogCollectSwitch bool   // on：开启设备侧日志收集功能 off：关闭设备侧日志收集开关
	EndTime          string // format yyyy-MM-dd'T'HH:mm:ss'Z'
}

func (lcc *LogCollectionConfig) SetLogCollectSwitch(switchFlag bool) {
	lcc.rw.Lock()
	defer lcc.rw.Unlock()
	lcc.LogCollectSwitch = switchFlag
}

func (lcc *LogCollectionConfig) SetEndTime(endTime string) {
	lcc.rw.Lock()
	defer lcc.rw.Unlock()
	lcc.EndTime = endTime
}

func (lcc *LogCollectionConfig) GetLogCollectSwitch() bool {
	lcc.rw.RLock()
	defer lcc.rw.RUnlock()
	return lcc.LogCollectSwitch
}

func (lcc *LogCollectionConfig) GetEndTime() string {
	lcc.rw.RLock()
	defer lcc.rw.RUnlock()
	return lcc.EndTime
}
