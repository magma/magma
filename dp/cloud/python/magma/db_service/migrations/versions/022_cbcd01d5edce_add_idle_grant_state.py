"""add_idle_grant_state

Revision ID: cbcd01d5edce
Revises: 467ad00fbc83
Create Date: 2022-09-20 13:19:39.743907

"""
from alembic import op

# revision identifiers, used by Alembic.
revision = 'cbcd01d5edce'
down_revision = '467ad00fbc83'
branch_labels = None
depends_on = None


def upgrade():
    """
    Run upgrade
    """
    op.execute("INSERT INTO grant_states (name) VALUES ('idle');")


def downgrade():
    """
    Run downgrade
    """
    op.execute("DELETE FROM grants WHERE state_id = (SELECT id FROM grant_states WHERE name = 'idle');")
    op.execute("DELETE FROM grant_states WHERE name = 'idle';")
