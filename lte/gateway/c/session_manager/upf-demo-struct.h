enum upfNodeType {
  IPv4 = 0,
  IPv6 = 1,
  FQDN = 2,
};
typedef struct tNodeId {
  upfNodeType node_id_type;
  char node_id[40];
} NodeId;

typedef struct Fseid {
  uint64_t f_seid;
  NodeId Nid;
} FSid;
