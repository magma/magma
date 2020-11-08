#pragma once

#include "messages_types.h"

/*
 * Sends N11_CREATE_PDU_SESSION_RESPONSE message to AMF.
 */
int send_n11_create_pdu_session_resp_itti(
    itti_n11_create_pdu_session_response_t* itti_msg);
