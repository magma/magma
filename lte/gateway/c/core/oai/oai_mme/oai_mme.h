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

/** @mainpage

  @section intro Introduction

  openair-mme project tends to provide an implementation of LTE core network.

  @section scope Scope


  @section design Design Philosophy

  Included protocol stacks:
  - SCTP RFC####
  - S1AP 3GPP TS 36.413 R10.5
  - S11 abstraction between MME and S-GW
  - 3GPP TS 23.401 R10.5
  - nw-gtpv1u for s1-u (http://amitchawre.net/)
  - freeDiameter project (http://www.freediameter.net/) 3GPP TS 29.272 R10.5

  @section applications Applications and Usage

  Please use the script to start LTE epc in root src directory

 */

/*! \file oai_mme.h
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/

#ifndef FILE_OAISIM_MME_SEEN
#define FILE_OAISIM_MME_SEEN

int main(int argc, char* argv[]);

#endif /* FILE_OAISIM_MME_SEEN */
