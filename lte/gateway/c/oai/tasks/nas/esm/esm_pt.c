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

/*****************************************************************************
  Source      esm_pt.c

  Version     0.1

  Date        2013/01/03

  Product     NAS stack

  Subsystem   EPS Session Management

  Author      Frederic Maurel

  Description Defines functions used to handle ESM procedure transactions.

*****************************************************************************/
#include <stdint.h>

#include "3gpp_24.007.h"
#include "esm_pt.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/*
   Minimal and maximal value of a procedure transaction identity:
   The Procedure Transaction Identity (PTI) identifies bi-directional
   messages flows
*/
#define ESM_PTI_MIN (PROCEDURE_TRANSACTION_IDENTITY_FIRST)
#define ESM_PTI_MAX (PROCEDURE_TRANSACTION_IDENTITY_LAST)

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    esm_pt_is_reserved()                                      **
 **                                                                        **
 ** Description: Check whether the given procedure transaction identity is **
 **      a reserved value                                          **
 **                                                                        **
 ** Inputs:  pti:       The identity of the procedure transaction  **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    true, false                                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_pt_is_reserved(int pti)
{
  return ((pti != ESM_PT_UNASSIGNED) && (pti > ESM_PTI_MAX));
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
