---
id: faq
title: FAQ
hide_title: true
---
# FAQ

1. **Issue:**  How can I restore a link if I get this error message at the prompt:
`  `/usr/local/go/pkg/tool/linux_amd64/link: signal: killed. ``

  **Resolution:** On Mac, the default memory limit of the Docker daemon is 2Gb.
  Increase the memory limit of the Docker daemon to 4Gb to build the images.
