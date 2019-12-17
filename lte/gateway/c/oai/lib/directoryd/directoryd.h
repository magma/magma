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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

/*
 * This enum should have the same definition as TableID in directoryd.proto .
 *
 * It's a bit difficult to directly use TableID here
 * due to some c++/c conversion when using directoryd.grpc.pb.h .
 * So we will define this type table_id_t here and later
 * cast it to magma::TableID in directoryd.cpp .
 */

bool directoryd_report_location(char* imsi);

bool directoryd_remove_location(char* imsi);

bool directoryd_update_location(char* imsi, char* location);

#ifdef __cplusplus
}
#endif
