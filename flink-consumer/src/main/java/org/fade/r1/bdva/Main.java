package org.fade.r1.bdva;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.connector.kafka.source.reader.deserializer.KafkaRecordDeserializationSchema;
import org.apache.flink.streaming.api.datastream.DataStreamSource;
import org.apache.flink.streaming.api.datastream.SingleOutputStreamOperator;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.connectors.influxdb.sink.InfluxDBSink;
import org.apache.kafka.common.TopicPartition;
import org.apache.kafka.common.serialization.StringDeserializer;

import java.io.IOException;
import java.io.InputStream;
import java.util.HashSet;
import java.util.Properties;
import java.util.Set;

public class Main {

    public static void main(String[] args) {
        // 读取配置
        Properties properties = new Properties();
        InputStream config = Main.class.getClassLoader().getResourceAsStream("configuration.properties");
        try {
            properties.load(config);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        Set<TopicPartition> partitionSet = new HashSet<>();
        partitionSet.add(new TopicPartition(properties.getProperty("kafka.topic"), Integer.parseInt(properties.getProperty("kafka.partition"))));
        KafkaSource<String> source = KafkaSource.<String>builder()
                .setBootstrapServers(properties.getProperty("kafka.address"))
                .setPartitions(partitionSet)
                .setStartingOffsets(OffsetsInitializer.earliest())
                .setDeserializer(KafkaRecordDeserializationSchema.valueOnly(StringDeserializer.class))
//                .setBounded(OffsetsInitializer.latest())
                .build();
        InfluxDBSink<BusStatus> influxDBSink = InfluxDBSink.builder()
                .setInfluxDBSchemaSerializer(new MySerializer())
                .setInfluxDBUrl(properties.getProperty("influxdb2.host"))
//                .setInfluxDBUsername(properties.getProperty("influxdb2.username"))
//                .setInfluxDBPassword(properties.getProperty("influxdb2.password"))
                .setInfluxDBBucket(properties.getProperty("influxdb2.bucket"))
                .setInfluxDBOrganization(properties.getProperty("influxdb2.organization"))
                .setInfluxDBToken(properties.getProperty("influxdb2.token"))
                .build();
        // 本地执行
//        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        try (StreamExecutionEnvironment env = StreamExecutionEnvironment.createRemoteEnvironment(properties.getProperty("flink.host"),
                Integer.parseInt(properties.getProperty("flink.port")),
                "E:\\GitHub\\r1-BDVA\\flink-consumer\\target\\flink-consumer-1.1.3.jar")) {
//                "D:\\apache-maven\\maven-repository\\org\\apache\\flink\\flink-connector-kafka\\3.4.0-1.20\\flink-connector-kafka-3.4.0-1.20.jar",
//                "D:\\apache-maven\\maven-repository\\org\\apache\\kafka\\kafka-clients\\3.4.0\\kafka-clients-3.4.0.jar")) {
//        try {
            // 与task slot有关
            env.setParallelism(2);
            DataStreamSource<String> kafka = env.fromSource(source, WatermarkStrategy.noWatermarks(), "kafka");
            SingleOutputStreamOperator<BusStatus> operator = kafka.map(json -> {
                ObjectMapper objectMapper = new ObjectMapper();
                return objectMapper.readValue(json, BusStatus.class);
            });
            operator.sinkTo(influxDBSink);
            env.execute();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

}
