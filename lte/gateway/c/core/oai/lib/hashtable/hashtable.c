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

/*! \file hashtable.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <pthread.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "hashtable.h"

#if TRACE_HASHTABLE
#define PRINT_HASHTABLE(hTbLe, ...)                                            \
  do {                                                                         \
    if (hTbLe->log_enabled) OAILOG_TRACE(LOG_UTIL, ##__VA_ARGS__);             \
  } while (0)
#else
#define PRINT_HASHTABLE(...)
#endif
//------------------------------------------------------------------------------
char* hashtable_rc_code2string(hashtable_rc_t rcP) {
  switch (rcP) {
    case HASH_TABLE_OK:
      return "HASH_TABLE_OK";
      break;

    case HASH_TABLE_INSERT_OVERWRITTEN_DATA:
      return "HASH_TABLE_INSERT_OVERWRITTEN_DATA";
      break;

    case HASH_TABLE_KEY_NOT_EXISTS:
      return "HASH_TABLE_KEY_NOT_EXISTS";
      break;

    case HASH_TABLE_KEY_ALREADY_EXISTS:
      return "HASH_TABLE_KEY_ALREADY_EXISTS";
      break;

    case HASH_TABLE_BAD_PARAMETER_HASHTABLE:
      return "HASH_TABLE_BAD_PARAMETER_HASHTABLE";
      break;

    default:
      return "UNKNOWN hashtable_rc_t";
  }
}

//------------------------------------------------------------------------------
/*
   free_wrapper int function
   hash_free_int_func() is used when this hashtable is used to store int values
   as data (pointer = value).
*/

void hash_free_int_func(void** memoryP) {}

//------------------------------------------------------------------------------
/*
   Default hash function
   def_hashfunc() is the default used by hashtable_create() when the user didn't
   specify one. This is a simple/naive hash function which adds the key's ASCII
   char values. It will probably generate lots of collisions on large hash
   tables.
*/

static inline hash_size_t def_hashfunc(const uint64_t keyP) {
  return (hash_size_t) keyP;
}

//------------------------------------------------------------------------------
/*
   Initialization
   hashtable_init() set up the initial structure of the hash table. The user
   specified size will be allocated and initialized to NULL. The user can also
   specify a hash function. If the hashfunc argument is NULL, a default hash
   function is used. If an error occurred, NULL is returned. All other values in
   the returned hash_table_t pointer should be released with
   hashtable_destroy().
*/
hash_table_t* hashtable_init(
    hash_table_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const hash_key_t), void (*freefuncP)(void**),
    bstring display_name_pP) {
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

  if (!(hashtblP->nodes = calloc(size, sizeof(hash_node_t*)))) {
    free_wrapper((void**) &hashtblP);
    return NULL;
  }
  hashtblP->log_enabled = true;

  PRINT_HASHTABLE(hashtblP, "allocated nodes\n");
  hashtblP->size = size;

  if (hashfuncP)
    hashtblP->hashfunc = hashfuncP;
  else
    hashtblP->hashfunc = def_hashfunc;

  if (freefuncP)
    hashtblP->freefunc = freefuncP;
  else
    hashtblP->freefunc = free_wrapper;

  if (display_name_pP) {
    bassign(hashtblP->name, display_name_pP);
  } else {
    hashtblP->name = bformat("hashtable%u@%p", size, hashtblP);
  }
  hashtblP->is_allocated_by_malloc = false;
  return hashtblP;
}

//------------------------------------------------------------------------------
/*
   Initialization
   hashtable_create() allocate and sets up the initial structure of the hash
   table. The user specified size will be allocated and initialized to NULL. The
   user can also specify a hash function. If the hashfunc argument is NULL, a
   default hash function is used. If an error occurred, NULL is returned. All
   other values in the returned hash_table_t pointer should be released with
   hashtable_destroy().
*/
hash_table_t* hashtable_create(
    const hash_size_t sizeP, hash_size_t (*hashfuncP)(const hash_key_t),
    void (*freefuncP)(void**), bstring display_name_pP) {
  hash_table_t* hashtbl = NULL;

  if (!(hashtbl = calloc(1, sizeof(hash_table_t)))) {
    return NULL;
  }
  hashtbl =
      hashtable_init(hashtbl, sizeP, hashfuncP, freefuncP, display_name_pP);
  hashtbl->is_allocated_by_malloc = true;
  return hashtbl;
}

//------------------------------------------------------------------------------
/*
   Initialization
   hashtable_ts_init() sets up the initial structure of the thread safe hash
   table. The user specified size will be allocated and initialized to NULL. The
   user can also specify a hash function. If the hashfunc argument is NULL, a
   default hash function is used. If an error occurred, NULL is returned. All
   other values in the returned hash_table_t pointer should be released with
   hashtable_destroy().
*/
hash_table_ts_t* hashtable_ts_init(
    hash_table_ts_t* const hashtblP, const hash_size_t sizeP,
    hash_size_t (*hashfuncP)(const hash_key_t), void (*freefuncP)(void**),
    bstring display_name_pP) {
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

  memset(hashtblP, 0, sizeof(*hashtblP));

  if (!(hashtblP->nodes = calloc(size, sizeof(hash_node_t*)))) {
    free_wrapper((void**) &hashtblP);
    return NULL;
  }

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

  hashtblP->size = size;

  if (hashfuncP)
    hashtblP->hashfunc = hashfuncP;
  else
    hashtblP->hashfunc = def_hashfunc;

  if (freefuncP)
    hashtblP->freefunc = freefuncP;
  else
    hashtblP->freefunc = free_wrapper;

  if (display_name_pP) {
    hashtblP->name = bstrcpy(display_name_pP);
  } else {
    hashtblP->name = bformat("hashtable@%p", hashtblP);
  }
  hashtblP->is_allocated_by_malloc = false;
  hashtblP->log_enabled            = true;
  return hashtblP;
}
//------------------------------------------------------------------------------
/*
   Initialization
   hashtable_ts_create() allocate and sets up the initial structure of the
   thread safe hash table. The user specified size will be allocated and
   initialized to NULL. The user can also specify a hash function. If the
   hashfunc argument is NULL, a default hash function is used. If an error
   occurred, NULL is returned. All other values in the returned hash_table_t
   pointer should be released with hashtable_destroy().
*/
hash_table_ts_t* hashtable_ts_create(
    const hash_size_t sizeP, hash_size_t (*hashfuncP)(const hash_key_t),
    void (*freefuncP)(void**), bstring display_name_pP) {
  hash_table_ts_t* hashtbl = NULL;

  if (!(hashtbl = calloc(1, sizeof(hash_table_ts_t)))) {
    return NULL;
  }
  hashtbl =
      hashtable_ts_init(hashtbl, sizeP, hashfuncP, freefuncP, display_name_pP);
  hashtbl->is_allocated_by_malloc = true;
  return hashtbl;
}

//------------------------------------------------------------------------------
/*
   Cleanup
   The hashtable_destroy() walks through the linked lists for each possible hash
   value, and releases the elements. It also releases the nodes array and the
   hash_table_t.
*/
hashtable_rc_t hashtable_destroy(hash_table_t* hashtblP) {
  hash_size_t n     = 0;
  hash_node_t *node = NULL, *oldnode = NULL;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  for (n = 0; n < hashtblP->size; ++n) {
    node = hashtblP->nodes[n];

    while (node) {
      oldnode = node;
      node    = node->next;

      if (oldnode->data) {
        hashtblP->freefunc(&oldnode->data);
      }

      free_wrapper((void**) &oldnode);
    }
  }

  free_wrapper((void**) &hashtblP->nodes);
  bdestroy_wrapper(&hashtblP->name);
  if (hashtblP->is_allocated_by_malloc) {
    free_wrapper((void**) &hashtblP);
  }
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   Cleanup
   The hashtable_destroy() walks through the linked lists for each possible
   hash value, and releases the elements. It also releases the nodes array and
   the hash_table_t.
*/
hashtable_rc_t hashtable_ts_destroy(hash_table_ts_t* hashtblP) {
  hash_size_t n     = 0;
  hash_node_t *node = NULL, *oldnode = NULL;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  for (n = 0; n < hashtblP->size; ++n) {
    pthread_mutex_lock(&hashtblP->lock_nodes[n]);
    node = hashtblP->nodes[n];

    while (node) {
      oldnode = node;
      node    = node->next;

      if (oldnode->data) {
        hashtblP->freefunc(&oldnode->data);
      }

      free_wrapper((void**) &oldnode);
    }

    pthread_mutex_unlock(&hashtblP->lock_nodes[n]);
    pthread_mutex_destroy(&hashtblP->lock_nodes[n]);
  }

  free_wrapper((void**) &hashtblP->nodes);
  bdestroy_wrapper(&hashtblP->name);
  free_wrapper((void**) &hashtblP->lock_nodes);
  free_wrapper((void**) &hashtblP->lock_attr);
  if (hashtblP->is_allocated_by_malloc) {
    free_wrapper((void**) &hashtblP);
  }
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
hashtable_rc_t hashtable_is_key_exists(
    const hash_table_t* const hashtblP, const hash_key_t keyP) {
  hash_node_t* node = NULL;
  hash_size_t hash  = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 ") return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
hashtable_rc_t hashtable_ts_is_key_exists(
    const hash_table_ts_t* const hashtblP, const hash_key_t keyP) {
  hash_node_t* node = NULL;
  hash_size_t hash  = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 ") return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
// may cost a lot CPU...
hashtable_key_array_t* hashtable_ts_get_keys(hash_table_ts_t* const hashtblP) {
  hash_node_t* node         = NULL;
  unsigned int i            = 0;
  hashtable_key_array_t* ka = NULL;

  ka = calloc(1, sizeof(hashtable_key_array_t));
  if (ka == NULL) return NULL;

  if (hashtblP->num_elements == 0) {
    free(ka);
    return NULL;
  }

  ka->keys = calloc(hashtblP->num_elements, sizeof(hash_key_t));
  if (ka->keys == NULL) {
    free(ka);
    return NULL;
  }

  while ((ka->num_keys < hashtblP->num_elements) && (i < hashtblP->size)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
      while (node) {
        ka->keys[ka->num_keys++] = node->key;
        node                     = node->next;
      }
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    i++;
  }
  return ka;
}

//------------------------------------------------------------------------------
// may cost a lot CPU...
hashtable_element_array_t* hashtable_ts_get_elements(
    hash_table_ts_t* const hashtblP) {
  hash_node_t* node             = NULL;
  unsigned int i                = 0;
  hashtable_element_array_t* ea = NULL;

  if ((!hashtblP) || !(hashtblP->num_elements)) {
    return NULL;
  }
  ea           = calloc(1, sizeof(hashtable_element_array_t));
  ea->elements = calloc(hashtblP->num_elements, sizeof(void*));

  while ((ea->num_elements < hashtblP->num_elements) && (i < hashtblP->size)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
      while (node) {
        ea->elements[ea->num_elements++] = node->data;
        node                             = node->next;
      }
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    i++;
  }
  return ea;
}

//------------------------------------------------------------------------------
// may cost a lot CPU...
// Also useful if we want to find an element in the collection based on compare
// criteria different than the single key The compare criteria in implemented in
// the funct_cb function
hashtable_rc_t hashtable_apply_callback_on_elements(
    hash_table_t* const hashtblP,
    bool funct_cb(
        hash_key_t keyP, void* dataP, void* parameterP, void** resultP),
    void* parameterP, void** resultP) {
  hash_node_t* node         = NULL;
  unsigned int i            = 0;
  unsigned int num_elements = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size)) {
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];

      while (node) {
        num_elements++;
        if (funct_cb(node->key, node->data, parameterP, resultP)) {
          return HASH_TABLE_OK;
        }
        node = node->next;
      }
    }
    i++;
  }

  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
// may cost a lot CPU...
// Also useful if we want to find an element in the collection based on compare
// criteria different than the single key The compare criteria in implemented
// in the funct_cb function
hashtable_rc_t hashtable_ts_apply_callback_on_elements(
    hash_table_ts_t* const hashtblP,
    bool funct_cb(
        const hash_key_t keyP, void* const dataP, void* parameterP,
        void** resultP),
    void* parameterP, void** resultP) {
  hash_node_t* node         = NULL;
  unsigned int i            = 0;
  unsigned int num_elements = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];

      while (node) {
        num_elements++;
        if (funct_cb(node->key, node->data, parameterP, resultP)) {
          pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
          return HASH_TABLE_OK;
        }
        node = node->next;
      }
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    i++;
  }

  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
hashtable_rc_t hashtable_dump_content(
    const hash_table_t* const hashtblP, bstring str) {
  hash_node_t* node = NULL;
  unsigned int i    = 0;

  if (!hashtblP) {
    bcatcstr(str, "HASH_TABLE_BAD_PARAMETER_HASHTABLE");
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  while (i < hashtblP->size) {
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];

      while (node) {
        bstring b0 = bformat(
            "Key 0x%" PRIx64 " Element %p Node %p\n", node->key, node->data,
            node);
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
hashtable_rc_t hashtable_ts_dump_content(
    const hash_table_ts_t* const hashtblP, bstring str) {
  hash_node_t* node = NULL;
  unsigned int i    = 0;

  if (!hashtblP) {
    bcatcstr(str, "HASH_TABLE_BAD_PARAMETER_HASHTABLE");
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  while (i < hashtblP->size) {
    if (hashtblP->nodes[i] != NULL) {
      pthread_mutex_lock(&hashtblP->lock_nodes[i]);
      node = hashtblP->nodes[i];

      while (node) {
        bstring b0 = bformat(
            "Key 0x%" PRIx64 " Element %p Node %p Next %p\n", node->key,
            node->data, node, node->next);
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
hashtable_rc_t hashtable_insert(
    hash_table_t* const hashtblP, const hash_key_t keyP, void* dataP) {
  hash_node_t* node = NULL;
  hash_size_t hash  = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if ((node->data) && (node->data != dataP)) {
        hashtblP->freefunc(&node->data);

        node->data = dataP;
        PRINT_HASHTABLE(
            hashtblP,
            "%s(%s,key 0x%" PRIx64 " data %p) return INSERT_OVERWRITTEN_DATA\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, dataP);
        return HASH_TABLE_INSERT_OVERWRITTEN_DATA;
      }
      node->data = dataP;
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 " data %p) return OK\n",
          __FUNCTION__, bdata(hashtblP->name), keyP, dataP);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }

  if (!(node = malloc(sizeof(hash_node_t)))) return -1;

  node->key  = keyP;
  node->data = dataP;

  if (hashtblP->nodes[hash]) {
    node->next = hashtblP->nodes[hash];
  } else {
    node->next = NULL;
  }

  hashtblP->nodes[hash] = node;
  hashtblP->num_elements += 1;

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 " data %p) return OK\n", __FUNCTION__,
      bdata(hashtblP->name), keyP, dataP);
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   Adding a new element
   To make sure the hash value is not bigger than size, the result of the user
   provided hash function is used modulo size.
*/
hashtable_rc_t hashtable_ts_insert(
    hash_table_ts_t* const hashtblP, const hash_key_t keyP, void* dataP) {
  hash_node_t* node = NULL;
  hash_size_t hash  = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if ((node->data) && (node->data != dataP)) {
        hashtblP->freefunc(&node->data);
        node->data = dataP;
        pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
        PRINT_HASHTABLE(
            hashtblP,
            "%s(%s,key 0x%" PRIx64 " data %p) return INSERT_OVERWRITTEN_DATA\n",
            __FUNCTION__, bdata(hashtblP->name), keyP, dataP);
        return HASH_TABLE_INSERT_OVERWRITTEN_DATA;
      }
      node->data = dataP;
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 " data %p) return OK\n",
          __FUNCTION__, bdata(hashtblP->name), keyP, dataP);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }

  if (!(node = malloc(sizeof(hash_node_t)))) return -1;

  node->key  = keyP;
  node->data = dataP;

  if (hashtblP->nodes[hash]) {
    node->next = hashtblP->nodes[hash];
  } else {
    node->next = NULL;
  }

  hashtblP->nodes[hash] = node;
  __sync_fetch_and_add(&hashtblP->num_elements, 1);
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 " data %p) next %p return OK\n",
      __FUNCTION__, bdata(hashtblP->name), keyP, dataP, node->next);
  return HASH_TABLE_OK;
}

//------------------------------------------------------------------------------
/*
   To free_wrapper an element from the hash table, we just search for it in the
   linked list for that hash value, and free_wrapper it if it is found. If it
   was not found, it is an error and -1 is returned.
*/
hashtable_rc_t hashtable_free(
    hash_table_t* const hashtblP, const hash_key_t keyP) {
  hash_node_t *node, *prevnode = NULL;
  hash_size_t hash = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if (prevnode)
        prevnode->next = node->next;
      else
        hashtblP->nodes[hash] = node->next;

      if (node->data) {
        hashtblP->freefunc(&node->data);
      }

      free_wrapper((void**) &node);
      __sync_fetch_and_sub(&hashtblP->num_elements, 1);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 ") return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   To free_wrapper an element from the hash table, we just search for it in the
   linked list for that hash value, and free_wrapper it if it is found. If it
   was not found, it is an error and -1 is returned.
*/
hashtable_rc_t hashtable_ts_free(
    hash_table_ts_t* const hashtblP, const hash_key_t keyP) {
  hash_node_t *node, *prevnode = NULL;
  hash_size_t hash = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if (prevnode)
        prevnode->next = node->next;
      else
        hashtblP->nodes[hash] = node->next;

      if (node->data) {
        hashtblP->freefunc(&node->data);
      }

      free_wrapper((void**) &node);
      __sync_fetch_and_sub(&hashtblP->num_elements, 1);
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 ") return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }

  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   To remove an element from the hash table, we just search for it in the linked
   list for that hash value, and remove it if it is found. If it was not found,
   it is an error and -1 is returned.
*/
hashtable_rc_t hashtable_remove(
    hash_table_t* const hashtblP, const hash_key_t keyP, void** dataP) {
  hash_node_t *node, *prevnode = NULL;
  hash_size_t hash = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if (prevnode)
        prevnode->next = node->next;
      else
        hashtblP->nodes[hash] = node->next;

      *dataP = node->data;
      free_wrapper((void**) &node);
      __sync_fetch_and_sub(&hashtblP->num_elements, 1);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 ") return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   To remove an element from the hash table, we just search for it in the linked
   list for that hash value, and remove it if it is found. If it was not found,
   it is an error and -1 is returned.
*/
hashtable_rc_t hashtable_ts_remove(
    hash_table_ts_t* const hashtblP, const hash_key_t keyP, void** dataP) {
  hash_node_t *node, *prevnode = NULL;
  hash_size_t hash = 0;

  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      if (prevnode)
        prevnode->next = node->next;
      else
        hashtblP->nodes[hash] = node->next;

      *dataP = node->data;
      free_wrapper((void**) &node);
      __sync_fetch_and_sub(&hashtblP->num_elements, 1);
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 ") return OK\n", __FUNCTION__,
          bdata(hashtblP->name), keyP);
      return HASH_TABLE_OK;
    }

    prevnode = node;
    node     = node->next;
  }
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   Searching for an element is easy. We just search through the linked list for
   the corresponding hash value. NULL is returned if we didn't find it.
*/
hashtable_rc_t hashtable_get(
    const hash_table_t* const hashtblP, const hash_key_t keyP, void** dataP) {
  hash_node_t* node = NULL;
  hash_size_t hash  = 0;

  *dataP = NULL;
  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      *dataP = node->data;
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 " data %p) return OK\n",
          __FUNCTION__, bdata(hashtblP->name), keyP, *dataP);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }

  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);
  return HASH_TABLE_KEY_NOT_EXISTS;
}

//------------------------------------------------------------------------------
/*
   Searching for an element is easy. We just search through the linked list for
   the corresponding hash value. NULL is returned if we didn't find it.
*/
hashtable_rc_t hashtable_ts_get(
    const hash_table_ts_t* const hashtblP, const hash_key_t keyP,
    void** dataP) {
  hash_node_t* node = NULL;
  hash_size_t hash  = 0;

  *dataP = NULL;
  if (!hashtblP) {
    return HASH_TABLE_BAD_PARAMETER_HASHTABLE;
  }

  hash = hashtblP->hashfunc(keyP) % hashtblP->size;

  pthread_mutex_lock(&hashtblP->lock_nodes[hash]);
  node = hashtblP->nodes[hash];

  while (node) {
    if (node->key == keyP) {
      *dataP = node->data;
      pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
      PRINT_HASHTABLE(
          hashtblP, "%s(%s,key 0x%" PRIx64 " data %p) return OK\n",
          __FUNCTION__, bdata(hashtblP->name), keyP, *dataP);
      return HASH_TABLE_OK;
    }

    node = node->next;
  }
  pthread_mutex_unlock(&hashtblP->lock_nodes[hash]);
  PRINT_HASHTABLE(
      hashtblP, "%s(%s,key 0x%" PRIx64 ") return KEY_NOT_EXISTS\n",
      __FUNCTION__, bdata(hashtblP->name), keyP);

#define TEMPORARY_DEBUG 1
#if TEMPORARY_DEBUG
  bstring b = bfromcstr(" ");
  hashtable_ts_dump_content(hashtblP, b);
  PRINT_HASHTABLE(hashtblP, "%s:%s\n", bdata(hashtblP->name), bdata(b));
  bdestroy(b);
#endif
  return HASH_TABLE_KEY_NOT_EXISTS;
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
   its new position. We create a temporary hash_table_t object (newtbl) to be
   used while building the new hashes. This allows us to reuse
   hashtable_insert() and hashtable_free(), when moving the elements to the new
   table. After that, we can just free_wrapper the old table and copy the
   elements from newtbl to hashtbl.
*/

hashtable_rc_t hashtable_resize(
    hash_table_t* const hashtblP, const hash_size_t sizeP) {
  hash_table_t newtbl;
  hash_size_t n;
  hash_node_t *node, *next;
  void* dummy = NULL;

  if (!hashtblP) {
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

  if (!(newtbl.nodes = calloc(size, sizeof(hash_node_t*)))) return -1;

  for (n = 0; n < hashtblP->size; ++n) {
    for (node = hashtblP->nodes[n]; node; node = next) {
      next = node->next;
      hashtable_remove(hashtblP, node->key, &dummy);
      hashtable_insert(&newtbl, node->key, node->data);
    }
  }

  free_wrapper((void**) &hashtblP->nodes);
  hashtblP->nodes = newtbl.nodes;
  hashtblP->size  = newtbl.size;
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
   its new position. We create a temporary hash_table_t object (newtbl) to be
   used while building the new hashes. This allows us to reuse
   hashtable_insert() and hashtable_free(), when moving the elements to the new
   table. After that, we can just free_wrapper the old table and copy the
   elements from newtbl to hashtbl. Dangerous not really thread safe.
*/

hashtable_rc_t hashtable_ts_resize(
    hash_table_ts_t* const hashtblP, const hash_size_t sizeP) {
  hash_table_ts_t newtbl = {.mutex = PTHREAD_MUTEX_INITIALIZER, 0};
  hash_size_t n          = 0;
  hash_node_t *node = NULL, *next = NULL;
  void* dummy = NULL;

  if (!hashtblP) {
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

  if (!(newtbl.nodes = calloc(size, sizeof(hash_node_t*))))
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
      hashtable_ts_remove(hashtblP, node->key, &dummy);
      hashtable_ts_insert(&newtbl, node->key, node->data);
    }
  }

  free_wrapper((void**) &hashtblP->nodes);
  free_wrapper((void**) &hashtblP->lock_nodes);
  hashtblP->size       = newtbl.size;
  hashtblP->nodes      = newtbl.nodes;
  hashtblP->lock_nodes = newtbl.lock_nodes;
  pthread_mutex_unlock(&hashtblP->mutex);
  return HASH_TABLE_OK;
}
