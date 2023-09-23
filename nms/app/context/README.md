This directory contains context providers for global state.

Most contexts have set or update methods which take the following parameters.

- networkID
- current state
- method to set current state
- key
- optional value

If the key and value is specified,

- if the key doesn't exist, then a new resource is created
- if the key exists, we update the resource
then we get the latest state and update the current state. The additional get might
appear redundant, we do this currently so that we can handle two cases
- One is the case where the POST/PUT (config) entity differs from GET(status) entity. E.g. mutable_gateway and gateway.
- Secondly, in case the backend instataneously updates some state, in that case doing a GET ensures that we have updated
state rather than relying on stale state.

If only key is specified,
then it is a delete operation and delete the resource and
we update the resource
