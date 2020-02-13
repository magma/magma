#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import os
import shutil
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
        self.external_id = "test_external_id"
        self.suffixes = ["txt", "pdf", "png"]
        self.tmpdir = tempfile.mkdtemp()
        self.location_1 = self.client.add_location(
            [("City", "Lima1")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        self.location_2 = self.client.add_location(
            [("City", "Lima2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        self.location_with_ext_id = self.client.add_location(
            [("City", "Lima3")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
            self.external_id,
        )
        self.location_child_1 = self.client.add_location(
            [("City", "parent"), ("City", "child1")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )
        self.location_child_2 = self.client.add_location(
            [("City", "parent"), ("City", "child2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )

    def tearDown(self):
        shutil.rmtree(self.tmpdir)
        self.client.delete_location(self.location_2)
        self.client.delete_location(self.location_1)
        self.client.delete_location(self.location_with_ext_id)
        self.client.delete_location(self.location_child_1)
        self.client.delete_location(self.location_child_2)
        for location_type in self.location_types_created:
            self.client.delete_location_type_with_locations(location_type)

    def test_location_created(self):
        fetch_location = self.client.get_location([("City", "Lima1")])
        self.assertEqual(self.location_1, fetch_location)

    def test_location_with_external_id_created(self):
        fetch_locations = self.client.get_locations_by_external_id(self.external_id)
        self.assertEqual(len(fetch_locations), 1)
        self.assertEqual(self.location_with_ext_id, fetch_locations[0])

    def test_location_edited(self):
        self.client.edit_location(
            self.location_1, "Lima4", 10, 20, None, {"Contact": "new_limacity@peru.pe"}
        )
        edited_location = self.client.get_location([("City", "Lima4")])
        # TODO update test to check updated properties
        self.assertEqual(self.location_1.id, edited_location.id)

    def test_location_add_file(self):
        temp_file_path = os.path.join(self.tmpdir, ".".join(["temp_file", "txt"]))
        with open(temp_file_path, "wb") as tmp_file:
            tmp_file.write(b"TEST DATA FILE")

        self.client.add_file(temp_file_path, "LOCATION", self.location_1.id)

        docs = self.client.get_location_documents(self.location_1)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            self.client.delete_document(doc)

    def test_location_add_file_with_category(self):
        temp_file_path = os.path.join(self.tmpdir, ".".join(["temp_file", "txt"]))
        with open(temp_file_path, "wb") as tmp_file:
            tmp_file.write(b"TEST DATA FILE")
        self.client.add_file(
            temp_file_path, "LOCATION", self.location_1.id, "test_category"
        )
        docs = self.client.get_location_documents(self.location_1)
        for doc in docs:
            self.assertEqual(doc.category, "test_category")
            self.client.delete_document(doc)

    def test_location_upload_folder(self):
        fetch_location = self.client.get_location([("City", "Lima1")])
        self.assertEqual(self.location_1, fetch_location)
        for suffix in self.suffixes:
            with open(
                os.path.join(self.tmpdir, ".".join(["temp_file", suffix])), "wb"
            ) as tmp:
                tmp.write(b"TEST DATA FILE")
        self.client.add_files(
            self.tmpdir, "LOCATION", self.location_1.id, "test_category"
        )

        docs = self.client.get_location_documents(self.location_1)
        self.assertEqual(len(docs), len(self.suffixes))
        for doc in docs:
            self.assertEqual(doc.category, "test_category")
            self.client.delete_document(doc)

    def test_location_moved(self):
        self.client.move_location(self.location_2.id, self.location_1.id)
        moved_location = self.client.get_location(
            [("City", "Lima1"), ("City", "Lima2")]
        )
        self.assertEqual(self.location_2.id, moved_location.id)

    def test_get_location_children(self):
        created_locations_arr = {self.location_child_1.id, self.location_child_2.id}
        parent_location = self.client.get_location([("City", "parent")])
        fetch_locations = self.client.get_location_children(parent_location.id)
        fetch_locations_arr = {location.id for location in fetch_locations}
        self.assertEqual(created_locations_arr, fetch_locations_arr)

    def test_delete_location_documents(self):
        with tempfile.NamedTemporaryFile() as fp:
            fp.write(b"DATA")
            self.client.add_location_image(fp.name, self.location_1)
        with self.assertRaises(LocationCannotBeDeletedWithDependency):
            self.client.delete_location(self.location_1)
        docs = self.client.get_location_documents(self.location_1)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            self.client.delete_document(doc)
