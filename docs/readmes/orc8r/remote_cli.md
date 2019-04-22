---
id: remote_cli
title: Orchestrator Remote CLI - Creating commands guide
sidebar_label: Orchestrator Remote CLI
hide_title: true
---
# Orchestrator Remote CLI - Creating commands guide
## Creating a Command

There are a couple steps needed to implement a command.

### Define rpc method in gateway

Create a new RPC method in the service protobuf that you want the method to live in.

```
// magma/orc8r/protos/magmad.proto

service Magmad {
  ...
  rpc Reboot (Void) returns (Void) {}
  ...
}
```

### Implement in gateway

Gateway services should have a gRPC server implementation located in `rpc_servicer.py`. Within the servicer, create a function that implements this RPC method.

```
# magma/orc8r/gateway/python/magma/magmad/rpc_servicer.py

class MagmadRpcServicer(magmad_pb2_grpc.MagmadServicer):
    ...
    def Reboot(self, _, context):
        """
        Reboot the gateway device
        """
        ...
```

### Implement gateway api in cloud

Within the cloud service, create `gateway_api.go` that will call service methods. We can create a function to get a connection to the service, and use the dispatcher to forward requests to gateway services, like so:

```
// magma/orc8r/cloud/go/services/magmad/gateway_api.go

func getGWMagmadClient(networkId string, gatewayId string) (protos.MagmadClient, *grpc.ClientConn, context.Context, error) {
    ...
    conn, ctx, err := gateway_registry.GetGatewayConnection("magmad", gwRecord.HwId.Id)
    ...
    return protos.NewMagmadClient(conn), conn, ctx, nil
}
```

Using this client, we can create a function that calls the method:

```
// magma/orc8r/cloud/go/services/magmad/gateway_api.go

func GatewayReboot(networkId string, gatewayId string) error {
    client, conn, ctx, err := getGWMagmadClient(networkId, gatewayId)
    if err != nil {
        return err
    }
    defer conn.Close()
    _, err = client.Reboot(ctx, new(protos.Void))
    return err
}
```

### Define REST API endpoint

Each cloud service should have a `swagger` folder that documents service paths and definitions in `swagger.yml`. Add your path (and parameters or definitions, if necessary) to this file:

```
# magma/orc8r/cloud/go/services/magmad/swagger/swagger.yml

...
paths:
  ...
  /networks/{network_id}/gateways/{gateway_id}/command/reboot:
    post:
      summary: Reboot gateway device
      tags:
      - Commands
      parameters:
      - $ref: './swagger-common.yml#/parameters/network_id'
      - $ref: './swagger-common.yml#/parameters/gateway_id'
      responses:
        '200':
          description: Success
        default:
          $ref: './swagger-common.yml#/responses/UnexpectedError'
```

### Handler implementation

Each cloud service should have an `obsidian` folder, which contains handlers in `obsidian/handlers` and generated models in `obsidian/models`. 

Create your handler function:

```
// magma/orc8r/cloud/go/services/magmad/obsidian/handlers/gateway_handlers.go

func rebootGateway(c echo.Context) error {
    ...
    return c.NoContent(http.StatusOK)
}
```

Add your handler to the list of handlers in `GetObsidianHandlers()`.

```
// magma/orc8r/cloud/go/services/magmad/obsidian/handlers/handlers.go

func GetObsidianHandlers() []handlers.Handler {
    ...
    return []handlers.Handler{
        ...
        {Path: RebootGateway, Methods: handlers.POST, HandlerFunc: rebootGateway},
        ...
    }
}
```

Build and see your new endpoint in the swagger UI.

## Adding commands to NMS UI
Use `MagmaAPIUrls.command()` to get the url to the command endpoint `/networks/{network_id}/gateways/{gateway_id}/command/{command_name}`.

We can then make a request using that url, for example:
```
const url = MagmaAPIUrls.command(match, id, commandName);

axios
  .post(url)
  .then(_resp => {
    this.props.alert('Success');
  })
  .catch(error => this.props.alert(error.response.data.message));
