#!/bin/bash
echo "Calling-Station-Id = aa:bb:33:44:cc:ff, Called-Station-Id = 11:11:11:11:11:12, User-Name = aa:bb:33:44:cc:ff, NAS-IP-Address = 192.168.1.1, NAS-Identifier = blabla" | radclient -x 127.0.0.1 auth 123456
