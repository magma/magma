#!/usr/bin/env python3

import os
from datetime import datetime
from typing import Optional

import filetype

from .. import SymphonyClient
from ..graphql.enum.image_entity import ImageEntity
from ..graphql.input.add_image import AddImageInput
from ..graphql.mutation.add_image import AddImageMutation
from ..graphql.mutation.delete_image import DeleteImageMutation


def add_image(
    client: SymphonyClient,
    local_file_path: str,
    entity_type: ImageEntity,
    entity_id: str,
    category: Optional[str] = None,
) -> None:
    file_type = filetype.guess(local_file_path)
    file_type = file_type.MIME if file_type is not None else ""
    img_key = client.store_file(local_file_path, file_type, False)
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


def delete_image(
    client: SymphonyClient, entity_type: ImageEntity, entity_id: str, image_id: str
) -> None:
    DeleteImageMutation.execute(
        client, entityType=entity_type, entityId=entity_id, id=image_id
    )
