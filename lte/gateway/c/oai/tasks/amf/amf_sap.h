/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once
#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"
#ifdef __cplusplus
};
#endif

#include "amf_app_ue_context_and_proc.h"
#include "amf_app_defs.h"

namespace magma5g {
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/
/*
 * AMFREG-SAP primitives
 */
typedef enum {
  _AMFREG_START = 0,
  _AMFREG_COMMON_PROC_REQ,     /* AMF common procedure requested   */
  _AMFREG_COMMON_PROC_CNF,     /* AMF common procedure successful  */
  _AMFREG_COMMON_PROC_REJ,     /* AMF common procedure failed, CN send REJECT */
  _AMFREG_COMMON_PROC_ABORT,   /* AMF common procedure aborted     */
  _AMFREG_REGISTRATION_CNF,    /* 5GS network REGISTRATION accepted      */
  _AMFREG_REGISTRATION_REJ,    /* 6GS network REGISTRATION rejected      */
  _AMFREG_REGISTRATION_ABORT,  /* 5GS network REGISTRATION aborted      */
  _AMFREG_DEREGISTRATION_INIT, /* Network DEREGISTRATION initiated         */
  _AMFREG_DEREGISTRATION_REQ,  /* Network DEREGISTRATION requested         */
  _AMFREG_DEREGISTRATION_FAILED, /* Network DEREGISTRATION attempt failed    */
  _AMFREG_DEREGISTRATION_CNF,    /* Network DEREGISTRATION accepted          */
  _AMFREG_SERVICE_REQ,
  _AMFREG_SERVICE_CNF,
  _AMFREG_SERVICE_REJ,
  _AMFREG_END
} amf_reg_primitive_t;

/*
 * 5GMM Mobility Management primitives
 * ----------------------------------
 * AMFREG-SAP provides registration services for location updating and
 * registration/deregistration procedures;
 * AMFSMF-SAP provides interlayer services to the 5GMM Session Management
 * NF for service registration and activate/deactivate PDU session context;
 * AMFAS-SAP provides services to the Access Stratum sublayer for NAS message
 * transfer;
 */
enum amf_primitive_t {
  AMFREG_COMMON_PROC_REQ = 1,
  AMFREG_COMMON_PROC_CNF,
  AMFREG_COMMON_PROC_REJ,
  AMFREG_COMMON_PROC_ABORT,
  /* AMFAS-SAP */
  AMFAS_SECURITY_REQ              = _AMFAS_SECURITY_REQ,
  AMFAS_SECURITY_IND              = _AMFAS_SECURITY_IND,
  AMFAS_SECURITY_RES              = _AMFAS_SECURITY_RES,
  AMFAS_SECURITY_REJ              = _AMFAS_SECURITY_REJ,
  AMFAS_ESTABLISH_REQ             = _AMFAS_ESTABLISH_REQ,
  AMFAS_ESTABLISH_CNF             = _AMFAS_ESTABLISH_CNF,
  AMFAS_ESTABLISH_REJ             = _AMFAS_ESTABLISH_REJ,
  AMFAS_RELEASE_REQ               = _AMFAS_RELEASE_REQ,
  AMFREG_REGISTRATION_REJ         = _AMFREG_REGISTRATION_REJ,
  AMFAS_DATA_IND                  = _AMFAS_DATA_IND,
  AMFREG_REGISTRATION_CNF         = _AMFREG_REGISTRATION_CNF,
  AMFREG_DEREGISTRATION_REQ       = _AMFREG_DEREGISTRATION_REQ,
  AMFAS_DATA_REQ                  = _AMFAS_DATA_REQ,
  AMFCN_CS_RESPONSE               = _AMFCN_SMC_PARAM_RES,
  AMFCN_AUTHENTICATION_PARAM_RES  = _AMFCN_AUTHENTICATION_PARAM_RES,
  AMFCN_AUTHENTICATION_PARAM_FAIL = _AMFCN_AUTHENTICATION_PARAM_FAIL
};

/*
 * Minimal identifier for AMF-SAP primitives
 */
#define AMFREG_PRIMITIVE_MIN _AMFREG_START
#define AMFAS_PRIMITIVE_MIN _AMFAS_START
#define AMFCN_PRIMITIVE_MIN _AMFCN_START

/*
 * Maximal identifier for AMF-SAP primitives
 */
#define AMFREG_PRIMITIVE_MAX _AMFREG_END
#define AMFAS_PRIMITIVE_MAX _AMFAS_END
#define AMFCN_PRIMITIVE_MAX _AMFCN_END

/*
 * AMFREG primitive for registration  procedure
 * -------------------------------------
 */
typedef struct amf_reg_register_s {
  nas_amf_registration_proc_t* proc;
} amf_reg_register_t;

/*
 * Structure of AMFREG-SAP primitive
 */
typedef struct sap_primitive_s {
  amf_reg_register_t registered;
  struct nas_amf_common_proc_s* common_proc;
} sap_primitive_t;

typedef struct amf_reg_s {
  amf_primitive_t primitive;
  amf_ue_ngap_id_t ue_id;
  amf_context_t* ctx;
  bool notify;
  bool free_proc;
  sap_primitive_t u;
} amf_reg_t;

typedef struct amf_cn_auth_res_s {
  /* UE identifier */
  amf_ue_ngap_id_t ue_id;

  /* For future use: nb of vectors provided */
  uint8_t nb_vectors;

  /* Consider only one E-UTRAN vector for the moment... */
  eutran_vector_t* vector[MAX_EPS_AUTH_VECTORS];
} amf_cn_auth_res_t;

// typedef struct amf_ul_s {
typedef struct amf_cn_s {
  amf_cn_primitive_t primitive;
  union {
    amf_cn_auth_res_t* auth_res;
  } u;
} amf_cn_t;

typedef struct primitive_s {
  amf_reg_t amf_reg; /* AMFREG-SAP primitives    */
  amf_as_t amf_as;   /* AMFAS-SAP primitives     */
  amf_cn_t amf_cn;   /* AMFCN-SAP primitives     */
} primitive_t;

/*
 * Structure of 5GMM Mobility Management primitive
 */
typedef struct amf_sap_s {
  amf_primitive_t primitive;
  primitive_t u;
} amf_sap_t;

// Message passing from Access to Service Access Point.
int amf_sap_send(amf_sap_t* msg);

// Functions to handle all UL message final actions*/
int amf_cn_send(const amf_cn_t* msg);
int amf_reg_send(amf_sap_t* const msg);

}  // namespace  magma5g
