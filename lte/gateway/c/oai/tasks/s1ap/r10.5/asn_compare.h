/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#ifndef _ASN_COMPARE_H_
#define _ASN_COMPARE_H_

#ifdef __cplusplus
extern "C" {
#endif

struct asn_TYPE_descriptor_s; /* Forward declaration */

typedef enum COMPARE_ERR_CODE_e {
  COMPARE_ERR_CODE_START = 0,
  COMPARE_ERR_CODE_NONE = COMPARE_ERR_CODE_START,
  COMPARE_ERR_CODE_NO_MATCH,
  COMPARE_ERR_CODE_TYPE_MISMATCH,
  COMPARE_ERR_CODE_TYPE_ARG_NULL,
  COMPARE_ERR_CODE_VALUE_NULL,
  COMPARE_ERR_CODE_VALUE_ARG_NULL,
  COMPARE_ERR_CODE_CHOICE_NUM,
  COMPARE_ERR_CODE_CHOICE_PRESENT,
  COMPARE_ERR_CODE_CHOICE_MALFORMED,
  COMPARE_ERR_CODE_SET_MALFORMED,
  COMPARE_ERR_CODE_COLLECTION_NUM_ELEMENTS,
  COMPARE_ERR_CODE_END
} COMPARE_ERR_CODE_t;

typedef struct asn_comp_rval_s {
  enum COMPARE_ERR_CODE_e err_code;
  char *
    name; // e_S1ap_ProtocolIE_ID not available for all ASN1 use (RRC vs S1AP, X2AP)
  void *structure1;
  void *structure2;
  struct asn_comp_rval_s *next;
} asn_comp_rval_t;

#define COMPARE_CHECK_ARGS(                                                    \
  aRg_tYpE_dEf1, aRg_tYpE_dEf2, aRg_vAl1, aRg_vAl2, rEsUlT)                    \
  do {                                                                         \
    if ((aRg_tYpE_dEf1) && (aRg_tYpE_dEf2)) {                                  \
      if ((aRg_tYpE_dEf1->name) && (aRg_tYpE_dEf2->name)) {                    \
        if (strcmp(aRg_tYpE_dEf1->name, aRg_tYpE_dEf2->name)) {                \
          rEsUlT = (asn_comp_rval_t *) calloc(1, sizeof(asn_comp_rval_t));     \
          rEsUlT->err_code = COMPARE_ERR_CODE_TYPE_MISMATCH;                   \
          rEsUlT->name = aRg_tYpE_dEf1->name;                                  \
          return rEsUlT;                                                       \
        }                                                                      \
      } else {                                                                 \
        if ((aRg_tYpE_dEf1->xml_tag) && (aRg_tYpE_dEf2->xml_tag)) {            \
          if (strcmp(aRg_tYpE_dEf1->xml_tag, aRg_tYpE_dEf2->xml_tag)) {        \
            rEsUlT = (asn_comp_rval_t *) calloc(1, sizeof(asn_comp_rval_t));   \
            rEsUlT->err_code = COMPARE_ERR_CODE_TYPE_MISMATCH;                 \
            rEsUlT->name = aRg_tYpE_dEf1->xml_tag;                             \
            return rEsUlT;                                                     \
          }                                                                    \
        }                                                                      \
      }                                                                        \
    } else {                                                                   \
      rEsUlT = (asn_comp_rval_t *) calloc(1, sizeof(asn_comp_rval_t));         \
      rEsUlT->name = aRg_tYpE_dEf1->name;                                      \
      rEsUlT->structure1 = aRg_vAl1;                                           \
      rEsUlT->structure2 = aRg_vAl2;                                           \
      rEsUlT->err_code = COMPARE_ERR_CODE_TYPE_ARG_NULL;                       \
      return rEsUlT;                                                           \
    }                                                                          \
    if ((NULL == aRg_vAl1) || (NULL == aRg_vAl2)) {                            \
      rEsUlT = (asn_comp_rval_t *) calloc(1, sizeof(asn_comp_rval_t));         \
      rEsUlT->name = aRg_tYpE_dEf1->name;                                      \
      rEsUlT->structure1 = aRg_vAl1;                                           \
      rEsUlT->structure2 = aRg_vAl2;                                           \
      rEsUlT->err_code = COMPARE_ERR_CODE_VALUE_ARG_NULL;                      \
      return rEsUlT;                                                           \
    }                                                                          \
  } while (0);

#ifdef __cplusplus
}
#endif

#endif /* _ASN_COMPARE_H_ */
