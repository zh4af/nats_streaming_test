Test nats streaming server cluster HA
-------
搭建3个node的nats streaming server集群，各node的配置文件见*.conf，以及启动方式见.sh

客户端分别在3个node机器上以publisher, subscriber, subscriber角色启动
```
测试效果较满意，在3个node中随意断开一个node，cluster从新选举new leader，
推送服务短暂超时后（据观察异常事件在数秒内），业务程序仍能正常运行，消息的推送和订阅不受影响，数据不丢失。
```
```
剩下的两个nats streaming server node会不断重连丢失的node，应当监控nats streaming server的日志，
或者通过http monitor port监控nats streaming server的状态，及时发现cluster的异常。
```
```
当发现node异常断开后，从新启动异常的node，3个node会自动发现对方，恢复为稳定的3 node cluster，业务程序不受影响。
```
```
断掉两个node，此时cluster无法工作，业务程序受影响，消息无法推送和接受；
恢复其中一个或两个node，cluster可以恢复并同步状态，但不能无感知地恢复业务程序，消息的推送和订阅此时不通，需要手动重启推送和订阅服务。
```