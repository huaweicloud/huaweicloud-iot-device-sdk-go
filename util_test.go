package iot

import (
	"testing"
)

func TestTimeStamp(t *testing.T) {
	timeStamp := timestamp()
	if len(timeStamp) != 10 {
		t.Error(`Time Stamp length must be 10`)
	}
}

func TestDataCollectionTime(t *testing.T) {
	if len(GetEventTimeStamp()) != 16 {
		t.Errorf(`Data Collection Time length must be 16,but is %d`, len(GetEventTimeStamp()))
	}
}

func TestHmacSha256(t *testing.T) {
	encodedPassword := "c0fefa1341fb0647290e93f641a9bcea74cd32111668cdc5f7418553640a55cc"
	if hmacSha256("123456789", "202012222200") != encodedPassword {
		t.Errorf("encoded password must be %s but is %s", encodedPassword, hmacSha256("123456789", "202012222200"))
	}
}

func TestInterface2JsonString(t *testing.T) {
	if Interface2JsonString(nil) != "" {
		t.Errorf("nill interface to json string must empty")
	}
}

func TestGetTopicRequestId(t *testing.T) {
	topic := "$os/device/down/request=123456789"
	if getTopicRequestId(topic) != "123456789" {
		t.Errorf("topic request id must be %s", "123456789")
	}
}

func TestFormatTopic(t *testing.T) {
	topic := "$os/device/{device_id}/up"
	deviceId := "123"
	formatTopicName := "$os/device/123/up"
	if formatTopicName != formatTopic(topic, deviceId) {
		t.Errorf("formated topic must be %s", formatTopicName)
	}

}

// 仅适用于windows系统
func TestSmartFileName(t *testing.T) {
	fileName := "D/go/sdk/test.log"
	name := "D:\\go\\sdk\\test.log"

	if name != smartFileName(fileName) {
		t.Errorf("in windows file system,smart file name must be %s", name)
	}
}
