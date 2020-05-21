#!/usr/bin/env python3

from pysymphony.api.image import add_image, delete_image
from pysymphony.graphql.enum.image_entity import ImageEntity

from .. import SymphonyClient
from ..common.data_class import SiteSurvey


def add_site_survey_image(
    client: SymphonyClient, local_file_path: str, id: str
) -> None:
    """This function adds image to existing site survey.

        Args:
            local_file_path (str): local system path to the file
            id (str): site survey ID

        Raises:
            FailedOperationException: on operation failure

        Example:
            ```
            client.add_site_survey_image(
                local_file_path="./document.pdf",
                id="123456"
            )
            ```
    """
    add_image(client, local_file_path, ImageEntity.SITE_SURVEY, id)


def delete_site_survey_image(client: SymphonyClient, survey: SiteSurvey) -> None:
    """This function deletes image from existing site survey.

        Args:
            survey ( `pyinventory.common.data_class.SiteSurvey` ): site survey object

        Raises:
            FailedOperationException: on operation failure

        Example:
            ```
            client.delete_site_survey_image(survey=survey)
            ```
    """
    source_file_key = survey.sourceFileKey
    source_file_id = survey.sourceFileId
    if source_file_key is not None:
        client.delete_file(source_file_key, False)
    if source_file_id is not None:
        delete_image(client, ImageEntity.SITE_SURVEY, survey.survey_id, source_file_id)
