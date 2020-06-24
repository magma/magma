/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/* This message asks for task initialization */
MESSAGE_DEF(
  INITIALIZE_MESSAGE,
  MESSAGE_PRIORITY_MED,
  IttiMsgEmpty,
  initialize_message)

/* This message asks for task activation */
MESSAGE_DEF(
  ACTIVATE_MESSAGE,
  MESSAGE_PRIORITY_MED,
  IttiMsgEmpty,
  activate_message)

/* This message asks for task deactivation */
MESSAGE_DEF(
  DEACTIVATE_MESSAGE,
  MESSAGE_PRIORITY_MED,
  IttiMsgEmpty,
  deactivate_message)

/* This message asks for task termination */
MESSAGE_DEF(
  TERMINATE_MESSAGE,
  MESSAGE_PRIORITY_MAX,
  IttiMsgEmpty,
  terminate_message)

/* Test message used for debug */
MESSAGE_DEF(MESSAGE_TEST, MESSAGE_PRIORITY_MED, IttiMsgEmpty, message_test)

/* Error message  */
MESSAGE_DEF(ERROR_LOG, MESSAGE_PRIORITY_MAX, IttiMsgEmpty, error_log)
/* Warning message  */
MESSAGE_DEF(WARNING_LOG, MESSAGE_PRIORITY_MAX, IttiMsgEmpty, warning_log)
/* Notice message  */
MESSAGE_DEF(NOTICE_LOG, MESSAGE_PRIORITY_MED, IttiMsgEmpty, notice_log)
/* Info message  */
MESSAGE_DEF(INFO_LOG, MESSAGE_PRIORITY_MED, IttiMsgEmpty, info_log)
/* Debug message  */
MESSAGE_DEF(DEBUG_LOG, MESSAGE_PRIORITY_MED, IttiMsgEmpty, debug_log)

/* Generic log message for text */
MESSAGE_DEF(GENERIC_LOG, MESSAGE_PRIORITY_MED, IttiMsgEmpty, generic_log)

// This message leads to recovery of timers for all tasks after MME restart
MESSAGE_DEF(
    RECOVERY_MESSAGE, MESSAGE_PRIORITY_MAX, IttiMsgEmpty, recovery_message)
