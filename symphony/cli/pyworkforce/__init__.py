#!/usr/bin/env python3

from gql.gql.reporter import DUMMY_REPORTER, Reporter
from pysymphony import SymphonyClient

from .common.constant import __version__


class WorkforceClient(SymphonyClient):

    from .api.site_survey import (
        delete_site_survey,
        export_to_excel,
        get_site_surveys,
        upload_site_survey,
    )

    def __init__(
        self,
        email: str,
        password: str,
        tenant: str = "fb-test",
        is_local_host: bool = False,
        is_dev_mode: bool = False,
        reporter: Reporter = DUMMY_REPORTER,
    ) -> None:

        super().__init__(
            email,
            password,
            tenant,
            f"Pyworkforce/{__version__}",
            is_local_host,
            is_dev_mode,
            reporter,
        )
