000287 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0588        SGi MTU (read)........: 1500
000288 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0591        NAT ..................: true
000289 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0594        User TCP MSS clamping : true
000290 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0597        User IP masquerading  : false
000291 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0600    - PCEF support ...........: false (in development)
000292 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0646    - DNS Configuration:
000293 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0649        IPv4 Primary Address ..........: 8.8.8.8
000294 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0652        IPv4 Secondary Address ..........: 8.8.4.4
000295 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0653    - Helpers:
000296 Mon Dec 21 17:58:58 2020 7F27EC323C80 INFO  SPGW-A tasks/sgw/pgw_config.c          :0656        Push PCO (DNS+MTU) ........: false
000297 Mon Dec 21 17:59:07 2020 7F27D1F26700 INFO  SCTP   tasks/sctp/sctp_itti_messaging.c:0081     ppid NGAP in sctp_itti_send_new_association 
000298 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 118 
000299 Mon Dec 21 17:59:07 2020 7F27D1F26700 INFO  SCTP   tasks/sctp/sctp_itti_messaging.c:0109    ppid NGAP in sctp_itti_send_new_message_ind 
000300 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 116 
000301 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0177    NGAP_AMF_TEST: decode new buffer
000302 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000303 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 3]  handler for procedureCode 15 in originating message
000304 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:3 
000305 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0094      NGAP_AMF_TEST: Received NGAP INITIAL_UE_MESSAGE GNB_UE_NGAP_ID 0x000001
000306 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0107      NGAP_AMF_TEST: New Initial UE message received with gNB UE NGAP ID: 0x000001
000307 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0171         Sending Initial UE Message to AMF_APP: ID: 96, NGAP_INITIAL_UE_MESSAGE: 96 
000308 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0175         Sending Initial UE Message to AMF_APP 
000309 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0225         ####ACL_TAG iniUEmsg sent to TASK_AMF_APP
000311 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0068    AMF_TEST: NGAP_INITIAL_UE_MESSAGE received
000312 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0461     AMF_TEST: AMF_APP_INITIAL_UE_MESSAGE from NGAP,without S-TMSI. 
000313 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0467    AMF_TEST: UE context doesn't exist -> create one
000310 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0181    NGAP_AMF_TEST: decode new buffer done
000314 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0516    AMF_TEST: Sending NAS Establishment Indication to NAS for ue_id = (1)
000315 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000316 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 26
000317 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000318 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_recv.cpp          :0065    AMF_TEST: Processing REGITRATION_REQUEST message
000319 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/nas_proc.cpp          :0245    AMF_TEST: Sending AS IDENTITY_REQUEST
000320 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0806    AMF_TEST: Sending IDENTITY_REQUEST to UE
000321 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0968    AMF_TEST: Start NAS encoding
000322 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0982    AMF_TEST: NAS Encoding Success
000323 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_transport.cpp :0139    AMF_TEST: sending downlink message to NGAP
000324 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 136 
000325 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0458    recieved message ID 136:AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION
000326 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 95 
000327 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0228    ########ACL_TAG :handle_message, NGAP_NAS_DL_DATA_REQ
000328 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:0488    ########ACL_TAG :ngap_generate_downlink_nas_transport, NGAP_NAS_DL_DATA_REQ
000329 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0263    NGAP_AMF_TEST: sending IDENTITY_RESPONSE to AMF_APP
000330 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0373    NGAP_AMF_TEST: decode of new buffer DONE
000331 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000332 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 106576]  handler for procedureCode 46 in originating message
000333 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:106576 
000334 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0302      NGAP_AMF_TEST: Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID 0x00000001
000335 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0085   [0]   Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (1) 
000336 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0101      ngap_amf_itti_nas_uplink_ind, ########send to AMF :101
000337 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0077    AMF_TEST: UPLINK_NAS_MESSAGE received
000338 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0531    AMF_TEST: Received NAS UPLINK DATA IND from NGAP
000339 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000340 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 21
000341 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000342 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_recv.cpp          :0272    AMF_TEST: Received IDENTITY_RESPONSE message
000343 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0323    imsi : 3 3
000344 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0325    imsi : 1 1
000345 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0327    imsi : 0 0
000346 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0329    imsi : 4 4
000347 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0331    imsi : 1 1
000348 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0333    imsi : 0 0
000349 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0335    imsi : 1 1
000350 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0337    imsi : 0 0
000351 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0339    imsi : 0 0
000352 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0341    imsi : 0 0
000353 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0343    imsi : 0 0
000354 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0345    imsi : 0 0
000355 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0347    imsi : 0 0
000356 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0349    imsi : 0 0
000357 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_recv.cpp          :0351    imsi : 1 1
000358 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_identity.cpp      :0237    AMF-TEST: Identification procedure complete for (ue_id=0x00000001)
000359 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/Registration.cpp      :0388    AMF_TEST: Identification procedure success
000360 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NGAP   tasks/amf/amf_authentication.cpp:0313    AMF_TEST: starting Authentication procedure
000361 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0820    AMF_TEST: Sending AUTHENTICATION_REQUEST to UE
000362 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0968    AMF_TEST: Start NAS encoding
000363 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0982    AMF_TEST: NAS Encoding Success
000364 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_transport.cpp :0139    AMF_TEST: sending downlink message to NGAP
000365 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 95 
000366 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0228    ########ACL_TAG :handle_message, NGAP_NAS_DL_DATA_REQ
000367 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:0488    ########ACL_TAG :ngap_generate_downlink_nas_transport, NGAP_NAS_DL_DATA_REQ
000368 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0278    NGAP_AMF_TEST: sending AUTHENTICATION_RESPONSE to AMF_APP
000369 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0373    NGAP_AMF_TEST: decode of new buffer DONE
000370 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000371 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 147824]  handler for procedureCode 46 in originating message
000372 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:147824 
000373 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0302      NGAP_AMF_TEST: Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID 0x00000001
000374 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0085   [0]   Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (1) 
000375 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0101      ngap_amf_itti_nas_uplink_ind, ########send to AMF :101
000376 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0077    AMF_TEST: UPLINK_NAS_MESSAGE received
000377 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0531    AMF_TEST: Received NAS UPLINK DATA IND from NGAP
000378 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000379 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 24
000380 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000381 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_recv.cpp          :0418    AMF_TEST: Received AUTHENTICATION_RESPONSE message
000382 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_authentication.cpp:0546    AMF_TEST: Authentication  procedures complete for (ue_id=0x00000001)
000383 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_authentication.cpp:0614    AMF_TEST: Successful authentication of the UE
000384 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/Registration.cpp      :0087    AMF_TEST Authentication procedure success and start  Security modecommand procedures
000385 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_security_mode_cont:0313    AMF_TEST: Initiating security mode control procedure, KSI = 0
000386 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0871    AMF_TEST: Sending SECURITY_MODE_COMMAND to UE
000387 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0968    AMF_TEST: Start NAS encoding
000388 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0982    AMF_TEST: NAS Encoding Success
000389 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_transport.cpp :0139    AMF_TEST: sending downlink message to NGAP
000390 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 95 
000391 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0228    ########ACL_TAG :handle_message, NGAP_NAS_DL_DATA_REQ
000392 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:0488    ########ACL_TAG :ngap_generate_downlink_nas_transport, NGAP_NAS_DL_DATA_REQ
000393 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0293    NGAP_AMF_TEST: sending SECURITY_MODE_COMPLETE to AMF_APP
000394 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0373    NGAP_AMF_TEST: decode of new buffer DONE
000395 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000396 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 133552]  handler for procedureCode 46 in originating message
000397 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:133552 
000398 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0302      NGAP_AMF_TEST: Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID 0x00000001
000399 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0085   [0]   Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (1) 
000400 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0101      ngap_amf_itti_nas_uplink_ind, ########send to AMF :101
000401 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0077    AMF_TEST: UPLINK_NAS_MESSAGE received
000402 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0531    AMF_TEST: Received NAS UPLINK DATA IND from NGAP
000403 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000404 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 6
000405 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000406 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_Security_Mode.cpp :0065    AMF_TEST: Security mode procedures complete for (ue_id=0x00000001)
000407 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/Registration.cpp      :0552    AMF_TEST: ue_id=0x00000001Start REGISTRATION_ACCEPT procedures for UE 
000408 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_as.cpp            :1032    AMF_TEST: Send AS connection establish confirmation for (ue_id = 1)
000409 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :1082    AMF_TEST: Sending REGISTRATION_ACCEPT to UE
000410 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :1185    AMF_TEST: start NAS encoding 
000411 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :1195    AMF_TEST: NAS encoding success
000412 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_transport.cpp :0139    AMF_TEST: sending downlink message to NGAP
000413 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 95 
000414 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0228    ########ACL_TAG :handle_message, NGAP_NAS_DL_DATA_REQ
000415 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:0488    ########ACL_TAG :ngap_generate_downlink_nas_transport, NGAP_NAS_DL_DATA_REQ
000416 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0316    NGAP_AMF_TEST: sending REGISTRATION_COMPLETE to AMF_APP
000417 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0373    NGAP_AMF_TEST: decode of new buffer DONE
000418 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000419 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 131856]  handler for procedureCode 46 in originating message
000420 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:131856 
000421 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0302      NGAP_AMF_TEST: Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID 0x00000001
000422 Mon Dec 21 17:59:07 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0085   [0]   Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (1) 
000423 Mon Dec 21 17:59:07 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0101      ngap_amf_itti_nas_uplink_ind, ########send to AMF :101
000424 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0077    AMF_TEST: UPLINK_NAS_MESSAGE received
000425 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0531    AMF_TEST: Received NAS UPLINK DATA IND from NGAP
000426 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000427 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 6
000428 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000429 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/Registration.cpp      :0950    AMFAS-SAP - Received Registration Complete message for ue_id = (1)
000430 Mon Dec 21 17:59:07 2020 7F27D5836700 INFO  NAS-AM tasks/amf/Registration.cpp      :1061     Sending AMF INFORMATION for ue_id = (1)
000431 Mon Dec 21 17:59:09 2020 7F27D1F26700 INFO  SCTP   tasks/sctp/sctp_itti_messaging.c:0134    ppid NGAP in sctp_itti_send_com_down_ind 
000432 Mon Dec 21 17:59:11 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0394    NGAP_AMF_TEST: sending PDU_SEESION_ESTABLISHMENT_REQUEST to AMF_APP
000433 Mon Dec 21 17:59:11 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0422    NGAP_AMF_TEST: decode new buffer DONE
000434 Mon Dec 21 17:59:11 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000435 Mon Dec 21 17:59:11 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 131856]  handler for procedureCode 46 in originating message
000436 Mon Dec 21 17:59:11 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:131856 
000437 Mon Dec 21 17:59:11 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0302      NGAP_AMF_TEST: Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID 0x00000004
000438 Mon Dec 21 17:59:11 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0085   [0]   Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (4) 
000439 Mon Dec 21 17:59:11 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0101      ngap_amf_itti_nas_uplink_ind, ########send to AMF :101
000440 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0077    AMF_TEST: UPLINK_NAS_MESSAGE received
000441 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0531    AMF_TEST: Received NAS UPLINK DATA IND from NGAP
000442 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000443 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 15
000444 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000445 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_as.cpp            :0248    AMF_TEST: Processing UL NAS Transport Message
000446 Mon Dec 21 17:59:11 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 119 
000447 Mon Dec 21 17:59:11 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0458    recieved message ID 119:SCTP_CLOSE_ASSOCIATION
000448 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_smf_send.cpp      :0049    AMF SMF Handler- Received PDN Connectivity Request message 
000449 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0041    AMF_TEST: check-1
000450 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0077    AMF_TEST: check-1
000451 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0052    AMF_TEST: check-1
000452 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0066    AMF_TEST: buff_ip:c0
000453 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0066    AMF_TEST: buff_ip:a8
000454 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0066    AMF_TEST: buff_ip:80
000455 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0066    AMF_TEST: buff_ip:6e
000456 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0072    AMF_TEST: call async set_smf_session()
000457 Mon Dec 21 17:59:11 2020 7F27D98E5700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0074    AMF_TEST: after set_smf_session()()
000458 Mon Dec 21 17:59:11 2020 7F27D5836700 INFO  AMF-AP tasks/amf/prepare_request_for_sm:0082    AMF_TEST: check-last
000459 Mon Dec 21 17:59:15 2020 7F27CFE81700 INFO  UTIL   tasks/grpc_service/AmfServiceImp:0053    Received  GRPC SetSMSessionContextAccess request
000460 Mon Dec 21 17:59:15 2020 7F27CFE81700 INFO  UTIL   tasks/grpc_service/amf_service_h:0026    Sending itti_n11_create_pdu_session_response to AMF 
000461 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0086    AMF_TEST: session created for imsi:310410100000001 with IP:À¨€n 
000462 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0089    AMF_TEST: IP:c0 
000463 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0089    AMF_TEST: IP:a8 
000464 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0089    AMF_TEST: IP:80 
000465 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0089    AMF_TEST: IP:6e 
000466 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0089    AMF_TEST: IP:0 
000467 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0089    AMF_TEST: IP:0 
000468 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_pdu_resource_s:0132    PDU session resource setup request message construction to NGAP
000469 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_pdu_resource_s:0235    Converting pdu_session_resource_setup_request_transfer_t to bstring 
000470 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0685    payload_container.len:30 
000471 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0686    encode success, sent dl packet to NGAP
000472 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0696    bytes:33 
000473 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0698    encode success, sent dl packet to NGAP
000474 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_transport.cpp :0139    AMF_TEST: sending downlink message to NGAP
000475 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 102 
000476 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0208     ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000477 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1028       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000478 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1038       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000479 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1058       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000480 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1068       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000481 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1074       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000482 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1081       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000483 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1097       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000484 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1109       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ items: 1
000485 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1116       #####  NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000486 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1119       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000487 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1127       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ :30
000488 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000489 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000490 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :03
000491 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000492 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :8b
000493 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000494 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :0a
000495 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :01
000496 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :f0
000497 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :7f
000498 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000499 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000500 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :02
000501 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000502 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000503 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000504 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :06
000505 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000506 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :86
000507 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000508 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :01
000509 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000510 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000511 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :88
000512 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000513 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :07
000514 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000515 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :01
000516 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000517 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000518 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :ff
000519 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000520 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1138       ##### ies: :00
000521 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1141       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000522 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1145       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000523 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0052      1
000524 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0085      ######ACL_TAG: ngap_amf_encode_initiating, 85  
000525 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0087      ######ACL_TAG: ngap_amf_encode_initiating, 87  
000526 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      01 
000527 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000528 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000529 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000530 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000531 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000532 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000533 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000534 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      1d 
000535 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000536 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000537 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000538 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000539 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000540 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000541 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000542 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000543 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000544 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000545 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000546 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000547 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000548 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000549 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000550 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      0c 
000551 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000552 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000553 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000554 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000555 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000556 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000557 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000558 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      b0 
000559 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      2c 
000560 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      0c 
000561 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000562 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      30 
000563 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      60 
000564 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000565 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000566 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      03 
000567 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000568 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000569 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000570 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      04 
000571 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000572 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000573 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000574 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000575 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000576 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000577 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000578 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000579 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000580 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000581 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000582 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000583 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000584 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000585 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000586 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000587 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000588 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000589 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000590 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000591 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000592 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000593 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000594 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000595 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000596 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000597 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000598 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000599 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000600 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000601 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000602 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000603 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000604 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000605 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000606 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000607 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000608 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000609 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000610 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000611 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000612 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000613 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000614 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000615 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000616 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000617 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000618 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000619 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000620 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000621 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000622 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000623 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000624 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000625 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000626 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000627 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000628 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000629 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000630 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000631 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000632 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000633 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000634 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000635 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000636 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000637 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000638 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000639 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000640 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000641 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000642 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000643 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000644 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000645 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000646 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000647 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000648 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000649 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000650 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000651 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000652 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000653 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000654 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000655 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000656 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000657 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000658 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000659 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000660 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000661 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000662 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000663 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000664 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000665 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000666 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000667 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000668 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000669 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000670 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000671 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000672 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000673 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000674 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000675 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000676 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000677 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000678 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000679 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000680 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000681 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000682 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000683 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000684 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000685 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000686 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000687 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000688 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000689 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000690 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000691 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000692 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000693 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000694 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000695 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000696 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000697 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000698 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000699 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000700 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000701 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000702 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0102      ######ACL_TAG: ngap_amf_encode_initiating, 102  
000703 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0123      buf:0x7f27d16f0430
000704 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0124      *buf:0x606000159b60
000705 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0125      l:62
000706 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000707 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      1d 
000708 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000709 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      3a 
000710 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000711 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000712 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      03 
000713 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000714 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      0a 
000715 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000716 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      02 
000717 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000718 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
000719 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000720 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      55 
000721 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000722 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      02 
000723 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000724 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
000725 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000726 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      4a 
000727 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000728 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      27 
000729 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000730 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000731 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
000732 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      06 
000733 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000734 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      21 
000735 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000736 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000737 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      03 
000738 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000739 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      8b 
000740 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000741 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      0a 
000742 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
000743 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      f0 
000744 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      7f 
000745 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000746 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000747 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      02 
000748 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000749 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000750 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000751 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      06 
000752 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000753 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      86 
000754 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000755 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
000756 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000757 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000758 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      88 
000759 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000760 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      07 
000761 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000762 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
000763 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000764 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000765 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      ff 
000766 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000767 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
000768 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0054      ####ACL_TAG
000769 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0074      ####ACL_TAG
000770 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0076      ####ACL_TAG
000771 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1153       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000772 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0062      ######ACL_TAG: ngap_amf_itti_send_sctp_request, 62 
000773 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1166       ##### NGAP_PDUSESSION_RESOURCE_SETUP_REQ
000774 Mon Dec 21 17:59:15 2020 7F27D5012700 INFO  SCTP   tasks/sctp/sctp_primitives_serve:0116    ppid NGAP in sctp_itti_send_lower_layer_conf 
000775 Mon Dec 21 17:59:15 2020 7F27D5012700 INFO  SCTP   tasks/sctp/sctp_itti_messaging.c:0048     ppid NGAP in sctp_itti_send_lower_layer_conf 
000776 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 95 
000777 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0228    ########ACL_TAG :handle_message, NGAP_NAS_DL_DATA_REQ
000778 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:0488    ########ACL_TAG :ngap_generate_downlink_nas_transport, NGAP_NAS_DL_DATA_REQ
000779 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0342    NGAP_AMF_TEST: sending PDU_SESSION_RELEASE_REQUEST to AMF_APP
000780 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0373    NGAP_AMF_TEST: decode of new buffer DONE
000781 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0181    ######ACL_TAG: ngap_amf_handle_message, 181  
000782 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0212    NGAP: [SCTP 236176]  handler for procedureCode 46 in originating message
000783 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_handlers.c  :0214    #####ACL_TAG :ASSOC_ID:236176 
000784 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_nas_procedur:0302      NGAP_AMF_TEST: Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID 0x00000004
000785 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf_itti_messagi:0085   [0]   Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (4) 
000786 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0101      ngap_amf_itti_nas_uplink_ind, ########send to AMF :101
000787 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_main.cpp      :0077    AMF_TEST: UPLINK_NAS_MESSAGE received
000788 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_handler.cpp   :0531    AMF_TEST: Received NAS UPLINK DATA IND from NGAP
000789 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0183    AMF_TEST: Decoding NAS Message
000790 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0187    AMF_TEST: rc = 11
000791 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_as.cpp            :0201    AMF_TEST: NAS Decode Success
000792 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  NAS-AM tasks/amf/amf_as.cpp            :0248    AMF_TEST: Processing UL NAS Transport Message
000793 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_pdu_resource_s:0258    PDU session resource release request message construction to NGAP
000794 Mon Dec 21 17:59:15 2020 7F27D5836700 INFO  AMF-AP tasks/amf/amf_app_pdu_resource_s:0287    Converting pdu_session_resource_release_command_transfer to bstring
000795 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 117 
000796 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0458    recieved message ID 117:SCTP_DATA_CNF
000797 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 104 
000798 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf.c           :0215     ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000799 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1178       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000800 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1195       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000801 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1212       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000802 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1227       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000803 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1240       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000804 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1250       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000805 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1272       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000806 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1274       ##### items:1
000807 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :01
000808 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000809 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000810 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000811 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000812 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000813 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000814 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000815 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000816 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000817 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000818 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000819 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000820 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000821 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000822 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000823 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000824 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000825 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000826 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000827 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000828 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000829 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000830 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000831 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000832 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000833 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000834 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000835 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000836 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000837 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000838 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000839 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000840 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000841 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000842 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000843 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000844 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000845 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000846 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1288       ##### ies: :00
000847 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1292       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
000848 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0052      1
000849 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0085      ######ACL_TAG: ngap_amf_encode_initiating, 85  
000850 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0087      ######ACL_TAG: ngap_amf_encode_initiating, 87  
000851 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      01 
000852 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000853 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000854 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000855 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000856 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000857 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000858 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000859 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      1c 
000860 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000861 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000862 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000863 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000864 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000865 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000866 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000867 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      01 
000868 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000869 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000870 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000871 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000872 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000873 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000874 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000875 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      0b 
000876 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000877 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000878 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000879 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000880 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000881 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000882 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000883 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      90 
000884 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      1c 
000885 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      0c 
000886 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000887 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      30 
000888 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      60 
000889 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000890 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000891 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      03 
000892 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000893 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000894 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000895 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      04 
000896 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000897 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000898 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000899 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000900 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000901 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000902 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000903 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000904 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000905 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000906 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000907 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000908 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000909 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000910 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000911 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000912 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000913 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000914 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000915 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000916 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000917 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000918 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000919 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000920 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000921 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000922 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000923 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000924 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000925 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000926 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000927 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000928 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000929 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000930 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000931 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000932 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000933 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000934 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000935 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000936 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000937 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000938 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000939 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000940 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000941 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000942 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000943 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000944 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000945 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000946 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000947 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000948 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000949 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000950 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000951 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000952 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000953 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000954 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000955 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000956 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000957 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000958 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000959 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000960 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000961 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000962 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000963 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000964 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000965 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000966 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000967 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000968 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000969 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000970 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000971 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000972 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000973 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000974 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000975 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000976 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000977 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000978 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000979 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000980 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000981 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000982 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000983 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000984 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000985 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000986 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000987 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000988 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000989 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000990 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000991 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000992 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000993 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000994 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000995 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000996 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000997 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000998 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
000999 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001000 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001001 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001002 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001003 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001004 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001005 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001006 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001007 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001008 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001009 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001010 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001011 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001012 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001013 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001014 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001015 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001016 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001017 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001018 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001019 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001020 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001021 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001022 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001023 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001024 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001025 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001026 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0091      00 
001027 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0102      ######ACL_TAG: ngap_amf_encode_initiating, 102  
001028 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0123      buf:0x7f27d16f03c0
001029 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0124      *buf:0x60c00006ca00
001030 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0125      l:67
001031 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001032 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      1c 
001033 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      40 
001034 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      3f 
001035 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001036 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001037 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      03 
001038 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001039 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      0a 
001040 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001041 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      02 
001042 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001043 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
001044 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001045 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      55 
001046 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001047 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      02 
001048 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001049 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
001050 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001051 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      4f 
001052 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001053 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      2c 
001054 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001055 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001056 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
001057 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      28 
001058 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      01 
001059 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001060 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001061 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001062 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001063 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001064 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001065 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001066 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001067 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001068 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001069 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001070 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001071 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001072 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001073 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001074 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001075 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001076 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001077 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001078 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001079 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001080 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001081 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001082 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001083 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001084 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001085 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001086 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001087 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001088 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001089 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001090 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001091 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001092 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001093 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001094 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001095 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001096 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001097 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0127      00 
001098 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0054      ####ACL_TAG
001099 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0074      ####ACL_TAG
001100 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_encoder.c   :0076      ####ACL_TAG
001101 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1298       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
001102 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1309       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
001103 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_itti_messagi:0062      ######ACL_TAG: ngap_amf_itti_send_sctp_request, 62 
001104 Mon Dec 21 17:59:15 2020 7F27D16F1700 ERROR NGAP   tasks/ngap/ngap_amf_nas_procedur:1312       ##### NGAP_PDUSESSIONRESOURCE_REL_REQ
001105 Mon Dec 21 17:59:15 2020 7F27D5012700 INFO  SCTP   tasks/sctp/sctp_primitives_serve:0116    ppid NGAP in sctp_itti_send_lower_layer_conf 
001106 Mon Dec 21 17:59:15 2020 7F27D5012700 INFO  SCTP   tasks/sctp/sctp_itti_messaging.c:0048     ppid NGAP in sctp_itti_send_lower_layer_conf 
001107 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0144    NGAP_AMF_TEST : inside handler received message with imsi 13744632839234567870  and message type 117 
001108 Mon Dec 21 17:59:15 2020 7F27D16F1700 INFO  NGAP   tasks/ngap/ngap_amf.c           :0458    recieved message ID 117:SCTP_DATA_CNF
