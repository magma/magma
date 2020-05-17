#!/usr/bin/env python3


import glob
import os.path
from typing import Generator, Optional

from pysymphony import SymphonyClient
from pysymphony.api.image import add_image, delete_image
from pysymphony.graphql.enum.image_entity import ImageEntity

from ..common.data_class import Document, Location


def list_dir(directory_path: str) -> Generator[str, None, None]:
    files = list(glob.glob(os.path.join(directory_path, "**/**"), recursive=True))
    for file_path in set(files):
        if os.path.isfile(file_path):
            yield file_path


def add_file(
    client: SymphonyClient,
    local_file_path: str,
    entity_type: str,
    entity_id: str,
    category: Optional[str] = None,
) -> None:
    """This function adds file to an entity of a given type.

        Args:
            local_file_path (str): local system path to the file
            entity_type (str): one of existing options ["LOCATION", "WORK_ORDER", "SITE_SURVEY", "EQUIPMENT"]
            entity_id (string): valid entity ID
            category (Optional[string]): file category name

        Raises:
            FailedOperationException: on operation failure

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            client.add_file(
                local_file_path="./document.pdf",
                entity_type="LOCATION",
                entity_id=location.id,
                category="category_name",
            )
            ```
    """
    entity = {
        "LOCATION": ImageEntity.LOCATION,
        "WORK_ORDER": ImageEntity.WORK_ORDER,
        "SITE_SURVEY": ImageEntity.SITE_SURVEY,
        "EQUIPMENT": ImageEntity.EQUIPMENT,
    }.get(entity_type, ImageEntity.LOCATION)
    add_image(client, local_file_path, entity, entity_id, category)


def add_files(
    client: SymphonyClient,
    local_directory_path: str,
    entity_type: str,
    entity_id: str,
    category: Optional[str] = None,
) -> None:
    """This function adds all files located in folder to an entity of a given type.

        Args:
            local_directory_path (str): local system path to the directory
            entity_type (str): one of existing options ["LOCATION", "WORK_ORDER", "SITE_SURVEY", "EQUIPMENT"]
            entity_id (string): valid entity ID
            category (Optional[string]): file category name

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            client.add_files(
                local_directory_path="./documents_folder/",
                entity_type="LOCATION",
                entity_id=location.id,
                category="category_name",
            )
            ```
    """
    for file in list_dir(local_directory_path):
        add_file(client, file, entity_type, entity_id, category)


def add_location_image(
    client: SymphonyClient, local_file_path: str, location: Location
) -> None:
    """This function adds image to existing location.

        Args:
            local_file_path (str): local system path to the file
            location ( `pyinventory.common.data_class.Location` ): existing location object

        Raises:
            FailedOperationException: on operation failure

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            client.add_location_image(
                local_file_path="./document.pdf",
                location=location,
            )
            ```
    """
    add_image(client, local_file_path, ImageEntity.LOCATION, location.id)


def delete_document(client: SymphonyClient, document: Document) -> None:
    """This function deletes existing document.

        Args:
            document ( `pyinventory.common.data_class.Document` ): document object

        Raises:
            FailedOperationException: on operation failure

        Example:
            ```
            client.delete_document(document=document)
            ```
    """
    delete_image(client, document.parentEntity, document.parentId, document.id)
