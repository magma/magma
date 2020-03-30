package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.resources;

import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.filter.auth.TenantIdHolder;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.restclient.RestClientManager;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model.SearchResult;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model.WorkflowSummary;
import com.sun.jersey.api.client.ClientResponse;
import com.sun.jersey.api.client.GenericType;
import com.sun.jersey.core.util.MultivaluedMapImpl;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import java.io.InputStream;
import java.lang.invoke.MethodHandles;
import java.util.Arrays;
import javax.inject.Inject;
import javax.inject.Singleton;
import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.Consumes;
import javax.ws.rs.DefaultValue;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.MultivaluedMap;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Api(value = "/workflow", produces = MediaType.APPLICATION_JSON, consumes = MediaType.APPLICATION_JSON, tags = "Workflow Management")
@Path("/workflow")
@Produces({MediaType.APPLICATION_JSON})
@Consumes({MediaType.APPLICATION_JSON})
@Singleton
public class WorkflowResource {
    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());
    private static final String QUERY = "query";

    private final RestClientManager manager;
    private final TenantIdHolder tenantIdHolder;



    @Inject
    public WorkflowResource(RestClientManager manager, TenantIdHolder tenantIdHolder) {
        this.manager = manager;
        this.tenantIdHolder = tenantIdHolder;
    }

    @ApiOperation(value = "Search for workflows based on payload and other parameters",
            notes = "use sort options as sort=<field>:ASC|DESC e.g. sort=name&sort=workflowId:DESC." +
                    " If order is not specified, defaults to ASC.")
    @GET
    @Consumes(MediaType.WILDCARD)
    @Produces(MediaType.APPLICATION_JSON)
    @Path("/search")
    public /*SearchResult<WorkflowSummary>*/InputStream search(@QueryParam("start") @DefaultValue("0") int start,
                                                               @QueryParam("size") @DefaultValue("100") int size,
                                                               @QueryParam("sort") String sort,
                                                               @QueryParam("freeText") @DefaultValue("*") String freeText,
                                                               @QueryParam(QUERY) String query,
                                                               @Context HttpServletRequest req
                                                ) {
        MultivaluedMap<String, String> params = manager.copyParams(req);
        // modify request
        // TODO: properly sanitize query
        String queryPrefix = "workflowType STARTS_WITH '" + tenantIdHolder.get() + "'";
        if (StringUtils.isBlank(query)) {
            query = queryPrefix;
        } else {
            query = queryPrefix + " AND " + query;
        }
        params.putSingle(QUERY, query);
        ClientResponse cr = manager.executeRequest(req, params);

        if (cr.getStatus() >= 400) {
            throw new IllegalStateException("Status:" + cr.getStatus());
        }

        if (cr.getStatus() != 204 && cr.hasEntity()) {
            // TODO: response headers, status code
            return cr.getEntityInputStream();

//            // extract body
//            SearchResult<WorkflowSummary> entity = cr.getEntity(new GenericType<SearchResult<WorkflowSummary>>() {
//            });
//            // modify response
//            entity.setTotalHits(999);
//            return entity;
        } else {
            return null;
        }
    }
}
