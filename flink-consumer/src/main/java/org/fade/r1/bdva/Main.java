package org.fade.r1.bdva;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.connector.kafka.source.reader.deserializer.KafkaRecordDeserializationSchema;
import org.apache.flink.streaming.api.datastream.DataStreamSource;
import org.apache.flink.streaming.api.datastream.SingleOutputStreamOperator;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
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
                .setBounded(OffsetsInitializer.latest())
                .build();
        // 本地执行
//        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        try (StreamExecutionEnvironment env = StreamExecutionEnvironment.createRemoteEnvironment(properties.getProperty("flink.host"),
                Integer.parseInt(properties.getProperty("flink.port")),
                "E:\\GitHub\\r1-BDVA\\flink-consumer\\target\\flink-consumer-1.0-SNAPSHOT.jar")) {
//                "D:\\apache-maven\\maven-repository\\org\\apache\\flink\\flink-connector-kafka\\3.4.0-1.20\\flink-connector-kafka-3.4.0-1.20.jar",
//                "D:\\apache-maven\\maven-repository\\org\\apache\\kafka\\kafka-clients\\3.4.0\\kafka-clients-3.4.0.jar")) {
//        try {
            DataStreamSource<String> kafka = env.fromSource(source, WatermarkStrategy.noWatermarks(), "kafka");
            SingleOutputStreamOperator<String> operator = kafka.map(String::toUpperCase);
            operator.print();
            env.execute();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

}
