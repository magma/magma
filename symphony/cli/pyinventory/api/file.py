#!/usr/bin/env python3
# pyre-strict


import glob
import os.path
from datetime import datetime
from typing import Generator, Optional

import filetype

from ..client import SymphonyClient
from ..consts import Document, Location, SiteSurvey
from ..graphql.add_image_input import AddImageInput
from ..graphql.add_image_mutation import AddImageMutation
from ..graphql.delete_image_mutation import DeleteImageMutation
from ..graphql.image_entity_enum import ImageEntity


def store_file(
    client: SymphonyClient, file_path: str, file_type: str, is_global: bool
) -> str:
    sign_response = client.session.get(
        client.put_endpoint,
        params={"contentType": file_type},
        headers={"Is-Global": str(is_global)},
    )
    sign_response = sign_response.json()
    signed_url = sign_response["URL"]
    with open(file_path, "rb") as f:
        file_data = f.read()
    response = client.session.put(
        signed_url, data=file_data, headers={"Content-Type": file_type}
    )
    response.raise_for_status()
    return sign_response["key"]


def delete_file(client: SymphonyClient, key: str, is_global: bool) -> None:
    sign_response = client.session.delete(
        client.delete_endpoint.format(key),
        headers={"Is-Global": str(is_global)},
        allow_redirects=False,
    )
    sign_response.raise_for_status()
    assert sign_response.status_code == 307
    signed_url = sign_response.headers["location"]
    response = client.session.delete(signed_url)
    response.raise_for_status()


def _add_image(
    client: SymphonyClient,
    local_file_path: str,
    entity_type: ImageEntity,
    entity_id: str,
    category: Optional[str] = None,
) -> None:
    file_type = filetype.guess(local_file_path)
    file_type = file_type.MIME if file_type is not None else ""
    img_key = store_file(client, local_file_path, file_type, False)
    file_size = os.path.getsize(local_file_path)

    AddImageMutation.execute(
        client,
        AddImageInput(
            entityType=entity_type,
            entityId=entity_id,
            imgKey=img_key,
            fileName=os.path.basename(local_file_path),
            fileSize=file_size,
            modified=datetime.utcnow(),
            contentType=file_type,
            category=category,
        ),
    )


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
            client (object):
                Client object
            local_file_path (str):
                local system path to the file
            entity_type (str):
                one of existing options ["LOCATION", "WORK_ORDER", "SITE_SURVEY", "EQUIPMENT"]
            entity_id (string):
                valid entity ID
            category (Optional[string]): file category name 

        Returns: None

        Example:
        ```
        client.add_file(client, './document.pdf', 'LOCATION', location.id, 'category_name')
        ```
    """
    entity = {
        "LOCATION": ImageEntity.LOCATION,
        "WORK_ORDER": ImageEntity.WORK_ORDER,
        "SITE_SURVEY": ImageEntity.SITE_SURVEY,
        "EQUIPMENT": ImageEntity.EQUIPMENT,
    }.get(entity_type, ImageEntity.LOCATION)
    _add_image(client, local_file_path, entity, entity_id, category)


def add_files(
    client: SymphonyClient,
    local_directory_path: str,
    entity_type: str,
    entity_id: str,
    category: Optional[str] = None,
) -> None:
    """This function adds all files located in folder to an entity of a given type.

        Args:
            client (object):
                Client object
            local_file_path (str):
                local system path to the file
            entity_type (str):
                one of existing options ["LOCATION", "WORK_ORDER", "SITE_SURVEY", "EQUIPMENT"]
            entity_id (string):
                valid entity ID
            category (Optional[string]): file category name

        Returns: None

        Example:
        ```
        client.add_files(client, './documents_folder/', 'LOCATION', location.id, 'category_name')
        ```
    """
    for file in list_dir(local_directory_path):
        add_file(client, file, entity_type, entity_id, category)


def add_location_image(
    client: SymphonyClient, local_file_path: str, location: Location
) -> None:
    _add_image(client, local_file_path, ImageEntity.LOCATION, location.id)


def add_site_survey_image(
    client: SymphonyClient, local_file_path: str, id: str
) -> None:
    _add_image(client, local_file_path, ImageEntity.SITE_SURVEY, id)


def _delete_image(
    client: SymphonyClient, entity_type: ImageEntity, entity_id: str, image_id: str
) -> None:
    DeleteImageMutation.execute(
        client, entityType=entity_type, entityId=entity_id, id=image_id
    )


def delete_site_survey_image(client: SymphonyClient, survey: SiteSurvey) -> None:
    source_file_key = survey.sourceFileKey
    source_file_id = survey.sourceFileId
    if source_file_key is not None:
        delete_file(client, source_file_key, False)
    if source_file_id is not None:
        _delete_image(client, ImageEntity.SITE_SURVEY, survey.id, source_file_id)


def delete_document(client: SymphonyClient, document: Document) -> None:
    _delete_image(client, document.parentEntity, document.parentId, document.id)
