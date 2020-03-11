// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

// ExportEntity is all entities that can be exportable
type ExportEntity string

const (
	ExportEntityEquipment  ExportEntity = "EQUIPMENT"
	ExportEntityService    ExportEntity = "SERVICE"
	ExportEntityLink       ExportEntity = "LINK"
	ExportEntityPort       ExportEntity = "PORT"
	ExportEntityLocation   ExportEntity = "LOCATION"
	ExportEntityWorkOrders ExportEntity = "WORK_ORDERS"
)
