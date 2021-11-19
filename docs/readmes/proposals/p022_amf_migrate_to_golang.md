---
id: p022_amf_migrate_to_golang
title: Golang Migration
hide_title: true
---


# Proposal: Go Migration

Last updated: 19/11/2021

## Abstract

This document describes the steps to migrate progressively, from C++ to Golang, the AMF core.

## Background

Initially, the development of AMF magma was done in C++. used ITTI for interprocess communication between SCTP - NGAP and NGAP - AMF.

## Proposal

The proposal is to migrate AMF into Golang, replace ITTI with Grpc.

## NGAP, AMF_APP, NAS migration

NGAP: define protos for SCTP-NGAP rpc calls, implement client and server functionality on sctp and ngap side.
      Procedures to be migrated :
        NG-SETUP request
        NG-SETUP response
        NG-SETUP failure
        NG-SETUP reset
        Initial UE message
        Uplink NAS message
        Initial context setup request 
        Initial context setup response
        
AMF : define protos for NGAP-AMF rpc calls, implement client and server functionality on ngap and amf side.
      Procedures to be migrated :
        Registration (request, accept and complete)
        Identification (request and response)
        Authentication (request and response)
        Security mode (command and complete)
        De-registration (request and accept)

NAS : Migrate encode/decode routines of the following messages from C++ to Golang :
        Registration reuqest
        Identity request
        Identity response
        Authentication request
        Authentication response
        Security mode command
        Security mode complete
        Registration accept
        Registration complete
        De-Registration request (UE-Originated)
        De-Registration accept (UE-Originated)
