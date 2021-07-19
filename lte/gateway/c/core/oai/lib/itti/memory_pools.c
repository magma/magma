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

#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>

#include "assertions.h"
#include "memory_pools.h"

/*------------------------------------------------------------------------------*/
const static int mp_debug = 0;

#define MP_DEBUG(x, args...)                                                   \
  do {                                                                         \
    if (mp_debug) fprintf(stdout, "[MP][D]" x, ##args);                        \
    fflush(stdout);                                                            \
  } while (0)

/*------------------------------------------------------------------------------*/
#define CHARS_TO_UINT32(c1, c2, c3, c4)                                        \
  ((((uint32_t) c1) << 24) | (((uint32_t) c2) << 16) |                         \
   (((uint32_t) c3) << 8) | ((uint32_t) c4))

#define MEMORY_POOL_ITEM_INFO_NUMBER 2

/*------------------------------------------------------------------------------*/
typedef int32_t items_group_position_t;
typedef int32_t items_group_index_t;

typedef union items_group_positions_u {
  uint64_t all;
  struct {
    items_group_position_t put;
    items_group_position_t get;
  } ind;
} items_group_positions_t;

typedef struct items_group_s {
  items_group_position_t number_plus_one;
  volatile uint32_t minimum;
  volatile items_group_positions_t positions;
  volatile items_group_index_t* indexes;
} items_group_t;

/*------------------------------------------------------------------------------*/
static const items_group_index_t ITEMS_GROUP_INDEX_INVALID = -1;

/*------------------------------------------------------------------------------*/
typedef uint32_t pool_item_start_mark_t;
typedef uint32_t pool_item_end_mark_t;

typedef uint32_t memory_pool_data_t;

typedef uint32_t pool_start_mark_t;

typedef uint32_t pools_start_mark_t;

typedef uint8_t pool_id_t;
typedef uint8_t item_status_t;

typedef struct memory_pool_item_start_s {
  pool_item_start_mark_t start_mark;

  pool_id_t pool_id;
  item_status_t item_status;
  uint16_t info[MEMORY_POOL_ITEM_INFO_NUMBER];
} memory_pool_item_start_t;

typedef struct memory_pool_item_end_s {
  pool_item_end_mark_t end_mark;
} memory_pool_item_end_t;

typedef struct memory_pool_item_s {
  memory_pool_item_start_t start;
  memory_pool_data_t data[0];
  memory_pool_item_end_t end;
} memory_pool_item_t;

typedef struct memory_pool_s {
  pool_start_mark_t start_mark;

  pool_id_t pool_id;
  uint32_t item_data_number;
  uint32_t pool_item_size;
  items_group_t items_group_free;
  memory_pool_item_t* items;
} memory_pool_t;

typedef struct memory_pools_s {
  pools_start_mark_t start_mark;

  uint32_t pools_number;
  uint32_t pools_defined;
  memory_pool_t* pools;
} memory_pools_t;

//------------------------------------------------------------------------------
static const uint32_t MAX_POOLS_NUMBER      = 20;
static const uint32_t MAX_POOL_ITEMS_NUMBER = 200 * 1000;
static const uint32_t MAX_POOL_ITEM_SIZE    = 100 * 1000;

static const pool_item_start_mark_t POOL_ITEM_START_MARK =
    CHARS_TO_UINT32('P', 'I', 's', 't');
static const pool_item_end_mark_t POOL_ITEM_END_MARK =
    CHARS_TO_UINT32('p', 'i', 'E', 'N');

static const item_status_t ITEM_STATUS_FREE      = 'F';
static const item_status_t ITEM_STATUS_ALLOCATED = 'a';

static const pool_start_mark_t POOL_START_MARK =
    CHARS_TO_UINT32('P', '_', 's', 't');

static const pools_start_mark_t POOLS_START_MARK =
    CHARS_TO_UINT32('P', 'S', 's', 't');

/*------------------------------------------------------------------------------*/
static inline uint32_t items_group_number_items(items_group_t* items_group) {
  return items_group->number_plus_one - 1;
}

//------------------------------------------------------------------------------
static inline uint32_t items_group_free_items(items_group_t* items_group) {
  items_group_positions_t positions;
  uint32_t free_items;

  positions.all = items_group->positions.all;
  free_items =
      items_group->number_plus_one + positions.ind.put - positions.ind.get;
  free_items %= items_group->number_plus_one;
  return free_items;
}

//------------------------------------------------------------------------------
static inline items_group_index_t items_group_get_free_item(
    items_group_t* items_group) {
  items_group_position_t get_raw;
  items_group_position_t put;
  items_group_position_t get;
  items_group_position_t free_items;
  items_group_index_t index = ITEMS_GROUP_INDEX_INVALID, tmp;

  /*
   * Get current put position
   */
  put = items_group->positions.ind.put % items_group->number_plus_one;
  /*
   * Get current get position and increase it
   */
  get_raw = __sync_fetch_and_add(&items_group->positions.ind.get, 1);
  get     = get_raw % items_group->number_plus_one;

  if (put == get) {
    /*
     * No more item free, restore previous position
     */
    __sync_fetch_and_sub(&items_group->positions.ind.get, 1);
  } else {
    /*
     * Get index at current get position,
     * atomically compare this index and index at the get position
     * to either make the index unavailable or receive an invalid
     * index as, due to race conditions, another thread acquired this index.
     */
    index = items_group->indexes[get];
    tmp   = __sync_val_compare_and_swap(
        &items_group->indexes[get], index, ITEMS_GROUP_INDEX_INVALID);
    if (tmp != index) {
      index = ITEMS_GROUP_INDEX_INVALID;
    }

    if (index <= ITEMS_GROUP_INDEX_INVALID) {
      /*
       * Index has not yet been completely freed, restore previous get position
       */
      __sync_fetch_and_sub(&items_group->positions.ind.get, 1);
    } else {
      if (get_raw == items_group->number_plus_one) {
        /*
         * Wrap get position
         */
        __sync_fetch_and_sub(
            &items_group->positions.ind.get, items_group->number_plus_one);
      }

      free_items = items_group_free_items(items_group);

      /*
       * Updates minimum free items if needed
       */
      while (items_group->minimum > free_items) {
        items_group->minimum = free_items;
      }
    }
  }

  return (index);
}

//------------------------------------------------------------------------------
static inline int items_group_put_free_item(
    items_group_t* items_group, items_group_index_t index) {
  items_group_position_t put_raw;
  items_group_position_t put;

  /*
   * Get current put position and increase it
   */
  put_raw = __sync_fetch_and_add(&items_group->positions.ind.put, 1);
  put     = put_raw % items_group->number_plus_one;

  if (put_raw == items_group->number_plus_one) {
    /*
     * Wrap position
     */
    __sync_fetch_and_sub(
        &items_group->positions.ind.put, items_group->number_plus_one);
  }

  AssertError(
      items_group->indexes[put] <= ITEMS_GROUP_INDEX_INVALID,
      return (EXIT_FAILURE),
      "Index at current put position (%d) is not marked as free (%d)!\n", put,
      items_group->number_plus_one);
  /*
   * Save freed item index at current put position
   */
  items_group->indexes[put] = index;
  return (EXIT_SUCCESS);
}

//------------------------------------------------------------------------------
static inline memory_pools_t* memory_pools_from_handler(
    memory_pools_handle_t memory_pools_handle) {
  memory_pools_t* memory_pools;

  /*
   * Recover memory_pools
   */
  memory_pools = (memory_pools_t*) memory_pools_handle;
  /*
   * Sanity check on passed handle
   */
  AssertError(
      memory_pools->start_mark == POOLS_START_MARK, memory_pools = NULL,
      "Handle %p is not a valid memory pools handle, start mark is missing!\n",
      memory_pools_handle);
  return (memory_pools);
}

//------------------------------------------------------------------------------
static inline memory_pool_item_t* memory_pool_item_from_handler(
    memory_pool_item_handle_t memory_pool_item_handle) {
  void* address;
  memory_pool_item_t* memory_pool_item;

  /*
   * Recover memory_pools
   */
  address          = memory_pool_item_handle - sizeof(memory_pool_item_start_t);
  memory_pool_item = (memory_pool_item_t*) address;
  /*
   * Sanity check on passed handle
   */
  // AssertError (memory_pool_item->start.start_mark == POOL_ITEM_START_MARK,
  // memory_pool_item = NULL, "Handle %p is not a valid memory pool item handle,
  // start mark is missing!\n", memory_pool_item);
  AssertFatal(
      memory_pool_item->start.start_mark == POOL_ITEM_START_MARK,
      "Handle %p is not a valid memory pool item handle, start mark is "
      "missing!\n",
      memory_pool_item);
  return (memory_pool_item);
}

//------------------------------------------------------------------------------
static inline memory_pool_item_t* memory_pool_item_from_index(
    memory_pool_t* memory_pool, items_group_index_t index) {
  void* address;

  address = (void*) memory_pool->items;
  address += index * memory_pool->pool_item_size;
  return (address);
}

//------------------------------------------------------------------------------
memory_pools_handle_t memory_pools_create(uint32_t pools_number) {
  memory_pools_t* memory_pools;
  pool_id_t pool;

  AssertFatal(
      pools_number <= MAX_POOLS_NUMBER,
      "Too many memory pools requested (%d/%d)!\n", pools_number,
      MAX_POOLS_NUMBER); /* Limit to a reasonable number of pools */
  /*
   * Allocate memory_pools
   */
  memory_pools = malloc(sizeof(memory_pools_t));
  AssertFatal(
      memory_pools != NULL, "Memory pools structure allocation failed!\n");
  /*
   * Initialize memory_pools
   */
  {
    memory_pools->start_mark    = POOLS_START_MARK;
    memory_pools->pools_number  = pools_number;
    memory_pools->pools_defined = 0;
    /*
     * Allocate pools
     */
    memory_pools->pools = calloc(pools_number, sizeof(memory_pool_t));
    AssertFatal(
        memory_pools->pools != NULL, "Memory pools allocation failed!\n");

    /*
     * Initialize pools
     */
    for (pool = 0; pool < pools_number; pool++) {
      memory_pools->pools[pool].start_mark = POOL_START_MARK;
    }
  }
  return ((memory_pools_handle_t) memory_pools);
}

//------------------------------------------------------------------------------
char* memory_pools_statistics(memory_pools_handle_t memory_pools_handle) {
  memory_pools_t* memory_pools;
  pool_id_t pool;
  char* statistics;
  int statistics_len;
  int printed_chars;
  uint32_t allocated_pool_memory;
  uint32_t allocated_pools_memory = 0;
  items_group_t* items_group;
  uint32_t pool_items_size;

  /*
   * Recover memory_pools
   */
  memory_pools = memory_pools_from_handler(memory_pools_handle);
  AssertFatal(
      memory_pools != NULL, "Failed to retrieve memory pool for handle %p!\n",
      memory_pools_handle);
  statistics_len = memory_pools->pools_defined * 200;
  statistics     = malloc(statistics_len);
  printed_chars  = snprintf(
      &statistics[0], statistics_len,
      "Pool:   size, number, minimum,   free, address space and memory used in "
      "Kbytes\n");

  for (pool = 0; pool < memory_pools->pools_defined; pool++) {
    items_group           = &memory_pools->pools[pool].items_group_free;
    allocated_pool_memory = items_group_number_items(items_group) *
                            memory_pools->pools[pool].pool_item_size;
    allocated_pools_memory += allocated_pool_memory;
    pool_items_size =
        memory_pools->pools[pool].item_data_number * sizeof(memory_pool_data_t);
    printed_chars += snprintf(
        &statistics[printed_chars], statistics_len - printed_chars,
        "  %2u: %6u, %6u,  %6u, %6u, [%p-%p] %6u\n", pool, pool_items_size,
        items_group_number_items(items_group), items_group->minimum,
        items_group_free_items(items_group), memory_pools->pools[pool].items,
        ((void*) memory_pools->pools[pool].items) + allocated_pool_memory,
        allocated_pool_memory / (1024));
  }

  printed_chars = snprintf(
      &statistics[printed_chars], statistics_len - printed_chars,
      "Pools memory %u Kbytes\n", allocated_pools_memory / (1024));
  return (statistics);
}

//------------------------------------------------------------------------------
int memory_pools_add_pool(
    memory_pools_handle_t memory_pools_handle, uint32_t pool_items_number,
    uint32_t pool_item_size) {
  memory_pools_t* memory_pools;
  memory_pool_t* memory_pool;
  pool_id_t pool;
  items_group_index_t item_index;
  memory_pool_item_t* memory_pool_item;

  AssertFatal(
      pool_items_number <= MAX_POOL_ITEMS_NUMBER,
      "Too many items for a memory pool (%u/%d)!\n", pool_items_number,
      MAX_POOL_ITEMS_NUMBER); /* Limit to a reasonable number of items */
  AssertFatal(
      pool_item_size <= MAX_POOL_ITEM_SIZE,
      "Item size is too big for memory pool items (%u/%d)!\n", pool_item_size,
      MAX_POOL_ITEM_SIZE); /* Limit to a reasonable item size */
  /*
   * Recover memory_pools
   */
  memory_pools = memory_pools_from_handler(memory_pools_handle);
  AssertFatal(
      memory_pools != NULL, "Failed to retrieve memory pool for handle %p!\n",
      memory_pools_handle);
  /*
   * Check number of already created pools
   */
  AssertFatal(
      memory_pools->pools_defined < memory_pools->pools_number,
      "Can not allocate more memory pool (%d)!\n", memory_pools->pools_number);
  /*
   * Select pool
   */
  pool        = memory_pools->pools_defined;
  memory_pool = &memory_pools->pools[pool];
  /*
   * Initialize pool
   */
  {
    memory_pool->pool_id = pool;
    /*
     * Item size in memory_pool_data_t items by excess
     */
    memory_pool->item_data_number =
        (pool_item_size + sizeof(memory_pool_data_t) - 1) /
        sizeof(memory_pool_data_t);
    memory_pool->pool_item_size =
        (memory_pool->item_data_number * sizeof(memory_pool_data_t)) +
        sizeof(memory_pool_item_t);
    memory_pool->items_group_free.number_plus_one   = pool_items_number + 1;
    memory_pool->items_group_free.minimum           = pool_items_number;
    memory_pool->items_group_free.positions.ind.put = pool_items_number;
    memory_pool->items_group_free.positions.ind.get = 0;
    /*
     * Allocate free indexes
     */
    memory_pool->items_group_free.indexes = malloc(
        memory_pool->items_group_free.number_plus_one *
        sizeof(items_group_index_t));
    AssertFatal(
        memory_pool->items_group_free.indexes != NULL,
        "Memory pool indexes allocation failed!\n");

    /*
     * Initialize free indexes
     */
    for (item_index = 0; item_index < pool_items_number; item_index++) {
      memory_pool->items_group_free.indexes[item_index] = item_index;
    }

    /*
     * Last index is not allocated
     */
    memory_pool->items_group_free.indexes[item_index] =
        ITEMS_GROUP_INDEX_INVALID;
    /*
     * Allocate items
     */
    memory_pool->items = calloc(pool_items_number, memory_pool->pool_item_size);
    AssertFatal(
        memory_pool->items != NULL, "Memory pool items allocation failed!\n");

    /*
     * Initialize items
     */
    for (item_index = 0; item_index < pool_items_number; item_index++) {
      memory_pool_item = memory_pool_item_from_index(memory_pool, item_index);
      memory_pool_item->start.start_mark  = POOL_ITEM_START_MARK;
      memory_pool_item->start.pool_id     = pool;
      memory_pool_item->start.item_status = ITEM_STATUS_FREE;
      memory_pool_item->data[memory_pool->item_data_number] =
          POOL_ITEM_END_MARK;
    }
  }
  memory_pools->pools_defined++;
  return (0);
}

//------------------------------------------------------------------------------
memory_pool_item_handle_t memory_pools_allocate(
    memory_pools_handle_t memory_pools_handle, uint32_t item_size,
    uint16_t info_0, uint16_t info_1) {
  memory_pools_t* memory_pools;
  memory_pool_item_t* memory_pool_item;
  memory_pool_item_handle_t memory_pool_item_handle = NULL;
  pool_id_t pool;
  items_group_index_t item_index = ITEMS_GROUP_INDEX_INVALID;

  /*
   * Recover memory_pools
   */
  memory_pools = memory_pools_from_handler(memory_pools_handle);
  AssertError(
      memory_pools != NULL, {},
      "Failed to retrieve memory pool for handle %p!\n", memory_pools_handle);

  for (pool = 0; pool < memory_pools->pools_defined; pool++) {
    if ((memory_pools->pools[pool].item_data_number *
         sizeof(memory_pool_data_t)) < item_size) {
      /*
       * This memory pool has too small items, skip it
       */
      continue;
    }

    item_index =
        items_group_get_free_item(&memory_pools->pools[pool].items_group_free);

    if (item_index <= ITEMS_GROUP_INDEX_INVALID) {
      /*
       * Allocation failed, skip this pool
       */
      continue;
    } else {
      /*
       * Allocation succeed, exit searching loop
       */
      break;
    }
  }

  if (item_index > ITEMS_GROUP_INDEX_INVALID) {
    /*
     * Convert item index into memory_pool_item address
     */
    memory_pool_item =
        memory_pool_item_from_index(&memory_pools->pools[pool], item_index);
    /*
     * Sanity check on item status, must be free
     */
    AssertFatal(
        memory_pool_item->start.item_status == ITEM_STATUS_FREE,
        "Item status is not set to free (%d) in pool %u, item %d!\n",
        memory_pool_item->start.item_status, pool, item_index);
    memory_pool_item->start.item_status = ITEM_STATUS_ALLOCATED;
    memory_pool_item->start.info[0]     = info_0;
    memory_pool_item->start.info[1]     = info_1;
    memory_pool_item_handle             = memory_pool_item->data;
    MP_DEBUG(
        " Alloc [%2u][%6d]{%6d}, %3u %3u, %6u, %p, %p, %p\n", pool, item_index,
        items_group_free_items(&memory_pools->pools[pool].items_group_free),
        info_0, info_1, item_size, memory_pools->pools[pool].items,
        memory_pool_item, memory_pool_item_handle);
  } else {
    MP_DEBUG(
        " Alloc [--][------]{------}, %3u %3u, %6u, failed!\n", info_0, info_1,
        item_size);
  }

  return memory_pool_item_handle;
}

//------------------------------------------------------------------------------
int memory_pools_free(
    memory_pools_handle_t memory_pools_handle,
    memory_pool_item_handle_t memory_pool_item_handle, uint16_t info_0) {
  memory_pools_t* memory_pools;
  memory_pool_item_t* memory_pool_item;
  pool_id_t pool;
  items_group_index_t item_index;
  uint32_t item_size;
  uint32_t pool_item_size;
  uint16_t info_1;
  int result;

  /*
   * Recover memory_pools
   */
  memory_pools = memory_pools_from_handler(memory_pools_handle);
  AssertError(
      memory_pools != NULL, return (EXIT_FAILURE),
      "Failed to retrieve memory pools for handle %p!\n", memory_pools_handle);
  /*
   * Recover memory pool item
   */
  memory_pool_item = memory_pool_item_from_handler(memory_pool_item_handle);
  AssertError(
      memory_pool_item != NULL, return (EXIT_FAILURE),
      "Failed to retrieve memory pool item for handle %p!\n",
      memory_pool_item_handle);
  info_1 = memory_pool_item->start.info[1];

  /*
   * Recover pool index
   */
  pool = memory_pool_item->start.pool_id;
  AssertFatal(
      pool < memory_pools->pools_defined, "Pool index is invalid (%u/%u)!\n",
      pool, memory_pools->pools_defined);
  item_size      = memory_pools->pools[pool].item_data_number;
  pool_item_size = memory_pools->pools[pool].pool_item_size;
  item_index =
      (((void*) memory_pool_item) - ((void*) memory_pools->pools[pool].items)) /
      pool_item_size;
  MP_DEBUG(
      " Free  [%2u][%6d]{%6d}, %3u %3u,         %p, %p, %p, %u\n", pool,
      item_index,
      items_group_free_items(&memory_pools->pools[pool].items_group_free),
      memory_pool_item->start.info[0], info_1, memory_pool_item_handle,
      memory_pool_item, memory_pools->pools[pool].items,
      ((uint32_t)(item_size * sizeof(memory_pool_data_t))));
  /*
   * Sanity check on calculated item index
   */
  AssertFatal(
      memory_pool_item ==
          memory_pool_item_from_index(&memory_pools->pools[pool], item_index),
      "Incorrect memory pool item address (%p, %p) for pool %u, item %d!\n",
      memory_pool_item,
      memory_pool_item_from_index(&memory_pools->pools[pool], item_index), pool,
      item_index);
  /*
   * Sanity check on end marker, must still be present (no write overflow)
   */
  AssertFatal(
      memory_pool_item->data[item_size] == POOL_ITEM_END_MARK,
      "Memory pool item is corrupted, end mark is not present for pool %u, "
      "item "
      "%d!\n",
      pool, item_index);
  /*
   * Sanity check on item status, must be allocated
   */
  AssertFatal(
      memory_pool_item->start.item_status == ITEM_STATUS_ALLOCATED,
      "Trying to free a non allocated (%x) memory pool item (pool %u, item "
      "%d)!\n",
      memory_pool_item->start.item_status, pool, item_index);
  memory_pool_item->start.item_status = ITEM_STATUS_FREE;
  result                              = items_group_put_free_item(
      &memory_pools->pools[pool].items_group_free, item_index);
  AssertError(
      result == EXIT_SUCCESS, {},
      "Failed to free memory pool item (pool %u, item %d)!\n", pool,
      item_index);

  return result;
}

//------------------------------------------------------------------------------
void memory_pools_set_info(
    memory_pools_handle_t memory_pools_handle,
    memory_pool_item_handle_t memory_pool_item_handle, int index,
    uint16_t info) {
  memory_pools_t* memory_pools;
  memory_pool_item_t* memory_pool_item;
  pool_id_t pool;
  items_group_index_t item_index;
  uint32_t item_size;
  uint32_t pool_item_size;

  AssertFatal(
      index < MEMORY_POOL_ITEM_INFO_NUMBER, "Incorrect info index (%d/%d)!\n",
      index, MEMORY_POOL_ITEM_INFO_NUMBER);
  /*
   * Recover memory pool item
   */
  memory_pool_item = memory_pool_item_from_handler(memory_pool_item_handle);
  AssertFatal(
      memory_pool_item != NULL,
      "Failed to retrieve memory pool item for handle %p!\n",
      memory_pool_item_handle);
  /*
   * Set info[1]
   */
  memory_pool_item->start.info[index] = info;

  /*
   * Check item validity and log (not mandatory)
   */
  if (1) {
    /*
     * Recover memory_pools
     */
    memory_pools = memory_pools_from_handler(memory_pools_handle);
    AssertFatal(
        memory_pools != NULL, "Failed to retrieve memory pool for handle %p!\n",
        memory_pools_handle);
    /*
     * Recover pool index
     */
    pool = memory_pool_item->start.pool_id;
    AssertFatal(
        pool < memory_pools->pools_defined, "Pool index is invalid (%u/%u)!\n",
        pool, memory_pools->pools_defined);
    item_size      = memory_pools->pools[pool].item_data_number;
    pool_item_size = memory_pools->pools[pool].pool_item_size;
    item_index     = (((void*) memory_pool_item) -
                  ((void*) memory_pools->pools[pool].items)) /
                 pool_item_size;
    MP_DEBUG(
        " Info  [%2u][%6d]{%6d}, %3u %3u,         %p, %p, %p, %u\n", pool,
        item_index,
        items_group_free_items(&memory_pools->pools[pool].items_group_free),
        memory_pool_item->start.info[0], memory_pool_item->start.info[1],
        memory_pool_item_handle, memory_pool_item,
        memory_pools->pools[pool].items,
        ((uint32_t)(item_size * sizeof(memory_pool_data_t))));
    /*
     * Sanity check on calculated item index
     */
    AssertFatal(
        memory_pool_item ==
            memory_pool_item_from_index(&memory_pools->pools[pool], item_index),
        "Incorrect memory pool item address (%p, %p) for pool %u, item %d!\n",
        memory_pool_item,
        memory_pool_item_from_index(&memory_pools->pools[pool], item_index),
        pool, item_index);
    /*
     * Sanity check on end marker, must still be present (no write overflow)
     */
    AssertFatal(
        memory_pool_item->data[item_size] == POOL_ITEM_END_MARK,
        "Memory pool item is corrupted, end mark is not present for pool %u, "
        "item %d!\n",
        pool, item_index);
    /*
     * Sanity check on item status, must be allocated
     */
    AssertFatal(
        memory_pool_item->start.item_status == ITEM_STATUS_ALLOCATED,
        "Trying to free a non allocated (%x) memory pool item (pool %u, item "
        "%d)\n",
        memory_pool_item->start.item_status, pool, item_index);
  }
}
