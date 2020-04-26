# 要订阅的channel

## 访客端

* {ent_id}_visitor_queue 访客的排队信息

* {trace_id}_conv_instruction 访客接受对话相关信令的通道,(关闭/转接/撤回)

* {trace_id} 以trace_id为名称的通道，访客通过订阅这个通道来接受自己和坐席发来的消息。

## 坐席端

* {conversation_id}_internal_messages 在参与同事的对话的时候，内部消息都会发送到这个通道里面

* {ent_id}_online_agents 在坐席分配的时候可以查询到在线的坐席，在坐席点击隐身的时候取消对这个通道的订阅

* {trace_id} 以trace_id为名称的通道，坐席通过订阅这个通道来接受自己和访客发来的消息。