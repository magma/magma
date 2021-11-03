#include "amf_map.h"

namespace magma5g {

/***************************************************************************
**                                                                        **
** Name:    map_rc_code2string()                                          **
**                                                                        **
** Description: This converts the map_rc_t, return code to string         **
**                                                                        **
***************************************************************************/

std::string map_rc_code2string(map_rc_t rc) {
  switch (rc) {
    case MAP_OK:
      return "MAP_OK";
      break;

    case MAP_KEY_NOT_EXISTS:
      return "MAP_KEY_NOT_EXISTS";
      break;

    case MAP_SEARCH_NO_RESULT:
      return "MAP_SEARCH_NO_RESULT";
      break;

    case MAP_KEY_ALREADY_EXISTS:
      return "MAP_KEY_ALREADY_EXISTS";
      break;

    case MAP_BAD_PARAMETER_KEY:
      return "MAP_BAD_PARAMETER_KEY";
      break;

    case MAP_BAD_PARAMETER_VALUE:
      return "MAP_BAD_PARAMETER_VALUE";
      break;

    case MAP_EMPTY:
      return "MAP_EMPTY";
      break;

    case MAP_DUMP_FAIL:
      return "MAP_DUMP_FAIL";
      break;

    default:
      return "UNKNOWN map_rc_t";
  }
}
}  // namespace magma5g
