sip-parser is a tool for extract SIP calls message flows from PCAP file.

**How to build**  
make sure your server have go environment  
if you dont have this, please see this tutorial  
https://linuxize.com/post/how-to-install-go-on-centos-7/  
go version is **go1.23.4**

git clone http://stash.denovolab.com/scm/~zeusho/sip-parser.git  
clone the repo into your server

enter the code dir  
use command  
`go build -o sip-parser`  
the command is build project with go  
then you will get a binary file named sip-parser  
start run the program  
`./sip-parser`

**How to use**

**Parse Call from pcap file**
For one pcap file  
`./sip_parser load -f data/test/202411140020.pcap`
when it finish you will get two files
in_data_test_202411140020.pcap.csv
out_data_test_202411140020.pcap.csv

For dir with pcap files  
`./sip_parser load -f data/test`  
when it finish you will get two files
in_test.csv
out_test.csv

**call_simulation from csv file**  
`./sip-parser get_cost -f in_test.csv -t 30`  
in_test.csv is in csv file path  
30 is telent threads number


**convert csv file**  
`./sip-parser convert_csv -f res_in.csv`  
you will get res_res_in.csv  
here is a final result file.  