/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/** @defgroup _intertask_interface_impl_ Intertask Interface Mechanisms
 * Implementation
 * @ingroup _ref_implementation_
 * @{
 */

#ifndef INTERTASK_INTERFACE_TYPES_H_
#define INTERTASK_INTERFACE_TYPES_H_

#include <time.h>
#include "itti_types.h"

#include "messages_types.h"

/* Defines to handle bit fields on unsigned long values */
#define UL_BIT_MASK(lENGTH) ((1UL << (lENGTH)) - 1UL)
#define UL_BIT_SHIFT(vALUE, oFFSET) ((vALUE) << (oFFSET))
#define UL_BIT_UNSHIFT(vALUE, oFFSET) ((vALUE) >> (oFFSET))

#define UL_FIELD_MASK(oFFSET, lENGTH)                                          \
  UL_BIT_SHIFT(UL_BIT_MASK(lENGTH), (oFFSET))
#define UL_FIELD_INSERT(vALUE, fIELD, oFFSET, lENGTH)                          \
  (((vALUE) & (~UL_FIELD_MASK(oFFSET, lENGTH))) |                              \
   UL_BIT_SHIFT(((fIELD) &UL_BIT_MASK(lENGTH)), oFFSET))
#define UL_FIELD_EXTRACT(vALUE, oFFSET, lENGTH)                                \
  (UL_BIT_UNSHIFT((vALUE), (oFFSET)) & UL_BIT_MASK(lENGTH))

/* Definitions of task ID fields */
#define TASK_THREAD_ID_OFFSET 8
#define TASK_THREAD_ID_LENGTH 8

#define TASK_SUB_TASK_ID_OFFSET 0
#define TASK_SUB_TASK_ID_LENGTH 8

/* Defines to extract task ID fields */
#define TASK_GET_THREAD_ID(tASKiD) (itti_desc.tasks_info[tASKiD].thread)
/* Extract the instance from a message */
#define ITTI_MESSAGE_GET_INSTANCE(mESSAGE) ((mESSAGE)->ittiMsgHeader.instance)

#include <messages_types.h>

/* This enum defines messages ids. Each one is unique. */
typedef enum {
#define MESSAGE_DEF(iD, sTRUCT, fIELDnAME) iD,
#include <messages_def.h>
#undef MESSAGE_DEF

  MESSAGES_ID_MAX,
} MessagesIds;

//! Thread id of each task
typedef enum {
  THREAD_NULL = 0,

#define TASK_DEF(tHREADiD) THREAD_##tHREADiD,
#include <tasks_def.h>
#undef TASK_DEF

  THREAD_MAX,
  THREAD_FIRST = 1,
} thread_id_t;

//! Sub-tasks id, to defined offset form thread id
typedef enum {
#define TASK_DEF(tHREADiD) tHREADiD##_THREAD = THREAD_##tHREADiD,
#include <tasks_def.h>
#undef TASK_DEF
} task_thread_id_t;

//! Tasks id of each task
typedef enum {
  TASK_UNKNOWN = 0,

#define TASK_DEF(tHREADiD) tHREADiD,
#include <tasks_def.h>
#undef TASK_DEF

  TASK_MAX,
  TASK_FIRST = 1,
} task_id_t;

typedef union msg_s {
#define MESSAGE_DEF(iD, sTRUCT, fIELDnAME) sTRUCT fIELDnAME;
#include <messages_def.h>
#undef MESSAGE_DEF
} msg_t;

typedef uint16_t MessageHeaderSize;

/** @struct MessageHeader
 *  @brief Message Header structure for inter-task communication.
 */
typedef struct MessageHeader_s {
  union {
    struct {
      MessagesIds
          messageId; /**< Unique message id as referenced in enum MessagesIds */

      task_id_t originTaskId;      /**< ID of the sender task */
      task_id_t destinationTaskId; /**< ID of the destination task */
      struct timespec timestamp;   /** Time msg got created */
      instance_t instance;         /**< Task instance for virtualization */
      imsi64_t imsi;               /** IMSI associated to sender task */
      long last_hop_latency;       /** Last hop zmq latency */

      MessageHeaderSize
          ittiMsgSize; /**< Message size (not including header size) */
    };
    // Add padding to avoid any holes in MessageDef object.
    uint8_t __pad[64];
  };
} MessageHeader;

/** @struct MessageDef
 *  @brief Message structure for inter-task communication.
 *  \internal
 * The memory allocation code expects \ref ittiMsg directly following \ref
 * ittiMsgHeader.
 */
typedef struct MessageDef_s {
  MessageHeader ittiMsgHeader; /**< Message header */
  msg_t
      ittiMsg; /**< Union of payloads as defined in x_messages_def.h headers */
} MessageDef;

#endif /* INTERTASK_INTERFACE_TYPES_H_ */
/* @} */
