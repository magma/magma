# Secure-by-Design Policy

## 1. Purpose and Scope

This policy defines the security principles and engineering requirements for designing, building, and operating the software so that components are secure by default, resilient under attack, and auditable.

It applies to all source code, services, APIs, containers, deployment manifests, CI/CD workflows, default configuration, and third-party dependencies maintained by the project — with particular attention to the network-facing control plane, data plane, and federation/roaming interfaces.

The intent is practical security a small team can actually sustain: maximize controls that live in code and CI; avoid process ceremony that cannot be reliably enforced.

---

## 2. Secure-by-Design Principles

1. **Security by default.** Components start in the most secure practical configuration. Insecure modes, demo credentials, debug ports, broad network exposure, permissive TLS, and unauthenticated APIs must not be enabled by default.
2. **Least privilege.** Each process, container, service account, and API client receives only the permissions it needs.
3. **Zero trust between components.** Internal traffic is not trusted merely because it originates from another component, namespace, or local network. Authentication, authorization, and integrity protection apply across trust boundaries.
4. **Defense in depth.** Do not rely on a single control. Layer authentication, authorization, network segmentation, cryptography, input validation, rate limiting, and logging.
5. **Fail secure.** If a dependency, database, CA, or policy service is unavailable, the component fails into a documented safe state rather than silently allowing unauthorized access or traffic.
6. **Explicit trust boundaries.** Document the boundaries between management plane, control plane, data plane, federation interfaces, and external integrations.
7. **Secure operability.** Operators can rotate credentials, revoke certificates, inspect logs, and confirm enforcement without unsafe manual workarounds.
8. **Secure supply chain.** Treat dependencies, build systems, container images, and CI credentials as part of the attack surface.

---

## 3. Ownership and Security Review

Every component must have a **named maintainer** responsible for its security review, vulnerability triage, dependency updates, and secure-configuration guidance. One person may own several components; the point is that responsibility is never unassigned.

A focused security review is required before merging changes that affect:

* Authentication or authorization.
* Certificate handling or secret storage.
* Subscriber identity, profile, or policy handling.
* Data-plane rule programming.
* Network-facing protocol parsers.
* External or northbound APIs.
* Container privileges or CI/CD release workflows.
* Logging of potentially sensitive information.

For new components, new network-facing interfaces, or changes to privilege boundaries, the reviewer should briefly reason through the threat surface: what the assets and entry points are, who the untrusted callers are, how the component fails, and what is logged. This is a short written note in the PR or design doc — not a separate formal artifact.

---

## 4. Architecture Security

### 4.1 Component Isolation

Each service runs as an isolated component with a dedicated runtime identity, minimal filesystem access, minimal Linux capabilities, explicit network exposure, and resource limits. Avoid privileged containers, host networking, and writable root filesystems unless required.

Privileged execution is permitted only when needed for packet processing, interface management, or kernel interaction. The justification must be documented in the deployment notes.

### 4.2 Management Plane

Management-plane components must enforce strong authentication, role-based access control, TLS for operator-facing APIs, session/token expiration, audit logging for configuration changes, brute-force protection, modern password hashing, and no default administrative credentials. Management interfaces must not be exposed to untrusted networks without an explicit operator decision.

### 4.3 Control Plane

Control-plane components must validate all protocol input, including malformed, oversized, replayed, duplicated, fragmented, and out-of-order messages. State machines must reject invalid transitions, enforce per-peer and per-subscriber limits, rate-limit expensive operations, avoid unbounded memory growth, and log security-relevant failures without leaking secrets or subscriber identifiers.

### 4.4 Data Plane

Data-plane components must enforce policy deterministically and must not silently bypass enforcement when dependent services fail. Required controls:

* Default-deny for unknown or invalid sessions.
* Validation of tunnel identifiers, subscriber mappings, and forwarding state.
* Cleanup of stale rules after session teardown, restart, or upgrade.
* Auditable correlation between session state and installed forwarding rules.
* Tests for rule collision, shadowing, priority conflicts, and stale state.

### 4.5 Federation and External Integration

Federation- and roaming-facing interfaces are high-risk boundaries. They must provide mutual authentication where supported, strict certificate validation, protocol-level input validation, peer allowlisting where practical, rate limits, replay protection where applicable, security logging, and negative tests for malformed messages.

---

## 5. Authentication and Authorization

### 5.1 Authentication

All administrative, northbound, and inter-component APIs require authentication unless explicitly documented as public health checks. Acceptable methods: mutual TLS, short-lived signed tokens, federated identity through a trusted provider, or rotatable service-account credentials.

Never acceptable: secrets embedded in source, default passwords, long-lived tokens without rotation, secrets in plaintext config, authentication disabled by default, or trusting client-supplied identity headers without upstream enforcement.

### 5.2 Authorization

Every privileged operation enforces authorization server-side — covering read/write access, tenant or network scope, subscriber and policy operations, gateway registration, certificate lifecycle, role management, upgrades, and log access. Hiding controls in the UI is not authorization; APIs must enforce access independently.

### 5.3 Service-to-Service Identity

Every service-to-service call across a trust boundary authenticates the caller and authorizes the operation. Service identity must be unique per workload class, rotatable, revocable, and auditable.

---

## 6. Cryptography, Certificates, and Secrets

### 6.1 Cryptographic Baseline

Use modern, widely reviewed libraries and protocols. TLS 1.2 minimum (1.3 preferred), strong cipher suites only, certificate and identity validation enabled by default, no self-signed production defaults without explicit operator action, no obsolete algorithms (MD5, SHA-1 signatures, RC4, DES, export ciphers), and no custom cryptography.

### 6.2 Certificate Lifecycle

Support certificate rotation without full redeployment where practical, expiration monitoring, revocation/replacement procedures, separate certificates per trust domain, a documented bootstrap process, and safe handling of expired, missing, malformed, or mismatched certificates.

### 6.3 Secret Handling

Secrets — private keys, API tokens, database passwords, subscriber authentication material, signing keys, cloud and CI credentials — must never be committed to source control. CI must run secret scanning and reject commits containing likely credentials.

---

## 7. Secure Coding

### 7.1 Input Validation

Validate all external input — API bodies, query parameters, headers, protocol messages, subscriber and network identifiers, configuration files, environment variables, CLI arguments, and records persisted by previous software versions. Validation covers type, length, encoding, range, required fields, allowed values, cross-field consistency, and state-machine validity.

### 7.2 Memory Safety

Components in memory-unsafe languages must use compiler hardening, static analysis, bounds checking, safe buffer wrappers, and strict review of pointer arithmetic and (de)serialization. Parsers and protocol handlers should have fuzz coverage and be tested under address sanitizer or equivalent where practical. Prefer memory-safe patterns for new protocol parsers.

### 7.3 Error Handling

Handle errors explicitly. The software must not ignore failed authorization checks, continue with partially initialized security state, fall back to insecure defaults, log secrets, expose stack traces to unauthenticated users, crash on malformed remote input, or return ambiguous success when enforcement failed.

### 7.4 Logging Safety

Logs must be useful for security operations without leaking sensitive data. Never log (unless masked/redacted): private keys, tokens, passwords, authentication vectors, subscriber secrets, session cookies, authorization headers, unnecessary PII, or connection strings containing credentials.

Security events should include enough metadata for investigation: timestamp, component, actor, action, result, target, source address where available, and a correlation ID.

---

## 8. API Security

All APIs must enforce: authentication by default, per-operation authorization, TLS in production, strict schema validation, versioned contracts, rate limiting for expensive operations, pagination on list operations, and audit logging for state-changing operations. They must not put sensitive data in URLs, ship unsafe CORS defaults, or expose debug/administrative endpoints by default.

API documentation should identify the authentication method, authorization scope, request/response schema, error behavior, and security-sensitive fields.

---

## 9. Protocol Parser and Network Interface Requirements

All network-facing protocol handlers are treated as hostile-input parsers. They must:

* Validate message length before parsing.
* Validate mandatory and optional fields; reject unknown critical fields where the protocol requires it.
* Handle duplicate fields safely.
* Enforce state-machine order and apply per-peer rate limits.
* Avoid unbounded allocations and recursion without depth limits.
* Avoid unsafe deserialization.
* Produce structured error logs and include fuzzing coverage where practical.

Protocol decoding failures must never crash the component.

---

## 10. Policy and Subscriber Enforcement

Policy and subscriber-handling components must keep runtime behavior consistent with configured intent:

* Validate subscriber profiles before activation.
* Reject incomplete or conflicting policy rules.
* Enforce default-deny for unknown subscribers.
* Prevent privilege escalation through malformed profiles.
* Track the policy version applied to active sessions.
* Audit subscriber and policy changes.
* Reconcile runtime state after restart and remove stale session/forwarding state.
* Test policy precedence and conflict resolution.

---

## 11. Secure Defaults and Configuration

Default deployments must be safe enough for evaluation without teaching unsafe production habits. The project must not ship with default administrative passwords, hard-coded private keys, TLS disabled on production paths, wildcard network exposure, unauthenticated dashboards, debug logging on by default, privileged containers unless necessary, unrestricted database or message-bus access, or open metrics endpoints exposing sensitive labels.

Configuration templates must clearly separate development, lab, and production settings, and any unsafe compatibility setting must carry an explicit warning.

---

## 12. Supply Chain and CI/CD

* Maintain awareness of third-party dependencies; pin versions or use lock files where practical, review new dependencies, remove unused ones, and update critical/high vulnerabilities promptly.
* Run automated dependency and container image scanning, and rebuild images when base-image vulnerabilities require it.
* Container images should use minimal base images, run as non-root unless justified, drop unnecessary capabilities, avoid embedded secrets, and be scanned before release.
* Build release artifacts through CI on protected branches with mandatory review. CI secrets must not be available to untrusted forks, and CI tokens follow least privilege.
* Heightened review for any change to build scripts, deployment scripts, or workflows that can publish artifacts or access signing keys.

(Artifact/image signing, SBOM generation, and provenance metadata are encouraged where the tooling exists, but are not gating requirements for this project.)

---

## 13. Testing

Security testing is part of normal engineering, not a final gate. Maintain, proportionate to each component's risk:

* Unit tests for authorization logic.
* Integration tests for authentication flows.
* Negative tests for malformed protocol messages.
* Regression tests for past vulnerabilities.
* Fuzz tests for parsers and decoders.
* Dependency, container, and secret scanning in CI.
* Data-plane policy enforcement tests.

Do not promote a release candidate with unresolved critical or high vulnerabilities unless the risk is explicitly documented and accepted.

---

## 14. Data Protection

Classify and protect sensitive data — subscriber identity, authentication material, policy and session state, location/mobility data, operator credentials, and network topology. Required controls: encryption in transit, encryption at rest where the deployment environment supports it, role-based access, redaction in logs, minimized retention, and access auditability.

---

## 15. Runtime Hardening Guidance

Production deployment documentation should cover network segmentation and firewalling, Kubernetes network policies or equivalent, container/host hardening, Linux capability minimization, file permissions, log forwarding, backup/restore, and certificate and secret rotation. Components requiring elevated host access must document why and how operators should restrict it.

---

## 16. Vulnerability Reporting and Exceptions

Maintain a current, monitored private channel for reporting suspected vulnerabilities (the repository's security policy), describing how to report, what response to expect, and which versions are supported. Triage each report for affected components and versions, exploitability, required privileges, impact, and available workarounds. Fix critical and high-severity issues promptly and provide mitigation guidance to operators when a fix will take time.

Any deliberate exception to this policy is recorded in a tracking issue noting the affected component, the justification, any compensating controls, and an owner. Exceptions should not be open-ended unless the risk is inherent to the component's documented function.

---

## 17. New Component Checklist

A new component should not merge until it has: a named owner, a documented purpose and trust boundaries, an authentication and authorization model, a secret-handling and logging approach, defined failure behavior, resource limits, security tests, and basic deployment hardening notes.
