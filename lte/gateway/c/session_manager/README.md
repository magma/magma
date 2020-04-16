
# Session Manager  
`sessiond` runs on the LTE gateway and CWAG and is used to track active sessions.
## Responsibilities
  
**Tracking Active Sessions**  
The main responsibility of `sessiond` is to track active sessions.
For each session, the currently installed rules are also tracked, 
both static and dynamic. Credit usage is also tracked for both charging 
and monitoring credit. Some additional configuration is also tracked for 
each session, and this can be found in `SessionState`.

**Enforcement Middleman**
`sessiond` receives requests from both the OCS and PCRF, and based on its 
internal view of sessions, is responsible for instructing pipelined for 
proper enforcement.

**Credit Usage Reporting**
`sessiond` will receive usage updates from `pipelined`. Updates on credit 
usage, and ended sessions are reported back up to the OCS and PCRF based on 
these updates.


## Stateless Operation
To operate sessiond as a stateless service, this feature should be enabled 
in `sessiond.yml`. The field is marked as `support_stateless`.
This will allow sessiond to be restarted without requiring sessions to be 
re-authenticated.

When stateless operation is enabled, sessiond will store session state in 
persistent storage. When a gRPC request request is received which requires 
acting on a session, its state will be read from storage, operated on and 
modified in memory, and then the updates to the session will be written back 
to storage before responding to the gRPC request.

`SessionStore` is the interface to storage required for stateless operation, 
but is still used even when session state is only stored in memory. 
Updates to stored session state are done with the `UpdateCriteria` passed 
into the `SessionStore`.

To ensure that `sessiond` can always be restarted properly, session state 
should always be consistent with services with which it interacts. 
On `sessiond` restart, session state will be synced to `mme` service, which 
it will use as the source of truth for active sessions. 
Based on the active sessions, `pipelined` will be instructed to 
allow/disallow traffic.
