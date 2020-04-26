# custm-chat

在线客服系统 backend

## 实现

基于Golang(go1.11.1)实现，对外通过http提供服务。

## 依赖

* MySQL5.7

* IMServer [centrifugo](https://github.com/centrifugal/centrifugo)

* Redis

* 地理位置 [GeoIP](https://www.ipip.net/)

* 邮件服务 [sendgrid](https://sendgrid.com/)

## IM 相关

* 客户端
  [js](https://github.com/centrifugal/centrifuge-js)
  [go](https://github.com/centrifugal/gocent)
  
* 基于channel的通信方式。
  在这个项目里就是坐席和访客在分配完坐席完成之后，通过客户端提供的api 订阅相应的通道来实现消息的发送和接受。
  
* 使用的通信协议
  websocket
  sockjs
  
* 在线访客和坐席的同步
  每隔一定时间，客户端需要调用im server的刷新接口，来同步自己的状态。在refresh接口需要一个connection token,
  这个token里包含有访客和坐席的id, 便于backend server可以通过后台提供的接口查询在线的坐席和访客。