/*-
 * Eurecom 2015.
 */
#ifndef _COMPARE_H_
#define _COMPARE_H_

#include <asn_application.h>

#ifdef __cplusplus
extern "C" {
#endif

struct asn_TYPE_descriptor_s; /* Forward declaration */

typedef asn_comp_rval_t *(type_compare_f)(
  struct asn_TYPE_descriptor_s *type_descriptor1,
  void *struct_ptr1,
  struct asn_TYPE_descriptor_s *type_descriptor2,
  void *struct_ptr2);

#ifdef __cplusplus
}
#endif

#endif /* _COMPARE_H_ */
