from magma.protocol_controller.plugins.cbsd_sas.tests.app_testcase import (
    AppTestCase,
)
from parameterized import parameterized

REGISTRATION = 'registration'
DEREGISTRATION = 'deregistration'
RELINQUISHMENT = 'relinquishment'
GRANT = 'grant'
HEARTBEAT = 'heartbeat'
SPECTRUM_INQUIRY = 'spectrumInquiry'
POST = 'post'
GET = 'get'
PUT = 'put'
PATCH = 'patch'
DELETE = 'delete'


class SASProtocolControllerTests(AppTestCase):

    @parameterized.expand([
        REGISTRATION,
        DEREGISTRATION,
        RELINQUISHMENT,
        GRANT,
        HEARTBEAT,
        SPECTRUM_INQUIRY,
    ])
    def test_routes_only_allow_post_calls(self, route):
        # Given / When
        response = self.client.options(f'/sas/v1/{route}')

        # Then
        self.assertListEqual(
            ['OPTIONS', 'POST'],
            sorted(response.allow._headers),
        )
