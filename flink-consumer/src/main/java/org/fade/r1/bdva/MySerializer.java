package org.fade.r1.bdva;

import com.influxdb.client.domain.WritePrecision;
import com.influxdb.client.write.Point;
import org.apache.flink.api.connector.sink.SinkWriter;
import org.apache.flink.streaming.connectors.influxdb.sink.writer.InfluxDBSchemaSerializer;

public class MySerializer implements InfluxDBSchemaSerializer<BusStatus> {

    @Override
    public Point serialize(BusStatus element, SinkWriter.Context context) {
        final Point dataPoint = new Point("BusStatus");
        dataPoint.addTag("BusId", element.busId());
        dataPoint.addField("Latitude", element.latitude());
        dataPoint.addField("Longitude", element.longitude());
        dataPoint.addField("Speed", element.speed());
        dataPoint.time(element.timestamp(), WritePrecision.NS);
        return dataPoint;
    }

}
