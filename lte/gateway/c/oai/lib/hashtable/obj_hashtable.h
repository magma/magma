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

/*! \file obj_hashtable.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_OBJ_HASHTABLE_SEEN
#define FILE_OBJ_HASHTABLE_SEEN

#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>

#include "hashtable.h"
#include "bstrlib.h"

#define FREE_OBJ_HASHTABLE_KEY_ARRAY(key_array_ptr)                            \
  do {                                                                         \
    AssertFatal(key_array_ptr, "Trying to free a NULL array pointer");         \
    free(*key_array_ptr);                                                      \
    free(key_array_ptr);                                                       \
  } while (0) /*Free the list of keys of an object hash table */

typedef struct obj_hash_node_s {
  int key_size;
  void* key;
  void* data;
  struct obj_hash_node_s* next;
} obj_hash_node_t;

typedef struct obj_hash_node_uint64_s {
  int key_size;
  void* key;
  uint64_t data;
  struct obj_hash_node_uint64_s* next;
} obj_hash_node_uint64_t;

typedef struct obj_hash_table_s {
  pthread_mutex_t mutex;
  hash_size_t size;
  hash_size_t num_elements;
  struct obj_hash_node_s** nodes;
  pthread_mutex_t* lock_nodes;
  hash_size_t (*hashfunc)(const void*, int);
  void (*freekeyfunc)(void**);
  void (*freedatafunc)(void**);
  bstring name;
  bool log_enabled;
} obj_hash_table_t;
typedef struct obj_hash_table_uint64_s {
  pthread_mutex_t mutex;
  hash_size_t size;
  hash_size_t num_elements;
  struct obj_hash_node_uint64_s** nodes;
  pthread_mutex_t* lock_nodes;
  hash_size_t (*hashfunc)(const void*, int);
  void (*freekeyfunc)(void**);
  bstring name;
  bool log_enabled;
} obj_hash_table_uint64_t;

void obj_hashtable_no_free_key_callback(void* param);
obj_hash_table_t* obj_hashtable_init(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const void*, int), void (*freekeyfuncP)(void**),
    void (*freedatafuncP)(void**), bstring display_name_pP);
obj_hash_table_t* obj_hashtable_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const void*, int),
    void (*freekeyfunc)(void**), void (*freedatafunc)(void**),
    bstring display_name_pP);
hashtable_rc_t obj_hashtable_destroy(obj_hash_table_t* const hashtblP);
hashtable_rc_t obj_hashtable_is_key_exists(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP) __attribute__((hot, warn_unused_result));
hashtable_rc_t obj_hashtable_insert(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void* dataP);
hashtable_rc_t obj_hashtable_dump_content(
    const obj_hash_table_t* const hashtblP, bstring str);
hashtable_rc_t obj_hashtable_free(
    obj_hash_table_t* hashtblP, const void* keyP, const int key_sizeP);
hashtable_rc_t obj_hashtable_remove(
    obj_hash_table_t* hashtblP, const void* keyP, const int key_sizeP,
    void** dataP);
hashtable_rc_t obj_hashtable_get(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void** dataP) __attribute__((hot));
hashtable_rc_t obj_hashtable_get_keys(
    const obj_hash_table_t* const hashtblP, void** keysP, unsigned int* sizeP);
hashtable_rc_t obj_hashtable_resize(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP);

// Thread-safe functions
obj_hash_table_t* obj_hashtable_ts_init(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const void*, int), void (*freekeyfuncP)(void**),
    void (*freedatafuncP)(void**), bstring display_name_pP);
obj_hash_table_t* obj_hashtable_ts_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const void*, int),
    void (*freekeyfunc)(void**), void (*freedatafunc)(void**),
    bstring display_name_pP);
hashtable_rc_t obj_hashtable_ts_destroy(obj_hash_table_t* const hashtblP);
hashtable_rc_t obj_hashtable_ts_is_key_exists(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP) __attribute__((hot, warn_unused_result));
hashtable_rc_t obj_hashtable_ts_insert(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void* dataP);
hashtable_rc_t obj_hashtable_ts_dump_content(
    const obj_hash_table_t* const hashtblP, bstring str);
hashtable_rc_t obj_hashtable_ts_free(
    obj_hash_table_t* hashtblP, const void* keyP, const int key_sizeP);
hashtable_rc_t obj_hashtable_ts_remove(
    obj_hash_table_t* hashtblP, const void* keyP, const int key_sizeP,
    void** dataP);
hashtable_rc_t obj_hashtable_ts_get(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void** dataP) __attribute__((hot));
hashtable_rc_t obj_hashtable_ts_get_keys(
    const obj_hash_table_t* const hashtblP, void** keysP, unsigned int* sizeP);
hashtable_rc_t obj_hashtable_ts_resize(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP);
obj_hash_table_uint64_t* obj_hashtable_uint64_init(
    obj_hash_table_uint64_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const void*, int), void (*freekeyfuncP)(void**),
    bstring display_name_pP);
obj_hash_table_uint64_t* obj_hashtable_uint64_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const void*, int),
    void (*freekeyfunc)(void**), bstring display_name_pP);
hashtable_rc_t obj_hashtable_uint64_destroy(
    obj_hash_table_uint64_t* const hashtblP);
hashtable_rc_t obj_hashtable_uint64_is_key_exists(
    const obj_hash_table_uint64_t* const hashtblP, const void* const keyP,
    const int key_sizeP) __attribute__((hot, warn_unused_result));
hashtable_rc_t obj_hashtable_uint64_insert(
    obj_hash_table_uint64_t* const hashtblP, const void* const keyP,
    const int key_sizeP, const uint64_t dataP);
hashtable_rc_t obj_hashtable_uint64_dump_content(
    const obj_hash_table_uint64_t* const hashtblP, bstring str);
hashtable_rc_t obj_hashtable_uint64_free(
    obj_hash_table_uint64_t* hashtblP, const void* keyP, const int key_sizeP);
hashtable_rc_t obj_hashtable_uint64_remove(
    obj_hash_table_uint64_t* hashtblP, const void* keyP, const int key_sizeP);
hashtable_rc_t obj_hashtable_uint64_get(
    const obj_hash_table_uint64_t* const hashtblP, const void* const keyP,
    const int key_sizeP, uint64_t* const dataP) __attribute__((hot));
hashtable_rc_t obj_hashtable_uint64_get_keys(
    const obj_hash_table_uint64_t* const hashtblP, void** keysP,
    unsigned int* sizeP);
hashtable_rc_t obj_hashtable_uint64_resize(
    obj_hash_table_uint64_t* const hashtblP, const hash_size_t sizeP);

// Thread-safe functions
obj_hash_table_uint64_t* obj_hashtable_uint64_ts_init(
    obj_hash_table_uint64_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const void*, int), void (*freekeyfuncP)(void**),
    bstring display_name_pP);
obj_hash_table_uint64_t* obj_hashtable_uint64_ts_create(
    const hash_size_t size, hash_size_t (*hashfunc)(const void*, int),
    void (*freekeyfunc)(void**), bstring display_name_pP);
hashtable_rc_t obj_hashtable_uint64_ts_destroy(
    obj_hash_table_uint64_t* const hashtblP);
hashtable_rc_t obj_hashtable_uint64_ts_is_key_exists(
    const obj_hash_table_uint64_t* const hashtblP, const void* const keyP,
    const int key_sizeP) __attribute__((hot, warn_unused_result));
hashtable_rc_t obj_hashtable_uint64_ts_insert(
    obj_hash_table_uint64_t* const hashtblP, const void* const keyP,
    const int key_sizeP, const uint64_t dataP);
hashtable_rc_t obj_hashtable_uint64_ts_dump_content(
    const obj_hash_table_uint64_t* const hashtblP, bstring str);
hashtable_rc_t obj_hashtable_uint64_ts_free(
    obj_hash_table_uint64_t* hashtblP, const void* keyP, const int key_sizeP);
hashtable_rc_t obj_hashtable_uint64_ts_remove(
    obj_hash_table_uint64_t* hashtblP, const void* keyP, const int key_sizeP);
hashtable_rc_t obj_hashtable_uint64_ts_get(
    const obj_hash_table_uint64_t* const hashtblP, const void* const keyP,
    const int key_sizeP, uint64_t* const dataP) __attribute__((hot));
hashtable_rc_t obj_hashtable_uint64_ts_get_keys(
    const obj_hash_table_uint64_t* const hashtblP, void*** keysP,
    unsigned int* sizeP);
hashtable_rc_t obj_hashtable_uint64_ts_resize(
    obj_hash_table_uint64_t* const hashtblP, const hash_size_t sizeP);

#endif
