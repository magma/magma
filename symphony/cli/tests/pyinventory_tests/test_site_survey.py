#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import os
from datetime import datetime

from pyinventory.api.location import add_location
from pyinventory.api.location_type import (
    add_location_type,
    delete_location_type_with_locations,
)
from pyinventory.site_survey import (
    delete_site_survey,
    get_site_surveys,
    upload_site_survey,
)

from .utils.base_test import BaseTest


class TestSiteSurvey(BaseTest):
    def setUp(self) -> None:
        super().setUp()
        self.location_types_created = []
        self.location_types_created.append(
            add_location_type(self.client, "City Center", [])
        )

    def tearDown(self) -> None:
        for location_type in self.location_types_created:
            delete_location_type_with_locations(self.client, location_type)

    def test_site_survey_created(self) -> None:
        location = add_location(
            self.client, [("City Center", "Lima Downtown")], {}, 10, 20
        )
        self.assertEqual(0, len(get_site_surveys(self.client, location)))
        completion_date = datetime.strptime("25-7-2019", "%d-%m-%Y")
        upload_site_survey(
            self.client,
            location,
            "My site survey",
            completion_date,
            os.path.join(
                os.path.dirname(__file__), "resources/city_center_site_survey.xlsx"
            ),
            os.path.join(
                os.path.dirname(__file__), "resources/city_center_site_survey.json"
            ),
        )
        surveys = get_site_surveys(self.client, location)
        self.assertEqual(1, len(surveys))
        survey = surveys[0]
        self.assertEqual("My site survey", survey.name)
        self.assertEqual(completion_date, survey.completionTime)

        self.assertIsNotNone(survey.sourceFileId)
        self.assertEqual(survey.sourceFileName, "city_center_site_survey.xlsx")

        delete_site_survey(self.client, survey)
        surveys = get_site_surveys(self.client, location)
        self.assertEqual(0, len(surveys))
