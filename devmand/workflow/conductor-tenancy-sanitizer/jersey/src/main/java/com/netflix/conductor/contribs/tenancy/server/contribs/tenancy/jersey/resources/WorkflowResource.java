package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.resources;

import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.restclient.RestClientManager;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model.SearchResult;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model.WorkflowSummary;
import com.sun.jersey.api.client.Client;
import com.sun.jersey.api.client.ClientResponse;
import com.sun.jersey.api.client.GenericType;
import com.sun.jersey.api.client.WebResource.Builder;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import java.lang.invoke.MethodHandles;
import javax.inject.Inject;
import javax.inject.Singleton;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.ws.rs.Consumes;
import javax.ws.rs.DefaultValue;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Api(value = "/workflow", produces = MediaType.APPLICATION_JSON, consumes = MediaType.APPLICATION_JSON, tags = "Workflow Management")
@Path("/workflow")
@Produces({MediaType.APPLICATION_JSON})
@Consumes({MediaType.APPLICATION_JSON})
@Singleton
public class WorkflowResource {
    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());
    private final RestClientManager manager;

    @Inject
    public WorkflowResource(RestClientManager manager) {
        this.manager = manager;
    }

    @ApiOperation(value = "Search for workflows based on payload and other parameters",
            notes = "use sort options as sort=<field>:ASC|DESC e.g. sort=name&sort=workflowId:DESC." +
                    " If order is not specified, defaults to ASC.")
    @GET
    @Consumes(MediaType.WILDCARD)
    @Produces(MediaType.APPLICATION_JSON)
    @Path("/search")
    public SearchResult<WorkflowSummary> search(@QueryParam("start") @DefaultValue("0") int start,
                                                @QueryParam("size") @DefaultValue("100") int size,
                                                @QueryParam("sort") String sort,
                                                @QueryParam("freeText") @DefaultValue("*") String freeText,
                                                @QueryParam("query") String query,
                                                @Context HttpServletRequest req,
                                                @Context HttpServletResponse res
                                                ) {

        String remoteAddress = req.getRemoteAddr();

        logger.debug("Remote address:{}", remoteAddress);

        res.addHeader("foo", "bar");

        // do http request...
        Client client = manager.getClient();
        String uri = "http://localhost:8080/api/workflow/search";
        String contentType = "application/json";
        Builder builder = client.resource(uri).type(contentType);
        // headers.forEach(builder::header);
        String accept = contentType;
        String method = "GET";
        // TODO modify request
        ClientResponse cr = builder.accept(accept).method(method, ClientResponse.class);

        if (cr.getStatus() >= 400) {
            throw new IllegalStateException("Status:" + cr.getStatus());
        }

        if (cr.getStatus() != 204 && cr.hasEntity()) {
            // extract body
            SearchResult<WorkflowSummary> entity = cr.getEntity(new GenericType<SearchResult<WorkflowSummary>>() {
            });
            // modify response
            entity.setTotalHits(999);
            return entity;
        } else {
            return null;
        }
    }
}
