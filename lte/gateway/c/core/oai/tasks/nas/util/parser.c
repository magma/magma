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

/*****************************************************************************
  Source    parser.c

  Version   0.1

  Date    2012/02/27

  Product   NAS stack

  Subsystem Utilities

  Author    Frederic Maurel

  Description Usefull command line parser

*****************************************************************************/

#include "parser.h"
#include "nas/commonDef.h"

#include <stdio.h>   // fprintf
#include <string.h>  // strcmp, strncpy

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:  parser_print_usage()                                      **
 **                                                                        **
 ** Description: Displays command line usage                               **
 **                                                                        **
 ** Inputs:  command_line:  Pointer to the command line structure      **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   Return:  None                                       **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
void parser_print_usage(const parser_command_line_t* command_line) {
  fprintf(stderr, "Usage: %s", command_line->name);

  for (int i = 0; i < command_line->nb_options; i++) {
    fprintf(
        stderr, " [%s %s]", command_line->options[i].name,
        command_line->options[i].argument);
  }

  fprintf(stderr, "\n");

  for (int i = 0; i < command_line->nb_options; i++) {
    fprintf(
        stderr, "\t%s\t%s\t(%s)\n", command_line->options[i].argument,
        command_line->options[i].usage, command_line->options[i].value);
  }
}

/****************************************************************************
 **                                                                        **
 ** Name:  parser_get_options()                                      **
 **                                                                        **
 ** Description: Parses the command line options                           **
 **                                                                        **
 ** Inputs:  argc:    Number of command line options             **
 **      argv:    Pointer to the command line                **
 **      command_line:  Pointer to the command line structure      **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   Return:  RETURNerror, RETURNok                      **
 **      command_line:  Pointer to the command line structure      **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int parser_get_options(
    int argc, const char** argv, parser_command_line_t* command_line) {
  int argument_not_found, option_not_found = 1;
  int option_length;

  /*
   * Initialize command line structure
   */
  strncpy(command_line->name, argv[0], PARSER_COMMAND_NAME_SIZE);

  for (int i = 0; i < command_line->nb_options; i++) {
    if (strcmp(command_line->options[i].value, "NULL")) {
      command_line->options[i].pvalue = command_line->options[i].value;
    }
  }

  /*
   * No command line options
   */
  if (argc == 1) {
    return RETURNok;
  }

  /*
   * Parse the command line looking for options
   */
  while (*++argv && **argv == '-') {
    argument_not_found = 1;
    option_not_found   = 1;

    for (int i = 0; i < command_line->nb_options; i++) {
      if (!strcmp(command_line->options[i].name, *argv)) {
        option_not_found = 0;

        if (argv[1] && *argv[1] != '-') {
          argument_not_found = 0;
          option_length      = (int) strlen(*++argv);

          if (option_length >= PARSER_OPTION_VALUE_SIZE) {
            fprintf(
                stderr,
                "%s: option name too long (%d), should be less than %d "
                "characters\n",
                *argv, option_length, PARSER_OPTION_VALUE_SIZE);
            return RETURNerror;
          }

          strcpy(command_line->options[i].value, *argv);

          if (command_line->options[i].pvalue == NULL) {
            command_line->options[i].pvalue = command_line->options[i].value;
          }
        }

        break;
      }
    }

    if (option_not_found) {
      if (strcmp(*argv, "-?") && strcmp(*argv, "-h") &&
          strcmp(*argv, "-help") && strcmp(*argv, "--help")) {
        fprintf(stderr, "%s: illegal option %s\n", command_line->name, *argv);
        return RETURNerror;
      }
    } else if (argument_not_found) {
      fprintf(
          stderr, "%s: option %s requires an argument\n", command_line->name,
          *argv);
      return RETURNerror;
    }
  }

  if (option_not_found) {
    return RETURNerror;
  }

  return RETURNok;
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
