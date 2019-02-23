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
/*! \file gtpv1u_task.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdlib.h>
#include <stdint.h>

#include "gtpv1u.h"

//#define GTPV1U_LINEAR_TEID_ALLOCATION 1

#if GTPV1U_LINEAR_TEID_ALLOCATION
static uint32_t g_gtpv1u_teid = 0;
#endif

// TODO (amar) Function is unused. See if we can remove.
uint32_t gtpv1u_new_teid(void)
{
#if GTPV1U_LINEAR_TEID_ALLOCATION
  g_gtpv1u_teid = g_gtpv1u_teid + 1;
  return g_gtpv1u_teid;
#else
  return random() + random() % (RAND_MAX - 1) + 1;
#endif
}
