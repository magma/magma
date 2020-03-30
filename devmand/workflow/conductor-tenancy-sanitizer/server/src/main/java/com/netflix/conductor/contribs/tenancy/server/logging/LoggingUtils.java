package com.netflix.conductor.contribs.tenancy.server.logging;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.slf4j.bridge.SLF4JBridgeHandler;

public class LoggingUtils {
    public static void setUpJULBridge() {
        SLF4JBridgeHandler.removeHandlersForRootLogger();
        SLF4JBridgeHandler.install();
        // https://stackoverflow.com/questions/9117030/jul-to-slf4j-bridge
        Logger.getLogger("").setLevel(Level.FINEST);
    }
}
