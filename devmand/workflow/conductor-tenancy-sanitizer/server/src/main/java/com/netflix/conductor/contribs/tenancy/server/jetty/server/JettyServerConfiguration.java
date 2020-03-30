package com.netflix.conductor.contribs.tenancy.server.jetty.server;

public interface JettyServerConfiguration {
    String ENABLED_PROPERTY_NAME = "conductor.jetty.server.enabled";
    public static boolean ENABLED_DEFAULT_VALUE = true;

    String PORT_PROPERTY_NAME = "conductor.jetty.server.port";
    public static int PORT_DEFAULT_VALUE = 8081;

    String JOIN_PROPERTY_NAME = "conductor.jetty.server.join";
    public static boolean JOIN_DEFAULT_VALUE = true;

    default boolean isEnabled(){
        return ENABLED_DEFAULT_VALUE;
    }

    default int getPort() {
        return PORT_DEFAULT_VALUE;
    }

    default boolean isJoin(){
        return JOIN_DEFAULT_VALUE;
    }
}
