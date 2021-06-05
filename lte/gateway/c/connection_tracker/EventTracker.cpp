/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <iostream>

#include <libmnl/libmnl.h>
#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_conntrack.h>

#include <linux/if_packet.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/socket.h>
#include <net/if.h>
#include <netinet/ether.h>
#include <linux/ip.h>
#include <memory>

#include "EventTracker.h"

#include "magma_logging.h"

static int data_cb(const struct nlmsghdr* nlh, void* data);

namespace magma {
namespace lte {

EventTracker::EventTracker(std::shared_ptr<PacketGenerator> pkt_gen, int zone)
    : pkt_gen_(pkt_gen), zone_(zone) {}

int EventTracker::init_conntrack_event_loop() {
  struct mnl_socket* nl;
  char buf[MNL_SOCKET_BUFFER_SIZE];
  int ret;

  nl = mnl_socket_open(NETLINK_NETFILTER);
  if (nl == NULL) {
    perror("mnl_socket_open");
    exit(EXIT_FAILURE);
  }

  if (mnl_socket_bind(
          nl,
          NF_NETLINK_CONNTRACK_NEW |
              // NF_NETLINK_CONNTRACK_UPDATE |
              NF_NETLINK_CONNTRACK_DESTROY,
          MNL_SOCKET_AUTOPID) < 0) {
    perror("mnl_socket_bind");
    exit(EXIT_FAILURE);
  }

  while (1) {
    ret = mnl_socket_recvfrom(nl, buf, sizeof(buf));
    if (ret == -1) {
      perror("mnl_socket_recvfrom");
      exit(EXIT_FAILURE);
    }
    ret = mnl_cb_run(buf, ret, 0, 0, data_cb, (void*) this);
    if (ret == -1) {
      perror("mnl_cb_run");
      exit(EXIT_FAILURE);
    }
  }

  mnl_socket_close(nl);

  return 0;
}
}  // namespace lte
}  // namespace magma

static int parse_ip_cb(const struct nlattr* attr, void* data) {
  const struct nlattr** tb = (const struct nlattr**) data;
  int type                 = mnl_attr_get_type(attr);

  if (mnl_attr_type_valid(attr, CTA_IP_MAX) < 0) return MNL_CB_OK;

  switch (type) {
    case CTA_IP_V4_SRC:
    case CTA_IP_V4_DST:
      if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
  }
  tb[type] = attr;
  return MNL_CB_OK;
}

static void parse_ip(const struct nlattr* nest, struct flow_information* flow) {
  struct nlattr* tb[CTA_IP_MAX + 1] = {};

  mnl_attr_parse_nested(nest, parse_ip_cb, tb);
  if (tb[CTA_IP_V4_SRC]) {
    struct in_addr* in =
        (struct in_addr*) mnl_attr_get_payload(tb[CTA_IP_V4_SRC]);
    flow->saddr = in->s_addr;
  }
  if (tb[CTA_IP_V4_DST]) {
    struct in_addr* in =
        (struct in_addr*) mnl_attr_get_payload(tb[CTA_IP_V4_DST]);
    flow->daddr = in->s_addr;
  }
}

static int parse_proto_cb(const struct nlattr* attr, void* data) {
  const struct nlattr** tb = (const struct nlattr**) data;
  int type                 = mnl_attr_get_type(attr);

  if (mnl_attr_type_valid(attr, CTA_PROTO_MAX) < 0) return MNL_CB_OK;

  switch (type) {
    case CTA_PROTO_NUM:
    case CTA_PROTO_ICMP_TYPE:
    case CTA_PROTO_ICMP_CODE:
      if (mnl_attr_validate(attr, MNL_TYPE_U8) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
    case CTA_PROTO_SRC_PORT:
    case CTA_PROTO_DST_PORT:
    case CTA_PROTO_ICMP_ID:
      if (mnl_attr_validate(attr, MNL_TYPE_U16) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
  }
  tb[type] = attr;
  return MNL_CB_OK;
}

static void parse_proto(
    const struct nlattr* nest, struct flow_information* flow) {
  struct nlattr* tb[CTA_PROTO_MAX + 1] = {};

  mnl_attr_parse_nested(nest, parse_proto_cb, tb);
  if (tb[CTA_PROTO_NUM]) {
    flow->l4_proto = mnl_attr_get_u8(tb[CTA_PROTO_NUM]);
  }
  if (tb[CTA_PROTO_SRC_PORT]) {
    flow->sport = mnl_attr_get_u8(tb[CTA_PROTO_NUM]);
  }
  if (tb[CTA_PROTO_DST_PORT]) {
    flow->dport = mnl_attr_get_u8(tb[CTA_PROTO_NUM]);
  }
}

static int parse_tuple_cb(const struct nlattr* attr, void* data) {
  const struct nlattr** tb = (const struct nlattr**) data;
  int type                 = mnl_attr_get_type(attr);

  if (mnl_attr_type_valid(attr, CTA_TUPLE_MAX) < 0) return MNL_CB_OK;

  switch (type) {
    case CTA_TUPLE_IP:
      if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
    case CTA_TUPLE_PROTO:
      if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
  }
  tb[type] = attr;
  return MNL_CB_OK;
}

static void print_tuple(
    const struct nlattr* nest, struct flow_information* flow) {
  struct nlattr* tb[CTA_TUPLE_MAX + 1] = {};

  mnl_attr_parse_nested(nest, parse_tuple_cb, tb);
  if (tb[CTA_TUPLE_IP]) {
    parse_ip(tb[CTA_TUPLE_IP], flow);
  }
  if (tb[CTA_TUPLE_PROTO]) {
    parse_proto(tb[CTA_TUPLE_PROTO], flow);
  }
}

static int data_attr_cb(const struct nlattr* attr, void* data) {
  const struct nlattr** tb = (const struct nlattr**) data;
  int type                 = mnl_attr_get_type(attr);

  if (mnl_attr_type_valid(attr, CTA_MAX) < 0) return MNL_CB_OK;

  switch (type) {
    case CTA_TUPLE_ORIG:
      if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
    case CTA_TIMEOUT:
    case CTA_MARK:
    case CTA_SECMARK:
      if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0) {
        perror("mnl_attr_validate");
        return MNL_CB_ERROR;
      }
      break;
  }
  tb[type] = attr;
  return MNL_CB_OK;
}

static int data_cb(const struct nlmsghdr* nlh, void* data) {
  struct nlattr* tb[CTA_MAX + 1] = {};
  struct nfgenmsg* nfg = (struct nfgenmsg*) mnl_nlmsg_get_payload(nlh);
  struct flow_information flow;
  struct in_addr src_ip;
  struct in_addr dst_ip;

  if ((nlh->nlmsg_type & 0xFF) != IPCTNL_MSG_CT_DELETE) {
    return 0;
  }

  mnl_attr_parse(nlh, sizeof(*nfg), data_attr_cb, tb);

  // If zone isn't set the event isn't coming from the OVS commited flow
  if (!tb[CTA_ZONE] || ((magma::lte::EventTracker*) data)->zone_ !=
                           ntohs(mnl_attr_get_u16(tb[CTA_ZONE]))) {
    return 0;
  }

  if (tb[CTA_TUPLE_ORIG]) {
    print_tuple(tb[CTA_TUPLE_ORIG], &flow);
  }

  src_ip.s_addr = flow.saddr;
  dst_ip.s_addr = flow.daddr;

  switch (nlh->nlmsg_type & 0xFF) {
    case IPCTNL_MSG_CT_NEW:
      if (nlh->nlmsg_flags & (NLM_F_CREATE | NLM_F_EXCL))
        MLOG(MINFO) << "     [NEW] src=" << inet_ntoa(src_ip) << ":"
                    << ntohs(flow.sport) << " dst=" << inet_ntoa(dst_ip) << ":"
                    << ntohs(flow.dport) << " proto=" << flow.l4_proto;
      else
        printf("%9s ", "[UPDATE] \n");
      break;
    case IPCTNL_MSG_CT_DELETE:
      MLOG(MINFO) << "[DESTROY] src=" << inet_ntoa(src_ip) << ":"
                  << ntohs(flow.sport) << " dst=" << inet_ntoa(dst_ip) << ":"
                  << ntohs(flow.dport) << " proto=" << flow.l4_proto;
      break;
  }

  if (tb[CTA_MARK]) {
    MLOG(MINFO) << "From zone " << mnl_attr_get_u16(tb[CTA_ZONE]);
  }

  ((magma::lte::EventTracker*) data)->pkt_gen_->send_packet(&flow);

  return MNL_CB_OK;
}
