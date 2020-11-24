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

/********************************************************************************
 *
 * !!! This header should only be included by the file that initialize
 *     the intertask interface module for the process !!!
 *
 * Other files should include "intertask_interface.h"
 *
 *******************************************************************************/

#ifndef INTERTASK_INTERFACE_INIT_H_
#define INTERTASK_INTERFACE_INIT_H_

#include "intertask_interface.h"

#ifndef CHECK_PROTOTYPE_ONLY

/* Map task id to printable name. */
const task_info_t tasks_info[] = {
    {0, "TASK_UNKNOWN", "ipc:///run/IPC_TASK_UNKNOWN"},
#define TASK_DEF(tHREADiD)                                                     \
  {tHREADiD##_THREAD, #tHREADiD, "ipc:///run/IPC_" #tHREADiD},
#include <tasks_def.h>
#undef TASK_DEF
};

/* Map message id to message information */
const message_info_t messages_info[] = {
#define MESSAGE_DEF(iD, sTRUCT, fIELDnAME) {iD, sizeof(sTRUCT), #iD},
#include <messages_def.h>
#undef MESSAGE_DEF
};

#endif

/** \brief Init function for the intertask interface. Init queues, Mutexes and
 * Cond vars.
 * \param thread_max Maximum number of threads
 * \param messages_id_max Maximum message id
 * \param threads_name Pointer on the threads name information as created by
 * this include file
 * \param messages_info Pointer on messages information as created by this
 * include file
 **/
int itti_init(
    task_id_t task_max, thread_id_t thread_max, MessagesIds messages_id_max,
    const task_info_t* tasks_info, const message_info_t* messages_info,
    const char* const messages_definition_xml,
    const char* const dump_file_name);

#endif /* INTERTASK_INTERFACE_INIT_H_ */
/* @} */
