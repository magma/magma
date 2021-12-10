---
id: version-1.2.0-p001_vpn_config_from_api
title: Configurable VPN from Orchestrator API
hide_title: true
original_id: p001_vpn_config_from_api
---

# Configurable VPN from Orchestrator API

- Feature owners: `@alexrod`, `@tdzik`
- Feedback requested from: `@apad`, `@xjtian`, `@shasan`

## Summary

Currently, the only option to gain remote access to any Access Gateway involves some manual steps to be able to introduce an OpenVPN connection to it. This proposal involves some updates to the process to allow a configurable VPN connection through the Orchestrator REST API. 

## Motivation

The current method of configuring and then creating a VPN connection to an AGW, involves some manual steps, detailed:

1. Download and install openvpn package, acting as the OpenVPN client from the AGW. (Which ultimately connects to an AWS instance - OpenVPN server)
2. Register magma AGW on Orchestrator using `show_gateway_info.py` information. (AGW needs to be bootstrapped)
3. Run a script `vpn_setup_cli` which contains fixed server endpoint and port, retrieves the certs from cloud and sets up client config.
4. End-user then uses updated files from setup script to connect to VPN.

We can automate some of this setup by configuring connection parameters from Orchestrator REST API, we also get more control of the setup in case of setup / connection goes wrong by enabling / disabling VPN, or changing the config parameters. Security is also another important topic to take account, by keeping the credentials updated with the gateway certificates / certifier and also revoked once they're not valid or expired.

## Goals

The goals of automating VPN setup and making it configurable are:

- Allow an easier method of configuring VPN connection setup for AGWs
- Automating enabling / disabling of the connection
- Improve security of the VPN workflow

## Proposal

```                                                                                         
|-------------------------------|                        +-------------------------------+
|                               |      VPN Credentials   |                               |
|                               |------------------------|                               |
|                               |    magmad connection   |         Access Gateway        |
|   Orchestrator Bootstrapper   |                        |                               |
|                               |                        |                               |
|                               |                        |                               |
|                               |                        |                               |
|                               |                        |                               |
|--------------------------------                        +-------------------------------+
                                                                          |               
                                                                          |               
                                                                          /               
                                                                         |                
                                                                         |                
                     +----------------------+                    +--------------------+   
                     |                      |                    |                    |   
                     |                      |                    |     UDP/TCP        |   
                     |    OpenVPN Server    |--------------------|  OpenVPN Client    |   
                     | module of Terraform  |                    |                    |   
                     |                      |                    |                    |   
                     |                      |                    +--------------------+   
                     +----------------------+                                          
```

For setting this up, we can take advantage of our terraform module configuration, to deploy and expose an OpenVPN server. 
We can deploy an OpenVPN server on kubernetes by using helm openvpn module on: https://hub.helm.sh/charts/stable/openvpn. This server should use persistent volume in k8s to store all the client keys mapping information. 

### Bootstrapper Updates

VPN credentials will be created on the client, and will be sent to cloud controller during bootstrapping process with a cert signing request (which is part of the magmad connection with the access gateway), these should be rotated along with the gateway certs and revoked if these become not valid or are expired. These credentials will be created with a default short-lived duration (e.g. 12 hours)


OpenVPN helm chart provides with configuration scripts to setup and manage the client keys. The interface can be extended with:
- requestVPNSignCert
  - Will be called by the client to request a signed certificate (.crt) by the server after the client keys were generated.
- revokeVPNClientCert
  - This will revoke the RSA vpn client cert, it should be aligned and rotated with an invalidation/expiration of the gateway certificates.

We can mount the same persistent volume that the server uses for storing the client keys onto the bootstrapper, which will be the communication and syncing between the bootstrapper and OpenVPN.

### REST API Endpoint

From Orchestrator, we can add a new controller app endpoint that will allow user to do multiple operations on the VPN connection config:
- `.../gateways/gateway_id/vpn_config`
  - Enable / Disable shell flag

From here, the AGW can spin off and enable an OpenVPN UDP client, we can wrap the client into a dynamic service that can be easier to manage and activate using magmad service. This implementation should give us more flexibility as deploying the OpenVPN server becomes a specification as a Terraform module, provides more security by maintaining the VPN credentials valid along with the Access Gateway certificates, and also provides an scenario for the user to configure and manage the VPN connections right from the API while also help with recovery options when issues arise. UDP is the preferred default configuration for the client due to performance, and should fallback to TCP when necessary.

### VPN Usage

For usage of the VPN, we will define as owner to the operator (e.g. support team member, maintainer) of the entire system as of now. In the future we can look into defining certificates with principal validity for a particular network, and the server to only accept certification who have that network registered as principal.
Current VPN server setup has `client-to-client` option enabled, which eases the SSH connection for gateways, but opens up vulnerabilities as other AGWs connected to VPN become open. If we disable this option, the `fab ssh` command should be extended with a step to first connect to the kube pod running openvpn, and then sshing into the client / gateway. 

## Timeline of Work

- Adding deployment of OpenVPN server through Terraform module
- Implement interface for bootstrapper RPC VPN client certificate management
- Update bootstrapper process to include provision / maintenance of VPN creds for enabled VPN gateways
- Add UDP/TPC OpenVPN client wrapper to AGWs
- Add cloud endpoints for VPN configuration / management 
