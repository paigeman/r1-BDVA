package config

const (
	ServerIP       = "127.0.0.1"         // 服务器IP
	ServerPort     = "8080"              // 服务器端口
	BufferSize     = 1024                // 缓冲区大小
	KafkaAddress   = "localhost:9092"    // Kafka服务器地址
	KafkaTopic     = "quickstart-events" // Kafka的topic
	KafkaPartition = 0                   // Kafka的partition
)
