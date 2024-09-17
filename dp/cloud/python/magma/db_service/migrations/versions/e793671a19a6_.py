"""empty message

Revision ID: e793671a19a6
Revises: 9cd338f28663
Create Date: 2022-04-22 14:41:00.038838

"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision = 'e793671a19a6'
down_revision = '9cd338f28663'
branch_labels = None
depends_on = None

cbsds = sa.table(
    'cbsds',
    sa.column('id', sa.Integer),
    sa.column('desired_state_id', sa.Integer),
)

active_mode_configs = sa.table(
    'active_mode_configs',
    sa.column('cbsd_id', sa.Integer),
    sa.column('desired_state_id', sa.Integer),
)

cbsd_states = sa.table(
    'cbsd_states',
    sa.column('id', sa.Integer),
    sa.column('name', sa.String),
)


def upgrade():
    """
    Run upgrade
    """
    op.add_column('cbsds', sa.Column('desired_state_id', sa.Integer()))
    op.execute(
        cbsds.update().
        values(
            desired_state_id=sa.select(cbsd_states.c.id).
            where(cbsd_states.c.name == 'unregistered').
            scalar_subquery(),
        ),
    )
    op.execute(
        cbsds.update().
        values(desired_state_id=active_mode_configs.c.desired_state_id).
        where(cbsds.c.id == active_mode_configs.c.cbsd_id),
    )
    op.alter_column('cbsds', 'desired_state_id', nullable=False)
    op.create_foreign_key(None, 'cbsds', 'cbsd_states', ['desired_state_id'], ['id'], ondelete='CASCADE')
    op.drop_table('active_mode_configs')
    # ### end Alembic commands ###


def downgrade():
    """
    Run downgrade
    """
    op.create_table(
        'active_mode_configs',
        sa.Column('id', sa.INTEGER(), autoincrement=True, nullable=False),
        sa.Column('cbsd_id', sa.INTEGER(), autoincrement=False, nullable=False),
        sa.Column('desired_state_id', sa.INTEGER(), autoincrement=False, nullable=False),
        sa.Column('created_date', postgresql.TIMESTAMP(timezone=True), server_default=sa.text('statement_timestamp()'), autoincrement=False, nullable=False),
        sa.Column('updated_date', postgresql.TIMESTAMP(timezone=True), server_default=sa.text('statement_timestamp()'), autoincrement=False, nullable=True),
        sa.ForeignKeyConstraint(['cbsd_id'], ['cbsds.id'], name='active_mode_configs_cbsd_id_fkey', ondelete='CASCADE'),
        sa.ForeignKeyConstraint(['desired_state_id'], ['cbsd_states.id'], name='active_mode_configs_desired_state_id_fkey', ondelete='CASCADE'),
        sa.PrimaryKeyConstraint('id', name='active_mode_configs_pkey'),
        sa.UniqueConstraint('cbsd_id', name='active_mode_configs_cbsd_id_key'),
    )
    op.execute(
        active_mode_configs.insert().
        from_select(['cbsd_id', 'desired_state_id'], sa.select(cbsds.c.id, cbsds.c.desired_state_id)),
    )
    op.drop_column('cbsds', 'desired_state_id')
    # ### end Alembic commands ###
