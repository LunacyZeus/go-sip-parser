import csv

# 假设文件名为 'data.csv'
file_path = 'res1.csv'

# 打开文件并读取
with open(file_path, mode='r', newline='', encoding='utf-8') as file:
    reader = csv.DictReader(file, delimiter=',')  # 用逗号作为分隔符
    for row in reader:
        #dict_keys(['Call-ID', 'ANI', 'DNIS', 'LRN', 'Via', 'RelatedCallId', 'OutVia', 'Invite Time', 'Ring Time', 'Answer Time', 'Hangup Time', 'Duration (msec)', 'InTrunkId', 'OutTrunkId', 'InRate', 'InRate Id', 'InCost', 'OutRate', 'OutRate Id', 'OutCost', 'SrcIP', 'DestIP', 'Command^C', 'Result'])
        Call_ID = row.get('Call-ID', '')
        ANI = row.get('ANI', '')
        DNIS = row.get('DNIS', '')
        LRN = row.get('LRN', '')
        Via = row.get('Via', '')
        RelatedCallId = row.get('RelatedCallId', '')
        OutVia = row.get('OutVia', '')
        Invite_Time = row.get('Invite Time', '')
        Ring_Time = row.get('Ring Time', '')
        Answer_Time = row.get('Answer Time', '')
        Hangup_Time = row.get('Hangup Time', '')
        Duration_msec = row.get('Duration (msec)', '')
        InTrunkId = row.get('InTrunkId', '')
        OutTrunkId = row.get('OutTrunkId', '')
        InRate = row.get('InRate', '')
        InRate_Id = row.get('InRate Id', '')
        InCost = row.get('InCost', '')
        OutRate = row.get('OutRate', '')
        OutRate_Id = row.get('OutRate Id', '')
        OutCost = row.get('OutCost', '')
        SrcIP = row.get('SrcIP', '')
        DestIP = row.get('DestIP', '')
        Command = row.get('Command^C', '')
        Result = row.get('Result', '')

        print(f"Call-ID: {Call_ID},ANI:{ANI}")
