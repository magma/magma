#!/usr/bin/env python3
# pyre-strict


import glob
import os.path
from datetime import datetime
from typing import Dict, Generator, Optional

import filetype

from ..consts import Document, ImageEntity, Location, SiteSurvey
from ..graphql.add_image_mutation import (
    AddImageInput,
    AddImageMutation,
    ImageEntity as AddImageEntity,
)
from ..graphql.delete_image_mutation import (
    DeleteImageMutation,
    ImageEntity as DeleteImageEntity,
)
from ..graphql_client import GraphqlClient


IMAGE_ENTITY_TO_DELETE_IMAGE_ENTITY: Dict[ImageEntity, DeleteImageEntity] = {
    ImageEntity.LOCATION: DeleteImageEntity.LOCATION,
    ImageEntity.WORK_ORDER: DeleteImageEntity.WORK_ORDER,
    ImageEntity.SITE_SURVEY: DeleteImageEntity.SITE_SURVEY,
    ImageEntity.EQUIPMENT: DeleteImageEntity.EQUIPMENT,
}


def store_file(
    client: GraphqlClient, file_path: str, file_type: str, is_global: bool
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


def delete_file(client: GraphqlClient, key: str, is_global: bool) -> None:
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
    client: GraphqlClient,
    local_file_path: str,
    entity_type: AddImageEntity,
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
    client: GraphqlClient,
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

        Returns: None

        Example:
            client.add_file(client, './document.pdf', 'LOCATION', location.id)
    """
    entity = {
        "LOCATION": AddImageEntity.LOCATION,
        "WORK_ORDER": AddImageEntity.WORK_ORDER,
        "SITE_SURVEY": AddImageEntity.SITE_SURVEY,
        "EQUIPMENT": AddImageEntity.EQUIPMENT,
    }.get(entity_type, AddImageEntity.LOCATION)
    _add_image(client, local_file_path, entity, entity_id, category)


def add_files(
    client: GraphqlClient,
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

        Returns: None

        Example:
            client.add_files(client, './documents_folder/', 'LOCATION', location.id)
    """
    for file in list_dir(local_directory_path):
        add_file(client, file, entity_type, entity_id, category)


def add_location_image(
    client: GraphqlClient, local_file_path: str, location: Location
) -> None:
    _add_image(client, local_file_path, AddImageEntity.LOCATION, location.id)


def add_site_survey_image(client: GraphqlClient, local_file_path: str, id: str) -> None:
    _add_image(client, local_file_path, AddImageEntity.SITE_SURVEY, id)


def _delete_image(
    client: GraphqlClient, entity_type: DeleteImageEntity, entity_id: str, image_id: str
) -> None:
    DeleteImageMutation.execute(
        client, entityType=entity_type, entityId=entity_id, id=image_id
    )


def delete_site_survey_image(client: GraphqlClient, survey: SiteSurvey) -> None:
    source_file_key = survey.sourceFileKey
    source_file_id = survey.sourceFileId
    if source_file_key is not None:
        delete_file(client, source_file_key, False)
    if source_file_id is not None:
        _delete_image(client, DeleteImageEntity.SITE_SURVEY, survey.id, source_file_id)


def delete_document(client: GraphqlClient, document: Document) -> None:
    _delete_image(
        client,
        IMAGE_ENTITY_TO_DELETE_IMAGE_ENTITY[document.parentEntity],
        document.parentId,
        document.id,
    )
