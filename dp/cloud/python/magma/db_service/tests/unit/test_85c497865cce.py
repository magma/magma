from magma.db_service.tests.alembic_testcase import AlembicTestCase

CBSDS = 'cbsds'
TEST_STATE_ID = 1
NEW_COLUMNS = [
    'single_step_registration',
    'cbsd_category',
    'latitude',
    'longitude',
    'height',
    'height_type',
    'horizontal_accuracy',
    'antenna_azimuth',
    'antenna_downtilt',
    'antenna_beamwidth',
    'antenna_model',
    'eirp_capability',
    'cpi_digital_signature',
    'indoor_deployment',
]


class Test85c497865cceTestCase(AlembicTestCase):
    down_revision = 'e793671a19a6'
    up_revision = '85c497865cce'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)
        cbsd_states = self.get_table('cbsd_states')
        self.engine.execute(cbsd_states.insert().values(id=TEST_STATE_ID, name='some_state'))
        self.cbsds_pre_upgrade = self.get_table(CBSDS)
        self.engine.execute(
            self.cbsds_pre_upgrade.insert().values(
                id=1, state_id=TEST_STATE_ID, desired_state_id=TEST_STATE_ID,
            ),
        )

    def test_columns_not_present_pre_upgrade(self):
        for col in NEW_COLUMNS:
            self.assertFalse(self.has_column(self.cbsds_pre_upgrade, col))

    def test_columns_present_post_upgrade(self):
        # given
        self.upgrade()

        # when
        cbsds = self.get_table(CBSDS)

        # then
        for col in NEW_COLUMNS:
            self.assertTrue(self.has_column(cbsds, col))

    def test_default_values_post_upgrade(self):
        # given
        self.upgrade()

        # when
        cbsds = self.get_table(CBSDS)
        cbsd = self.engine.execute(cbsds.select()).fetchall()[0]

        # then
        self.assertFalse(cbsd.single_step_registration)
        self.assertFalse(cbsd.indoor_deployment)
        self.assertEqual('B', cbsd.cbsd_category)

    def test_downgrade(self):
        # given
        self.upgrade()

        # when
        self.downgrade()
        for col in NEW_COLUMNS:
            self.assertFalse(self.has_column(self.cbsds_pre_upgrade, col))
