# (WIP) trans-broker

curl -XPOST -d '{"log_name":"tb-log","token":"123456789","client":"10.12.88.99","mid":"abc123","event_id":50101,"event_time":"2018-09-20 12:02:30.453","computer_name":"domacli-PC0.tencent.com","event_data":{"Image":"D\\Program Files (X86)\\QQMgr.exe","ProcessMd5":"testmd5123","ProcessId":1024,"CommandLine":"xxx.bat","Op":"click","SourceIp":"10.12.88.99","OpTime":"2018-09-20 12:02:30.453","FilePath":"D\\Program Files (X86)\\QQMgr.exe","FileMd5":"testmd5123"}}'  'http://127.0.0.1:28083/data/push'


curl -XPOST -d '{"message":"sss","log_name":"tb-log","token":"123456789","client":"10.12.88.99","mid":"abc123","event_id":50101,"event_time":"2018-09-20T12:02:30.453Z","computer_name":"domacli-PC0.tencent.com","event_data":{"Image":"D\\Program Files (X86)\\QQMgr.exe","ProcessMd5":"testmd5123","ProcessId":1024,"CommandLine":"xxx.bat","Op":"click","SourceIp":"10.12.88.99","FilePath":"D:\\Program Files (X86)\\QQMgr.exe","FileMd5":"testmd5125"}}'  'http://127.0.0.1:28083/data/push'