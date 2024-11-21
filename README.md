Usage: call_simulation caller_ip,caller_port,ani,dnis,[test_shaken],[signature]

wget -O go_sip_parser https://raw.githubusercontent.com/LunacyZeus/go-sip-parser/refs/heads/main/go_sip_parser_amd64_linux


./go_sip_parser_amd64_linux telnet --cip "68.68.120.215" --cport "5060" --ani "15033478582" --dnis "+17472292998"


call_simulation 69.85.136.200,5060,17083744774,+12243462475

./go_sip_parser_amd64_linux telnet --cip "69.85.136.200" --cport "5060" --ani "17083744774" --dnis "+12243462475"
