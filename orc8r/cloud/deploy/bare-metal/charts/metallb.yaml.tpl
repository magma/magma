configInline:
  address-pools:
  - name: default
    protocol: layer2
    addresses:
      - ${metallb_addresses}
prometheus:
  scrapeAnnotations: true
psp:
  create: false
