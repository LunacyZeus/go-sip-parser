package csv_utils

import (
	"sip-parser/pkg/sip"
)

func ConvertRow(row *PcapCsv) (new_row *CostPcapCsv, err error) {
	row.ANI = sip.GetSipPart(row.ANI)
	row.DNIS = sip.GetSipPart(row.DNIS)

	new_row = &CostPcapCsv{
		CallId:        row.CallId,
		ANI:           row.ANI,
		DNIS:          row.DNIS,
		LRN:           row.LRN,
		RelatedCallId: row.RelatedCallId,
		InviteTime:    row.InviteTime,
		RingTime:      row.RingTime,
		AnswerTime:    row.AnswerTime,
		HangupTime:    row.HangupTime,
		Duration:      row.Duration, // in milliseconds
		InTrunkId:     row.InTrunkId,
		OutTrunkId:    row.OutTrunkId,
		InRate:        row.InRate,
		InRateID:      row.InRateID,
		InCost:        row.InCost,
		OutRate:       row.OutRate,
		OutRateID:     row.OutRateID,
		OutCost:       row.OutCost,
		SrcIP:         row.SrcIP,
		DestIP:        row.DestIP,
		Command:       row.Command,
		Result:        row.Result,
	}
	return
}
