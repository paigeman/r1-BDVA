package org.fade.r1.bdva;

import com.fasterxml.jackson.annotation.JsonProperty;

public record BusStatus (
        @JsonProperty("bus_id") String busId,
        double latitude,
        double longitude,
        double speed,
        long timestamp) {}
