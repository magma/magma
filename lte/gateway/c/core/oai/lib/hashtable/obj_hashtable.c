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

/*! \file obj_hashtable.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <pthread.h>

#include "bstrlib.h"
#include "obj_hashtable.h"
#include "dynamic_memory_check.h"

#if TRACE_HASHTABLE
#define PRINT_HASHTABLE(hTbLe, ...)                                            \
  do {                                                                         \
    if (hTbLe->log_enabled) OAILOG_TRACE(LOG_UTIL, ##__VA_ARGS__);             \
  } while (0)
#else
#define PRINT_HASHTABLE(...)
#endif

//------------------------------------------------------------------------------
// Free function selected if we do not want to free_wrapper the key when
// removing an entry
void obj_hashtable_no_free_key_callback(void* param) {
  ;  // volountary do nothing
}

//------------------------------------------------------------------------------
/*
   Default hash function
   def_hashfunc() is the default used by hashtable_create() when the user didn't
   specify one. This is a simple/naive hash function which adds the key's ASCII
   char values. It will probably generate lots of collisions on large hash
   tables.
*/

static hash_size_t def_hashfunc(const void* const keyP, const int key_sizeP) {
  hash_size_t hash = 0;
  int key_size     = key_sizeP;

  // may use MD4 ?
  while (key_size > 0) {
    uint32_t val = 0;
    int size     = sizeof(val);
    while ((size > 0) && (key_size > 0)) {
      val = val << 8;
      val |= ((uint8_t*) keyP)[key_size - 1];
      size--;
      key_size--;
    }
    hash ^= val;
  }

  return hash;
}

//------------------------------------------------------------------------------
/*
 *    Initialization
 *    obj_hashtable_init() sets up the initial structure of the hash table. The
 * user specified size will be allocated and initialized to NULL. The user can
 * also specify a hash function. If the hashfunc argument is NULL, a default
 * hash function is used. If an error occurred, NULL is returned. All other
 * values in the returned obj_hash_table_t pointer should be released with
 * hashtable_destroy().
 *
 */
obj_hash_table_t* obj_hashtable_init(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const void*, int), void (*freekeyfuncP)(void**),
    void (*freedatafuncP)(void**), bstring display_name_pP) {
  hash_size_t size = sizeP;
  // upper power of two:
  // http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2Float
  //  By Sean Eron Anderson
  // seander@cs.stanford.edu
  // Individually, the code snippets here are in the public domain (unless
  // otherwise noted) — feel free to use them however you please. The aggregate
  // collection and descriptions are © 1997-2005 Sean Eron Anderson. The code
  // and descriptions are distributed in the hope that they will be useful, but
  // WITHOUT ANY WARRANTY and without even the implied warranty of
  // merchantability or fitness for a particular purpose. As of May 5, 2005, all
  // the code has been tested thoroughly. Thousands of people have read it.
  // Moreover, Professor Randal Bryant, the Dean of Computer Science at Carnegie
  // Mellon University, has personally tested almost everything with his Uclid
  // code verification system. What he hasn't tested, I have checked against all
  // possible inputs on a 32-bit machine. To the first person to inform me of a
  // legitimate bug in the code, I'll pay a bounty of US$10 (by check or
  // Paypal). If directed to a charity, I'll pay US$20.
  size--;
  size |= size >> 1;
  size |= size >> 2;
  size |= size >> 4;
  size |= size >> 8;
  size |= size >> 16;
  size++;

  if (!(hashtblP->nodes = calloc(size, sizeof(obj_hash_node_t*)))) {
    free_wrapper((void**) &hashtblP);
    return NULL;
  }

  hashtblP->size = size;

  if (hashfuncP)
    hashtblP->hashfunc = hashfuncP;
  else
    hashtblP->hashfunc = def_hashfunc;

  if (freekeyfuncP)
    hashtblP->freekeyfunc = freekeyfuncP;
  else
    hashtblP->freekeyfunc = free_wrapper;

  if (freedatafuncP)
    hashtblP->freedatafunc = freedatafuncP;
  else
    hashtblP->freedatafunc = free_wrapper;

  if (display_name_pP) {
    hashtblP->name = bstrcpy(display_name_pP);
  } else {
    hashtblP->name = bfromcstr("");
    btrunc(hashtblP->name, 0);
    bassignformat(hashtblP->name, "obj_hashtable%u@%p", size, hashtblP);
  }
  hashtblP->log_enabled = true;
  return hashtblP;
}

//------------------------------------------------------------------------------
/*
   Initialization
   obj_hashtable_create() allocate and set up the initial structure of the hash
   table. The user specified size will be allocated and initialized to NULL. The
   user can also specify a hash function. If the hashfunc argument is NULL, a
   default hash function is used. If an error occurred, NULL is returned. All
   other values in the returned obj_hash_table_t pointer should be released with
   hashtable_destroy().
*/
obj_hash_table_t* obj_hashtable_create(
    const hash_size_t sizeP, hash_size_t (*hashfuncP)(const void*, int),
    void (*freekeyfuncP)(void**), void (*freedatafuncP)(void**),
    bstring display_name_pP) {
  obj_hash_table_t* hashtbl = NULL;

  if (!(hashtbl = calloc(1, sizeof(obj_hash_table_t)))) return NULL;

  return obj_hashtable_init(
      hashtbl, sizeP, hashfuncP, freekeyfuncP, freedatafuncP, display_name_pP);
}

//------------------------------------------------------------------------------
/*
   Initialization
   obj_hashtable_ts_init() sets up the initial structure of the hash table. The
   user specified size will be allocated and initialized to NULL. The user can
   also specify a hash function. If the hashfunc argument is NULL, a default
   hash function is used. If an error occurred, NULL is returned. All other
   values in the returned obj_hash_table_t pointer should be released with
   hashtable_destroy().
*/
obj_hash_table_t* obj_hashtable_ts_init(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const void*, int), void (*freekeyfuncP)(void**),
    void (*freedatafuncP)(void**), bstring display_name_pP) {
  hash_size_t size = sizeP;
  // upper power of two:
  // http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2Float
  //  By Sean Eron Anderson
  // seander@cs.stanford.edu
  // Individually, the code snippets here are in the public domain (unless
  // otherwise noted) — feel free to use them however you please. The aggregate
  // collection and descriptions are © 1997-2005 Sean Eron Anderson. The code
  // and descriptions are distributed in the hope that they will be useful, but
  // WITHOUT ANY WARRANTY and without even the implied warranty of
  // merchantability or fitness for a particular purpose. As of May 5, 2005, all
  // the code has been tested thoroughly. Thousands of people have read it.
  // Moreover, Professor Randal Bryant, the Dean of Computer Science at Carnegie
  // Mellon University, has personally tested almost everything with his Uclid
  // code verification system. What he hasn't tested, I have checked against all
  // possible inputs on a 32-bit machine. To the first person to inform me of a
  // legitimate bug in the code, I'll pay a bounty of US$10 (by check or
  // Paypal). If directed to a charity, I'll pay US$20.
  size--;
  size |= size >> 1;
  size |= size >> 2;
  size |= size >> 4;
  size |= size >> 8;
  size |= size >> 16;
  size++;

  if (!(hashtblP->lock_nodes = calloc(size, sizeof(pthread_mutex_t)))) {
    free_wrapper((void**) &hashtblP->nodes);
    free_wrapper((void**) &hashtblP->name);
    free_wrapper((void**) &hashtblP);
    return NULL;
  }

  pthread_mutex_init(&hashtblP->mutex, NULL);
  for (int i = 0; i < size; i++) {
    pthread_mutex_init(&hashtblP->lock_nodes[i], NULL);
  }

  hashtblP->log_enabled = true;
  return hashtblP;
}

//------------------------------------------------------------------------------
/*
   Initialisation
   obj_hashtable_ts_create() allocate and sets up the initial structure of the
   hash table. The user specified size will be allocated and initialized to
   NULL. The user can also specify a hash function. If the hashfunc argument is
   NULL, a default hash function is used. If an error occurred, NULL is
   returned. All other values in the returned obj_hash_table_t pointer should be
   released with hashtable_destroy().
*/
obj_hash_table_t* obj_hashtable_ts_create(
    const hash_size_t sizeP, hash_size_t (*hashfuncP)(const void*, int),
    void (*freekeyfuncP)(void**), void (*freedatafuncP)(void**),
    bstring display_name_pP) {
  obj_hash_table_t* hashtbl = NULL;

  hash_size_t size = sizeP;
  // upper power of two:
  // http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2Float
  //  By Sean Eron Anderson
  // seander@cs.stanford.edu
  // Individually, the code snippets here are in the public domain (unless
  // otherwise noted) — feel free to use them however you please. The aggregate
  // collection and descriptions are © 1997-2005 Sean Eron Anderson. The code
  // and descriptions are distributed in the hope that they will be useful, but
  // WITHOUT ANY WARRANTY and without even the implied warranty of
  // merchantability or fitness for a particular purpose. As of May 5, 2005, all
  // the code has been tested thoroughly. Thousands of people have read it.
  // Moreover, Professor Randal Bryant, the Dean of Computer Science at Carnegie
  // Mellon University, has personally tested almost everything with his Uclid
  // code verification system. What he hasn't tested, I have checked against all
  // possible inputs on a 32-bit machine. To the first person to inform me of a
  // legitimate bug in the code, I'll pay a bounty of US$10 (by check or
  // Paypal). If directed to a charity, I'll pay US$20.
  size--;
  size |= size >> 1;
  size |= size >> 2;
  size |= size >> 4;
  size |= size >> 8;
  size |= size >> 16;
  size++;

  if (!(hashtbl = obj_hashtable_create(
            size, hashfuncP, freekeyfuncP, freedatafuncP, display_name_pP))) {
    return NULL;
  }

  return obj_hashtable_ts_init(
      hashtbl, size, hashfuncP, freekeyfuncP, freedatafuncP, display_name_pP);
}

//------------------------------------------------------------------------------
/*
   Cleanup
   The hashtable_destroy() walks through the linked lists for each possible hash
   value, and releases the elements. It also releases the nodes array and the
   obj_hash_table_t.
*/
hashtable_rc_t obj_hashtable_destroy(obj_hash_table_t* const hashtblP) {
  hash_size_t n;
  obj_hash_node_t *node, *oldnode;

  for (n = 0; n < hashtblP->size; ++n) {
    node = hashtblP->nodes[n];

    while (node) {
      oldnode = node;
      node    = node->next;
      hashtblP->freekeyfunc(&oldnode->key);
      hashtblP->freedatafunc(&oldnode->data);
      free_wrapper((void**) &oldnode);
    }
  }

  free_wrapper((void**) &hashtblP->nodes);
  free_wrapper((void**) &hashtblP->lock_nodes);  // mmm....
  bdestroy_wrapper(&hashtblP->name);
  free_wrapper((void**) &hashtblP);
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   Cleanup
   The hashtable_destroy() walks through the linked lists for each possible hash
   value, and releases the elements. It also releases the nodes array and the
   obj_hash_table_t.
*/
hashtable_rc_t obj_hashtable_ts_destroy(obj_hash_table_t* const hashtblP) {
  hash_size_t n;
  obj_hash_node_t *node, *oldnode;

  for (n = 0; n < hashtblP->size; ++n) {
    pthread_mutex_lock(&hashtblP->lock_nodes[n]);
    node = hashtblP->nodes[n];

    while (node) {
      oldnode = node;
      node    = node->next;
      hashtblP->freekeyfunc(&oldnode->key);
      hashtblP->freedatafunc(&oldnode->data);
      free_wrapper((void**) &oldnode);
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[n]);
    pthread_mutex_destroy(&hashtblP->lock_nodes[n]);
  }

  free_wrapper((void**) &hashtblP->nodes);
  free_wrapper((void**) &hashtblP->lock_nodes);
  bdestroy_wrapper(&hashtblP->name);
  free_wrapper((void**) &hashtblP);
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
hashtable_rc_t obj_hashtable_is_key_exists(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP) {
  obj_hash_node_t* node;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p klen %u) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, key_sizeP, hash);
      return HASH_TABLE_OK;
    } else if (node->key_size == key_sizeP) {
      if (memcmp(node->key, keyP, key_sizeP) == 0) {
        PRINT_HASHTABLE(
            hashtblP, "%s(%s,key %p klen %u) hash %lx return OK\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, hash);
        return HASH_TABLE_OK;
      }
    }

    node = node->next;
  }

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key %p klen %u) hash %lx return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, hash);
  return HASH_TABLE_KEY_NOT_EXISTS;
}
//------------------------------------------------------------------------------
hashtable_rc_t obj_hashtable_ts_is_key_exists(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP) {
  obj_hash_node_t* node;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p klen %u) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, key_sizeP, hash);
      return HASH_TABLE_OK;
    } else if (node->key_size == key_sizeP) {
      if (memcmp(node->key, keyP, key_sizeP) == 0) {
        pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
        PRINT_HASHTABLE(
            hashtblP, "%s(%s,key %p klen %u) hash %lx return OK\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, hash);
        return HASH_TABLE_OK;
      }
    }

    node = node->next;
  }
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key %p klen %u) hash %lx return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, hash);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
hashtable_rc_t obj_hashtable_dump_content(
    const obj_hash_table_t* const hashtblP, bstring str) {
  obj_hash_node_t* node = NULL;
  unsigned int i        = 0;

  if (hashtblP == NULL) {
    bcatcstr(str, "HASH_TABLE_BAD_PARAMETER_HASHTABLE");
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  while (i < hashtblP->size) {
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];

      while (node) {
        bstring b0 = bformat(
            "Hash %x Key %p Key length %d Element %p\n", i, node->key,
            node->key_size, node->data);
        if (!b0) {
          PRINT_HASHTABLE(hashtblP, "Error while dumping hashtable content");
        } else {
          bconcat(str, b0);
          bdestroy_wrapper(&b0);
        }
        node = node->next;
      }
    }
    i += 1;
  }

  return HASH_TABLE_OK;
}
//------------------------------------------------------------------------------
hashtable_rc_t obj_hashtable_ts_dump_content(
    const obj_hash_table_t* const hashtblP, bstring str) {
  obj_hash_node_t* node = NULL;
  unsigned int i        = 0;

  if (hashtblP == NULL) {
    bcatcstr(str, "HASH_TABLE_BAD_PARAMETER_HASHTABLE");
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  while (i < hashtblP->size) {
    if (hashtblP->nodes[i] != NULL) {
      pthread_mutex_lock(&hashtblP->lock_nodes[i]);
      node = hashtblP->nodes[i];

      while (node) {
        bstring b0 = bformat(
            "Hash %x Key %p Key length %d Element %p\n", i, node->key,
            node->key_size, node->data);
        if (!b0) {
          PRINT_HASHTABLE(hashtblP, "Error while dumping hashtable content");
        } else {
          bconcat(str, b0);
          bdestroy_wrapper(&b0);
        }
        node = node->next;
      }
      pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    }
    i += 1;
  }
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   Adding a new element
   To make sure the hash value is not bigger than size, the result of the user
   provided hash function is used modulo size.
*/
hashtable_rc_t obj_hashtable_insert(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void* dataP) {
  obj_hash_node_t* node;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if (node->data) {
        hashtblP->freedatafunc(&node->data);
      }

      node->data     = dataP;
      node->key_size = key_sizeP;
      // waste of memory here (keyP is lost) we should free_wrapper it now
      PRINT_HASHTABLE(
          hashtblP,
          "%s(%s,key %p data %p) hash %lx return INSERT_OVERWRITTEN_DATA\n",
          __FUNCTION__, bdata(hashtblP->name), keyP, dataP, hash);
      return HASH_TABLE_INSERT_OVERWRITTEN_DATA;
    }

    node = node->next;
  }

  if (!(node = malloc(sizeof(obj_hash_node_t)))) {
    PRINT_HASHTABLE(
        hashtblP, "%s(%s,key %p) hash %lx return SYSTEM_ERROR\n", __FUNCTION__,
        bdata(hashtblP->name), keyP, hash);
    return HASH_TABLE_SYSTEM_ERROR;
  }

  if (!(node->key = malloc(key_sizeP))) {
    PRINT_HASHTABLE(
        hashtblP, "%s(%s,key %p) hash %lx return SYSTEM_ERROR\n", __FUNCTION__,
        bdata(hashtblP->name), keyP, hash);
    free_wrapper((void**) &node);
    return -1;
  }

  memcpy(node->key, keyP, key_sizeP);
  node->data     = dataP;
  node->key_size = key_sizeP;

  if (hashtblP->nodes[hash]) {
    node->next = hashtblP->nodes[hash];
  } else {
    node->next = NULL;
  }

  hashtblP->nodes[hash] = node;
  __sync_fetch_and_add(&hashtblP->num_elements, 1);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key %p klen %u data %p) hash %lx return OK\n",
      __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, dataP, hash);
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   Adding a new element
   To make sure the hash value is not bigger than size, the result of the user
   provided hash function is used modulo size.
*/
hashtable_rc_t obj_hashtable_ts_insert(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void* dataP) {
  obj_hash_node_t* node;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if ((node->data) && (node->data != dataP)) {
        hashtblP->freedatafunc(&node->data);

        node->data     = dataP;
        node->key_size = key_sizeP;
        // no waste of memory here because if node->key == keyP, it is a reuse
        // of the same key
        pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
        PRINT_HASHTABLE(
            hashtblP,
            "%s(%s,key %p data %p) hash %lx return INSERT_OVERWRITTEN_DATA\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, dataP, hash);
        return HASH_TABLE_INSERT_OVERWRITTEN_DATA;
      }
      node->data = dataP;
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p data %p) hash %lx return ok\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, dataP, hash);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }

  if (!(node = calloc(1, sizeof(obj_hash_node_t)))) {
    pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
    PRINT_HASHTABLE(
        hashtblP, "%s(%s,key %p) hash %lx return SYSTEM_ERROR\n", __FUNCTION__,
        bdata(hashtblP->name), keyP, hash);
    return HASH_TABLE_SYSTEM_ERROR;
  }

  if (!(node->key = calloc(1, key_sizeP))) {
    free_wrapper((void**) &node);
    pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
    PRINT_HASHTABLE(
        hashtblP, "%s(%s,key %p) hash %lx return SYSTEM_ERROR\n", __FUNCTION__,
        bdata(hashtblP->name), keyP, hash);
    return HASH_TABLE_SYSTEM_ERROR;
  }

  memcpy(node->key, keyP, key_sizeP);
  node->data     = dataP;
  node->key_size = key_sizeP;

  if (hashtblP->nodes[hash]) {
    node->next = hashtblP->nodes[hash];
  } else {
    node->next = NULL;
  }

  hashtblP->nodes[hash] = node;
  __sync_fetch_and_add(&hashtblP->num_elements, 1);
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key %p klen %u data %p) hash %lx return OK\n",
      __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, dataP, hash);
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   To remove an element from the hash table, we just search for it in the linked
   list for that hash value, and remove it if it is found. If it was not found,
   it is an error and -1 is returned.
*/
hashtable_rc_t obj_hashtable_free(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP) {
  obj_hash_node_t *node, *prevnode = NULL;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if ((node->key == keyP) || ((node->key_size == key_sizeP) &&
                                (memcmp(node->key, keyP, key_sizeP) == 0))) {
      if (prevnode) {
        prevnode->next = node->next;
      } else {
        hashtblP->nodes[hash] = node->next;
      }

      hashtblP->freekeyfunc(&node->key);
      hashtblP->freedatafunc(&node->data);
      free_wrapper((void**) &node);
      hashtblP->num_elements -= 1;
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, hash);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }

  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   To remove an element from the hash table, we just search for it in the linked
   list for that hash value, and remove it if it is found. If it was not found,
   it is an error and -1 is returned.
*/
hashtable_rc_t obj_hashtable_ts_free(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP) {
  obj_hash_node_t *node, *prevnode = NULL;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if ((node->key == keyP) || ((node->key_size == key_sizeP) &&
                                (memcmp(node->key, keyP, key_sizeP) == 0))) {
      if (prevnode) {
        prevnode->next = node->next;
      } else {
        hashtblP->nodes[hash] = node->next;
      }

      hashtblP->freekeyfunc(&node->key);
      hashtblP->freedatafunc(&node->data);
      free_wrapper((void**) &node);
      __sync_fetch_and_sub(&hashtblP->num_elements, 1);
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, hash);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);

  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   To remove an element from the hash table, we just search for it in the linked
   list for that hash value, and remove it if it is found. If it was not found,
   it is an error and -1 is returned.
*/
hashtable_rc_t obj_hashtable_remove(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void** dataP) {
  obj_hash_node_t *node, *prevnode = NULL;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if ((node->key == keyP) || ((node->key_size == key_sizeP) &&
                                (memcmp(node->key, keyP, key_sizeP) == 0))) {
      if (prevnode) {
        prevnode->next = node->next;
      } else {
        hashtblP->nodes[hash] = node->next;
      }

      hashtblP->freekeyfunc(&node->key);
      *dataP = node->data;
      free_wrapper((void**) &node);
      hashtblP->num_elements -= 1;
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, hash);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }

  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   To remove an element from the hash table, we just search for it in the linked
   list for that hash value, and remove it if it is found. If it was not found,
   it is an error and -1 is returned.
*/
hashtable_rc_t obj_hashtable_ts_remove(
    obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void** dataP) {
  obj_hash_node_t *node, *prevnode = NULL;
  hash_size_t hash;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if ((node->key == keyP) || ((node->key_size == key_sizeP) &&
                                (memcmp(node->key, keyP, key_sizeP) == 0))) {
      if (prevnode) {
        prevnode->next = node->next;
      } else {
        hashtblP->nodes[hash] = node->next;
      }

      hashtblP->freekeyfunc(&node->key);
      *dataP = node->data;
      free_wrapper((void**) &node);
      __sync_fetch_and_sub(&hashtblP->num_elements, 1);
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, hash);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);

  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   Searching for an element is easy. We just search through the linked list for
   the corresponding hash value. NULL is returned if we didn't find it.
*/
hashtable_rc_t obj_hashtable_get(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void** dataP) {
  obj_hash_node_t* node;
  hash_size_t hash;

  if (hashtblP == NULL) {
    *dataP = NULL;
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      *dataP = node->data;
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p data %p) hash %lx return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP, *dataP, hash);
      return HASH_TABLE_OK;
    } else if (node->key_size == key_sizeP) {
      if (memcmp(node->key, keyP, key_sizeP) == 0) {
        *dataP = node->data;
        PRINT_HASHTABLE(
            hashtblP, "%s(%s,key %p data %p) hash %lx return OK\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, *dataP, hash);
        return HASH_TABLE_OK;
      }
    }

    node = node->next;
  }

  *dataP = NULL;
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key %p) hash %lx return KEY_NOT_EXISTS\n", __FUNCTION__,
      bdata(hashtblP->name), keyP, hash);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   Searching for an element is easy. We just search through the linked list for
   the corresponding hash value. NULL is returned if we didn't find it.
*/
hashtable_rc_t obj_hashtable_ts_get(
    const obj_hash_table_t* const hashtblP, const void* const keyP,
    const int key_sizeP, void** dataP) {
  obj_hash_node_t* node;
  hash_size_t hash;

  if (hashtblP == NULL) {
    *dataP = NULL;
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  if (keyP == NULL) {
    PRINT_HASHTABLE(hashtblP, "return HASH_TABLE_BAD_PARAMETER_KEY\n");
    return HASH_TABLE_BAD_PARAMETER_KEY;
  }

  hash = hashtblP->hashfunc(keyP, key_sizeP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      *dataP = node->data;
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key %p klen %u data %p) hash %lx return OK\n",
          __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, *dataP, hash);
      return HASH_TABLE_OK;
    } else if (node->key_size == key_sizeP) {
      if (memcmp(node->key, keyP, key_sizeP) == 0) {
        *dataP = node->data;
        pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
        PRINT_HASHTABLE(
            hashtblP, "%s(%s,key %p klen %u data %p) hash %lx return OK\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, *dataP, hash);
        return HASH_TABLE_OK;
      }
    }

    node = node->next;
  }

  *dataP = NULL;
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key %p klen %u) hash %lx return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP, key_sizeP, hash);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   Function to return all keys of an object hash table
*/
hashtable_rc_t obj_hashtable_get_keys(
    const obj_hash_table_t* const hashtblP, void** keysP, unsigned int* sizeP) {
  size_t n              = 0;
  obj_hash_node_t* node = NULL;
  obj_hash_node_t* next = NULL;

  if (hashtblP == NULL) {
    keysP = NULL;
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  *sizeP = 0;
  keysP  = calloc(hashtblP->num_elements, sizeof(void*));

  if (keysP) {
    for (n = 0; n < hashtblP->size; ++n) {
      for (node = hashtblP->nodes[n]; node; node = next) {
        keysP[*sizeP++] = node->key;
        next            = node->next;
      }
    }

    PRINT_HASHTABLE(hashtblP, "return OK\n");
    return HASH_TABLE_OK;
  }
  PRINT_HASHTABLE(hashtblP, "return SYSTEM_ERROR\n");
  return HASH_TABLE_SYSTEM_ERROR;
}

//------------------------------------------------------------------------------
/*
   Function to return all keys of an object hash table
*/
hashtable_rc_t obj_hashtable_ts_get_keys(
    const obj_hash_table_t* const hashtblP, void** keysP, unsigned int* sizeP) {
  size_t n              = 0;
  obj_hash_node_t* node = NULL;
  obj_hash_node_t* next = NULL;

  if (hashtblP == NULL) {
    keysP = NULL;
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  *sizeP = 0;
  keysP  = calloc(hashtblP->num_elements, sizeof(void*));

  if (keysP) {
    for (n = 0; n < hashtblP->size; ++n) {
      pthread_mutex_lock(&hashtblP->lock_nodes[n]);
      for (node = hashtblP->nodes[n]; node; node = next) {
        keysP[*sizeP++] = node->key;
        next            = node->next;
      }
      pthread_mutex_unlock(&hashtblP->lock_nodes[n]);
    }

    PRINT_HASHTABLE(hashtblP, "return OK\n");
    return HASH_TABLE_OK;
  }
  PRINT_HASHTABLE(hashtblP, "return SYSTEM_ERROR\n");
  return HASH_TABLE_SYSTEM_ERROR;
}

//------------------------------------------------------------------------------
/*
   Resizing
   The number of elements in a hash table is not always known when creating the
   table. If the number of elements grows too large, it will seriously reduce
   the performance of most hash table operations. If the number of elements are
   reduced, the hash table will waste memory. That is why we provide a function
   for resizing the table. Resizing a hash table is not as easy as a realloc().
   All hash values must be recalculated and each element must be inserted into
   its new position. We create a temporary obj_hash_table_t object (newtbl) to
   be used while building the new hashes. This allows us to reuse
   hashtable_insert() and hashtable_free(), when moving the elements to the new
   table. After that, we can just free_wrapper the old table and copy the
   elements from newtbl to hashtbl.
*/
hashtable_rc_t obj_hashtable_resize(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP) {
  obj_hash_table_t newtbl = {.mutex = PTHREAD_MUTEX_INITIALIZER, 0};
  hash_size_t n;
  obj_hash_node_t *node, *next;
  void* dummy = NULL;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }
  hash_size_t size = sizeP;
  // upper power of two:
  // http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2Float
  //  By Sean Eron Anderson
  // seander@cs.stanford.edu
  // Individually, the code snippets here are in the public domain (unless
  // otherwise noted) — feel free to use them however you please. The aggregate
  // collection and descriptions are © 1997-2005 Sean Eron Anderson. The code
  // and descriptions are distributed in the hope that they will be useful, but
  // WITHOUT ANY WARRANTY and without even the implied warranty of
  // merchantability or fitness for a particular purpose. As of May 5, 2005, all
  // the code has been tested thoroughly. Thousands of people have read it.
  // Moreover, Professor Randal Bryant, the Dean of Computer Science at Carnegie
  // Mellon University, has personally tested almost everything with his Uclid
  // code verification system. What he hasn't tested, I have checked against all
  // possible inputs on a 32-bit machine. To the first person to inform me of a
  // legitimate bug in the code, I'll pay a bounty of US$10 (by check or
  // Paypal). If directed to a charity, I'll pay US$20.
  size--;
  size |= size >> 1;
  size |= size >> 2;
  size |= size >> 4;
  size |= size >> 8;
  size |= size >> 16;
  size++;

  newtbl.size     = size;
  newtbl.hashfunc = hashtblP->hashfunc;

  if (!(newtbl.nodes = calloc(size, sizeof(obj_hash_node_t*))))
    return HASH_TABLE_SYSTEM_ERROR;

  for (n = 0; n < hashtblP->size; ++n) {
    for (node = hashtblP->nodes[n]; node; node = next) {
      next = node->next;
      obj_hashtable_remove(hashtblP, node->key, node->key_size, &dummy);
      obj_hashtable_insert(&newtbl, node->key, node->key_size, node->data);
    }
  }

  free_wrapper((void**) &hashtblP->nodes);
  hashtblP->size  = newtbl.size;
  hashtblP->nodes = newtbl.nodes;
  PRINT_HASHTABLE(hashtblP, "return OK\n");
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   Resizing
   The number of elements in a hash table is not always known when creating the
   table. If the number of elements grows too large, it will seriously reduce
   the performance of most hash table operations. If the number of elements are
   reduced, the hash table will waste memory. That is why we provide a function
   for resizing the table. Resizing a hash table is not as easy as a realloc().
   All hash values must be recalculated and each element must be inserted into
   its new position. We create a temporary obj_hash_table_t object (newtbl) to
   be used while building the new hashes. This allows us to reuse
   hashtable_insert() and hashtable_free(), when moving the elements to the new
   table. After that, we can just free_wrapper the old table and copy the
   elements from newtbl to hashtbl.
*/
hashtable_rc_t obj_hashtable_ts_resize(
    obj_hash_table_t* const hashtblP, const hash_size_t sizeP) {
  obj_hash_table_t newtbl = {.mutex = PTHREAD_MUTEX_INITIALIZER, 0};
  hash_size_t n;
  obj_hash_node_t *node, *next;
  void* dummy = NULL;

  if (hashtblP == NULL) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }
  hash_size_t size = sizeP;
  // upper power of two:
  // http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2Float
  //  By Sean Eron Anderson
  // seander@cs.stanford.edu
  // Individually, the code snippets here are in the public domain (unless
  // otherwise noted) — feel free to use them however you please. The aggregate
  // collection and descriptions are © 1997-2005 Sean Eron Anderson. The code
  // and descriptions are distributed in the hope that they will be useful, but
  // WITHOUT ANY WARRANTY and without even the implied warranty of
  // merchantability or fitness for a particular purpose. As of May 5, 2005, all
  // the code has been tested thoroughly. Thousands of people have read it.
  // Moreover, Professor Randal Bryant, the Dean of Computer Science at Carnegie
  // Mellon University, has personally tested almost everything with his Uclid
  // code verification system. What he hasn't tested, I have checked against all
  // possible inputs on a 32-bit machine. To the first person to inform me of a
  // legitimate bug in the code, I'll pay a bounty of US$10 (by check or
  // Paypal). If directed to a charity, I'll pay US$20.
  size--;
  size |= size >> 1;
  size |= size >> 2;
  size |= size >> 4;
  size |= size >> 8;
  size |= size >> 16;
  size++;

  newtbl.size     = size;
  newtbl.hashfunc = hashtblP->hashfunc;

  if (!(newtbl.nodes = calloc(size, sizeof(obj_hash_node_t*))))
    return HASH_TABLE_SYSTEM_ERROR;

  if (!(newtbl.lock_nodes = calloc(size, sizeof(pthread_mutex_t)))) {
    free_wrapper((void**) &newtbl.nodes);
    return HASH_TABLE_SYSTEM_ERROR;
  }
  for (n = 0; n < hashtblP->size; ++n) {
    pthread_mutex_init(&newtbl.lock_nodes[n], NULL);
  }

  pthread_mutex_lock(&hashtblP->mutex);
  for (n = 0; n < hashtblP->size; ++n) {
    for (node = hashtblP->nodes[n]; node; node = next) {
      next = node->next;
      obj_hashtable_ts_remove(hashtblP, node->key, node->key_size, &dummy);
      obj_hashtable_insert(&newtbl, node->key, node->key_size, node->data);
    }
  }

  free_wrapper((void**) &hashtblP->nodes);
  free_wrapper((void**) &hashtblP->nodes);
  hashtblP->size       = newtbl.size;
  hashtblP->nodes      = newtbl.nodes;
  hashtblP->lock_nodes = newtbl.lock_nodes;
  pthread_mutex_unlock(&hashtblP->mutex);
  PRINT_HASHTABLE(hashtblP, "return OK\n");
  return HASH_TABLE_OK;
}
