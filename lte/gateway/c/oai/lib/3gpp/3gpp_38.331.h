
#pragma once

// could be extracted with asn1 tool

typedef enum m5g_EstablishmentCause {
  M5G_EMERGENCY = 1,
  M5G_HIGH_PRIORITY_ACCESS,
  M5G_MT_ACCESS,
  M5G_MO_SIGNALLING,
  M5G_MO_DATA,
  M5G_MO_VIDEOCALL,
  M5G_MO_SMS,
  M5G_MPS_PRIORITYACCESS,
  M5G_MCS_PRIORITYACCESS,
  M5G_SPARE6,
  M5G_SPARE5,
  M5G_SPARE4,
  M5G_SPARE3,
  M5G_SPARE2,
  M5G_SPARE1,
} m5g_rrc_establishment_cause_t;

/*EstablishmentCause ::= ENUMERATED {
emergency, highPriorityAccess, mt-Access, mo-Signalling,
mo-Data, mo-VoiceCall, mo-VideoCall, mo-SMS, mps-PriorityAccess,
mcs-PriorityAccess, spare6, spare5, spare4, spare3, spare2, spare1}*/
