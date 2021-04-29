---
id: nms_arch_overview
title: Overview
hide_title: true
---
# Overview

Magma’s NMS provides a single pane of glass for managing Magma based networks. NMS provides the ability to configure gateways and associated eNodeBs, provides visibility into status, events and metrics observed in these networks and finally
ability to configure and receive alerts.

## Big Picture

![nms](assets/nms/userguide/nms.png)

### NMS UI

NMS UI is a React App that uses hooks to manage the internal state. The NMS UI contains various react components to implement NMS page, admin page and the master page. NMS UI primarily uses most of the front end components from Material UI framework.

**Master page**

Master page contains components for displaying organizations, features, metrics and users in the master portal.

**NMS Page**

NMS page contains components for displaying dashboard, equipment, networks, policies/APNs, call tracing, subscribers, metrics and alerts. It also contains a network selector button to switch between various networks.

**Admin Page**

Admin page consists of UI components providing the ability to add networks and users to the respective organizations.


### Magmalte

Magmalte is a microservice built using express framework. It contains set of application and router level middlewares.
It uses sequelize ORM to connect to the NMS DB for servicing any routes involving DB interaction.

We will describe the middlewares used in the magmalte app below.

**AuthenticationMiddleware**

Magmalte uses passport authentication middleware for authenticating incoming requests. Passport authentication within magmalte currently supports

* local strategy(where username and password is validated against the information stored in db)

* SAML strategy(for enabling enterprise SSO). For enabling this, operator will have to configure the organization entrypoint, issuer name and cert.


**SessionMiddleware**

Session middleware helps with creating the session, setting the session cookie(use secure cookie in prod deployments) and creating the session object in req. The session token(from environment or hardcoded value) is used to sign the session cookie.  The session information is stored in the NMS DB through sequelize.  We only serialize userID as part of the session information.

**csrfMiddleware**

This protects against any CSRF attacks.  For more understanding, check this (https://github.com/pillarjs/understanding-csrf#csrf-tokens). Magmalte server includes the csrf tokens when it responds to the client and client submits the form with the token. (nms/app/packages/magmalte/app/common/axiosConfig.js)


**appMiddleware** sets some generic app level configuration. Including expecting the body to be json, body not to exceed 1mb and adds compression to the responses.

**userMiddleware**

Router middlewares to handle login routes for the user

**masterOrgMiddleware**

Router middleware to ensure that only requests from “master” organization is handled by master routes.

**Notes on Routes handling**

**apiControllerRouter**

All Orc8r API calls are handled by apiControllerRouter. ApiControllerRouter acts as a proxy and makes calls to orc8r using th e admin certs present on the container. ApiController Router has networkID filter to ensure that the networkID invoked in the request is part of the organization making the request. On the response side, there is a similar decorator to ensure that networkIDs which belong to the organization is only passed. Additionally, the apicontroller router has a audit log decorator which logs the mutating requests in a auditlogtable in the NMS DB.

**networkRouter**

Network creation calls are handled by this router and the network is added to the organization and passed on to the apiController.

**grafana router**

Grafana router is used for displaying the grafana component within iframe in UI. It handles all requests and proxies it to the underlying grafana service. It attempts to synchronize default grafana state when the route is invoked. The default state includes syncing the tenants, user information, datasource and dashboards. The datasource URL is a Orc8r tenants URL(magma/v1/tenants/<org_name>/metrics). NMS certs is also passed along with the datasource information so that grafana can use this to query orc8r and retrieve relevant metrics.
The default state synchronized helps display set of default dashboards to the users. The default dashboards include dashboards for gateways, subscribers, internal metrics etc.
Finally grafana router uses the tenant information synced to ensure that we only retrieve information for networks associated with a particular organization.


### NMS DB

NMS DB is the underlying SQL store for magmalte service. NMS DB contains following tables with associated schema

* Users

```
    id
    email character varying(255),
    organization character varying(255),
    password character varying(255),
    role integer,
    "networkIDs" json DEFAULT '[]'::json NOT NULL,
    tabs json,
    "createdAt" timestamp with time zone NOT NULL,
    "updatedAt" timestamp with time zone NOT NULL
```

* Organization

```
    id integer NOT NULL,
    name character varying(255),
    tabs json DEFAULT '[]'::json NOT NULL,
    "csvCharset" character varying(255),
    "customDomains" json DEFAULT '[]'::json NOT NULL,
    "networkIDs" json DEFAULT '[]'::json NOT NULL,
    "ssoSelectedType" public."enum_Organizations_ssoSelectedType" DEFAULT 'none'::public."enum_Organizations_ssoSelectedType" NOT NULL,
    "ssoCert" text DEFAULT ''::text NOT NULL,
    "ssoEntrypoint" character varying(255) DEFAULT ''::character varying NOT NULL,
    "ssoIssuer" character varying(255) DEFAULT ''::character varying NOT NULL,
    "ssoOidcClientID" character varying(255) DEFAULT ''::character varying NOT NULL,
    "ssoOidcClientSecret" character varying(255) DEFAULT ''::character varying NOT NULL,
    "ssoOidcConfigurationURL" character varying(255) DEFAULT ''::character varying NOT NULL,
    "createdAt" timestamp with time zone NOT NULL,
    "updatedAt" timestamp with time zone NOT NULL
```

* FeatureFlags - features enabled for the NMS

```
    id integer NOT NULL,
    "featureId" character varying(255) NOT NULL,
    organization character varying(255) NOT NULL,
    enabled boolean NOT NULL,
    "createdAt" timestamp with time zone NOT NULL,
    "updatedAt" timestamp with time zone NOT NULL
```

* AuditLogEntries - Table containing the mutations(POST/PUT/DELETE) made on NMS

```
    id integer NOT NULL,
    "actingUserId" integer NOT NULL,
    organization character varying(255) NOT NULL,
    "mutationType" character varying(255) NOT NULL,
    "objectId" character varying(255) NOT NULL,
    "objectDisplayName" character varying(255) NOT NULL,
    "objectType" character varying(255) NOT NULL,
    "mutationData" json NOT NULL,
    url character varying(255) NOT NULL,
    "ipAddress" character varying(255) NOT NULL,
    status character varying(255) NOT NULL,
    "statusCode" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone NOT NULL,
    "updatedAt" timestamp with time zone NOT NULL
```

Additionally the DB also contains session table to hold current set of sessions.

```
    sid character varying(36) NOT NULL,
    expires timestamp with time zone,
    data text,
    "createdAt" timestamp with time zone NOT NULL,
    "updatedAt" timestamp with time zone NOT NULL
```


