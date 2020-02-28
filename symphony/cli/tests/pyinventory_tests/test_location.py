#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import os
import shutil
import tempfile

from pyinventory.api.file import (
    add_file,
    add_files,
    add_location_image,
    delete_document,
)
from pyinventory.api.location import (
    add_location,
    delete_location,
    edit_location,
    get_location,
    get_location_children,
    get_location_documents,
    get_locations_by_external_id,
    move_location,
)
from pyinventory.api.location_type import (
    add_location_type,
    delete_location_type_with_locations,
)
from pyinventory.exceptions import LocationCannotBeDeletedWithDependency

from .utils.base_test import BaseTest


class TestLocation(BaseTest):
    def setUp(self) -> None:
        super().setUp()
        self.location_types_created = []
        self.location_types_created.append(
            add_location_type(
                self.client,
                "City",
                [("Mayor", "string", None, True), ("Contact", "email", None, True)],
            )
        )
        self.external_id = "test_external_id"
        self.suffixes = ["txt", "pdf", "png"]
        self.tmpdir = tempfile.mkdtemp()
        self.location_1 = add_location(
            self.client,
            [("City", "Lima1")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        self.location_2 = add_location(
            self.client,
            [("City", "Lima2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        self.location_with_ext_id = add_location(
            self.client,
            [("City", "Lima3")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
            self.external_id,
        )
        self.location_child_1 = add_location(
            self.client,
            [("City", "parent"), ("City", "child1")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )
        self.location_child_2 = add_location(
            self.client,
            [("City", "parent"), ("City", "child2")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )

    def tearDown(self) -> None:
        shutil.rmtree(self.tmpdir)
        delete_location(self.client, self.location_2)
        delete_location(self.client, self.location_1)
        delete_location(self.client, self.location_with_ext_id)
        delete_location(self.client, self.location_child_1)
        delete_location(self.client, self.location_child_2)
        for location_type in self.location_types_created:
            delete_location_type_with_locations(self.client, location_type)

    def test_location_created(self) -> None:
        fetch_location = get_location(self.client, [("City", "Lima1")])
        self.assertEqual(self.location_1, fetch_location)

    def test_location_with_external_id_created(self) -> None:
        fetch_locations = get_locations_by_external_id(self.client, self.external_id)
        self.assertEqual(len(fetch_locations), 1)
        self.assertEqual(self.location_with_ext_id, fetch_locations[0])

    def test_location_edited(self) -> None:
        edit_location(
            self.client,
            self.location_1,
            "Lima4",
            10,
            20,
            None,
            {"Contact": "new_limacity@peru.pe"},
        )
        edited_location = get_location(self.client, [("City", "Lima4")])
        # TODO(T63055774): update test to check updated properties
        self.assertEqual(self.location_1.id, edited_location.id)

    def test_location_add_file(self) -> None:
        temp_file_path = os.path.join(self.tmpdir, ".".join(["temp_file", "txt"]))
        with open(temp_file_path, "wb") as tmp_file:
            tmp_file.write(b"TEST DATA FILE")

        add_file(self.client, temp_file_path, "LOCATION", self.location_1.id)

        docs = get_location_documents(self.client, self.location_1)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            delete_document(self.client, doc)

    def test_location_add_file_with_category(self) -> None:
        temp_file_path = os.path.join(self.tmpdir, ".".join(["temp_file", "txt"]))
        with open(temp_file_path, "wb") as tmp_file:
            tmp_file.write(b"TEST DATA FILE")
        add_file(
            self.client, temp_file_path, "LOCATION", self.location_1.id, "test_category"
        )
        docs = get_location_documents(self.client, self.location_1)
        for doc in docs:
            self.assertEqual(doc.category, "test_category")
            delete_document(self.client, doc)

    def test_location_upload_folder(self) -> None:
        fetch_location = get_location(self.client, [("City", "Lima1")])
        self.assertEqual(self.location_1, fetch_location)
        for suffix in self.suffixes:
            with open(
                os.path.join(self.tmpdir, ".".join(["temp_file", suffix])), "wb"
            ) as tmp:
                tmp.write(b"TEST DATA FILE")
        add_files(
            self.client, self.tmpdir, "LOCATION", self.location_1.id, "test_category"
        )

        docs = get_location_documents(self.client, self.location_1)
        self.assertEqual(len(docs), len(self.suffixes))
        for doc in docs:
            self.assertEqual(doc.category, "test_category")
            delete_document(self.client, doc)

    def test_location_moved(self) -> None:
        move_location(self.client, self.location_2.id, self.location_1.id)
        moved_location = get_location(
            self.client, [("City", "Lima1"), ("City", "Lima2")]
        )
        self.assertEqual(self.location_2.id, moved_location.id)

    def test_get_location_children(self) -> None:
        created_locations_arr = {self.location_child_1.id, self.location_child_2.id}
        parent_location = get_location(self.client, [("City", "parent")])
        fetch_locations = get_location_children(self.client, parent_location.id)
        fetch_locations_arr = {location.id for location in fetch_locations}
        self.assertEqual(created_locations_arr, fetch_locations_arr)

    def test_delete_location_documents(self) -> None:
        with tempfile.NamedTemporaryFile() as fp:
            fp.write(b"DATA")
            add_location_image(self.client, fp.name, self.location_1)
        with self.assertRaises(LocationCannotBeDeletedWithDependency):
            delete_location(self.client, self.location_1)
        docs = get_location_documents(self.client, self.location_1)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            delete_document(self.client, doc)
