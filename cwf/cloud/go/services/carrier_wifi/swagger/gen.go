/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate cp $SWAGGER_ROOT/$SWAGGER_COMMON $SWAGGER_COMMON
//go:generate swagger generate model -f swagger.yml -t ../obsidian/ -C $SWAGGER_TEMPLATE
//go:generate rm ./$SWAGGER_COMMON

package swagger
