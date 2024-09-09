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
	"github.com/go-co-op/gocron"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/callback"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"strings"
	"time"
)

type RuleManageService struct {
	RuleIdList       map[string]bool
	RuleInfoMap      map[string]model.RuleInfo
	TimerRuleMap     map[string]TimerRuleInstance
	ConditionExecute ConditionExecute
}

func (ruleService *RuleManageService) ModifyRule(service model.DevicePropertyDownRequestEntry, ruleDelete callback.ReportRuleDelete) {
	if service.Properties == nil {
		return
	}
	var ruleProperties map[string]model.RuleVersion
	err := json.Unmarshal([]byte(iot.Interface2JsonString(service.Properties)), &ruleProperties)
	if err != nil {
		glog.Warning("modify rule failed. properties is : %s", service.Properties)
		return
	}
	var ruleIdDel []string
	for key, value := range ruleProperties {
		version := value.Version
		if version != -1 {
			if !ruleService.RuleIdList[key] {
				ruleService.RuleIdList[key] = true
			}
		} else {
			if ruleService.RuleIdList[key] {
				delete(ruleService.RuleIdList, key)
			}
			_, ok := ruleService.RuleInfoMap[key]
			if ok {
				delete(ruleService.RuleInfoMap, key)
			}
			timerRule, ok := ruleService.TimerRuleMap[key]
			if ok {
				timerRule.ShutdownTimer()
			}
			ruleIdDel = append(ruleIdDel, key)
		}
	}
	ruleService.reportRuleEvent(ruleIdDel, ruleDelete)
}

func (ruleService *RuleManageService) reportRuleEvent(ruleIdDel []string, ruleDelete callback.ReportRuleDelete) {
	deviceRuleEvent := model.DeviceRuleEvent{}
	deviceRuleEvent.ServiceId = "$device_rule"
	deviceRuleEvent.EventType = "device_rule_config_request"
	deviceRuleEvent.EventTime = iot.GetEventTimeStamp()

	params := model.DeviceRuleRequestEventParams{
		RuleIds: make([]string, 0),
		DelIds:  ruleIdDel,
	}
	for key := range ruleService.RuleIdList {
		params.RuleIds = append(params.RuleIds, key)
	}
	deviceRuleEvent.Paras = params
	var eventList []model.DeviceRuleEvent
	eventList = append(eventList, deviceRuleEvent)
	deviceEvents := model.DeviceEvents{
		Services: eventList,
	}
	b := ruleDelete(deviceEvents)
	glog.Infof("modify rule result: %v", b)
}

func (ruleService *RuleManageService) QueryRuleResponse(ruleInfos []model.RuleInfo, handler callback.RuleActionHandler) {
	if len(ruleInfos) <= 0 {
		glog.Warningf("rule info length is below 0.")
		return
	}
	for _, ruleInfo := range ruleInfos {
		_, ruleIdExist := ruleService.RuleIdList[ruleInfo.RuleId]
		oldRuleInfo, ruleInfoExist := ruleService.RuleInfoMap[ruleInfo.RuleId]
		if ruleInfo.Status == "inactive" {
			if ruleIdExist {
				delete(ruleService.RuleIdList, ruleInfo.RuleId)
			}
			if ruleInfoExist {
				delete(ruleService.RuleInfoMap, ruleInfo.RuleId)
			}
			continue
		}
		ruleVersion := ruleInfo.RuleVersionInShadow
		if ruleInfoExist && oldRuleInfo.RuleVersionInShadow >= ruleVersion {
			glog.Infof("rule version is not change. no need to refresh.")
			continue
		}
		ruleService.RuleInfoMap[ruleInfo.RuleId] = ruleInfo
		ruleService.submitTimerRule(ruleInfo, handler)
	}
}

func (ruleService *RuleManageService) HandleRule(services []model.DevicePropertyEntry, handler callback.RuleActionHandler) {
	for _, ruleInfo := range ruleService.RuleInfoMap {
		if !checkTimeRange(ruleInfo.TimeRange) {
			glog.Warningf("rule not match the time.")
			continue
		}
		conditions := ruleInfo.Conditions
		logic := ruleInfo.Logic
		if strings.EqualFold("or", logic) {
			for _, condition := range conditions {
				satisfied := ruleService.ConditionExecute.isConditionSatisfied(condition, services)
				if satisfied {
					handler(ruleInfo.Actions)
				}
			}
		} else if strings.EqualFold("and", logic) {
			isSatisfied := true
			for _, condition := range conditions {
				satisfied := ruleService.ConditionExecute.isConditionSatisfied(condition, services)
				if !satisfied {
					isSatisfied = false
				}
			}
			if isSatisfied {
				handler(ruleInfo.Actions)
			}
		} else {
			glog.Warningf("rule logic is not match. logic: %s", logic)
		}
	}
}

func (ruleService *RuleManageService) submitTimerRule(ruleInfo model.RuleInfo, handler callback.RuleActionHandler) {
	conditionList := ruleInfo.Conditions
	isTimerRule := false
	for _, condition := range conditionList {
		if strings.EqualFold("DAILY_TIMER", condition.Type) || strings.EqualFold("SIMPLE_TIMER", condition.Type) {
			isTimerRule = true
			break
		}
	}
	if !isTimerRule {
		return
	}
	if isTimerRule && len(conditionList) > 1 && strings.EqualFold("and", ruleInfo.Logic) {
		glog.Warningf("multy timer rule only support or logic. ruleId: %s", ruleInfo.RuleId)
		return
	}
	timerRule, exist := ruleService.TimerRuleMap[ruleInfo.RuleId]
	if exist {
		timerRule.ShutdownTimer()
		delete(ruleService.TimerRuleMap, ruleInfo.RuleId)
	}
	timerRule = TimerRuleInstance{
		schedule: gocron.NewScheduler(time.UTC),
		jobMap:   make(map[string]*gocron.Job),
	}
	timerRule.submitRule(ruleInfo, handler)
	timerRule.Start()
	ruleService.TimerRuleMap[ruleInfo.RuleId] = timerRule
}
