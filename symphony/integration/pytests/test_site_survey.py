#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from datetime import datetime

from utils.base_test import BaseTest


class TestSiteSurvey(BaseTest):
    def setUp(self):
        super().setUp()
        self.location_types_created = []
        self.location_types_created.append(
            self.client.add_location_type("City Center", [])
        )

    def tearDown(self):
        for location_type in self.location_types_created:
            self.client.delete_location_type_with_locations(location_type)

    def test_site_survey_created(self):
        location = self.client.add_location(
            [("City Center", "Lima Downtown")], {}, 10, 20
        )
        self.assertEqual(0, len(self.client.get_site_surveys(location)))
        completion_date = datetime.strptime("25-7-2019", "%d-%m-%Y")
        self.client.upload_site_survey(
            location,
            "My site survey",
            completion_date,
            "resources/city_center_site_survey.xlsx",
            "resources/city_center_site_survey.json",
        )
        surveys = self.client.get_site_surveys(location)
        self.assertEqual(1, len(surveys))
        survey = surveys[0]
        self.assertEqual("My site survey", survey.name)
        self.assertEqual(completion_date, survey.completionTime)

        self.assertIsNotNone(survey.sourceFileId)
        self.assertEqual(survey.sourceFileName, "city_center_site_survey.xlsx")

        self.client.delete_site_survey(survey)
        surveys = self.client.get_site_surveys(location)
        self.assertEqual(0, len(surveys))
