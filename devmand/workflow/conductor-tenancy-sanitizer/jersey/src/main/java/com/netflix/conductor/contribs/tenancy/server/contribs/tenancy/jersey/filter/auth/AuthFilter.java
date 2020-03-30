package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.filter.auth;

import java.io.IOException;
import java.lang.invoke.MethodHandles;
import javax.inject.Inject;
import javax.inject.Singleton;
import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Expect that tenant id in x-auth-organization header
 */
@Singleton
public class AuthFilter implements Filter {
    private static final String X_AUTH_ORG = "x-auth-organization";
    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());

    private final TenantIdHolder tenantIdHolder;

    @Inject
    public AuthFilter(TenantIdHolder tenantIdHolder) {
        this.tenantIdHolder = tenantIdHolder;
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {

    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) throws IOException, ServletException {
        HttpServletRequest hReq = (HttpServletRequest) request;
        String tenantId = hReq.getHeader(X_AUTH_ORG);
        if (StringUtils.isNotBlank(tenantId)) {
            logger.debug("Setting tenant id to '{}'", tenantId);
            tenantIdHolder.set(tenantId);
            chain.doFilter(request, response);
            tenantIdHolder.set(null);
        } else {
            logger.error("Cannot find auth header");
            throw new AuthHeaderNotFoundException();
        }
    }

    @Override
    public void destroy() {

    }
}
