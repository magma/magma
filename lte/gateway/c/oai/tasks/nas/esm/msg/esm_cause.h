/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
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
Source    esm_cause.h

Version   0.1

Date    2013/02/06

Product   NAS stack

Subsystem EPS Session Management

Author    Frederic Maurel

Description Defines error cause code returned upon receiving unknown,
    unforeseen, and erroneous EPS session management protocol
    data.

*****************************************************************************/
#ifndef FILE_ESM_CAUSE_SEEN
#define FILE_ESM_CAUSE_SEEN

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * Cause code used to notify that the EPS session management procedure
 * has been successfully processed
 */
#define ESM_CAUSE_SUCCESS (esm_cause_t)(-1)

#endif /* FILE_ESM_CAUSE_SEEN*/
