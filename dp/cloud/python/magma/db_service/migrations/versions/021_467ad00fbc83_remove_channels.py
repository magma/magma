"""Remove channels

Revision ID: 467ad00fbc83
Revises: fa12c537244a
Create Date: 2022-09-07 12:29:28.162120

"""
import json

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = '467ad00fbc83'
down_revision = '37bd12af762a'
branch_labels = None
depends_on = None


def upgrade():
    """ Upgrade with data migration. """
    op.add_column('cbsds', sa.Column('channels', sa.JSON(), server_default=sa.text("'[]'::json"), nullable=False))

    # migrate existing channels
    conn = op.get_bind()
    conn.execute(
        "UPDATE channels SET max_eirp = %s WHERE max_eirp IS NULL", 37,
    )
    for cbsd in conn.execute('SELECT cbsd_id FROM channels').fetchall():
        cbsd_id = cbsd[0]
        channels = conn.execute(
            "SELECT low_frequency, high_frequency, max_eirp "
            "FROM channels WHERE cbsd_id = %s",
            cbsd_id,
        )
        channels_data = [dict(c) for c in channels.mappings().all()]
        conn.execute("UPDATE cbsds SET channels = %s WHERE id = %s", (json.dumps(channels_data), cbsd_id))

    op.drop_index('ix_channels_cbsd_id', table_name='channels')
    op.drop_table('channels')


def downgrade():
    """ Downgrade with data migration. """
    op.create_table(
        'channels',
        sa.Column('id', sa.INTEGER(), primary_key=True, autoincrement=True, nullable=False),
        sa.Column('cbsd_id', sa.INTEGER(), autoincrement=False, nullable=True),
        sa.Column('low_frequency', sa.BIGINT(), autoincrement=False, nullable=False),
        sa.Column('high_frequency', sa.BIGINT(), autoincrement=False, nullable=False),
        sa.Column('channel_type', sa.VARCHAR(), autoincrement=False, nullable=True),
        sa.Column('rule_applied', sa.VARCHAR(), autoincrement=False, nullable=True),
        sa.Column(
            'max_eirp', postgresql.DOUBLE_PRECISION(precision=53), autoincrement=False, nullable=True,
        ),
        sa.Column(
            'created_date', postgresql.TIMESTAMP(timezone=True),
            server_default=sa.text('statement_timestamp()'), autoincrement=False, nullable=False,
        ),
        sa.Column(
            'updated_date', postgresql.TIMESTAMP(timezone=True),
            server_default=sa.text('statement_timestamp()'), autoincrement=False, nullable=True,
        ),
        sa.ForeignKeyConstraint(
            ['cbsd_id'], ['cbsds.id'], name='channels_cbsd_id_fkey', ondelete='CASCADE',
        ),
        sa.PrimaryKeyConstraint('id', name='channels_pkey'),
    )
    op.create_index('ix_channels_cbsd_id', 'channels', ['cbsd_id'], unique=False)

    # migrate existing channels
    conn = op.get_bind()
    for cbsd_id, channels in conn.execute("SELECT id, channels FROM cbsds WHERE channels::text <> '[]'").fetchall():
        for channel in channels:
            conn.execute(
                "INSERT INTO channels (cbsd_id, low_frequency, high_frequency, max_eirp) "
                "VALUES (%(cbsd_id)s, %(low_frequency)s, %(high_frequency)s, %(max_eirp)s)",
                dict(cbsd_id=cbsd_id, **channel),
            )

    op.drop_column('cbsds', 'channels')
