/*
 * This source file is part of the bstring string library.  This code was
 * written by Paul Hsieh in 2002-2015, and is covered by the BSD open source
 * license and the GPL. Refer to the accompanying documentation for details
 * on usage and license.
 */

/*
 * bstrlib.h
 *
 * This file is the interface for the core bstring functions.
 */

#ifndef BSTRLIB_INCLUDE
#define BSTRLIB_INCLUDE

#ifdef __cplusplus
extern "C" {
#endif

#include <stdarg.h>
#include <string.h>
#include <limits.h>
#include <ctype.h>

struct bStream;

#if !defined(BSTRLIB_VSNP_OK) && !defined(BSTRLIB_NOVSNP)
#if defined(__TURBOC__) && !defined(__BORLANDC__)
#define BSTRLIB_NOVSNP
#endif
#endif

#define BSTR_ERR (-1)
#define BSTR_OK (0)
#define BSTR_BS_BUFF_LENGTH_GET (0)

typedef struct tagbstring* bstring;
typedef const struct tagbstring* const_bstring;

/* Copy functions */
extern bstring bfromcstr(const char* str);
extern bstring bfromcstr_with_str_len(const char* str, int len);
extern bstring bfromcstralloc(int mlen, const char* str);
extern bstring bfromcstrrangealloc(int minl, int maxl, const char* str);
extern bstring blk2bstr(const void* blk, int len);
extern char* bstr2cstr(const_bstring s, char z);
extern int bcstrfree(char* s);
extern bstring bstrcpy(const_bstring b1);
extern int bassign(bstring a, const_bstring b);
extern int bassigncstr(bstring a, const char* str);
extern int bassignblk(bstring a, const void* s, int len);

/* Destroy function */
extern int bdestroy(bstring b);

/* Space allocation hinting functions */
extern int balloc(bstring s, int len);

/* Substring extraction */
extern bstring bmidstr(const_bstring b, int left, int len);

/* Various standard manipulations */
extern int bconcat(bstring b0, const_bstring b1);
extern int bconchar(bstring b0, char c);
extern int bcatcstr(bstring b, const char* s);
extern int bcatblk(bstring b, const void* s, int len);
extern int binsert(bstring s1, int pos, const_bstring s2, unsigned char fill);
extern int binsertblk(
    bstring s1, int pos, const void* s2, int len, unsigned char fill);
extern int breplace(
    bstring b1, int pos, int len, const_bstring b2, unsigned char fill);
extern int bdelete(bstring s1, int pos, int len);
extern int bsetstr(bstring b0, int pos, const_bstring b1, unsigned char fill);
extern int btrunc(bstring b, int n);

/* Scan/search functions */
extern int bstricmp(const_bstring b0, const_bstring b1);
extern int bstrnicmp(const_bstring b0, const_bstring b1, int n);
extern int biseqcaselessblk(const_bstring b, const void* blk, int len);
extern int bisstemeqcaselessblk(const_bstring b0, const void* blk, int len);
extern int biseqblk(const_bstring b, const void* blk, int len);
extern int bisstemeqblk(const_bstring b0, const void* blk, int len);
extern int biseqcstrcaseless(const_bstring b, const char* s);
extern int binstr(const_bstring s1, int pos, const_bstring s2);
extern int bstrchrp(const_bstring b, int c, int pos);
#define bstrchr(b, c) bstrchrp((b), (c), 0)
extern int binchr(const_bstring b0, int pos, const_bstring b1);
extern int bfindreplace(
    bstring b, const_bstring find, const_bstring repl, int pos);

/* List of string container functions */
struct bstrList {
  int qty, mlen;
  bstring* entry;
};
extern int bstrListDestroy(struct bstrList* sl);

/* String split and join functions */
extern struct bstrList* bsplit(const_bstring str, unsigned char splitChar);
extern bstring bjoinblk(const struct bstrList* bl, const void* s, int len);
extern int bsplitcb(
    const_bstring str, unsigned char splitChar, int pos,
    int (*cb)(void* parm, int ofs, int len), void* parm);

/* Miscellaneous functions */
extern int bpattern(bstring b, int len);
extern int btoupper(bstring b);
extern int btolower(bstring b);
extern int btrimws(bstring b);

#if !defined(BSTRLIB_NOVSNP)
extern bstring bformat(const char* fmt, ...);
extern int bformata(bstring b, const char* fmt, ...);
extern int bassignformat(bstring b, const char* fmt, ...);
extern int bvcformata(bstring b, int count, const char* fmt, va_list arglist);

#endif

typedef int (*bNgetc)(void* parm);
typedef size_t (*bNread)(void* buff, size_t elsize, size_t nelem, void* parm);

/* Input functions */
extern bstring bread(bNread readPtr, void* parm);
extern int breada(bstring b, bNread readPtr, void* parm);

/* Stream functions */
extern struct bStream* bsopen(bNread readPtr, void* parm);
extern void* bsclose(struct bStream* s);
extern int bsbufflength(struct bStream* s, int sz);
extern int bsread(bstring b, struct bStream* s, int n);
extern int bsreada(bstring b, struct bStream* s, int n);
extern int bsunread(struct bStream* s, const_bstring b);

struct tagbstring {
  int mlen;
  int slen;
  unsigned char* data;
};

/* Accessor macros */
#define blengthe(b, e)                                                         \
  (((b) == (void*) 0 || (b)->slen < 0) ? (int) (e) : ((b)->slen))
#define blength(b) (blengthe((b), 0))
#define bdataofse(b, o, e)                                                     \
  (((b) == (void*) 0 || (b)->data == (void*) 0) ? (char*) (e) :                \
                                                  ((char*) (b)->data) + (o))
#define bdataofs(b, o) (bdataofse((b), (o), (void*) 0))
#define bdatae(b, e) (bdataofse(b, 0, e))
#define bdata(b) (bdataofs(b, 0))
#define bchare(b, p, e)                                                        \
  ((((unsigned) (p)) < (unsigned) blength(b)) ? ((b)->data[(p)]) : (e))
#define bchar(b, p) bchare((b), (p), '\0')

/* Static constant string initialization macro */
#define bsStaticMlen(q, m)                                                     \
  { (m), (int) sizeof(q) - 1, (unsigned char*) ("" q "") }
#if defined(_MSC_VER)
#define bsStatic(q) bsStaticMlen(q, -32)
#endif
#ifndef bsStatic
#define bsStatic(q) bsStaticMlen(q, -__LINE__)
#endif

/* Static constant block parameter pair */
#define bsStaticBlkParms(q) ((void*) ("" q "")), ((int) sizeof(q) - 1)

#define bcatStatic(b, s) ((bcatblk)((b), bsStaticBlkParms(s)))
#define bfromStatic(s) ((blk2bstr)(bsStaticBlkParms(s)))

/* Reference building macros */
#define blk2tbstr(t, s, l)                                                     \
  {                                                                            \
    (t).data = (unsigned char*) (s);                                           \
    (t).slen = l;                                                              \
    (t).mlen = -1;                                                             \
  }
#define bmid2tbstr(t, b, p, l)                                                 \
  {                                                                            \
    const_bstring bstrtmp_s = (b);                                             \
    if (bstrtmp_s && bstrtmp_s->data && bstrtmp_s->slen >= 0) {                \
      int bstrtmp_left = (p);                                                  \
      int bstrtmp_len  = (l);                                                  \
      if (bstrtmp_left < 0) {                                                  \
        bstrtmp_len += bstrtmp_left;                                           \
        bstrtmp_left = 0;                                                      \
      }                                                                        \
      if (bstrtmp_len > bstrtmp_s->slen - bstrtmp_left)                        \
        bstrtmp_len = bstrtmp_s->slen - bstrtmp_left;                          \
      if (bstrtmp_len <= 0) {                                                  \
        (t).data = (unsigned char*) "";                                        \
        (t).slen = 0;                                                          \
      } else {                                                                 \
        (t).data = bstrtmp_s->data + bstrtmp_left;                             \
        (t).slen = bstrtmp_len;                                                \
      }                                                                        \
    } else {                                                                   \
      (t).data = (unsigned char*) "";                                          \
      (t).slen = 0;                                                            \
    }                                                                          \
    (t).mlen = -__LINE__;                                                      \
  }

#ifdef __cplusplus
}
#endif

#endif
