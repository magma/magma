#!/usr/bin/env python3
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
from pyinventory.api.location_type import add_location_type
from pyinventory.exceptions import LocationCannotBeDeletedWithDependency
from pysymphony import SymphonyClient

from ..utils.base_test import BaseTest
from ..utils.grpc.rpc_pb2_grpc import TenantServiceStub


class TestLocation(BaseTest):
    def __init__(
        self, testName: str, client: SymphonyClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
        add_location_type(
            client=self.client,
            name="City",
            properties=[
                ("Mayor", "string", None, True),
                ("Contact", "email", None, True),
            ],
        )
        self.external_id = "test_external_id"
        self.suffixes = ["txt", "pdf", "png"]
        self.tmpdir = tempfile.mkdtemp()
        self.location_1 = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima1")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
        )
        self.location_2 = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima2")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
        )
        self.location_with_ext_id = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima3")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
            externalID=self.external_id,
        )
        self.location_child_1 = add_location(
            client=self.client,
            location_hirerchy=[("City", "parent"), ("City", "child1")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )
        self.location_child_2 = add_location(
            client=self.client,
            location_hirerchy=[("City", "parent"), ("City", "child2")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )

    def tearDown(self) -> None:
        shutil.rmtree(self.tmpdir)

    def test_location_created(self) -> None:
        fetch_location = get_location(
            client=self.client, location_hirerchy=[("City", "Lima1")]
        )
        self.assertEqual(self.location_1, fetch_location)

    def test_location_created_already_exists(self) -> None:
        fetched_location = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima1")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
        )
        self.assertEqual(self.location_1, fetched_location)

    def test_location_hierarchy_created(self) -> None:
        fetched_parent = get_location(
            client=self.client, location_hirerchy=[("City", "parent")]
        )
        self.assertEqual(fetched_parent.name, "parent")
        fetched_child_location1 = get_location(
            client=self.client, location_hirerchy=[("City", "child1")]
        )
        self.assertEqual(self.location_child_1, fetched_child_location1)

        fetcheed_child_location2 = get_location(
            client=self.client, location_hirerchy=[("City", "child2")]
        )
        self.assertEqual(self.location_child_2, fetcheed_child_location2)

    def test_location_created_hierarchy_already_exists(self) -> None:
        fetched_location = add_location(
            client=self.client,
            location_hirerchy=[("City", "parent"), ("City", "child1")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
        )
        self.assertEqual(self.location_child_1, fetched_location)

    def test_location_with_external_id_created(self) -> None:
        fetch_locations = get_locations_by_external_id(
            client=self.client, external_id=self.external_id
        )
        self.assertEqual(len(fetch_locations), 1)
        self.assertEqual(self.location_with_ext_id, fetch_locations[0])

    def test_location_edited(self) -> None:
        edit_location(
            client=self.client,
            location=self.location_1,
            new_name="Lima4",
            new_lat=10,
            new_long=20,
            new_external_id=None,
            new_properties={"Contact": "new_limacity@peru.pe"},
        )
        edited_location = get_location(
            client=self.client, location_hirerchy=[("City", "Lima4")]
        )
        # TODO(T63055774): update test to check updated properties
        self.assertEqual(self.location_1.id, edited_location.id)

    def test_location_add_file(self) -> None:
        temp_file_path = os.path.join(self.tmpdir, ".".join(["temp_file", "txt"]))
        with open(temp_file_path, "wb") as tmp_file:
            tmp_file.write(b"TEST DATA FILE")

        add_file(
            client=self.client,
            local_file_path=temp_file_path,
            entity_type="LOCATION",
            entity_id=self.location_1.id,
        )

        docs = get_location_documents(client=self.client, location=self.location_1)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            delete_document(self.client, doc)

    def test_location_add_file_with_category(self) -> None:
        temp_file_path = os.path.join(self.tmpdir, ".".join(["temp_file", "txt"]))
        with open(temp_file_path, "wb") as tmp_file:
            tmp_file.write(b"TEST DATA FILE")
        add_file(
            client=self.client,
            local_file_path=temp_file_path,
            entity_type="LOCATION",
            entity_id=self.location_1.id,
            category="test_category",
        )
        docs = get_location_documents(client=self.client, location=self.location_1)
        for doc in docs:
            self.assertEqual(doc.category, "test_category")
            delete_document(self.client, doc)

    def test_location_upload_folder(self) -> None:
        fetch_location = get_location(
            client=self.client, location_hirerchy=[("City", "Lima1")]
        )
        self.assertEqual(self.location_1, fetch_location)
        for suffix in self.suffixes:
            with open(
                os.path.join(self.tmpdir, ".".join(["temp_file", suffix])), "wb"
            ) as tmp:
                tmp.write(b"TEST DATA FILE")
        add_files(
            client=self.client,
            local_directory_path=self.tmpdir,
            entity_type="LOCATION",
            entity_id=self.location_1.id,
            category="test_category",
        )

        docs = get_location_documents(client=self.client, location=self.location_1)
        self.assertEqual(len(docs), len(self.suffixes))
        for doc in docs:
            self.assertEqual(doc.category, "test_category")
            delete_document(client=self.client, document=doc)

    def test_location_moved(self) -> None:
        move_location(
            client=self.client,
            location_id=self.location_2.id,
            new_parent_id=self.location_1.id,
        )
        moved_location = get_location(
            client=self.client, location_hirerchy=[("City", "Lima1"), ("City", "Lima2")]
        )
        self.assertEqual(self.location_2.id, moved_location.id)

    def test_get_location_children(self) -> None:
        created_locations_arr = {self.location_child_1.id, self.location_child_2.id}
        parent_location = get_location(
            client=self.client, location_hirerchy=[("City", "parent")]
        )
        fetch_locations = get_location_children(
            client=self.client, location_id=parent_location.id
        )
        fetch_locations_arr = {location.id for location in fetch_locations}
        self.assertEqual(created_locations_arr, fetch_locations_arr)

    def test_delete_location_documents(self) -> None:
        with tempfile.NamedTemporaryFile() as fp:
            fp.write(b"DATA")
            add_location_image(
                client=self.client, local_file_path=fp.name, location=self.location_1
            )
        with self.assertRaises(LocationCannotBeDeletedWithDependency):
            delete_location(client=self.client, location=self.location_1)
        docs = get_location_documents(client=self.client, location=self.location_1)
        self.assertEqual(len(docs), 1)
        for doc in docs:
            delete_document(client=self.client, document=doc)
