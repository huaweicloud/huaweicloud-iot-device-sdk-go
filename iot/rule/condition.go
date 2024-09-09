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

package rule

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"strconv"
	"strings"
)

type ConditionExecute struct {
}

func (execute *ConditionExecute) isConditionSatisfied(condition model.Condition, services []model.DevicePropertyEntry) bool {
	if !strings.EqualFold(condition.Type, "DEVICE_DATA") {
		return false
	}
	deviceInfo := condition.DeviceInfo
	pathArr := strings.Split(deviceInfo.Path, "/")
	if len(pathArr) == 0 {
		glog.Warning("rule condition path is invalid. path: %s", deviceInfo.Path)
		return false
	}
	serviceIdPath := pathArr[0]
	propertyPath := pathArr[len(pathArr)-1]
	operate := condition.Operator
	if strings.EqualFold(operate, ">") {
		return execute.operationMoreThan(condition.Value, serviceIdPath, propertyPath, services)
	}
	if strings.EqualFold(operate, ">=") {
		return execute.operationMoreEqual(condition.Value, serviceIdPath, propertyPath, services)
	}
	if strings.EqualFold(operate, "<") {
		return execute.operationLessThan(condition.Value, serviceIdPath, propertyPath, services)
	}
	if strings.EqualFold(operate, "<=") {
		return execute.operationLessEqual(condition.Value, serviceIdPath, propertyPath, services)
	}
	if strings.EqualFold(operate, "=") {
		return execute.operationEquals(condition.Value, serviceIdPath, propertyPath, services)
	}
	if strings.EqualFold(operate, "between") {
		return execute.operationBetween(condition.Value, serviceIdPath, propertyPath, services)
	}
	if strings.EqualFold(operate, "in") {
		return execute.operationIn(condition.InValue, serviceIdPath, propertyPath, services)
	}
	glog.Warningf("operate is not match. operate: %s", operate)
	return false
}

func (execute *ConditionExecute) operationMoreThan(value, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		ok, target, current := getCompareFloatValue(value, propertyPath, properties)
		if ok && strings.EqualFold(serviceId, service.ServiceId) && current > target {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, target)
			return true
		}
	}
	return false
}

func (execute *ConditionExecute) operationMoreEqual(value, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		ok, target, current := getCompareFloatValue(value, propertyPath, properties)
		if ok && strings.EqualFold(serviceId, service.ServiceId) && current >= target {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, target)
			return true
		}
	}
	return false
}

func (execute *ConditionExecute) operationLessThan(value, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		ok, target, current := getCompareFloatValue(value, propertyPath, properties)
		if ok && strings.EqualFold(serviceId, service.ServiceId) && current < target {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, target)
			return true
		}
	}
	return false
}

func (execute *ConditionExecute) operationLessEqual(value, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		ok, target, current := getCompareFloatValue(value, propertyPath, properties)
		if ok && strings.EqualFold(serviceId, service.ServiceId) && current <= target {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, target)
			return true
		}
	}
	return false
}

func (execute *ConditionExecute) operationEquals(value, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		if len(value) == 0 {
			glog.Warningf("rule condition value is invalid. value: %s", value)
			continue
		}
		ok, target, current := getCompareFloatValue(value, propertyPath, properties)
		if ok && strings.EqualFold(serviceId, service.ServiceId) && current == target {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, value)
			return true
		} else {
			propertyMap := make(map[string]string)
			if json.Unmarshal([]byte(iot.Interface2JsonString(properties)), &propertyMap) != nil {
				continue
			}
			current, exist := propertyMap[propertyPath]
			if !exist {
				continue
			}
			if exist && strings.EqualFold(serviceId, service.ServiceId) && strings.EqualFold(current, value) {
				glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, value)
				return true
			}
		}
	}
	return false
}

func (execute *ConditionExecute) operationBetween(value, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		valueList := strings.Split(value, ",")
		if len(valueList) != 2 {
			glog.Warningf("rule condition value is invalid. value: %s", value)
			continue
		}
		ok, target, current := getCompareFloatValue(valueList[0], propertyPath, properties)
		target2, err := strconv.ParseFloat(valueList[1], 64)
		if err != nil {
			glog.Warningf("condition value convert to float failed. value: %s", value)
			continue
		}
		if ok && strings.EqualFold(serviceId, service.ServiceId) && current >= target && current <= target2 {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, target)
			return true
		}
	}
	return false
}

func (execute *ConditionExecute) operationIn(value []string, serviceId, propertyPath string, services []model.DevicePropertyEntry) bool {
	for _, service := range services {
		properties := service.Properties
		if len(value) == 0 {
			glog.Warningf("rule condition value is invalid. value: %s", value)
			continue
		}
		ok, _, current := getCompareFloatValue(value[0], propertyPath, properties)
		if ok && strings.EqualFold(serviceId, service.ServiceId) && valueInList(current, value) {
			glog.Infof("match condition for service. service_id: %s, value: %s", serviceId, value)
			return true
		}
	}
	return false
}

func getCompareFloatValue(value, propertyPath string, properties interface{}) (bool, float64, float64) {
	propertyMap := make(map[string]float64)
	if json.Unmarshal([]byte(iot.Interface2JsonString(properties)), &propertyMap) != nil {
		return false, 0, 0
	}
	target, err := strconv.ParseFloat(value, 64)
	if err != nil {
		glog.Warningf("condition value convert to float failed. value: %s", value)
		return false, 0, 0
	}
	current, exist := propertyMap[propertyPath]
	if !exist {
		return false, 0, 0
	}
	return true, target, current
}

func valueInList(current float64, valueList []string) bool {
	for _, value := range valueList {
		target, err := strconv.ParseFloat(value, 64)
		if err != nil {
			glog.Warningf("condition value convert to float failed. value: %s", value)
			continue
		}
		if target == current {
			return true
		}
	}
	return false
}
