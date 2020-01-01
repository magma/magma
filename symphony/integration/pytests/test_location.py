#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import tempfile

from pyinventory.exceptions import LocationCannotBeDeletedWithDependency
from utils.base_test import BaseTest


class TestLocation(BaseTest):
    def setUp(self):
        super().setUp()
        self.location_types_created = []
        self.location_types_created.append(
            self.client.add_location_type(
                "City",
                [("Mayor", "string", None, True), ("Contact", "email", None, True)],
            )
        )

    def tearDown(self):
        for location_type in self.location_types_created:
            self.client.delete_location_type_with_locations(location_type)

    def test_location_created(self):
        location = self.client.add_location(
            [("City", "Lima")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        fetch_location = self.client.get_location([("City", "Lima")])
        self.assertEqual(location, fetch_location)

    def test_location_with_external_id_created(self):
        external_id = "test_external_id"
        location = self.client.add_location(
            [("City", "Lima2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
            external_id,
        )
        fetch_locations = self.client.get_locations_by_external_id(external_id)
        self.assertEquals(len(fetch_locations), 1)
        self.assertEqual(location, fetch_locations[0])

    def test_location_edited(self):
        created_location = self.client.add_location(
            [("City", "Lima")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )

        self.client.edit_location(created_location, "Lima2", 10, 20, None)

        edited_location = self.client.get_location([("City", "Lima2")])
        self.assertEqual(created_location.id, edited_location.id)

    def test_location_moved(self):
        created_location_1 = self.client.add_location(
            [("City", "Lima1")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        created_location_2 = self.client.add_location(
            [("City", "Lima2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )

        self.client.move_location(created_location_2.id, created_location_1.id)

        moved_location = self.client.get_location(
            [("City", "Lima1"), ("City", "Lima2")]
        )
        self.assertEqual(created_location_2.id, moved_location.id)

        # cleanup, otherwise teardown won't delete the locations in the right order,
        # and will fail
        self.client.delete_location(moved_location)

    def test_get_location_children(self):
        location_1 = self.client.add_location(
            [("City", "parent"), ("City", "child1")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )
        location_2 = self.client.add_location(
            [("City", "parent"), ("City", "child2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )
        created_locations_arr = [location_1.id, location_2.id]
        parent_location = self.client.get_location([("City", "parent")])
        fetch_locations = self.client.get_location_children(parent_location.id)
        fetch_locations_arr = [location.id for location in fetch_locations]

        if created_locations_arr[0] == fetch_locations_arr[0]:
            self.assertEqual(created_locations_arr[0], created_locations_arr[0])
            self.assertEqual(created_locations_arr[1], created_locations_arr[1])
        else:
            self.assertEqual(created_locations_arr[0], created_locations_arr[1])
            self.assertEqual(created_locations_arr[1], created_locations_arr[0])

        # cleanup, otherwise teardown won't delete the locations in the right order,
        # and will fail
        self.client.delete_location(location_1)
        self.client.delete_location(location_2)

    def test_delete_location_documents(self):
        location = self.client.add_location(
            [("City", "Lima")], {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"}
        )
        with tempfile.NamedTemporaryFile() as fp:
            fp.write(b"DATA")
            self.client.add_location_image(fp.name, location)
        with self.assertRaises(LocationCannotBeDeletedWithDependency):
            self.client.delete_location(location)
        docs = self.client.get_location_documents(location)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            self.client.delete_document(doc)
        self.client.delete_location(location)
