package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.resources;

import com.sun.jersey.spi.container.ContainerRequest;
import com.sun.jersey.spi.container.ContainerRequestFilter;
import java.lang.invoke.MethodHandles;
import javax.inject.Singleton;
import javax.ws.rs.ext.Provider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Provider
@Singleton
public class LogFilter implements ContainerRequestFilter {
    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());

    @Override
    public ContainerRequest filter(ContainerRequest request) {
        logger.debug("URL: {}", request.getRequestUri());
        return request;
    }
}
