package tests

import (
	"sip-parser/pkg/sip"
	"testing"
)

func TestParseAniDnis(t *testing.T) {
	ani := "<sip:18562149721;rn=16096359982;npdi=yes@172.241.26.20>"
	dnis := "<sip:16073647476@192.40.216.65>;tag=sansay2051231311rdb49831"
	aniSip := sip.GetSipPart(ani)
	dnisSip := sip.GetSipPart(dnis)
	t.Logf("%s->%s %s->%s", ani, aniSip, dnis, dnisSip)

	ani = "<sip:2603302130@172.241.26.21>"
	dnis = "\"+12602762470\" <sip:+12602762470@207.223.71.229>;tag=gK082c9c8f"
	aniSip = sip.GetSipPart(ani)
	dnisSip = sip.GetSipPart(dnis)
	t.Logf("%s->%s %s->%s", ani, aniSip, dnis, dnisSip)
}
