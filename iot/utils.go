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

package iot

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/config"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/constants"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

// 时间戳：为设备连接平台时的UTC时间，格式为YYYYMMDDHH，如UTC 时间2018/7/24 17:56:20 则应表示为2018072417。
func TimeStamp() string {
	strFormatTime := time.Now().Format("2006-01-02 15:04:05")
	strFormatTime = strings.ReplaceAll(strFormatTime, "-", "")
	strFormatTime = strings.ReplaceAll(strFormatTime, " ", "")
	strFormatTime = strFormatTime[0:10]
	return strFormatTime
}

// GetEventTimeStamp 设备采集数据UTC时间（格式：yyyyMMdd'T'HHmmss'Z'），如：20161219T114920Z。
// 设备上报数据不带该参数或参数格式错误时，则数据上报时间以平台时间为准。
func GetEventTimeStamp() string {
	now := time.Now().UTC()
	return now.Format("20060102T150405Z")
}

func GetDateTime(timeStr string) (time.Time, error) {
	format := "2006-01-02 15:04:05"
	parse, err := time.Parse(format, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

func HmacSha256(data string, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		glog.Warningf("hmac password failed. err: %s", err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Sha256FromFile(filePath string) (int, string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		glog.Warningf("read file failed. path : %s", filePath)
		return 0, "", errors.New("read file failed")
	}
	hash := sha256.New()
	_, err = hash.Write(file)
	if err != nil {
		return 0, "", nil
	}
	return len(file), hex.EncodeToString(hash.Sum(nil)), nil
}

func Interface2JsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	byteData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(byteData)
}

func GetTopicRequestId(topic string) string {
	return strings.Split(topic, "=")[1]
}

func FormatTopic(topic, deviceId string) string {
	return strings.ReplaceAll(topic, "{device_id}", deviceId)
}

// 根据当前运行的操作系统重新修改文件路径以适配操作系统
func SmartFileName(filename string) string {
	// Windows操作系统适配
	if strings.Contains(OsName(), "windows") {
		if strings.Contains(filename, "/") {
			pathParts := strings.Split(filename, "/")
			pathParts[0] = pathParts[0] + ":"
			return strings.Join(pathParts, "\\")
		}
	}

	return filename
}

func CreateMqttClientId(deviceId string, scopeConfig *config.ScopeConfig) string {
	if scopeConfig == nil || scopeConfig.ScopeId == "" {
		return fmt.Sprintf("%s_%s_%s_%s", deviceId, "0", "0", TimeStamp())
	}
	if scopeConfig.ScopeType == constants.AuthTypePassword {
		return fmt.Sprintf("%s_%s_%s_%s_%s", deviceId, "0", scopeConfig.ScopeId, "0", TimeStamp())
	}
	return fmt.Sprintf("%s_%s_%s", deviceId, "0", scopeConfig.ScopeId)
}

func OsName() string {
	return runtime.GOOS
}

func SdkInfo() map[string]string {
	f, err := os.Open("sdk_info")
	if err != nil {
		glog.Warning("read sdk info failed")
		return map[string]string{}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			glog.Warningf("close stream failed. err: %s", err.Error())
		}
	}(f)

	// 文件很小
	info := make(map[string]string)
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			glog.Warningf("read sdk info failed or end")
			break
		}
		line := string(b)
		if len(line) != 0 {
			parts := strings.Split(line, "=")
			info[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	return info
}

func GetServer() *model.ServerInfo {
	f, err := os.Open("server_info.txt")
	if err != nil {
		glog.Warning("read sdk info failed")
		return nil
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			glog.Warningf("close stream failed. err: %s", err.Error())
		}
	}(f)

	// 文件很小
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			glog.Warningf("read sdk info failed or end")
			break
		}
		if len(b) != 0 {
			var serverInfo model.ServerInfo
			err := json.Unmarshal(b, &serverInfo)
			if err != nil {
				glog.Warningf("unmarshal server info failed. err: %s", err)
			}
			return &serverInfo
		}
	}
	return nil
}

func SavePwd(server model.ServerInfo) error {
	marshal, err := json.Marshal(server)
	if err != nil {
		glog.Warningf("marshal server info failed. err: %s", err)
		return err
	}
	return ioutil.WriteFile("server_info.txt", marshal, fs.ModePerm)
}

// VerifyConnection 使用自定义方式校验证书，忽略主机名校验
func VerifyConnection(serverCaPool *x509.CertPool) func(stats tls.ConnectionState) error {
	return func(stats tls.ConnectionState) error {
		certificates := stats.PeerCertificates
		if len(certificates) == 0 {
			glog.Warningf("ssl cert is empty.")
			return errors.New("cert is not suit")
		}
		opts := x509.VerifyOptions{
			Roots:         serverCaPool,
			Intermediates: x509.NewCertPool(),
		}
		for _, cert := range certificates {
			opts.Intermediates.AddCert(cert)
		}
		for _, cert := range certificates {
			_, err := cert.Verify(opts)
			if err == nil {
				return nil
			}
		}
		return errors.New("cert is not suit")
	}
}
