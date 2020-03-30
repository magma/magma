package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.filter.auth;

import com.google.common.base.Preconditions;
import javax.inject.Singleton;

@Singleton
public class TenantIdHolder {
    private final ThreadLocal<String> tenantIdHolder = new ThreadLocal<>();

    void set(String tenantId) {
        tenantIdHolder.set(tenantId);
    }

    public String get() {
        return Preconditions.checkNotNull(tenantIdHolder.get(), "TenantId should not be null");
    }
}
