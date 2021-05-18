---
id: p017_mme_migrate_to_c++
title: C++ Migration
hide_title: true
---


# Proposal: C++ Migration

Author(s): [Author Name, Co-Author Name]

Last updated: 04/19/2021

Discussion at
[https://github.com/magma/magma/issues/4074](https://github.com/magma/magma/issues/4074).

## Abstract

This document describes the steps to migrate progressively, from C to C++, the MME core.

## Background

Initially and historically, the development of MME magma was done only in C.

Many third party libraries and features that were integrated and or developed around the MME core are written in C++.

We are now to the point where, for an easier and cleaner way of integrating new C++ libraries, for the ease of the maintenance of the code, a C to C++ migration seems unavoidable.

## Proposal

The proposal is constrained by a lack of performance of the current context serialization library (protobuf).
 
Because the replacement libraries may have a big impact on the MME core code, this library migration will be studied first.


### State Serialization
The serialization performance of protobuf is insufficient, so a serialization library that allows better performance must be selected.

The bad performances are on the global serialization process (translation, memory allocation, copy).

The replacement libraries considered are "cap'n proto" and "Flat buffers”, they support C/C++ standart API, binary output format and do not need memory copy (zero copy).


The states that have to be serialized are (To be detailled):
- S1AP, S1AP_IMSI_MAP, S1AP_UE
- MME_UE, MME_NAS
- SGW, SGW_UE, SPGW, SPGW_UE
- S11/GTPV2-C
- NGAP, NGAP_IMSI_MAP, NGAP_UE

#### S1AP, MME_APP, NAS migration

 S1AP seems to have no state dependency with other tasks.
 
 Should be an opportunity to start doing a first migration with S1AP.

 MME_APP and NAS tasks have contexts in common.
 

////### Migrate some MM and NAS IEs
////For the ease of manipulating IEs, migrate a selected list of IEs (operator ==, copy, conversion, etc)


### 3GPP protocols (OAI)
This work is necessary only for the deployment of Magma MME with CUPS architecture.

S6a Interface implemented in C should be migrated to C++.
A C++ port of S6a interface is already available in OAI HSS repo (https://github.com/OPENAIRINTERFACE/openair-hss) could serve as a basis.

S11 Interface and its GTPV2-C protocol implemented in C should be migrated to C++.
A C++ port of S11 and GTPV2-C interface done in OAI SPGW-C repo (https://github.com/OPENAIRINTERFACE/openair-spgwc) could serve as a basis.


## Rationale

### Get/Set state attributes

For a progressive migration, and a decoupling of MME core code from the serialization library, macros, inline, helpers functions/C++ patterns may be used instead of directly using serialization library API in MME core code.

Some MM and NAS IEs, for the ease of their manipulation (comparison operators, etc) could be migrated with wrapper classes.

[A discussion of alternate approaches and the trade offs, advantages, and
disadvantages of the specified approach.]

## Compatibility

Endianness: Whatever is the HW platform, endianness is fixed.

Backward compatibility of the new serialization library with protobuf is not envisioned.

## Observability and Debug

TODO
[A description, of how issues with this design would be observed and debugged
in various stages from development through production.]

## Implementation

  Description of steps in the C++ migration:
 
* Serialization library selection
    * Build 2 realistics Prototypes with MME_APP, NAS protos (cap’n proto, Flat buffers)
        * Migrate hashtables (find functions)
    * Get measurements and evaluate performances, maintainability.
    * Select library
* MME_APP/NAS Serialization migration
    * Write proto
    * Write redis client libraries
    * Migrate hashtables (find functions, evaluate performance)
    * Design classes for wrapping generated POD classes (?)
    	Not all context attribute need to be serialized
    * Migrate states (update core code with C++ patterns)
    * Write Unit tests
* S1AP Serialization migration
    * Write proto
    * Write redis client libraries
    * Migrate hashtables (find functions, evaluate performance)
    * Design classes for wrapping generated POD classes 
    	Not all context attribute need to be serialized
    * Migrate states (update core code)
    * Write Unit tests
* Intertask Messaging migration
    * Study message copy avoidance
* S6a migration (OAI only)
* S11 migration (OAI only)
 
[A description of the steps in the implementation, who will do them, and when.] TODO

## Open issues (if applicable)

[A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none.]


