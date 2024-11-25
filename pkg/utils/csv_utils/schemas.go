package csv_utils

type NotUsed struct {
	Name string
}

type PcapCsv struct { // Our example struct, you can use "-" to ignore a field
	CallId        string `csv:"Call-ID"`
	ANI           string `csv:"ANI"`
	DNIS          string `csv:"DNIS"`
	Via           string `csv:"Via"`
	RelatedCallId string `csv:"RelatedCallId"`
	OutVia        string `csv:"OutVia"`
	InviteTime    string `csv:"Invite Time"`
	RingTime      string `csv:"Ring Time"`
	AnswerTime    string `csv:"Answer Time"`
	HangupTime    string `csv:"Hangup Time"`
	Duration      string `csv:"Duration (msec)"`
	InTrunkId     string `csv:"InTrunkId"`
	OutTrunkId    string `csv:"OutTrunkId"`
	InRate        string `csv:"InRate"`
	InRateID      string `csv:"InRate Id"`
	InCost        string `csv:"InCost"`
	OutRate       string `csv:"OutRate"`
	OutRateID     string `csv:"OutRate Id"`
	OutCost       string `csv:"OutCost"`
	Command       string `csv:"Command"`
	Result        string `csv:"Result"`
	//NotUsedString string  `csv:"-"`
}
