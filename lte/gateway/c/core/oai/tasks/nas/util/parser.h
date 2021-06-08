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
Source    parser.h

Version   0.1

Date    2012/02/27

Product   NAS stack

Subsystem Utilities

Author    Frederic Maurel

Description Usefull command line parser

*****************************************************************************/
#ifndef FILE_PARSER_SEEN
#define FILE_PARSER_SEEN

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/* -----------
 * Option type
 * -----------
 *  An option is defined with a name, an argument following the option's
 *  name, the usage message and a value
 */
typedef struct {
  const char* name;     /* Option name         */
  const char* argument; /* Argument following the option   */
  const char* usage;    /* Option and Argument usage     */
#define PARSER_OPTION_VALUE_SIZE 32
  char value[PARSER_OPTION_VALUE_SIZE]; /* Option value      */
  char* pvalue;
} parser_option_t;

/* -----------------
 * Command line type
 * -----------------
 *  An command line is defined with a name, the number of options and the
 *  list of command's options
 */
typedef struct {
#define PARSER_COMMAND_NAME_SIZE 32
  char name[PARSER_COMMAND_NAME_SIZE]; /* Command name      */
  const int nb_options;                /* Number of options     */
  parser_option_t options[];           /* Command line options    */
} parser_command_line_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

void parser_print_usage(const parser_command_line_t* commamd_line);
int parser_get_options(
    int argc, const char** argv, parser_command_line_t* commamd_line);

#endif /* FILE_PARSER_SEEN*/
