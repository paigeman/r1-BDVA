package org.fade.r1.bdva;

import com.influxdb.client.write.Point;
import org.apache.flink.api.connector.sink.SinkWriter;
import org.apache.flink.streaming.connectors.influxdb.sink.writer.InfluxDBSchemaSerializer;

public class MySerializer implements InfluxDBSchemaSerializer<Long> {

    @Override
    public Point serialize(Long element, SinkWriter.Context context) {
        final Point dataPoint = new Point("test");
        dataPoint.addTag("longValue", String.valueOf(element));
        dataPoint.addField("fieldKey", "fieldValue");
        return dataPoint;
    }

}
