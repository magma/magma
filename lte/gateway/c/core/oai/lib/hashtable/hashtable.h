/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file hashtable.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_HASH_TABLE_SEEN
#define FILE_HASH_TABLE_SEEN

#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>
#include <stddef.h>

#include "bstrlib.h"

typedef size_t hash_size_t;
typedef uint64_t hash_key_t;

#define HASHTABLE_NOT_A_KEY_VALUE ((uint64_t) -1)

typedef enum hashtable_return_code_e {
  HASH_TABLE_OK = 0,
  HASH_TABLE_INSERT_OVERWRITTEN_DATA,
  HASH_TABLE_KEY_NOT_EXISTS,
  HASH_TABLE_SEARCH_NO_RESULT,
  HASH_TABLE_KEY_ALREADY_EXISTS,
  HASH_TABLE_SAME_KEY_VALUE_EXISTS,
  HASH_TABLE_BAD_PARAMETER_HASHTABLE,
  HASH_TABLE_BAD_PARAMETER_KEY,
  HASH_TABLE_SYSTEM_ERROR,
  HASH_TABLE_CODE_MAX
} hashtable_rc_t;

#define HASH_TABLE_DEFAULT_HASH_FUNC NULL
#define HASH_TABLE_DEFAULT_free_wrapper_FUNC NULL
#define FREE_HASHTABLE_KEY_ARRAY(key_array_ptr)                                \
  do {                                                                         \
    AssertFatal(key_array_ptr, "Trying to free a NULL array pointer");         \
    free(key_array_ptr->keys);                                                 \
    free(key_array_ptr);                                                       \
  } while (0) /*Free the list of keys of a hash table */

typedef struct hash_node_s {
  hash_key_t key;
  void* data;
  struct hash_node_s* next;

} hash_node_t;

typedef struct hash_node_uint64_s {
  hash_key_t key;
  uint64_t data;
  struct hash_node_uint64_s* next;
} hash_node_uint64_t;

typedef struct hash_table_s {
  hash_size_t size;
  hash_size_t num_elements;
  struct hash_node_s** nodes;
  hash_size_t (*hashfunc)(const hash_key_t);
  void (*freefunc)(void**);
  bstring name;
  bool is_allocated_by_malloc;
  bool log_enabled;
} hash_table_t;

typedef struct hash_table_ts_s {
  pthread_mutex_t mutex;
  hash_size_t size;
  hash_size_t num_elements;
  struct hash_node_s** nodes;
  pthread_mutex_t* lock_nodes;
  pthread_mutexattr_t* lock_attr;
  hash_size_t (*hashfunc)(const hash_key_t);
  void (*freefunc)(void**);
  bstring name;
  bool is_allocated_by_malloc;
  bool log_enabled;
} hash_table_ts_t;
typedef struct hash_table_uint64_s {
  hash_size_t size;
  hash_size_t num_elements;
  struct hash_node_uint64_s** nodes;
  hash_size_t (*hashfunc)(const hash_key_t);
  bstring name;
  bool is_allocated_by_malloc;
  bool log_enabled;
} hash_table_uint64_t;

typedef struct hash_table_uint64_ts_s {
  pthread_mutex_t mutex;
  hash_size_t size;
  hash_size_t num_elements;
  struct hash_node_uint64_s** nodes;
  pthread_mutex_t* lock_nodes;
  hash_size_t (*hashfunc)(const hash_key_t);
  bstring name;
  bool is_allocated_by_malloc;
  bool log_enabled;
} hash_table_uint64_ts_t;

typedef struct hashtable_key_array_s {
  int num_keys;
  hash_key_t* keys;
} hashtable_key_array_t;

typedef struct hashtable_element_array_s {
  int num_elements;
  void** elements;
} hashtable_element_array_t;

typedef struct hashtable_uint64_element_array_s {
  int num_elements;
  uint64_t* elements;
} hashtable_uint64_element_array_t;

char* hashtable_rc_code2string(hashtable_rc_t rc);
void hash_free_int_func(void** memory);
hash_table_t* hashtable_init(
    hash_table_t* const hashtbl, const hash_size_t size,
    hash_size_t (*hashfunc)(const hash_key_t), void (*freefunc)(void**),
    bstring display_name_p);
__attribute__((malloc)) hash_table_t* hashtable_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const hash_key_t),
    void (*freefunc)(void**), bstring name_p);
hashtable_rc_t hashtable_destroy(hash_table_t* hashtbl);
hashtable_rc_t hashtable_is_key_exists(
    const hash_table_t* const hashtbl, const hash_key_t key)
    __attribute__((hot, warn_unused_result));
hashtable_rc_t hashtable_apply_callback_on_elements(
    hash_table_t* const hashtbl,
    bool func_cb(hash_key_t key, void* element, void* parameter, void** result),
    void* parameter, void** result);
hashtable_rc_t hashtable_dump_content(
    const hash_table_t* const hashtbl, bstring str);
hashtable_rc_t hashtable_insert(
    hash_table_t* const hashtbl, const hash_key_t key, void* element);
hashtable_rc_t hashtable_free(
    hash_table_t* const hashtbl, const hash_key_t key);
hashtable_rc_t hashtable_remove(
    hash_table_t* const hashtbl, const hash_key_t key, void** element);
hashtable_rc_t hashtable_get(
    const hash_table_t* const hashtbl, const hash_key_t key, void** element)
    __attribute__((hot));
hashtable_rc_t hashtable_resize(
    hash_table_t* const hashtbl, const hash_size_t size);

// Thread-safe functions
hash_table_ts_t* hashtable_ts_init(
    hash_table_ts_t* const hashtbl, const hash_size_t size,
    hash_size_t (*hashfunc)(const hash_key_t), void (*freefunc)(void**),
    bstring display_name_p);
__attribute__((malloc)) hash_table_ts_t* hashtable_ts_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const hash_key_t),
    void (*freefunc)(void**), bstring name_p);
hashtable_rc_t hashtable_ts_destroy(hash_table_ts_t* hashtbl);
hashtable_rc_t hashtable_ts_is_key_exists(
    const hash_table_ts_t* const hashtbl, const hash_key_t key)
    __attribute__((hot, warn_unused_result));
hashtable_key_array_t* hashtable_ts_get_keys(hash_table_ts_t* const hashtblP);
hashtable_element_array_t* hashtable_ts_get_elements(
    hash_table_ts_t* const hashtblP);
hashtable_rc_t hashtable_ts_apply_callback_on_elements(
    hash_table_ts_t* const hashtbl,
    bool func_cb(
        const hash_key_t key, void* const element, void* parameter,
        void** result),
    void* parameter, void** result);
hashtable_rc_t hashtable_ts_dump_content(
    const hash_table_ts_t* const hashtbl, bstring str);
hashtable_rc_t hashtable_ts_insert(
    hash_table_ts_t* const hashtbl, const hash_key_t key, void* element);
hashtable_rc_t hashtable_ts_free(
    hash_table_ts_t* const hashtbl, const hash_key_t key);
hashtable_rc_t hashtable_ts_remove(
    hash_table_ts_t* const hashtbl, const hash_key_t key, void** element);
hashtable_rc_t hashtable_ts_get(
    const hash_table_ts_t* const hashtbl, const hash_key_t key, void** element)
    __attribute__((hot));
hashtable_rc_t hashtable_ts_resize(
    hash_table_ts_t* const hashtbl, const hash_size_t size);
hash_table_uint64_ts_t* hashtable_uint64_ts_init(
    hash_table_uint64_ts_t* const hashtbl, const hash_size_t size,
    hash_size_t (*hashfunc)(const hash_key_t), bstring display_name_p);
__attribute__((malloc)) hash_table_uint64_ts_t* hashtable_uint64_ts_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const hash_key_t),
    bstring name_p);
hashtable_rc_t hashtable_uint64_ts_destroy(hash_table_uint64_ts_t* hashtbl);
hashtable_rc_t hashtable_uint64_ts_is_key_exists(
    const hash_table_uint64_ts_t* const hashtbl, const hash_key_t key)
    __attribute__((hot, warn_unused_result));
hashtable_key_array_t* hashtable_uint64_ts_get_keys(
    hash_table_uint64_ts_t* const hashtblP);
hashtable_uint64_element_array_t* hashtable_uint64_ts_get_elements(
    hash_table_uint64_ts_t* const hashtblP);
hashtable_rc_t hashtable_uint64_ts_apply_callback_on_elements(
    hash_table_uint64_ts_t* const hashtbl,
    bool func_cb(
        const hash_key_t key, const uint64_t element, void* parameter,
        void** result),
    void* parameter, void** result);
hashtable_rc_t hashtable_uint64_ts_dump_content(
    const hash_table_uint64_ts_t* const hashtbl, bstring str);
hashtable_rc_t hashtable_uint64_ts_insert(
    hash_table_uint64_ts_t* const hashtbl, const hash_key_t key,
    const uint64_t dataP);
hashtable_rc_t hashtable_uint64_ts_remove(
    hash_table_uint64_ts_t* const hashtbl, const hash_key_t key);
hashtable_rc_t hashtable_uint64_ts_get(
    const hash_table_uint64_ts_t* const hashtbl, const hash_key_t key,
    uint64_t* const dataP) __attribute__((hot));
hashtable_rc_t hashtable_uint64_ts_resize(
    hash_table_uint64_ts_t* const hashtbl, const hash_size_t size);

#endif
