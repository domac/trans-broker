package app

import (
	"encoding/json"
)

//上报数据结构
type PushData struct {
	Token         string      `json:"token"`
	ClientIp      string      `json:"client_ip"`
	Mid           string      `json:"mid"`
	Guid          string      `json:"guid"`
	EventTime     string      `json:"event_time"`
	EventId       int         `json:"event_id"`
	EventName     int         `json:"event_name"`
	ComputerName  string      `json:"computer_name"`
	EventBusiness []string    `json:"event_business"`
	EventData     interface{} `json:"event_data"`
}

//event_data 通用结构
type CommonData struct {
	Image       string `json:"Image"`
	ProcessMd5  string `json:"ProcessMd5"`
	ProcessId   int    `json:"ProcessId"`
	CommandLine string `json:"CommandLine"`
}

func (p *PushData) Marshal() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return b
}

type Packet struct {
	body []byte
}

func NewPacket(b []byte) *Packet {
	return &Packet{
		body: b,
	}
}
