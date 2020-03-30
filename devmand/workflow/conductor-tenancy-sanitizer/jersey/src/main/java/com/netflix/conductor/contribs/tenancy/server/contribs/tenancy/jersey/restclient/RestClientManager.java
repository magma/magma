package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.restclient;

import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.config.ProxyConfiguration;
import com.sun.jersey.api.client.Client;
import com.sun.jersey.api.client.ClientResponse;
import com.sun.jersey.api.client.WebResource;
import com.sun.jersey.api.client.WebResource.Builder;
import com.sun.jersey.api.client.filter.LoggingFilter;
import com.sun.jersey.core.util.MultivaluedMapImpl;
import java.lang.invoke.MethodHandles;
import java.util.ArrayList;
import java.util.Arrays;
import javax.inject.Inject;
import javax.inject.Singleton;
import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.MultivaluedMap;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Singleton
public class RestClientManager {
    private static final java.util.logging.Logger julClientLogger =
        java.util.logging.Logger.getLogger(MethodHandles.lookup().lookupClass().getName() + ".httpclient");
    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());
    static final int DEFAULT_READ_TIMEOUT = 150;
    static final int DEFAULT_CONNECT_TIMEOUT = 100;

    private final ThreadLocal<Client> threadLocalClient;
    private final int defaultReadTimeout;
    private final int defaultConnectTimeout;
    private final ProxyConfiguration proxyConfiguration;

    @Inject
    public RestClientManager(ProxyConfiguration proxyConfiguration) {
        this.proxyConfiguration = proxyConfiguration;
        this.threadLocalClient = ThreadLocal.withInitial(Client::create);
        this.defaultReadTimeout = DEFAULT_READ_TIMEOUT;
        this.defaultConnectTimeout = DEFAULT_CONNECT_TIMEOUT;
    }

    public Client getClient() {
        Client client = threadLocalClient.get();
        client.setReadTimeout(defaultReadTimeout);
        client.setConnectTimeout(defaultConnectTimeout);
        client.addFilter(new LoggingFilter(julClientLogger));
        return client;
    }

    public MultivaluedMap<String, String> copyParams(HttpServletRequest req) {
        MultivaluedMap<String, String> params = new MultivaluedMapImpl();
        // wrap in ArrayList because replacing values includes clearing the list
        req.getParameterMap().forEach((key, vals) -> params.put(key, new ArrayList<>(Arrays.asList(vals))));
        return params;
    }

    public ClientResponse executeRequest(HttpServletRequest req, MultivaluedMap<String, String> params) {
        String uri = proxyConfiguration.getBaseConductorURL() + req.getRequestURI();
        WebResource webResource = getClient().resource(uri);
        webResource = webResource.queryParams(params);

        Builder builder = webResource.getRequestBuilder()
            .type(MediaType.APPLICATION_JSON)
            .accept(MediaType.APPLICATION_JSON);
        return builder.method(req.getMethod(), ClientResponse.class);
    }
}
