# Robot For Higo 

### Feature
    
    提供 手机号+密码，配置本地Redis的IP:Port即可实现抓取Higo全球购Tab中最热的数据（其他Tab也可）,自动让你的账号进入该买手的群里
    
    因为发消息设计消息签名，无法伪造，故无法提供发送消息的业务
    
### QuickStart

    sh build.sh
    
    ./robot -mobile=12322 -password=12344 -redis=localhost:6379 -robotType=higo/ymt -message="hi"
    
### Note:
    
    可以通过绑定的redis中的zset结构获取当前已经发送广告的卖家ID
    
    洋码头  _ymt_send_message_
    
    Higo  _higo_group_followed
      
    
