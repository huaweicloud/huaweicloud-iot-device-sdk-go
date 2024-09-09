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
	"github.com/golang/glog"
	"github.com/huaweicloud/huaweicloud-iot-device-sdk-go/iot/model"
	"strconv"
	"strings"
	"time"
)

func convertString2intList(str string) ([]int, error) {
	split := strings.Split(str, ":")
	var timeList []int
	for _, time := range split {
		atoi, err := strconv.Atoi(time)
		if err != nil {
			return nil, err
		}
		timeList = append(timeList, atoi)
	}
	return timeList, nil
}

func getRuleWeek(weekList []int) {
	for index, week := range weekList {
		curWeek := week - 1
		weekList[index] = curWeek
	}
}

func checkTimeRange(timeRange model.TimeRange) bool {
	startTime := timeRange.StartTime
	endTime := timeRange.EndTime
	weekStr := timeRange.DaysOfWeek
	if len(startTime) == 0 || len(endTime) == 0 || len(weekStr) == 0 {
		return true
	}
	now := time.Now().UTC()
	nowWeek := (now.Weekday() + 1) % 7
	weekList := strings.Split(weekStr, ",")
	// 开始结束时间分割
	beginTimeList := strings.Split(startTime, ":")
	endTimeList := strings.Split(endTime, ":")
	beginHour, err := strconv.Atoi(beginTimeList[0])
	if err != nil {
		glog.Warningf("start time format is invalid. startTime: %s", startTime)
		return false
	}
	beginMinute, err := strconv.Atoi(beginTimeList[1])
	if err != nil {
		glog.Warningf("start time format is invalid. startTime: %s", startTime)
		return false
	}

	endHour, err := strconv.Atoi(endTimeList[0])
	if err != nil {
		glog.Warningf("end time format is invalid. startTime: %s", startTime)
		return false
	}
	endMinute, err := strconv.Atoi(endTimeList[1])
	if err != nil {
		glog.Warningf("end time format is invalid. startTime: %s", startTime)
		return false
	}

	nowHour := now.Hour()
	nowMinute := now.Minute()

	//8:00 - 9:00 形式
	if (beginHour*60 + beginMinute) < (endHour*60 + endMinute) {
		return ((beginHour*60+beginMinute) <= (nowHour*60+nowMinute) && (nowHour*60+nowMinute) <= (endHour*60+endMinute)) && arraysContains(strconv.Itoa(int(nowWeek)), weekList)
	}
	// 23:00 -01:00形式， 处于23:00-00:00之间的形式
	if (beginHour*60+beginMinute) <= (nowHour*60+nowMinute) && (nowHour*60+nowMinute) <= (24*60+00) &&
		arraysContains(strconv.Itoa(int(nowWeek)), weekList) {
		return true
	} else if (nowHour*60 + nowMinute) <= (endHour*60 + endMinute) {
		nowWeek = nowWeek - 1
		if 0 == nowWeek {
		}
		return arraysContains(strconv.Itoa(int(nowWeek)), weekList)
	}
	return false
}

func arraysContains(str string, arrays []string) bool {
	for _, arr := range arrays {
		if strings.EqualFold(arr, str) {
			return true
		}
	}
	return false
}
