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
#include <arpa/inet.h>
#include <stdio.h>

#include "RpcClient.h"

int main(int argc, char **argv)
{
  int status;
  struct in_addr tmp;
  struct in_addr ipv4_addr1;
  struct in_addr ipv4_addr2;
  struct in_addr ipv4_addr3;

  char apn[10] = "oai.ipv4";
  {
    char str[INET_ADDRSTRLEN];
    printf("Allocating IP address...\n");

    status = allocate_ipv4_address("0001", apn, &ipv4_addr1);
    if (status) {
      printf("allocate_ipv4_address error %d for sid %s for apn %s\n",
             status, "0001", apn);
      return -1;
    }
    tmp.s_addr = htonl(ipv4_addr1.s_addr);
    inet_ntop(AF_INET, &tmp, str, INET_ADDRSTRLEN);
    printf("IP allocated: %s\n", str);

    status = allocate_ipv4_address("0002", apn, &ipv4_addr2);
    if (status) {
      printf("allocate_ipv4_address error %d for sid %s for apn %s\n",
              status, "0002", apn);
      return -1;
    }
    tmp.s_addr = htonl(ipv4_addr2.s_addr);
    inet_ntop(AF_INET, &tmp, str, INET_ADDRSTRLEN);
    printf("IP allocated: %s\n", str);

    status = get_ipv4_address_for_subscriber("0002", apn, &ipv4_addr3);
    if (status) {
      printf(
        "get_ipv4_address_for_subscriber error %d for sid %s\n",
        status,
        "0002");
      return -1;
    }
    tmp.s_addr = htonl(ipv4_addr2.s_addr);
    inet_ntop(AF_INET, &tmp, str, INET_ADDRSTRLEN);
    printf("Retrieved ip address for user 0002: %s\n", str);
  }

  {
    printf("Releasing IP address...\n");

    status = release_ipv4_address("0001", apn, &ipv4_addr1);
    if (status) {
      printf("release_ipv4_address error %d for sid %s\n", status, "0001");
      return -1;
    }
    status = release_ipv4_address("0002", apn, &ipv4_addr2);
    if (status) {
      printf("release_ipv4_address error %d for sid %s\n", status, "0002");
      return -1;
    }
  }

  printf("All tests passed...\n");
}
