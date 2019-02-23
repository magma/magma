/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
 */

INSERT INTO mmeidentity (`idmmeidentity`,`mmehost`,`mmerealm`,`UE-reachability`)
SELECT * FROM (SELECT '7','magma-oai.openair4G.eur','openair4G.eur','0') AS tmp
WHERE NOT EXISTS (
  SELECT * FROM mmeidentity
  WHERE mmehost = 'magma-oai.openair4G.eur'
) LIMIT 1;
