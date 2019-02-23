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

#ifndef PROCEDURE_TRANSACTION_IDENTITY_H_
#define PROCEDURE_TRANSACTION_IDENTITY_H_
#include <stdint.h>

#define PROCEDURE_TRANSACTION_IDENTITY_MINIMUM_LENGTH 1
#define PROCEDURE_TRANSACTION_IDENTITY_MAXIMUM_LENGTH 1

typedef pti_t ProcedureTransactionIdentity;

int encode_procedure_transaction_identity(
  ProcedureTransactionIdentity *proceduretransactionidentity,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

void dump_procedure_transaction_identity_xml(
  ProcedureTransactionIdentity *proceduretransactionidentity,
  uint8_t iei);

int decode_procedure_transaction_identity(
  ProcedureTransactionIdentity *proceduretransactionidentity,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

#endif /* PROCEDURE TRANSACTION IDENTITY_H_ */
