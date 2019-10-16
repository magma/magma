/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate go run generate_main.go -o generated.go /usr/share/freeradius/dictionary.rfc2865 /usr/share/freeradius/dictionary.rfc2866 /usr/share/freeradius/dictionary.rfc2867 /usr/share/freeradius/dictionary.rfc2869 /usr/share/freeradius/dictionary.rfc3162 /usr/share/freeradius/dictionary.rfc3576 /usr/share/freeradius/dictionary.rfc5176

package debug
