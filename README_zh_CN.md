# r1-BDVA
研一课程“大数据可视化分析”的期末大作业

## 要求

> 模拟多个设备，具体场景自己讨论，比如车辆，每秒都在产生数据，服务端接收数据，udp协议就可以，建议使用Go语言，用中间件消息队列kafka，flink作kafka的消费者，kafka作基础的配置就行，flink处理好写至数据库，grafana显示数据库数据。

## 城市公交车实时运行状态监控系统

### 环境

- Java 17
- Docker
- WSL / Linux distro (已在 Kali Linux 中进行测试)
- Go 1.22.5
- Grafana 11.6.0

### 如何使用

首先运行脚本：

```shell
./start.sh up
```

把 `dashboard.json` 导入Grafana

然后运行flink-consumer目录下的Main.java程序（可能需要更改jar文件的路径）

### 业务场景

本系统模拟采集城市公交车的运行状态（如经纬度、速度、站点信息等），进行数据采集、清洗、分析与可视化

### 系统架构

客户端->服务端->kafka->flink->数据库(InfluxDB2)->Grafana

### 需求

#### 客户端模拟采集城市公交车的运行状态信息

要求收集：

- 经度
- 维度
- 速度

因为是模拟采集，可以随机生成上述信息

#### 客户端发送至服务端

以JSON格式发送，采用UDP协议

#### 服务端把数据写入Kafka

以JSON格式写入

#### Flink处理消费的消息

消费Kafka的消息，对JSON格式的消息反序列化，把数据写入InfluxDB

#### Grafana对数据进行可视化

比如显示：

- 车辆的运动轨迹
- 车辆的速度随时间的变化趋势
- ...

### 可视化结果

![](resources/images/Snipaste_2025-05-20_17-20-52.png)

![](resources/images/Snipaste_2025-05-20_17-21-39.png)

![](resources/images/Snipaste_2025-05-20_17-22-16.png)

![](resources/images/Snipaste_2025-05-20_17-22-33.png)

![](resources/images/Snipaste_2025-05-20_17-22-47.png)