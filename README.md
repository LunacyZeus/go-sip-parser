Usage: call_simulation caller_ip,caller_port,ani,dnis,[test_shaken],[signature]

wget -O go_sip_parser https://raw.githubusercontent.com/LunacyZeus/go-sip-parser/refs/heads/main/go_sip_parser_amd64_linux


./go_sip_parser_amd64_linux telnet --cip "68.68.120.215" --cport "5060" --ani "15033478582" --dnis "+17472292998"


call_simulation 69.85.136.200,5060,17083744774,+12243462475


call_simulation 88.151.130.19,5060,17862944218,+17252997974
call_simulation 88.151.132.30,5060,5482#+13236778193,9093237141

./go_sip_parser_amd64_linux telnet --cip "88.151.130.19" --cport "5060" --ani "17862944218" --dnis "+17252997974"

./go_sip_parser_amd64_linux telnet --cip "88.151.132.30" --cport "5060" --ani "5482#+13236778193" --dnis "9093237141"

./go_sip_parser_amd64_linux telnet --cip "208.79.54.183" --cport "5060" --ani "+16787784146" --dnis "+17065423030"

call_simulation 208.79.54.183,5060,+16787784146,+17065423030