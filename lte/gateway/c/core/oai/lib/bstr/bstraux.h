/*
 * This source file is part of the bstring string library.  This code was
 * written by Paul Hsieh in 2002-2015, and is covered by the BSD open source
 * license and the GPL. Refer to the accompanying documentation for details
 * on usage and license.
 */

/*
 * bstraux.h
 *
 * This file is not a necessary part of the core bstring library itself, but
 * is just an auxilliary module which includes miscellaneous or trivial
 * functions.
 */

#ifndef BSTRAUX_INCLUDE
#define BSTRAUX_INCLUDE

#include <time.h>
#include <string.h>

#include "bstrlib.h"

struct bStream;
struct bwriteStream;

#ifdef __cplusplus
extern "C" {
#endif

/* Backward compatibilty with previous versions of Bstrlib */
#if !defined(BSTRLIB_REDUCE_NAMESPACE_POLLUTION)
#define bAssign(a, b) ((bassign)((a), (b)))
#define bTrunc(b, n) ((btrunc)((b), (n)))
#endif

/* Unusual functions */
extern int bSetChar(bstring b, int pos, char c);
extern int bReplicate(bstring b, int n);
extern int bReverse(bstring b);
extern bstring bStrfTime(const char* fmt, const struct tm* timeptr);

/* Writable stream */
typedef int (*bNwrite)(
    const void* buf, size_t elsize, size_t nelem, void* parm);

int bwsWriteBstr(struct bwriteStream* stream, const_bstring b);

#ifdef __cplusplus
}
#endif

#endif
