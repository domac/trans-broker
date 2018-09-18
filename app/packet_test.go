package app

import (
	"encoding/json"
	"testing"
)

var testData = `
{
	"token":"123456",
	"client_ip":"127.0.0.1",
	"mid":"a001",
	"guid":"g001",
	"event_time":"2017",
	"event_id":1001,
	"computer_name":"domac-PC0.tencent.com",
	"event_business":["testxxx"],
	"event_data": {
		"Image":"/tmp/test.exe",
		"ProcessMd5":"277623723hdhsd",
		"ProcessId":1024,
		"CommandLine":"sudo test",
		"Op":"hhh"
	}
}
`

func TestUnmarshalPacketJson(t *testing.T) {
	pd := PushData{}
	err := json.Unmarshal([]byte(testData), &pd)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	js, err := json.Marshal(pd.EventData)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	commonData := CommonData{}
	json.Unmarshal(js, &commonData)

	if commonData.Image != "/tmp/test.exe" {
		t.Fail()
	}
}
