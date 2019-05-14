---
id: faq
title: FAQ
hide_title: true
---
# FAQ

1. I'm seeing `/usr/local/go/pkg/tool/linux_amd64/link: signal: killed`. How do I fix this?

   - If you are running on Mac, the memory limit of the docker daemon is 2GB by default.
   You may need to increase the memory limit of the docker daemon to about 4GB 
   to build the images.
    