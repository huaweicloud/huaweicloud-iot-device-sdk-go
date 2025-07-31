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

package file

import (
	"bytes"
	"crypto/tls"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 仅用于设备上传文件
type HttpClient interface {
	UploadFile(filename, uri string) bool
	DownloadFile(filename, uri, token string) bool
	OTADownloadFile(upgradeType byte, filename, uri, token string) bool
}

type httpClient struct {
	client *http.Client
}

func (client *httpClient) OTADownloadFile(upgradeType byte, fileName, downloadUrl, token string) bool {
	glog.Infof("begin to download file %s, url = %s", fileName, downloadUrl)
	fileName = iot.SmartFileName(fileName)

	originalUri, err := url.ParseRequestURI(downloadUrl)
	if err != nil {
		glog.Errorf("parse request uri failed %v", err)
		return false
	}
	if strings.Contains(downloadUrl, "https") {
		client.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	request, err := http.NewRequest("GET", downloadUrl, nil)

	if err != nil {
		glog.Errorf("download file request failed %v", err)
		return false
	}
	if upgradeType < 2 {
		request.Header.Set("Content-Type", "text/plain")
	}
	request.Header.Set("Host", originalUri.Host)
	if len(token) != 0 && upgradeType < 2 {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.client.Do(request)
	if err != nil {
		glog.Errorf("download file request failed %v", err)
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("download file request failed %v", err)
		return false
	}
	err = ioutil.WriteFile(fileName, body, fs.ModePerm)
	if err != nil {
		glog.Errorf("write file failed")
		return false
	}

	return true
}

func (client *httpClient) DownloadFile(fileName, downloadUrl, token string) bool {
	return client.OTADownloadFile(0, fileName, downloadUrl, token)
}

func (client *httpClient) UploadFile(filename, uri string) bool {
	filename = iot.SmartFileName(filename)
	fileBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		glog.Errorf("read file failed %v", err)
		return false
	}

	originalUri, err := url.ParseRequestURI(uri)
	if err != nil {
		glog.Errorf("parse request uri failed %v", err)
		return false
	}

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(fileBytes))

	if err != nil {
		glog.Errorf("upload file request failed %v", err)
		return false
	}

	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Host", originalUri.Host)

	if err != nil {
		glog.Errorf("upload request failed %v", err)
	}
	resp, err := client.client.Do(request)
	if err != nil {
		glog.Errorf("download file request failed %v", err)
		return false
	}
	return resp.StatusCode == 200
}

func CreateHttpClient() HttpClient {
	client := &http.Client{}

	httpClient := &httpClient{
		client: client,
	}

	return httpClient

}
