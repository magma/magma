"""empty message

Revision ID: 98f7ccfbd2f8
Revises: 58a1b16ef73c
Create Date: 2022-07-26 14:46:13.599502

"""
from alembic import op

# revision identifiers, used by Alembic.
revision = '98f7ccfbd2f8'
down_revision = '58a1b16ef73c'
branch_labels = None
depends_on = None


def upgrade():
    """
    Run upgrade
    """
    op.execute("DELETE FROM grants WHERE state_id = (SELECT id FROM grant_states WHERE name = 'idle');")
    op.execute("DELETE FROM grant_states WHERE name = 'idle';")


def downgrade():
    """
    Run downgrade
    """
    op.execute("INSERT INTO grant_states (name) VALUES ('idle');")
