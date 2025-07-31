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
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"io/ioutil"
	"sync"
	"time"
)

type BootstrapClient interface {
	Boot(properties *model.BootStrapProperties) (string, string)
	Close()
}

func NewBootstrapClient(authConfig config.ConnectAuthConfig) (BootstrapClient, error) {
	client := &bsClient{
		id:             authConfig.Id,
		Secret:         authConfig.Secret,
		bsServer:       authConfig.Servers,
		serverCaPath:   authConfig.BsServerCaPath,
		connectTimeOut: authConfig.ConnectTimeOut,
		scopeConfig: &config.ScopeConfig{
			ScopeId:      authConfig.ScopeId,
			ScopeType:    authConfig.AuthType,
			CertFilePath: authConfig.CertFilePath,
			CertKeyPath:  authConfig.CertKeyFilePath,
		},
		iotdaServer: newResult(),
	}

	res, err := client.init()
	if res {
		return client, nil
	}

	return nil, err
}

type bsClient struct {
	id             string
	Secret         string
	bsServer       string // 设备发放接入地址
	serverCaPath   string
	scopeConfig    *config.ScopeConfig // 注册组的scope信息
	client         mqtt.Client         // 使用的MQTT客户端
	iotdaServer    *Result             // 设备接入平台地址
	connectTimeOut time.Duration
}

func (bs *bsClient) init() (bool, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(bs.bsServer)
	options.SetClientID(iot.CreateMqttClientId(bs.id, bs.scopeConfig))
	options.SetUsername(bs.id)
	pwd, err := bs.generatePwd()
	if err != nil {
		return false, err
	}
	options.SetPassword(pwd)
	var clientCerts []tls.Certificate
	// 注册组证书模式和非注册组证书模式
	if bs.scopeConfig != nil && bs.scopeConfig.ScopeType == constants.AuthTypeX509 {
		options.SetPassword("")
		deviceCert, err := tls.LoadX509KeyPair(bs.scopeConfig.CertFilePath, bs.scopeConfig.CertKeyPath)
		if err != nil {
			glog.Error("load device cert failed")
			panic("load device cert failed")
		}
		clientCerts = append(clientCerts, deviceCert)
	}
	options.SetKeepAlive(250 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)

	ca, err := ioutil.ReadFile(bs.serverCaPath)
	if err != nil {
		glog.Error("load bs server ca failed\n")
		panic(err)
	}
	serverCaPool := x509.NewCertPool()
	serverCaPool.AppendCertsFromPEM(ca)

	tlsConfig := &tls.Config{
		RootCAs:            serverCaPool,
		Certificates:       clientCerts,
		InsecureSkipVerify: true,
		MaxVersion:         tls.VersionTLS13,
		MinVersion:         tls.VersionTLS12,
		CipherSuites:       constants.CipherSuites,
		VerifyConnection:   iot.VerifyConnection(serverCaPool),
	}
	options.SetTLSConfig(tlsConfig)

	bs.client = mqtt.NewClient(options)
	if token := bs.client.Connect(); token.WaitTimeout(bs.connectTimeOut) && token.Error() != nil {
		glog.Warningf("device %s create bootstrap client failed,error = %v", bs.id, token.Error())
		return false, token.Error()
	}
	downTopic := fmt.Sprintf("$oc/devices/%s/sys/bootstrap/down", bs.id)
	subRes := bs.client.Subscribe(downTopic, 0, bs.handleSubscribeHandler())
	if subRes.WaitTimeout(bs.connectTimeOut) && subRes.Error() != nil {
		glog.Warningf("sub topic %s failed,error is %s\n", downTopic, subRes.Error())
		panic("sub bs topic failed.")
	} else {
		glog.Infof("sub topic %s success\n", downTopic)
	}
	return true, nil
}

func (bs *bsClient) handleSubscribeHandler() func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		go func() {
			glog.Infof("get message from bs server")
			serverResponse := &serverResponse{}
			err := json.Unmarshal(message.Payload(), serverResponse)
			if err != nil {
				glog.Warningf("subscribe bootstrap topic failed. err: %s", err.Error())
				bs.iotdaServer.CompleteError(err)
			} else {
				glog.Infof("bootstrap success. address is : %s", serverResponse.Address)
				bs.iotdaServer.Complete(serverResponse.Address, serverResponse.DeviceSecret)
			}
		}()
	}
}

func (bs *bsClient) generatePwd() (string, error) {
	// 注册组秘钥模式
	if bs.scopeConfig != nil && bs.scopeConfig.ScopeId != "" && bs.scopeConfig.ScopeType == constants.AuthTypePassword {
		password, err := base64.StdEncoding.DecodeString(bs.Secret)
		if err != nil {
			glog.Warningf("decode Secret failed.")
			return "", err
		}
		newPwd, err := iot.HmacSha256(string(password), bs.id)
		if err != nil {
			return "", err
		}
		newPwd, err = iot.HmacSha256(newPwd, iot.TimeStamp())
		if err != nil {
			return "", err
		}
		return newPwd, nil
	}
	if bs.Secret != "" {
		pwd, err := iot.HmacSha256(bs.Secret, iot.TimeStamp())
		if err != nil {
			return "", err
		}
		return pwd, nil
	}
	return "", nil
}

func (bs *bsClient) Boot(properties *model.BootStrapProperties) (string, string) {
	upTopic := fmt.Sprintf("$oc/devices/%s/sys/bootstrap/up", bs.id)
	propertiesData := iot.Interface2JsonString(properties)
	pubRes := bs.client.Publish(upTopic, 0, false, propertiesData)
	if pubRes.Wait() && pubRes.Error() != nil {
		glog.Warningf("bootstrap failed. err: %s", pubRes.Error())
		return "", ""
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	success, err := bs.iotdaServer.WaitTimeout(timeout)
	if err != nil {
		glog.Warningf("wait iotdaServer response err. err: %s", err.Error())
		return "", ""
	}
	if success {
		return "tls://" + bs.iotdaServer.Value(), bs.iotdaServer.secret
	}
	glog.Warningf("receive response from server failed. err: %s", err.Error())
	return "", ""
}

func (bs *bsClient) Close() {
	bs.client.Disconnect(1000)
}

type serverResponse struct {
	Address      string `json:"address"`
	DeviceSecret string `json:"deviceSecret"`
}

type Result struct {
	Flag   chan int
	err    error
	mErr   sync.RWMutex
	res    string
	secret string
	mRes   sync.RWMutex
}

type Response struct {
	address      string
	deviceSecret string
}

func (b *Result) Value() string {
	b.mRes.RLock()
	defer b.mRes.RUnlock()
	return b.res
}

func (b *Result) Wait() bool {
	<-b.Flag
	return true
}

func (b *Result) WaitTimeout(ctx context.Context) (bool, error) {
	for {
		select {
		case <-b.Flag:
			return true, nil

		case <-ctx.Done():
			return false, ctx.Err()
		}
	}
}

func (b *Result) Complete(res, deviceSecret string) {
	b.mRes.Lock()
	defer b.mRes.Unlock()
	b.res = res
	b.secret = deviceSecret
	b.Flag <- 1
}

func (b *Result) CompleteError(err error) {
	b.mErr.Lock()
	defer b.mErr.Unlock()
	b.err = err
	b.Flag <- 1
}

func (b *Result) Error() error {
	b.mErr.RLock()
	defer b.mErr.RUnlock()
	return b.err
}

func newResult() *Result {
	return &Result{
		Flag: make(chan int),
	}
}
