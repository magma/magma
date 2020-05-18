// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/ent/schema/field"
)

var (
	// ActionsRulesColumns holds the columns for the "actions_rules" table.
	ActionsRulesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "trigger_id", Type: field.TypeString},
		{Name: "rule_filters", Type: field.TypeJSON},
		{Name: "rule_actions", Type: field.TypeJSON},
	}
	// ActionsRulesTable holds the schema information for the "actions_rules" table.
	ActionsRulesTable = &schema.Table{
		Name:        "actions_rules",
		Columns:     ActionsRulesColumns,
		PrimaryKey:  []*schema.Column{ActionsRulesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// CheckListCategoriesColumns holds the columns for the "check_list_categories" table.
	CheckListCategoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "work_order_check_list_categories", Type: field.TypeInt, Nullable: true},
	}
	// CheckListCategoriesTable holds the schema information for the "check_list_categories" table.
	CheckListCategoriesTable = &schema.Table{
		Name:       "check_list_categories",
		Columns:    CheckListCategoriesColumns,
		PrimaryKey: []*schema.Column{CheckListCategoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "check_list_categories_work_orders_check_list_categories",
				Columns: []*schema.Column{CheckListCategoriesColumns[5]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// CheckListCategoryDefinitionsColumns holds the columns for the "check_list_category_definitions" table.
	CheckListCategoryDefinitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "work_order_type_check_list_category_definitions", Type: field.TypeInt, Nullable: true},
	}
	// CheckListCategoryDefinitionsTable holds the schema information for the "check_list_category_definitions" table.
	CheckListCategoryDefinitionsTable = &schema.Table{
		Name:       "check_list_category_definitions",
		Columns:    CheckListCategoryDefinitionsColumns,
		PrimaryKey: []*schema.Column{CheckListCategoryDefinitionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "check_list_category_definitions_work_order_types_check_list_category_definitions",
				Columns: []*schema.Column{CheckListCategoryDefinitionsColumns[5]},

				RefColumns: []*schema.Column{WorkOrderTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// CheckListItemsColumns holds the columns for the "check_list_items" table.
	CheckListItemsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "checked", Type: field.TypeBool, Nullable: true},
		{Name: "string_val", Type: field.TypeString, Nullable: true},
		{Name: "enum_values", Type: field.TypeString, Nullable: true},
		{Name: "enum_selection_mode_value", Type: field.TypeEnum, Nullable: true, Enums: []string{"single", "multiple"}},
		{Name: "selected_enum_values", Type: field.TypeString, Nullable: true},
		{Name: "yes_no_val", Type: field.TypeEnum, Nullable: true, Enums: []string{"YES", "NO"}},
		{Name: "help_text", Type: field.TypeString, Nullable: true},
		{Name: "check_list_category_check_list_items", Type: field.TypeInt, Nullable: true},
	}
	// CheckListItemsTable holds the schema information for the "check_list_items" table.
	CheckListItemsTable = &schema.Table{
		Name:       "check_list_items",
		Columns:    CheckListItemsColumns,
		PrimaryKey: []*schema.Column{CheckListItemsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "check_list_items_check_list_categories_check_list_items",
				Columns: []*schema.Column{CheckListItemsColumns[11]},

				RefColumns: []*schema.Column{CheckListCategoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// CheckListItemDefinitionsColumns holds the columns for the "check_list_item_definitions" table.
	CheckListItemDefinitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "title", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "enum_values", Type: field.TypeString, Nullable: true},
		{Name: "enum_selection_mode_value", Type: field.TypeEnum, Nullable: true, Enums: []string{"single", "multiple"}},
		{Name: "help_text", Type: field.TypeString, Nullable: true},
		{Name: "check_list_category_definition_check_list_item_definitions", Type: field.TypeInt, Nullable: true},
	}
	// CheckListItemDefinitionsTable holds the schema information for the "check_list_item_definitions" table.
	CheckListItemDefinitionsTable = &schema.Table{
		Name:       "check_list_item_definitions",
		Columns:    CheckListItemDefinitionsColumns,
		PrimaryKey: []*schema.Column{CheckListItemDefinitionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "check_list_item_definitions_check_list_category_definitions_check_list_item_definitions",
				Columns: []*schema.Column{CheckListItemDefinitionsColumns[9]},

				RefColumns: []*schema.Column{CheckListCategoryDefinitionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// CommentsColumns holds the columns for the "comments" table.
	CommentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "text", Type: field.TypeString},
		{Name: "comment_author", Type: field.TypeInt, Nullable: true},
		{Name: "project_comments", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_comments", Type: field.TypeInt, Nullable: true},
	}
	// CommentsTable holds the schema information for the "comments" table.
	CommentsTable = &schema.Table{
		Name:       "comments",
		Columns:    CommentsColumns,
		PrimaryKey: []*schema.Column{CommentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "comments_users_author",
				Columns: []*schema.Column{CommentsColumns[4]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "comments_projects_comments",
				Columns: []*schema.Column{CommentsColumns[5]},

				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "comments_work_orders_comments",
				Columns: []*schema.Column{CommentsColumns[6]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// CustomersColumns holds the columns for the "customers" table.
	CustomersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "external_id", Type: field.TypeString, Unique: true, Nullable: true},
	}
	// CustomersTable holds the schema information for the "customers" table.
	CustomersTable = &schema.Table{
		Name:        "customers",
		Columns:     CustomersColumns,
		PrimaryKey:  []*schema.Column{CustomersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// EquipmentColumns holds the columns for the "equipment" table.
	EquipmentColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "future_state", Type: field.TypeString, Nullable: true},
		{Name: "device_id", Type: field.TypeString, Nullable: true},
		{Name: "external_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "equipment_type", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_work_order", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_position_attachment", Type: field.TypeInt, Unique: true, Nullable: true},
		{Name: "location_equipment", Type: field.TypeInt, Nullable: true},
	}
	// EquipmentTable holds the schema information for the "equipment" table.
	EquipmentTable = &schema.Table{
		Name:       "equipment",
		Columns:    EquipmentColumns,
		PrimaryKey: []*schema.Column{EquipmentColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "equipment_equipment_types_type",
				Columns: []*schema.Column{EquipmentColumns[7]},

				RefColumns: []*schema.Column{EquipmentTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_work_orders_work_order",
				Columns: []*schema.Column{EquipmentColumns[8]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_equipment_positions_attachment",
				Columns: []*schema.Column{EquipmentColumns[9]},

				RefColumns: []*schema.Column{EquipmentPositionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_locations_equipment",
				Columns: []*schema.Column{EquipmentColumns[10]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// EquipmentCategoriesColumns holds the columns for the "equipment_categories" table.
	EquipmentCategoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// EquipmentCategoriesTable holds the schema information for the "equipment_categories" table.
	EquipmentCategoriesTable = &schema.Table{
		Name:        "equipment_categories",
		Columns:     EquipmentCategoriesColumns,
		PrimaryKey:  []*schema.Column{EquipmentCategoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// EquipmentPortsColumns holds the columns for the "equipment_ports" table.
	EquipmentPortsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "equipment_ports", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_port_definition", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_port_link", Type: field.TypeInt, Nullable: true},
	}
	// EquipmentPortsTable holds the schema information for the "equipment_ports" table.
	EquipmentPortsTable = &schema.Table{
		Name:       "equipment_ports",
		Columns:    EquipmentPortsColumns,
		PrimaryKey: []*schema.Column{EquipmentPortsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "equipment_ports_equipment_ports",
				Columns: []*schema.Column{EquipmentPortsColumns[3]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_ports_equipment_port_definitions_definition",
				Columns: []*schema.Column{EquipmentPortsColumns[4]},

				RefColumns: []*schema.Column{EquipmentPortDefinitionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_ports_links_link",
				Columns: []*schema.Column{EquipmentPortsColumns[5]},

				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "equipmentport_equipment_port_definition_equipment_ports",
				Unique:  true,
				Columns: []*schema.Column{EquipmentPortsColumns[4], EquipmentPortsColumns[3]},
			},
		},
	}
	// EquipmentPortDefinitionsColumns holds the columns for the "equipment_port_definitions" table.
	EquipmentPortDefinitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "bandwidth", Type: field.TypeString, Nullable: true},
		{Name: "visibility_label", Type: field.TypeString, Nullable: true},
		{Name: "equipment_port_definition_equipment_port_type", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_type_port_definitions", Type: field.TypeInt, Nullable: true},
	}
	// EquipmentPortDefinitionsTable holds the schema information for the "equipment_port_definitions" table.
	EquipmentPortDefinitionsTable = &schema.Table{
		Name:       "equipment_port_definitions",
		Columns:    EquipmentPortDefinitionsColumns,
		PrimaryKey: []*schema.Column{EquipmentPortDefinitionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "equipment_port_definitions_equipment_port_types_equipment_port_type",
				Columns: []*schema.Column{EquipmentPortDefinitionsColumns[7]},

				RefColumns: []*schema.Column{EquipmentPortTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_port_definitions_equipment_types_port_definitions",
				Columns: []*schema.Column{EquipmentPortDefinitionsColumns[8]},

				RefColumns: []*schema.Column{EquipmentTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// EquipmentPortTypesColumns holds the columns for the "equipment_port_types" table.
	EquipmentPortTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// EquipmentPortTypesTable holds the schema information for the "equipment_port_types" table.
	EquipmentPortTypesTable = &schema.Table{
		Name:        "equipment_port_types",
		Columns:     EquipmentPortTypesColumns,
		PrimaryKey:  []*schema.Column{EquipmentPortTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// EquipmentPositionsColumns holds the columns for the "equipment_positions" table.
	EquipmentPositionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "equipment_positions", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_position_definition", Type: field.TypeInt, Nullable: true},
	}
	// EquipmentPositionsTable holds the schema information for the "equipment_positions" table.
	EquipmentPositionsTable = &schema.Table{
		Name:       "equipment_positions",
		Columns:    EquipmentPositionsColumns,
		PrimaryKey: []*schema.Column{EquipmentPositionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "equipment_positions_equipment_positions",
				Columns: []*schema.Column{EquipmentPositionsColumns[3]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "equipment_positions_equipment_position_definitions_definition",
				Columns: []*schema.Column{EquipmentPositionsColumns[4]},

				RefColumns: []*schema.Column{EquipmentPositionDefinitionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "equipmentposition_equipment_position_definition_equipment_positions",
				Unique:  true,
				Columns: []*schema.Column{EquipmentPositionsColumns[4], EquipmentPositionsColumns[3]},
			},
		},
	}
	// EquipmentPositionDefinitionsColumns holds the columns for the "equipment_position_definitions" table.
	EquipmentPositionDefinitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "visibility_label", Type: field.TypeString, Nullable: true},
		{Name: "equipment_type_position_definitions", Type: field.TypeInt, Nullable: true},
	}
	// EquipmentPositionDefinitionsTable holds the schema information for the "equipment_position_definitions" table.
	EquipmentPositionDefinitionsTable = &schema.Table{
		Name:       "equipment_position_definitions",
		Columns:    EquipmentPositionDefinitionsColumns,
		PrimaryKey: []*schema.Column{EquipmentPositionDefinitionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "equipment_position_definitions_equipment_types_position_definitions",
				Columns: []*schema.Column{EquipmentPositionDefinitionsColumns[6]},

				RefColumns: []*schema.Column{EquipmentTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// EquipmentTypesColumns holds the columns for the "equipment_types" table.
	EquipmentTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "equipment_type_category", Type: field.TypeInt, Nullable: true},
	}
	// EquipmentTypesTable holds the schema information for the "equipment_types" table.
	EquipmentTypesTable = &schema.Table{
		Name:       "equipment_types",
		Columns:    EquipmentTypesColumns,
		PrimaryKey: []*schema.Column{EquipmentTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "equipment_types_equipment_categories_category",
				Columns: []*schema.Column{EquipmentTypesColumns[4]},

				RefColumns: []*schema.Column{EquipmentCategoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// FilesColumns holds the columns for the "files" table.
	FilesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "type", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "size", Type: field.TypeInt, Nullable: true},
		{Name: "modified_at", Type: field.TypeTime, Nullable: true},
		{Name: "uploaded_at", Type: field.TypeTime, Nullable: true},
		{Name: "content_type", Type: field.TypeString},
		{Name: "store_key", Type: field.TypeString},
		{Name: "category", Type: field.TypeString, Nullable: true},
		{Name: "check_list_item_files", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_files", Type: field.TypeInt, Nullable: true},
		{Name: "location_files", Type: field.TypeInt, Nullable: true},
		{Name: "survey_question_photo_data", Type: field.TypeInt, Nullable: true},
		{Name: "survey_question_images", Type: field.TypeInt, Nullable: true},
		{Name: "user_profile_photo", Type: field.TypeInt, Unique: true, Nullable: true},
		{Name: "work_order_files", Type: field.TypeInt, Nullable: true},
	}
	// FilesTable holds the schema information for the "files" table.
	FilesTable = &schema.Table{
		Name:       "files",
		Columns:    FilesColumns,
		PrimaryKey: []*schema.Column{FilesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "files_check_list_items_files",
				Columns: []*schema.Column{FilesColumns[11]},

				RefColumns: []*schema.Column{CheckListItemsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "files_equipment_files",
				Columns: []*schema.Column{FilesColumns[12]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "files_locations_files",
				Columns: []*schema.Column{FilesColumns[13]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "files_survey_questions_photo_data",
				Columns: []*schema.Column{FilesColumns[14]},

				RefColumns: []*schema.Column{SurveyQuestionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "files_survey_questions_images",
				Columns: []*schema.Column{FilesColumns[15]},

				RefColumns: []*schema.Column{SurveyQuestionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "files_users_profile_photo",
				Columns: []*schema.Column{FilesColumns[16]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "files_work_orders_files",
				Columns: []*schema.Column{FilesColumns[17]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// FloorPlansColumns holds the columns for the "floor_plans" table.
	FloorPlansColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "floor_plan_location", Type: field.TypeInt, Nullable: true},
		{Name: "floor_plan_reference_point", Type: field.TypeInt, Nullable: true},
		{Name: "floor_plan_scale", Type: field.TypeInt, Nullable: true},
		{Name: "floor_plan_image", Type: field.TypeInt, Nullable: true},
	}
	// FloorPlansTable holds the schema information for the "floor_plans" table.
	FloorPlansTable = &schema.Table{
		Name:       "floor_plans",
		Columns:    FloorPlansColumns,
		PrimaryKey: []*schema.Column{FloorPlansColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "floor_plans_locations_location",
				Columns: []*schema.Column{FloorPlansColumns[4]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "floor_plans_floor_plan_reference_points_reference_point",
				Columns: []*schema.Column{FloorPlansColumns[5]},

				RefColumns: []*schema.Column{FloorPlanReferencePointsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "floor_plans_floor_plan_scales_scale",
				Columns: []*schema.Column{FloorPlansColumns[6]},

				RefColumns: []*schema.Column{FloorPlanScalesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "floor_plans_files_image",
				Columns: []*schema.Column{FloorPlansColumns[7]},

				RefColumns: []*schema.Column{FilesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// FloorPlanReferencePointsColumns holds the columns for the "floor_plan_reference_points" table.
	FloorPlanReferencePointsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "x", Type: field.TypeInt},
		{Name: "y", Type: field.TypeInt},
		{Name: "latitude", Type: field.TypeFloat64},
		{Name: "longitude", Type: field.TypeFloat64},
	}
	// FloorPlanReferencePointsTable holds the schema information for the "floor_plan_reference_points" table.
	FloorPlanReferencePointsTable = &schema.Table{
		Name:        "floor_plan_reference_points",
		Columns:     FloorPlanReferencePointsColumns,
		PrimaryKey:  []*schema.Column{FloorPlanReferencePointsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// FloorPlanScalesColumns holds the columns for the "floor_plan_scales" table.
	FloorPlanScalesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "reference_point1_x", Type: field.TypeInt},
		{Name: "reference_point1_y", Type: field.TypeInt},
		{Name: "reference_point2_x", Type: field.TypeInt},
		{Name: "reference_point2_y", Type: field.TypeInt},
		{Name: "scale_in_meters", Type: field.TypeFloat64},
	}
	// FloorPlanScalesTable holds the schema information for the "floor_plan_scales" table.
	FloorPlanScalesTable = &schema.Table{
		Name:        "floor_plan_scales",
		Columns:     FloorPlanScalesColumns,
		PrimaryKey:  []*schema.Column{FloorPlanScalesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// HyperlinksColumns holds the columns for the "hyperlinks" table.
	HyperlinksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "url", Type: field.TypeString},
		{Name: "name", Type: field.TypeString, Nullable: true},
		{Name: "category", Type: field.TypeString, Nullable: true},
		{Name: "equipment_hyperlinks", Type: field.TypeInt, Nullable: true},
		{Name: "location_hyperlinks", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_hyperlinks", Type: field.TypeInt, Nullable: true},
	}
	// HyperlinksTable holds the schema information for the "hyperlinks" table.
	HyperlinksTable = &schema.Table{
		Name:       "hyperlinks",
		Columns:    HyperlinksColumns,
		PrimaryKey: []*schema.Column{HyperlinksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "hyperlinks_equipment_hyperlinks",
				Columns: []*schema.Column{HyperlinksColumns[6]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "hyperlinks_locations_hyperlinks",
				Columns: []*schema.Column{HyperlinksColumns[7]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "hyperlinks_work_orders_hyperlinks",
				Columns: []*schema.Column{HyperlinksColumns[8]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// LinksColumns holds the columns for the "links" table.
	LinksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "future_state", Type: field.TypeString, Nullable: true},
		{Name: "link_work_order", Type: field.TypeInt, Nullable: true},
	}
	// LinksTable holds the schema information for the "links" table.
	LinksTable = &schema.Table{
		Name:       "links",
		Columns:    LinksColumns,
		PrimaryKey: []*schema.Column{LinksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "links_work_orders_work_order",
				Columns: []*schema.Column{LinksColumns[4]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// LocationsColumns holds the columns for the "locations" table.
	LocationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "external_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "latitude", Type: field.TypeFloat64},
		{Name: "longitude", Type: field.TypeFloat64},
		{Name: "site_survey_needed", Type: field.TypeBool, Nullable: true},
		{Name: "location_type", Type: field.TypeInt, Nullable: true},
		{Name: "location_children", Type: field.TypeInt, Nullable: true},
	}
	// LocationsTable holds the schema information for the "locations" table.
	LocationsTable = &schema.Table{
		Name:       "locations",
		Columns:    LocationsColumns,
		PrimaryKey: []*schema.Column{LocationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "locations_location_types_type",
				Columns: []*schema.Column{LocationsColumns[8]},

				RefColumns: []*schema.Column{LocationTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "locations_locations_children",
				Columns: []*schema.Column{LocationsColumns[9]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "location_name_location_type_location_children",
				Unique:  true,
				Columns: []*schema.Column{LocationsColumns[3], LocationsColumns[8], LocationsColumns[9]},
			},
		},
	}
	// LocationTypesColumns holds the columns for the "location_types" table.
	LocationTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "site", Type: field.TypeBool},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "map_type", Type: field.TypeString, Nullable: true},
		{Name: "map_zoom_level", Type: field.TypeInt, Nullable: true, Default: 7},
		{Name: "index", Type: field.TypeInt},
	}
	// LocationTypesTable holds the schema information for the "location_types" table.
	LocationTypesTable = &schema.Table{
		Name:        "location_types",
		Columns:     LocationTypesColumns,
		PrimaryKey:  []*schema.Column{LocationTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// PermissionsPoliciesColumns holds the columns for the "permissions_policies" table.
	PermissionsPoliciesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "is_global", Type: field.TypeBool, Nullable: true},
		{Name: "inventory_policy", Type: field.TypeJSON, Nullable: true},
		{Name: "workforce_policy", Type: field.TypeJSON, Nullable: true},
	}
	// PermissionsPoliciesTable holds the schema information for the "permissions_policies" table.
	PermissionsPoliciesTable = &schema.Table{
		Name:        "permissions_policies",
		Columns:     PermissionsPoliciesColumns,
		PrimaryKey:  []*schema.Column{PermissionsPoliciesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// ProjectsColumns holds the columns for the "projects" table.
	ProjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "project_location", Type: field.TypeInt, Nullable: true},
		{Name: "project_creator", Type: field.TypeInt, Nullable: true},
		{Name: "project_type_projects", Type: field.TypeInt, Nullable: true},
	}
	// ProjectsTable holds the schema information for the "projects" table.
	ProjectsTable = &schema.Table{
		Name:       "projects",
		Columns:    ProjectsColumns,
		PrimaryKey: []*schema.Column{ProjectsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "projects_locations_location",
				Columns: []*schema.Column{ProjectsColumns[5]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "projects_users_creator",
				Columns: []*schema.Column{ProjectsColumns[6]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "projects_project_types_projects",
				Columns: []*schema.Column{ProjectsColumns[7]},

				RefColumns: []*schema.Column{ProjectTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "project_name_project_type_projects",
				Unique:  true,
				Columns: []*schema.Column{ProjectsColumns[3], ProjectsColumns[7]},
			},
		},
	}
	// ProjectTypesColumns holds the columns for the "project_types" table.
	ProjectTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
	}
	// ProjectTypesTable holds the schema information for the "project_types" table.
	ProjectTypesTable = &schema.Table{
		Name:        "project_types",
		Columns:     ProjectTypesColumns,
		PrimaryKey:  []*schema.Column{ProjectTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// PropertiesColumns holds the columns for the "properties" table.
	PropertiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "int_val", Type: field.TypeInt, Nullable: true},
		{Name: "bool_val", Type: field.TypeBool, Nullable: true},
		{Name: "float_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "latitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "range_from_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "range_to_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "string_val", Type: field.TypeString, Nullable: true},
		{Name: "equipment_properties", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_port_properties", Type: field.TypeInt, Nullable: true},
		{Name: "link_properties", Type: field.TypeInt, Nullable: true},
		{Name: "location_properties", Type: field.TypeInt, Nullable: true},
		{Name: "project_properties", Type: field.TypeInt, Nullable: true},
		{Name: "property_type", Type: field.TypeInt, Nullable: true},
		{Name: "property_equipment_value", Type: field.TypeInt, Nullable: true},
		{Name: "property_location_value", Type: field.TypeInt, Nullable: true},
		{Name: "property_service_value", Type: field.TypeInt, Nullable: true},
		{Name: "property_work_order_value", Type: field.TypeInt, Nullable: true},
		{Name: "property_user_value", Type: field.TypeInt, Nullable: true},
		{Name: "service_properties", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_properties", Type: field.TypeInt, Nullable: true},
	}
	// PropertiesTable holds the schema information for the "properties" table.
	PropertiesTable = &schema.Table{
		Name:       "properties",
		Columns:    PropertiesColumns,
		PrimaryKey: []*schema.Column{PropertiesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "properties_equipment_properties",
				Columns: []*schema.Column{PropertiesColumns[11]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_equipment_ports_properties",
				Columns: []*schema.Column{PropertiesColumns[12]},

				RefColumns: []*schema.Column{EquipmentPortsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_links_properties",
				Columns: []*schema.Column{PropertiesColumns[13]},

				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_locations_properties",
				Columns: []*schema.Column{PropertiesColumns[14]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_projects_properties",
				Columns: []*schema.Column{PropertiesColumns[15]},

				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_property_types_type",
				Columns: []*schema.Column{PropertiesColumns[16]},

				RefColumns: []*schema.Column{PropertyTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_equipment_equipment_value",
				Columns: []*schema.Column{PropertiesColumns[17]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_locations_location_value",
				Columns: []*schema.Column{PropertiesColumns[18]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_services_service_value",
				Columns: []*schema.Column{PropertiesColumns[19]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_work_orders_work_order_value",
				Columns: []*schema.Column{PropertiesColumns[20]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_users_user_value",
				Columns: []*schema.Column{PropertiesColumns[21]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_services_properties",
				Columns: []*schema.Column{PropertiesColumns[22]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "properties_work_orders_properties",
				Columns: []*schema.Column{PropertiesColumns[23]},

				RefColumns: []*schema.Column{WorkOrdersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PropertyTypesColumns holds the columns for the "property_types" table.
	PropertyTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "type", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "external_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "category", Type: field.TypeString, Nullable: true},
		{Name: "int_val", Type: field.TypeInt, Nullable: true},
		{Name: "bool_val", Type: field.TypeBool, Nullable: true},
		{Name: "float_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "latitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "string_val", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "range_from_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "range_to_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "is_instance_property", Type: field.TypeBool, Default: true},
		{Name: "editable", Type: field.TypeBool, Default: true},
		{Name: "mandatory", Type: field.TypeBool},
		{Name: "deleted", Type: field.TypeBool},
		{Name: "node_type", Type: field.TypeString, Nullable: true},
		{Name: "equipment_port_type_property_types", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_port_type_link_property_types", Type: field.TypeInt, Nullable: true},
		{Name: "equipment_type_property_types", Type: field.TypeInt, Nullable: true},
		{Name: "location_type_property_types", Type: field.TypeInt, Nullable: true},
		{Name: "project_type_properties", Type: field.TypeInt, Nullable: true},
		{Name: "service_type_property_types", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_type_property_types", Type: field.TypeInt, Nullable: true},
	}
	// PropertyTypesTable holds the schema information for the "property_types" table.
	PropertyTypesTable = &schema.Table{
		Name:       "property_types",
		Columns:    PropertyTypesColumns,
		PrimaryKey: []*schema.Column{PropertyTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "property_types_equipment_port_types_property_types",
				Columns: []*schema.Column{PropertyTypesColumns[21]},

				RefColumns: []*schema.Column{EquipmentPortTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "property_types_equipment_port_types_link_property_types",
				Columns: []*schema.Column{PropertyTypesColumns[22]},

				RefColumns: []*schema.Column{EquipmentPortTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "property_types_equipment_types_property_types",
				Columns: []*schema.Column{PropertyTypesColumns[23]},

				RefColumns: []*schema.Column{EquipmentTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "property_types_location_types_property_types",
				Columns: []*schema.Column{PropertyTypesColumns[24]},

				RefColumns: []*schema.Column{LocationTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "property_types_project_types_properties",
				Columns: []*schema.Column{PropertyTypesColumns[25]},

				RefColumns: []*schema.Column{ProjectTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "property_types_service_types_property_types",
				Columns: []*schema.Column{PropertyTypesColumns[26]},

				RefColumns: []*schema.Column{ServiceTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "property_types_work_order_types_property_types",
				Columns: []*schema.Column{PropertyTypesColumns[27]},

				RefColumns: []*schema.Column{WorkOrderTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "propertytype_name_location_type_property_types",
				Unique:  true,
				Columns: []*schema.Column{PropertyTypesColumns[4], PropertyTypesColumns[24]},
			},
			{
				Name:    "propertytype_name_equipment_port_type_property_types",
				Unique:  true,
				Columns: []*schema.Column{PropertyTypesColumns[4], PropertyTypesColumns[21]},
			},
			{
				Name:    "propertytype_name_equipment_type_property_types",
				Unique:  true,
				Columns: []*schema.Column{PropertyTypesColumns[4], PropertyTypesColumns[23]},
			},
			{
				Name:    "propertytype_name_equipment_port_type_link_property_types",
				Unique:  true,
				Columns: []*schema.Column{PropertyTypesColumns[4], PropertyTypesColumns[22]},
			},
			{
				Name:    "propertytype_name_work_order_type_property_types",
				Unique:  true,
				Columns: []*schema.Column{PropertyTypesColumns[4], PropertyTypesColumns[27]},
			},
		},
	}
	// ReportFiltersColumns holds the columns for the "report_filters" table.
	ReportFiltersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "entity", Type: field.TypeEnum, Enums: []string{"WORK_ORDER", "PORT", "EQUIPMENT", "LINK", "LOCATION", "SERVICE"}},
		{Name: "filters", Type: field.TypeString, Size: 2147483647, Default: "[]"},
	}
	// ReportFiltersTable holds the schema information for the "report_filters" table.
	ReportFiltersTable = &schema.Table{
		Name:        "report_filters",
		Columns:     ReportFiltersColumns,
		PrimaryKey:  []*schema.Column{ReportFiltersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
		Indexes: []*schema.Index{
			{
				Name:    "reportfilter_name_entity",
				Unique:  true,
				Columns: []*schema.Column{ReportFiltersColumns[3], ReportFiltersColumns[4]},
			},
		},
	}
	// ServicesColumns holds the columns for the "services" table.
	ServicesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "external_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "status", Type: field.TypeString},
		{Name: "service_type", Type: field.TypeInt, Nullable: true},
	}
	// ServicesTable holds the schema information for the "services" table.
	ServicesTable = &schema.Table{
		Name:       "services",
		Columns:    ServicesColumns,
		PrimaryKey: []*schema.Column{ServicesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "services_service_types_type",
				Columns: []*schema.Column{ServicesColumns[6]},

				RefColumns: []*schema.Column{ServiceTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ServiceEndpointsColumns holds the columns for the "service_endpoints" table.
	ServiceEndpointsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "service_endpoints", Type: field.TypeInt, Nullable: true},
		{Name: "service_endpoint_port", Type: field.TypeInt, Nullable: true},
		{Name: "service_endpoint_equipment", Type: field.TypeInt, Nullable: true},
		{Name: "service_endpoint_definition_endpoints", Type: field.TypeInt, Nullable: true},
	}
	// ServiceEndpointsTable holds the schema information for the "service_endpoints" table.
	ServiceEndpointsTable = &schema.Table{
		Name:       "service_endpoints",
		Columns:    ServiceEndpointsColumns,
		PrimaryKey: []*schema.Column{ServiceEndpointsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "service_endpoints_services_endpoints",
				Columns: []*schema.Column{ServiceEndpointsColumns[3]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "service_endpoints_equipment_ports_port",
				Columns: []*schema.Column{ServiceEndpointsColumns[4]},

				RefColumns: []*schema.Column{EquipmentPortsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "service_endpoints_equipment_equipment",
				Columns: []*schema.Column{ServiceEndpointsColumns[5]},

				RefColumns: []*schema.Column{EquipmentColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "service_endpoints_service_endpoint_definitions_endpoints",
				Columns: []*schema.Column{ServiceEndpointsColumns[6]},

				RefColumns: []*schema.Column{ServiceEndpointDefinitionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ServiceEndpointDefinitionsColumns holds the columns for the "service_endpoint_definitions" table.
	ServiceEndpointDefinitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "role", Type: field.TypeString, Nullable: true},
		{Name: "name", Type: field.TypeString},
		{Name: "index", Type: field.TypeInt},
		{Name: "equipment_type_service_endpoint_definitions", Type: field.TypeInt, Nullable: true},
		{Name: "service_type_endpoint_definitions", Type: field.TypeInt, Nullable: true},
	}
	// ServiceEndpointDefinitionsTable holds the schema information for the "service_endpoint_definitions" table.
	ServiceEndpointDefinitionsTable = &schema.Table{
		Name:       "service_endpoint_definitions",
		Columns:    ServiceEndpointDefinitionsColumns,
		PrimaryKey: []*schema.Column{ServiceEndpointDefinitionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "service_endpoint_definitions_equipment_types_service_endpoint_definitions",
				Columns: []*schema.Column{ServiceEndpointDefinitionsColumns[6]},

				RefColumns: []*schema.Column{EquipmentTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "service_endpoint_definitions_service_types_endpoint_definitions",
				Columns: []*schema.Column{ServiceEndpointDefinitionsColumns[7]},

				RefColumns: []*schema.Column{ServiceTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "serviceendpointdefinition_index_service_type_endpoint_definitions",
				Unique:  true,
				Columns: []*schema.Column{ServiceEndpointDefinitionsColumns[5], ServiceEndpointDefinitionsColumns[7]},
			},
			{
				Name:    "serviceendpointdefinition_name_service_type_endpoint_definitions",
				Unique:  true,
				Columns: []*schema.Column{ServiceEndpointDefinitionsColumns[4], ServiceEndpointDefinitionsColumns[7]},
			},
		},
	}
	// ServiceTypesColumns holds the columns for the "service_types" table.
	ServiceTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "has_customer", Type: field.TypeBool},
		{Name: "is_deleted", Type: field.TypeBool},
		{Name: "discovery_method", Type: field.TypeEnum, Nullable: true, Enums: []string{"INVENTORY"}},
	}
	// ServiceTypesTable holds the schema information for the "service_types" table.
	ServiceTypesTable = &schema.Table{
		Name:        "service_types",
		Columns:     ServiceTypesColumns,
		PrimaryKey:  []*schema.Column{ServiceTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// SurveysColumns holds the columns for the "surveys" table.
	SurveysColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "owner_name", Type: field.TypeString, Nullable: true},
		{Name: "creation_timestamp", Type: field.TypeTime, Nullable: true},
		{Name: "completion_timestamp", Type: field.TypeTime},
		{Name: "survey_location", Type: field.TypeInt, Nullable: true},
		{Name: "survey_source_file", Type: field.TypeInt, Nullable: true},
	}
	// SurveysTable holds the schema information for the "surveys" table.
	SurveysTable = &schema.Table{
		Name:       "surveys",
		Columns:    SurveysColumns,
		PrimaryKey: []*schema.Column{SurveysColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "surveys_locations_location",
				Columns: []*schema.Column{SurveysColumns[7]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "surveys_files_source_file",
				Columns: []*schema.Column{SurveysColumns[8]},

				RefColumns: []*schema.Column{FilesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SurveyCellScansColumns holds the columns for the "survey_cell_scans" table.
	SurveyCellScansColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "network_type", Type: field.TypeString},
		{Name: "signal_strength", Type: field.TypeInt},
		{Name: "timestamp", Type: field.TypeTime, Nullable: true},
		{Name: "base_station_id", Type: field.TypeString, Nullable: true},
		{Name: "network_id", Type: field.TypeString, Nullable: true},
		{Name: "system_id", Type: field.TypeString, Nullable: true},
		{Name: "cell_id", Type: field.TypeString, Nullable: true},
		{Name: "location_area_code", Type: field.TypeString, Nullable: true},
		{Name: "mobile_country_code", Type: field.TypeString, Nullable: true},
		{Name: "mobile_network_code", Type: field.TypeString, Nullable: true},
		{Name: "primary_scrambling_code", Type: field.TypeString, Nullable: true},
		{Name: "operator", Type: field.TypeString, Nullable: true},
		{Name: "arfcn", Type: field.TypeInt, Nullable: true},
		{Name: "physical_cell_id", Type: field.TypeString, Nullable: true},
		{Name: "tracking_area_code", Type: field.TypeString, Nullable: true},
		{Name: "timing_advance", Type: field.TypeInt, Nullable: true},
		{Name: "earfcn", Type: field.TypeInt, Nullable: true},
		{Name: "uarfcn", Type: field.TypeInt, Nullable: true},
		{Name: "latitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "survey_cell_scan_checklist_item", Type: field.TypeInt, Nullable: true},
		{Name: "survey_cell_scan_survey_question", Type: field.TypeInt, Nullable: true},
		{Name: "survey_cell_scan_location", Type: field.TypeInt, Nullable: true},
	}
	// SurveyCellScansTable holds the schema information for the "survey_cell_scans" table.
	SurveyCellScansTable = &schema.Table{
		Name:       "survey_cell_scans",
		Columns:    SurveyCellScansColumns,
		PrimaryKey: []*schema.Column{SurveyCellScansColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "survey_cell_scans_check_list_items_checklist_item",
				Columns: []*schema.Column{SurveyCellScansColumns[23]},

				RefColumns: []*schema.Column{CheckListItemsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "survey_cell_scans_survey_questions_survey_question",
				Columns: []*schema.Column{SurveyCellScansColumns[24]},

				RefColumns: []*schema.Column{SurveyQuestionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "survey_cell_scans_locations_location",
				Columns: []*schema.Column{SurveyCellScansColumns[25]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SurveyQuestionsColumns holds the columns for the "survey_questions" table.
	SurveyQuestionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "form_name", Type: field.TypeString, Nullable: true},
		{Name: "form_description", Type: field.TypeString, Nullable: true},
		{Name: "form_index", Type: field.TypeInt},
		{Name: "question_type", Type: field.TypeString, Nullable: true},
		{Name: "question_format", Type: field.TypeString, Nullable: true},
		{Name: "question_text", Type: field.TypeString, Nullable: true},
		{Name: "question_index", Type: field.TypeInt},
		{Name: "bool_data", Type: field.TypeBool, Nullable: true},
		{Name: "email_data", Type: field.TypeString, Nullable: true},
		{Name: "latitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "location_accuracy", Type: field.TypeFloat64, Nullable: true},
		{Name: "altitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "phone_data", Type: field.TypeString, Nullable: true},
		{Name: "text_data", Type: field.TypeString, Nullable: true},
		{Name: "float_data", Type: field.TypeFloat64, Nullable: true},
		{Name: "int_data", Type: field.TypeInt, Nullable: true},
		{Name: "date_data", Type: field.TypeTime, Nullable: true},
		{Name: "survey_question_survey", Type: field.TypeInt, Nullable: true},
	}
	// SurveyQuestionsTable holds the schema information for the "survey_questions" table.
	SurveyQuestionsTable = &schema.Table{
		Name:       "survey_questions",
		Columns:    SurveyQuestionsColumns,
		PrimaryKey: []*schema.Column{SurveyQuestionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "survey_questions_surveys_survey",
				Columns: []*schema.Column{SurveyQuestionsColumns[21]},

				RefColumns: []*schema.Column{SurveysColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SurveyTemplateCategoriesColumns holds the columns for the "survey_template_categories" table.
	SurveyTemplateCategoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "category_title", Type: field.TypeString},
		{Name: "category_description", Type: field.TypeString},
		{Name: "location_type_survey_template_categories", Type: field.TypeInt, Nullable: true},
	}
	// SurveyTemplateCategoriesTable holds the schema information for the "survey_template_categories" table.
	SurveyTemplateCategoriesTable = &schema.Table{
		Name:       "survey_template_categories",
		Columns:    SurveyTemplateCategoriesColumns,
		PrimaryKey: []*schema.Column{SurveyTemplateCategoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "survey_template_categories_location_types_survey_template_categories",
				Columns: []*schema.Column{SurveyTemplateCategoriesColumns[5]},

				RefColumns: []*schema.Column{LocationTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SurveyTemplateQuestionsColumns holds the columns for the "survey_template_questions" table.
	SurveyTemplateQuestionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "question_title", Type: field.TypeString},
		{Name: "question_description", Type: field.TypeString},
		{Name: "question_type", Type: field.TypeString},
		{Name: "index", Type: field.TypeInt},
		{Name: "survey_template_category_survey_template_questions", Type: field.TypeInt, Nullable: true},
	}
	// SurveyTemplateQuestionsTable holds the schema information for the "survey_template_questions" table.
	SurveyTemplateQuestionsTable = &schema.Table{
		Name:       "survey_template_questions",
		Columns:    SurveyTemplateQuestionsColumns,
		PrimaryKey: []*schema.Column{SurveyTemplateQuestionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "survey_template_questions_survey_template_categories_survey_template_questions",
				Columns: []*schema.Column{SurveyTemplateQuestionsColumns[7]},

				RefColumns: []*schema.Column{SurveyTemplateCategoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "surveytemplatequestion_index_survey_template_category_survey_template_questions",
				Unique:  true,
				Columns: []*schema.Column{SurveyTemplateQuestionsColumns[6], SurveyTemplateQuestionsColumns[7]},
			},
		},
	}
	// SurveyWiFiScansColumns holds the columns for the "survey_wi_fi_scans" table.
	SurveyWiFiScansColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "ssid", Type: field.TypeString, Nullable: true},
		{Name: "bssid", Type: field.TypeString},
		{Name: "timestamp", Type: field.TypeTime},
		{Name: "frequency", Type: field.TypeInt},
		{Name: "channel", Type: field.TypeInt},
		{Name: "band", Type: field.TypeString, Nullable: true},
		{Name: "channel_width", Type: field.TypeInt, Nullable: true},
		{Name: "capabilities", Type: field.TypeString, Nullable: true},
		{Name: "strength", Type: field.TypeInt},
		{Name: "latitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude", Type: field.TypeFloat64, Nullable: true},
		{Name: "survey_wi_fi_scan_checklist_item", Type: field.TypeInt, Nullable: true},
		{Name: "survey_wi_fi_scan_survey_question", Type: field.TypeInt, Nullable: true},
		{Name: "survey_wi_fi_scan_location", Type: field.TypeInt, Nullable: true},
	}
	// SurveyWiFiScansTable holds the schema information for the "survey_wi_fi_scans" table.
	SurveyWiFiScansTable = &schema.Table{
		Name:       "survey_wi_fi_scans",
		Columns:    SurveyWiFiScansColumns,
		PrimaryKey: []*schema.Column{SurveyWiFiScansColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "survey_wi_fi_scans_check_list_items_checklist_item",
				Columns: []*schema.Column{SurveyWiFiScansColumns[14]},

				RefColumns: []*schema.Column{CheckListItemsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "survey_wi_fi_scans_survey_questions_survey_question",
				Columns: []*schema.Column{SurveyWiFiScansColumns[15]},

				RefColumns: []*schema.Column{SurveyQuestionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "survey_wi_fi_scans_locations_location",
				Columns: []*schema.Column{SurveyWiFiScansColumns[16]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "auth_id", Type: field.TypeString, Unique: true},
		{Name: "first_name", Type: field.TypeString, Nullable: true},
		{Name: "last_name", Type: field.TypeString, Nullable: true},
		{Name: "email", Type: field.TypeString, Nullable: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"ACTIVE", "DEACTIVATED"}, Default: "ACTIVE"},
		{Name: "role", Type: field.TypeEnum, Enums: []string{"USER", "ADMIN", "OWNER"}, Default: "USER"},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:        "users",
		Columns:     UsersColumns,
		PrimaryKey:  []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// UsersGroupsColumns holds the columns for the "users_groups" table.
	UsersGroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"ACTIVE", "DEACTIVATED"}, Default: "ACTIVE"},
	}
	// UsersGroupsTable holds the schema information for the "users_groups" table.
	UsersGroupsTable = &schema.Table{
		Name:        "users_groups",
		Columns:     UsersGroupsColumns,
		PrimaryKey:  []*schema.Column{UsersGroupsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// WorkOrdersColumns holds the columns for the "work_orders" table.
	WorkOrdersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "status", Type: field.TypeString, Default: "PLANNED"},
		{Name: "priority", Type: field.TypeString, Default: "NONE"},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "install_date", Type: field.TypeTime, Nullable: true},
		{Name: "creation_date", Type: field.TypeTime},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "close_date", Type: field.TypeTime, Nullable: true},
		{Name: "project_work_orders", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_type", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_location", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_owner", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_assignee", Type: field.TypeInt, Nullable: true},
	}
	// WorkOrdersTable holds the schema information for the "work_orders" table.
	WorkOrdersTable = &schema.Table{
		Name:       "work_orders",
		Columns:    WorkOrdersColumns,
		PrimaryKey: []*schema.Column{WorkOrdersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "work_orders_projects_work_orders",
				Columns: []*schema.Column{WorkOrdersColumns[11]},

				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "work_orders_work_order_types_type",
				Columns: []*schema.Column{WorkOrdersColumns[12]},

				RefColumns: []*schema.Column{WorkOrderTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "work_orders_locations_location",
				Columns: []*schema.Column{WorkOrdersColumns[13]},

				RefColumns: []*schema.Column{LocationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "work_orders_users_owner",
				Columns: []*schema.Column{WorkOrdersColumns[14]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "work_orders_users_assignee",
				Columns: []*schema.Column{WorkOrdersColumns[15]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// WorkOrderDefinitionsColumns holds the columns for the "work_order_definitions" table.
	WorkOrderDefinitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "project_type_work_orders", Type: field.TypeInt, Nullable: true},
		{Name: "work_order_definition_type", Type: field.TypeInt, Nullable: true},
	}
	// WorkOrderDefinitionsTable holds the schema information for the "work_order_definitions" table.
	WorkOrderDefinitionsTable = &schema.Table{
		Name:       "work_order_definitions",
		Columns:    WorkOrderDefinitionsColumns,
		PrimaryKey: []*schema.Column{WorkOrderDefinitionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "work_order_definitions_project_types_work_orders",
				Columns: []*schema.Column{WorkOrderDefinitionsColumns[4]},

				RefColumns: []*schema.Column{ProjectTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:  "work_order_definitions_work_order_types_type",
				Columns: []*schema.Column{WorkOrderDefinitionsColumns[5]},

				RefColumns: []*schema.Column{WorkOrderTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// WorkOrderTypesColumns holds the columns for the "work_order_types" table.
	WorkOrderTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
	}
	// WorkOrderTypesTable holds the schema information for the "work_order_types" table.
	WorkOrderTypesTable = &schema.Table{
		Name:        "work_order_types",
		Columns:     WorkOrderTypesColumns,
		PrimaryKey:  []*schema.Column{WorkOrderTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// ServiceUpstreamColumns holds the columns for the "service_upstream" table.
	ServiceUpstreamColumns = []*schema.Column{
		{Name: "service_id", Type: field.TypeInt},
		{Name: "downstream_id", Type: field.TypeInt},
	}
	// ServiceUpstreamTable holds the schema information for the "service_upstream" table.
	ServiceUpstreamTable = &schema.Table{
		Name:       "service_upstream",
		Columns:    ServiceUpstreamColumns,
		PrimaryKey: []*schema.Column{ServiceUpstreamColumns[0], ServiceUpstreamColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "service_upstream_service_id",
				Columns: []*schema.Column{ServiceUpstreamColumns[0]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:  "service_upstream_downstream_id",
				Columns: []*schema.Column{ServiceUpstreamColumns[1]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// ServiceLinksColumns holds the columns for the "service_links" table.
	ServiceLinksColumns = []*schema.Column{
		{Name: "service_id", Type: field.TypeInt},
		{Name: "link_id", Type: field.TypeInt},
	}
	// ServiceLinksTable holds the schema information for the "service_links" table.
	ServiceLinksTable = &schema.Table{
		Name:       "service_links",
		Columns:    ServiceLinksColumns,
		PrimaryKey: []*schema.Column{ServiceLinksColumns[0], ServiceLinksColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "service_links_service_id",
				Columns: []*schema.Column{ServiceLinksColumns[0]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:  "service_links_link_id",
				Columns: []*schema.Column{ServiceLinksColumns[1]},

				RefColumns: []*schema.Column{LinksColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// ServiceCustomerColumns holds the columns for the "service_customer" table.
	ServiceCustomerColumns = []*schema.Column{
		{Name: "service_id", Type: field.TypeInt},
		{Name: "customer_id", Type: field.TypeInt},
	}
	// ServiceCustomerTable holds the schema information for the "service_customer" table.
	ServiceCustomerTable = &schema.Table{
		Name:       "service_customer",
		Columns:    ServiceCustomerColumns,
		PrimaryKey: []*schema.Column{ServiceCustomerColumns[0], ServiceCustomerColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "service_customer_service_id",
				Columns: []*schema.Column{ServiceCustomerColumns[0]},

				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:  "service_customer_customer_id",
				Columns: []*schema.Column{ServiceCustomerColumns[1]},

				RefColumns: []*schema.Column{CustomersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UsersGroupMembersColumns holds the columns for the "users_group_members" table.
	UsersGroupMembersColumns = []*schema.Column{
		{Name: "users_group_id", Type: field.TypeInt},
		{Name: "user_id", Type: field.TypeInt},
	}
	// UsersGroupMembersTable holds the schema information for the "users_group_members" table.
	UsersGroupMembersTable = &schema.Table{
		Name:       "users_group_members",
		Columns:    UsersGroupMembersColumns,
		PrimaryKey: []*schema.Column{UsersGroupMembersColumns[0], UsersGroupMembersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "users_group_members_users_group_id",
				Columns: []*schema.Column{UsersGroupMembersColumns[0]},

				RefColumns: []*schema.Column{UsersGroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:  "users_group_members_user_id",
				Columns: []*schema.Column{UsersGroupMembersColumns[1]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UsersGroupPoliciesColumns holds the columns for the "users_group_policies" table.
	UsersGroupPoliciesColumns = []*schema.Column{
		{Name: "users_group_id", Type: field.TypeInt},
		{Name: "permissions_policy_id", Type: field.TypeInt},
	}
	// UsersGroupPoliciesTable holds the schema information for the "users_group_policies" table.
	UsersGroupPoliciesTable = &schema.Table{
		Name:       "users_group_policies",
		Columns:    UsersGroupPoliciesColumns,
		PrimaryKey: []*schema.Column{UsersGroupPoliciesColumns[0], UsersGroupPoliciesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "users_group_policies_users_group_id",
				Columns: []*schema.Column{UsersGroupPoliciesColumns[0]},

				RefColumns: []*schema.Column{UsersGroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:  "users_group_policies_permissions_policy_id",
				Columns: []*schema.Column{UsersGroupPoliciesColumns[1]},

				RefColumns: []*schema.Column{PermissionsPoliciesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ActionsRulesTable,
		CheckListCategoriesTable,
		CheckListCategoryDefinitionsTable,
		CheckListItemsTable,
		CheckListItemDefinitionsTable,
		CommentsTable,
		CustomersTable,
		EquipmentTable,
		EquipmentCategoriesTable,
		EquipmentPortsTable,
		EquipmentPortDefinitionsTable,
		EquipmentPortTypesTable,
		EquipmentPositionsTable,
		EquipmentPositionDefinitionsTable,
		EquipmentTypesTable,
		FilesTable,
		FloorPlansTable,
		FloorPlanReferencePointsTable,
		FloorPlanScalesTable,
		HyperlinksTable,
		LinksTable,
		LocationsTable,
		LocationTypesTable,
		PermissionsPoliciesTable,
		ProjectsTable,
		ProjectTypesTable,
		PropertiesTable,
		PropertyTypesTable,
		ReportFiltersTable,
		ServicesTable,
		ServiceEndpointsTable,
		ServiceEndpointDefinitionsTable,
		ServiceTypesTable,
		SurveysTable,
		SurveyCellScansTable,
		SurveyQuestionsTable,
		SurveyTemplateCategoriesTable,
		SurveyTemplateQuestionsTable,
		SurveyWiFiScansTable,
		UsersTable,
		UsersGroupsTable,
		WorkOrdersTable,
		WorkOrderDefinitionsTable,
		WorkOrderTypesTable,
		ServiceUpstreamTable,
		ServiceLinksTable,
		ServiceCustomerTable,
		UsersGroupMembersTable,
		UsersGroupPoliciesTable,
	}
)

func init() {
	CheckListCategoriesTable.ForeignKeys[0].RefTable = WorkOrdersTable
	CheckListCategoryDefinitionsTable.ForeignKeys[0].RefTable = WorkOrderTypesTable
	CheckListItemsTable.ForeignKeys[0].RefTable = CheckListCategoriesTable
	CheckListItemDefinitionsTable.ForeignKeys[0].RefTable = CheckListCategoryDefinitionsTable
	CommentsTable.ForeignKeys[0].RefTable = UsersTable
	CommentsTable.ForeignKeys[1].RefTable = ProjectsTable
	CommentsTable.ForeignKeys[2].RefTable = WorkOrdersTable
	EquipmentTable.ForeignKeys[0].RefTable = EquipmentTypesTable
	EquipmentTable.ForeignKeys[1].RefTable = WorkOrdersTable
	EquipmentTable.ForeignKeys[2].RefTable = EquipmentPositionsTable
	EquipmentTable.ForeignKeys[3].RefTable = LocationsTable
	EquipmentPortsTable.ForeignKeys[0].RefTable = EquipmentTable
	EquipmentPortsTable.ForeignKeys[1].RefTable = EquipmentPortDefinitionsTable
	EquipmentPortsTable.ForeignKeys[2].RefTable = LinksTable
	EquipmentPortDefinitionsTable.ForeignKeys[0].RefTable = EquipmentPortTypesTable
	EquipmentPortDefinitionsTable.ForeignKeys[1].RefTable = EquipmentTypesTable
	EquipmentPositionsTable.ForeignKeys[0].RefTable = EquipmentTable
	EquipmentPositionsTable.ForeignKeys[1].RefTable = EquipmentPositionDefinitionsTable
	EquipmentPositionDefinitionsTable.ForeignKeys[0].RefTable = EquipmentTypesTable
	EquipmentTypesTable.ForeignKeys[0].RefTable = EquipmentCategoriesTable
	FilesTable.ForeignKeys[0].RefTable = CheckListItemsTable
	FilesTable.ForeignKeys[1].RefTable = EquipmentTable
	FilesTable.ForeignKeys[2].RefTable = LocationsTable
	FilesTable.ForeignKeys[3].RefTable = SurveyQuestionsTable
	FilesTable.ForeignKeys[4].RefTable = SurveyQuestionsTable
	FilesTable.ForeignKeys[5].RefTable = UsersTable
	FilesTable.ForeignKeys[6].RefTable = WorkOrdersTable
	FloorPlansTable.ForeignKeys[0].RefTable = LocationsTable
	FloorPlansTable.ForeignKeys[1].RefTable = FloorPlanReferencePointsTable
	FloorPlansTable.ForeignKeys[2].RefTable = FloorPlanScalesTable
	FloorPlansTable.ForeignKeys[3].RefTable = FilesTable
	HyperlinksTable.ForeignKeys[0].RefTable = EquipmentTable
	HyperlinksTable.ForeignKeys[1].RefTable = LocationsTable
	HyperlinksTable.ForeignKeys[2].RefTable = WorkOrdersTable
	LinksTable.ForeignKeys[0].RefTable = WorkOrdersTable
	LocationsTable.ForeignKeys[0].RefTable = LocationTypesTable
	LocationsTable.ForeignKeys[1].RefTable = LocationsTable
	ProjectsTable.ForeignKeys[0].RefTable = LocationsTable
	ProjectsTable.ForeignKeys[1].RefTable = UsersTable
	ProjectsTable.ForeignKeys[2].RefTable = ProjectTypesTable
	PropertiesTable.ForeignKeys[0].RefTable = EquipmentTable
	PropertiesTable.ForeignKeys[1].RefTable = EquipmentPortsTable
	PropertiesTable.ForeignKeys[2].RefTable = LinksTable
	PropertiesTable.ForeignKeys[3].RefTable = LocationsTable
	PropertiesTable.ForeignKeys[4].RefTable = ProjectsTable
	PropertiesTable.ForeignKeys[5].RefTable = PropertyTypesTable
	PropertiesTable.ForeignKeys[6].RefTable = EquipmentTable
	PropertiesTable.ForeignKeys[7].RefTable = LocationsTable
	PropertiesTable.ForeignKeys[8].RefTable = ServicesTable
	PropertiesTable.ForeignKeys[9].RefTable = WorkOrdersTable
	PropertiesTable.ForeignKeys[10].RefTable = UsersTable
	PropertiesTable.ForeignKeys[11].RefTable = ServicesTable
	PropertiesTable.ForeignKeys[12].RefTable = WorkOrdersTable
	PropertyTypesTable.ForeignKeys[0].RefTable = EquipmentPortTypesTable
	PropertyTypesTable.ForeignKeys[1].RefTable = EquipmentPortTypesTable
	PropertyTypesTable.ForeignKeys[2].RefTable = EquipmentTypesTable
	PropertyTypesTable.ForeignKeys[3].RefTable = LocationTypesTable
	PropertyTypesTable.ForeignKeys[4].RefTable = ProjectTypesTable
	PropertyTypesTable.ForeignKeys[5].RefTable = ServiceTypesTable
	PropertyTypesTable.ForeignKeys[6].RefTable = WorkOrderTypesTable
	ServicesTable.ForeignKeys[0].RefTable = ServiceTypesTable
	ServiceEndpointsTable.ForeignKeys[0].RefTable = ServicesTable
	ServiceEndpointsTable.ForeignKeys[1].RefTable = EquipmentPortsTable
	ServiceEndpointsTable.ForeignKeys[2].RefTable = EquipmentTable
	ServiceEndpointsTable.ForeignKeys[3].RefTable = ServiceEndpointDefinitionsTable
	ServiceEndpointDefinitionsTable.ForeignKeys[0].RefTable = EquipmentTypesTable
	ServiceEndpointDefinitionsTable.ForeignKeys[1].RefTable = ServiceTypesTable
	SurveysTable.ForeignKeys[0].RefTable = LocationsTable
	SurveysTable.ForeignKeys[1].RefTable = FilesTable
	SurveyCellScansTable.ForeignKeys[0].RefTable = CheckListItemsTable
	SurveyCellScansTable.ForeignKeys[1].RefTable = SurveyQuestionsTable
	SurveyCellScansTable.ForeignKeys[2].RefTable = LocationsTable
	SurveyQuestionsTable.ForeignKeys[0].RefTable = SurveysTable
	SurveyTemplateCategoriesTable.ForeignKeys[0].RefTable = LocationTypesTable
	SurveyTemplateQuestionsTable.ForeignKeys[0].RefTable = SurveyTemplateCategoriesTable
	SurveyWiFiScansTable.ForeignKeys[0].RefTable = CheckListItemsTable
	SurveyWiFiScansTable.ForeignKeys[1].RefTable = SurveyQuestionsTable
	SurveyWiFiScansTable.ForeignKeys[2].RefTable = LocationsTable
	WorkOrdersTable.ForeignKeys[0].RefTable = ProjectsTable
	WorkOrdersTable.ForeignKeys[1].RefTable = WorkOrderTypesTable
	WorkOrdersTable.ForeignKeys[2].RefTable = LocationsTable
	WorkOrdersTable.ForeignKeys[3].RefTable = UsersTable
	WorkOrdersTable.ForeignKeys[4].RefTable = UsersTable
	WorkOrderDefinitionsTable.ForeignKeys[0].RefTable = ProjectTypesTable
	WorkOrderDefinitionsTable.ForeignKeys[1].RefTable = WorkOrderTypesTable
	ServiceUpstreamTable.ForeignKeys[0].RefTable = ServicesTable
	ServiceUpstreamTable.ForeignKeys[1].RefTable = ServicesTable
	ServiceLinksTable.ForeignKeys[0].RefTable = ServicesTable
	ServiceLinksTable.ForeignKeys[1].RefTable = LinksTable
	ServiceCustomerTable.ForeignKeys[0].RefTable = ServicesTable
	ServiceCustomerTable.ForeignKeys[1].RefTable = CustomersTable
	UsersGroupMembersTable.ForeignKeys[0].RefTable = UsersGroupsTable
	UsersGroupMembersTable.ForeignKeys[1].RefTable = UsersTable
	UsersGroupPoliciesTable.ForeignKeys[0].RefTable = UsersGroupsTable
	UsersGroupPoliciesTable.ForeignKeys[1].RefTable = PermissionsPoliciesTable
}
