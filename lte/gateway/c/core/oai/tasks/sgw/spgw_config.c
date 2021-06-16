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

/*! \file spgw_config.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define SGW
#define SPGW_CONFIG_C

#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <libconfig.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "spgw_config.h"

void spgw_config_display(spgw_config_t* config_p);

//------------------------------------------------------------------------------
void spgw_config_init(spgw_config_t* config_pP) {
  sgw_config_init(&config_pP->sgw_config);
  pgw_config_init(&config_pP->pgw_config);

  /*
   * Set service 303
   */
  config_pP->service303_config.name = bfromcstr(SERVICE303_SPGW_PACKAGE_NAME);
  config_pP->service303_config.version =
      bfromcstr(SERVICE303_SPGW_PACKAGE_VERSION);
}

//------------------------------------------------------------------------------
static int spgw_config_process(spgw_config_t* config_pP) {
#if (!EMBEDDED_SGW)
  async_system_command(
      TASK_ASYNC_SYSTEM, SPGW_WARN_ON_ERROR, "sysctl -w net.ipv4.ip_forward=1");
  async_system_command(TASK_ASYNC_SYSTEM, SPGW_WARN_ON_ERROR, "sync");
#endif

  if (RETURNok != sgw_config_process(&config_pP->sgw_config)) {
    return RETURNerror;
  }

  if (RETURNok != pgw_config_process(&config_pP->pgw_config)) {
    return RETURNerror;
  }
  return RETURNok;
}

//------------------------------------------------------------------------------
int spgw_config_parse_file(spgw_config_t* config_pP) {
  config_t cfg = {0};
  config_init(&cfg);

  if (config_pP->config_file) {
    /*
     * Read the file. If there is an error, report it and exit.
     */
    if (!config_read_file(&cfg, bdata(config_pP->config_file))) {
      OAILOG_ERROR(
          LOG_SPGW_APP, "%s:%d - %s\n", bdata(config_pP->config_file),
          config_error_line(&cfg), config_error_text(&cfg));
      config_destroy(&cfg);
      Fatal(
          "Failed to parse SP-GW configuration file %s!\n",
          bdata(config_pP->config_file));
    }
  } else {
    OAILOG_ERROR(LOG_SPGW_APP, "No SP-GW configuration file provided!\n");
    config_destroy(&cfg);
    Fatal("No SP-GW configuration file provided!\n");
  }

  OAILOG_INFO(
      LOG_SPGW_APP, "Parsing configuration file provided %s\n",
      bdata(config_pP->config_file));
  if (sgw_config_parse_file(&config_pP->sgw_config) != 0) {
    config_destroy(&cfg);
    return RETURNerror;
  }

  if (pgw_config_parse_file(&config_pP->pgw_config) != 0) {
    config_destroy(&cfg);
    return RETURNerror;
  }

  if (spgw_config_process(config_pP) != 0) {
    return RETURNerror;
  }

  config_destroy(&cfg);
  return RETURNok;
}

//------------------------------------------------------------------------------
void spgw_config_display(spgw_config_t* config_p) {
  sgw_config_display(&config_p->sgw_config);
  pgw_config_display(&config_p->pgw_config);
}

//------------------------------------------------------------------------------
static void usage(char* target) {
  OAILOG_INFO(
      LOG_CONFIG, "==== EURECOM %s version: %s ====\n", PACKAGE_NAME,
      PACKAGE_VERSION);
  OAILOG_INFO(LOG_CONFIG, "Please report any bug to: %s\n", PACKAGE_BUGREPORT);
  OAILOG_INFO(LOG_CONFIG, "Usage: %s [options]\n", target);
  OAILOG_INFO(LOG_CONFIG, "Available options:\n");
  OAILOG_INFO(LOG_CONFIG, "-h      Print this help and return\n");
  OAILOG_INFO(LOG_CONFIG, "-c <path>\n");
  OAILOG_INFO(LOG_CONFIG, "        Set the configuration file for S/P-GW\n");
  OAILOG_INFO(LOG_CONFIG, "        See template in ETC\n");
  OAILOG_INFO(LOG_CONFIG, "-K <file>\n");
  OAILOG_INFO(
      LOG_CONFIG, "        Output intertask messages to provided file\n");
  OAILOG_INFO(
      LOG_CONFIG, "-V      Print %s version and return\n", PACKAGE_NAME);
}
//------------------------------------------------------------------------------
int spgw_config_parse_opt_line(
    int argc, char* argv[], spgw_config_t* spgw_config_p) {
  int c;

  spgw_config_init(spgw_config_p);

  /*
   * Parsing command line
   */
  while ((c = getopt(argc, argv, "c:hi:K:V")) != -1) {
    switch (c) {
      case 'c': {
        /*
         * Store the given configuration file. If no file is given,
         * * * * then the default values will be used.
         */
        spgw_config_p->config_file = blk2bstr(optarg, strlen(optarg));
        spgw_config_p->pgw_config.config_file =
            bstrcpy(spgw_config_p->config_file);
        spgw_config_p->sgw_config.config_file =
            bstrcpy(spgw_config_p->config_file);
        OAILOG_DEBUG(
            LOG_CONFIG, "spgw_config.config_file %s\n",
            bdata(spgw_config_p->config_file));
      } break;

      case 'V': {
        OAILOG_DEBUG(
            LOG_CONFIG,
            "==== EURECOM %s v%s ===="
            "Please report any bug to: %s\n",
            PACKAGE_NAME, PACKAGE_VERSION, PACKAGE_BUGREPORT);
      } break;

      case 'K':
        spgw_config_p->sgw_config.itti_config.log_file =
            blk2bstr(optarg, strlen(optarg));
        OAILOG_DEBUG(
            LOG_CONFIG, "spgw_config.sgw_config.itti_config.log_file %s\n",
            bdata(spgw_config_p->sgw_config.itti_config.log_file));
        break;

      case 'h': /* Fall through */
      default:
        usage(argv[0]);
        exit(0);
    }
  }

  /*
   * Parse the configuration file using libconfig
   */
  if (!spgw_config_p->config_file) {
    spgw_config_p->config_file = bfromcstr("/usr/local/etc/oai/spgw.conf");
    spgw_config_p->pgw_config.config_file =
        bfromcstr("/usr/local/etc/oai/spgw.conf");
    spgw_config_p->sgw_config.config_file =
        bfromcstr("/usr/local/etc/oai/spgw.conf");
  }
  if (spgw_config_parse_file(spgw_config_p) != 0) {
    return RETURNerror;
  }

  /*
   * Display the configuration
   */
  spgw_config_display(spgw_config_p);
  return RETURNok;
}

void free_spgw_config(spgw_config_t* spgw_config_p) {
  OAI_FPRINTF_INFO("Cleaning up SPGW configs");
  free_pgw_config(&spgw_config_p->pgw_config);
  bdestroy_wrapper(&spgw_config_p->config_file);
  bdestroy_wrapper(&spgw_config_p->sgw_config.config_file);
  bdestroy_wrapper(&spgw_config_p->pgw_config.config_file);
}
