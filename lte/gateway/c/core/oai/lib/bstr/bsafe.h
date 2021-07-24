/*
 * This source file is part of the bstring string library.  This code was
 * written by Paul Hsieh in 2002-2004, and is covered by the BSD open source
 * license. Refer to the accompanying documentation for details on usage and
 * license.
 */

/*
 * bsafe.h
 *
 * This is an optional module that can be used to help enforce a safety
 * standard based on pervasive usage of bstrlib.  This file is not necessarily
 * portable, however, it has been tested to work correctly with Intel's C/C++
 * compiler, WATCOM C/C++ v11.x and Microsoft Visual C++.
 */

#ifndef BSTRLIB_BSAFE_INCLUDE
#define BSTRLIB_BSAFE_INCLUDE

#ifdef __cplusplus
extern "C" {
#endif

#if !defined(__GNUC__) && (!defined(_MSC_VER) || (_MSC_VER <= 1310))
/* This is caught in the linker, so its not necessary for gcc. */
extern char*(gets)(char* buf);
#endif

extern char*(strncpy)(char* dst, const char* src, size_t n);
extern char*(strncat)(char* dst, const char* src, size_t n);
extern char*(strtok)(char* s1, const char* s2);
extern char*(strdup)(const char* s);

#undef strcpy
#undef strcat
#define strcpy(a, b) bsafe_strcpy(a, b)
#define strcat(a, b) bsafe_strcat(a, b)

#ifdef __cplusplus
}
#endif

#endif
