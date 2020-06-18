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

#include <stdio.h>
#include <stdint.h>
#include <unistd.h>

#include "bstrlib.h"
#include "intertask_interface.h"
#include "intertask_interface_init.h"

#include "common_defs.h"
#include "mme_config.h"
#include "mock_task.h"
#include "log.h"
#include "shared_ts_log.h"
#include "s1ap_mme.h"

// minus one from all string sizeof's due to trailing null byte
#define STATIC_BUF_LEN(buf) (sizeof(buf) - 1)

#define FUZZ_BUF_SIZE 4096

#define FUZZ_STREAM 1
#define FUZZ_ASSOC_ID 1
#define FUZZ_INSTREAMS 1
#define FUZZ_OUTSTREAMS 1

char s1ap_header[] = "\x00\x0c@";
char initial_header[] = "\x00\x00\x05\x00\x08\x00\x02\x00\x01\x00\x1a\x00";
char initial_footer[] =
  "\x00C\x00\x06\x00\x00\xf1\x10\x00\x01\x00d@"
  "\x08\x00\x00\xf1\x10\x00\x00\x00\xa0\x00\x86@\x010";

int nas_to_s1ap(char *buf, char *nas, int nas_len)
{
  int initial_len;
  int n;

  initial_len = STATIC_BUF_LEN(initial_header) + 2 + nas_len +
                STATIC_BUF_LEN(initial_footer);
  n = 0;

  memcpy(&buf[n], s1ap_header, STATIC_BUF_LEN(s1ap_header));
  n += STATIC_BUF_LEN(s1ap_header);

  buf[n] = (uint8_t) initial_len;
  n++;

  memcpy(&buf[n], initial_header, STATIC_BUF_LEN(initial_header));
  n += STATIC_BUF_LEN(initial_header);

  buf[n] = (uint8_t) nas_len + 1;
  n++;

  buf[n] = (uint8_t) nas_len;
  n++;

  memcpy(&buf[n], nas, nas_len);
  n += nas_len;

  memcpy(&buf[n], initial_footer, STATIC_BUF_LEN(initial_footer));
  n += STATIC_BUF_LEN(initial_footer);

  return n;
}

MessageDef *generate_fuzz_msg(char *packet, int len)
{
  MessageDef *msg;

  msg = itti_alloc_new_message(TASK_SCTP, SCTP_DATA_IND);

  SCTP_DATA_IND(msg).payload = blk2bstr(packet, len);
  SCTP_DATA_IND(msg).stream = FUZZ_STREAM;
  SCTP_DATA_IND(msg).assoc_id = FUZZ_ASSOC_ID;
  SCTP_DATA_IND(msg).instreams = FUZZ_INSTREAMS;
  SCTP_DATA_IND(msg).outstreams = FUZZ_OUTSTREAMS;

  return msg;
}

MessageDef *generate_s1ap_fuzz(int fd)
{
  char buf[FUZZ_BUF_SIZE];
  int n_read;

  n_read = read(fd, buf, sizeof(buf));

  return generate_fuzz_msg(buf, n_read);
}

MessageDef *generate_nas_fuzz(int fd)
{
  char s1ap[256];
  char nas[127];
  int nas_len;
  int s1ap_len;

  nas_len = read(fd, nas, sizeof(nas));

  s1ap_len = nas_to_s1ap(s1ap, nas, nas_len);

  return generate_fuzz_msg(s1ap, s1ap_len);
}

int fuzz(char *target)
{
  mock_task_t *sctp_task;
  MessageDef *msg;

  mme_config.max_enbs = 1;
  mme_config.max_ues = 32;

  if (
    itti_init(
      TASK_MAX,
      THREAD_MAX,
      MESSAGES_ID_MAX,
      tasks_info,
      messages_info,
      NULL,
      NULL) != RETURNok)
    return -1;

  if (
    OAILOG_INIT(
      MME_CONFIG_STRING_MME_CONFIG, OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS) !=
    RETURNok)
    return -1;
  if (shared_log_init(MAX_LOG_PROTOS) != RETURNok) return -1;

  sctp_task = mock_task_new(TASK_SCTP);
  (void) sctp_task;

  if (s1ap_mme_init(&mme_config) != RETURNok) return -1;

  if (strcmp(target, "nas") == 0) {
    msg = generate_nas_fuzz(STDIN_FILENO);
  } else if (strcmp(target, "s1ap") == 0) {
    msg = generate_s1ap_fuzz(STDIN_FILENO);
  } else {
    return -1;
  }

  itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, msg);

  usleep(500);

  return 0;
}

int main(int argc, char **argv)
{
  if (argc != 2) {
    printf("Usage: %s <s1ap|nas>\n", argv[0]);
    return -1;
  }

  if (fuzz(argv[1]) < 0) {
    printf("Invalid fuzzing target -- specify s1ap or nas\n");
    return -1;
  }

  return 0;
}
