/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */
#ifndef FILE_SGI_FORWARD_MESSAGES_TYPES_SEEN
#define FILE_SGI_FORWARD_MESSAGES_TYPES_SEEN

#include "common_types.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
typedef enum SGIStatus_e {
  SGI_STATUS_OK                                   = 16,
  SGI_STATUS_ERROR_CONTEXT_NOT_FOUND              = 64,
  SGI_STATUS_ERROR_INVALID_MESSAGE_FORMAT         = 65,
  SGI_STATUS_ERROR_SERVICE_NOT_SUPPORTED          = 68,
  SGI_STATUS_ERROR_SYSTEM_FAILURE                 = 72,
  SGI_STATUS_ERROR_NO_RESOURCES_AVAILABLE         = 73,
  SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED = 84,
  SGI_STATUS_ERROR_NO_MEMORY_AVAILABLE            = 91,
  SGI_STATUS_ERROR_FAILED_TO_PROCESS_PCO          = 94,
  SGI_STATUS_MAX,
} SGIStatus_t;

typedef struct {
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  teid_t sgw_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
  pdn_type_t pdn_type;  ///< PDN Type
  paa_t paa;            ///< PDN Address Allocation
} itti_sgi_create_end_point_request_t;

typedef struct {
  SGIStatus_t status;   ///< Status of  endpoint creation (Failed = 0xFF or
                        ///< Success = 0x0)
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  teid_t sgw_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
  pdn_type_t pdn_type;  ///< PDN Type
  paa_t paa;            ///< PDN Address Allocation
  protocol_configuration_options_t pco;  ///< Protocol configuration options
} itti_sgi_create_end_point_response_t;

typedef struct {
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  teid_t sgw_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  teid_t enb_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
} itti_sgi_update_end_point_request_t;

typedef struct {
  teid_t sgw_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  teid_t enb_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
} bearer_cxt_t;

typedef struct {
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  uint8_t num_bearers_modified;
  bearer_cxt_t bearer_contexts_to_be_modified[BEARERS_PER_UE];
  uint8_t num_bearers_removed;
  ebi_t bearer_contexts_to_be_removed[BEARERS_PER_UE];
  uint8_t num_bearers_not_found;
  ebi_t bearer_contexts_not_found[BEARERS_PER_UE];
} itti_sgi_update_end_point_response_t;

typedef struct {
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  teid_t sgw_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
  pdn_type_t pdn_type;  ///< PDN Type
  paa_t paa;            ///< PDN Address Allocation
} itti_sgi_delete_end_point_request_t;

typedef struct {
  SGIStatus_t status;   ///< Status of  endpoint deletion (Failed = 0xFF or
                        ///< Success = 0x0)
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  teid_t sgw_S1u_teid;  ///< Tunnel Endpoint Identifier S1-U
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
  pdn_type_t pdn_type;  ///< PDN Type
  paa_t paa;            ///< PDN Address Allocation
} itti_sgi_delete_end_point_response_t;

//-----------------------------------------------------------------------------
/**
 * IP Allocation response from MobilityServiceClient, sent from GRPC task to
 * SPGW as part of IP allocation process, notifying the SPGW task if it was
 * successful or not.
 */
typedef struct itti_ip_allocation_response_s {
  SGIStatus_t status;  ///< Status of  endpoint creation (Failed = 0xFF or
  ///< Success = 0x0)
  teid_t context_teid;  ///< Tunnel Endpoint Identifier S11
  ebi_t eps_bearer_id;  ///< EPS bearer identifier
  paa_t paa;            ///< PDN Address Allocation
} itti_ip_allocation_response_t;

#endif /* FILE_SGI_FORWARD_MESSAGES_TYPES_SEEN */
