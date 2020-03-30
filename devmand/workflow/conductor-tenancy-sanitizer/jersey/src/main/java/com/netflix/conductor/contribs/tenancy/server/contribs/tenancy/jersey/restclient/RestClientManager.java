package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.restclient;

import com.sun.jersey.api.client.Client;
import javax.inject.Inject;
import javax.inject.Singleton;

@Singleton
public class RestClientManager {

	static final int DEFAULT_READ_TIMEOUT = 150;
	static final int DEFAULT_CONNECT_TIMEOUT = 100;

	private final ThreadLocal<Client> threadLocalClient;
	private final int defaultReadTimeout;
	private final int defaultConnectTimeout;

	@Inject
	public RestClientManager() {
		this.threadLocalClient = ThreadLocal.withInitial(Client::create);
		this.defaultReadTimeout = DEFAULT_READ_TIMEOUT;
		this.defaultConnectTimeout = DEFAULT_CONNECT_TIMEOUT;
	}

	public Client getClient() {
		Client client = threadLocalClient.get();
		client.setReadTimeout(defaultReadTimeout);
		client.setConnectTimeout(defaultConnectTimeout);
		return client;
	}
}
