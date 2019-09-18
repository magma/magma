/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

//go:generate cp $SWAGGER_ROOT/$SWAGGER_COMMON $SWAGGER_COMMON
//go:generate swagger generate model -f swagger.yml -t ../obsidian/ -C $SWAGGER_TEMPLATE
//go:generate rm ./$SWAGGER_COMMON

package swagger
