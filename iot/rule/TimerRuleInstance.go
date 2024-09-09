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
	"github.com/go-co-op/gocron"
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/callback"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type TimerRuleInstance struct {
	schedule *gocron.Scheduler
	jobMap   map[string]*gocron.Job
}

func (timerRule *TimerRuleInstance) submitRule(ruleInfo model.RuleInfo, handler callback.RuleActionHandler) {
	conditions := ruleInfo.Conditions
	for _, condition := range conditions {
		if strings.EqualFold("DAILY_TIMER", condition.Type) {
			timerRule.handlerTimerRule(condition, ruleInfo, handler)
		} else if strings.EqualFold("SIMPLE_TIMER", condition.Type) {
			timerRule.handlerSimpleRule(condition, ruleInfo, handler)
		}
	}
}

func (timerRule *TimerRuleInstance) handlerTimerRule(condition model.Condition, ruleInfo model.RuleInfo, handler callback.RuleActionHandler) {
	executeTime := condition.Time
	daysOfWeek := condition.DaysOfWeek
	if len(executeTime) == 0 || len(daysOfWeek) == 0 {
		glog.Warning("time or days of week is empty. time: %s, days of week: %s", executeTime, daysOfWeek)
		return
	}
	timeList, err := convertString2intList(executeTime)
	if err != nil || len(timeList) != 2 {
		glog.Warningf("time format is invalid. time: %s", timeList)
		return
	}
	weekList, err := convertString2intList(daysOfWeek)
	if err != nil || len(timeList) == 0 {
		glog.Warningf("week format is invalid. time: %s", timeList)
		return
	}

	time.Sleep(1 * time.Second)
	getRuleWeek(weekList)
	for _, week := range weekList {
		_, err := timerRule.schedule.Every(1).Weekday(time.Weekday(week)).At(executeTime).Do(func() {
			if checkTimeRange(ruleInfo.TimeRange) {
				handler(ruleInfo.Actions)
			}
		})
		if err != nil {
			glog.Warningf("create schedule failed. err: %s", err)
			continue
		}
		glog.Infof("add rule schedule daily timer job success.")
	}
}

func (timerRule *TimerRuleInstance) handlerSimpleRule(condition model.Condition, ruleInfo model.RuleInfo, handler callback.RuleActionHandler) {
	interval := condition.RepeatInterval
	count := condition.RepeatCount
	startTime, err := iot.GetDateTime(condition.StartTime)
	if err != nil {
		glog.Warningf("rule start time is invalid. time: %s, err: %s", condition.StartTime, err)
		return
	}
	jobId := uuid.NewV4().String()
	var jobStartTime time.Time
	job, err := timerRule.schedule.Every(interval).Second().LimitRunsTo(count).StartAt(startTime).Do(func() {
		if !checkTimeRange(ruleInfo.TimeRange) {
			return
		}
		currentJob, exist := timerRule.jobMap[jobId]
		if !exist {
			glog.Warningf("job not found. jobId: %s", jobId)
			return
		}
		nextTime := currentJob.LastRun()
		runCount := currentJob.RunCount()
		if runCount == 1 {
			jobStartTime = nextTime
		}
		if jobStartTime.After(startTime) {
			glog.Warningf("job was timeout. jobName: %s", jobId)
			return
		}
		handler(ruleInfo.Actions)
	})
	if err != nil {
		glog.Warningf("create schedule failed. err: %s", err)
		return
	}
	job.Name(jobId)
	timerRule.jobMap[jobId] = job
	glog.Infof("add rule schedule simple timer job success.")
}

func (timerRule *TimerRuleInstance) Start() {
	timerRule.schedule.StartAsync()
}

func (timerRule *TimerRuleInstance) ShutdownTimer() {
	timerRule.schedule.Stop()
}
