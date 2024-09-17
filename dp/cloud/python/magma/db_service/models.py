"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from sqlalchemy import (
    JSON,
    BigInteger,
    Boolean,
    Column,
    DateTime,
    Float,
    ForeignKey,
    Integer,
    String,
)
from sqlalchemy import text as sa_text
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship

Base = declarative_base()


def now():
    """
    Return a function for setting actual time of datetime columns
    """
    return sa_text('statement_timestamp()')


class DBRequestType(Base):
    """
    SAS DB Request type class
    """
    __tablename__ = "request_types"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)

    requests = relationship("DBRequest", back_populates="type")

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', name='{self.name}')>"


class DBRequest(Base):
    """
    SAS DB Request class
    """
    __tablename__ = "requests"
    id = Column(Integer, primary_key=True, autoincrement=True)
    type_id = Column(Integer, ForeignKey("request_types.id", ondelete="CASCADE"))
    cbsd_id = Column(Integer, ForeignKey("cbsds.id", ondelete="CASCADE"), index=True)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )
    payload = Column(JSON)

    type = relationship("DBRequestType", back_populates="requests")
    cbsd = relationship("DBCbsd", back_populates="requests")

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
            f"type_id='{self.type_id}', " \
            f"cbsd_id='{self.cbsd_id}' " \
            f"created_date='{self.created_date}' " \
            f"updated_date='{self.updated_date}' " \
            f"payload='{self.payload}')>"


class DBGrantState(Base):
    """
    SAS DB Grant state class
    """
    __tablename__ = "grant_states"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)

    grants = relationship("DBGrant", back_populates="state")

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
            f"name='{self.name}'>"


class DBGrant(Base):
    """
    SAS DB Grant class
    """
    __tablename__ = "grants"
    id = Column(Integer, primary_key=True, autoincrement=True)
    state_id = Column(
        Integer, ForeignKey(
            "grant_states.id", ondelete="CASCADE",
        ), nullable=False,
    )
    cbsd_id = Column(Integer, ForeignKey("cbsds.id", ondelete="CASCADE"), index=True)
    grant_id = Column(String, nullable=False)
    grant_expire_time = Column(DateTime(timezone=True))
    transmit_expire_time = Column(DateTime(timezone=True))
    heartbeat_interval = Column(Integer)
    last_heartbeat_request_time = Column(DateTime(timezone=True))
    channel_type = Column(String)
    low_frequency = Column(BigInteger, nullable=False)
    high_frequency = Column(BigInteger, nullable=False)
    max_eirp = Column(Float, nullable=False)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )

    state = relationship("DBGrantState", back_populates="grants")
    cbsd = relationship("DBCbsd", back_populates="grants")

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        state_name = self.state.name
        return f"<{class_name}(id='{self.id}', " \
            f"state='{state_name}', " \
            f"cbsd_id='{self.cbsd_id}', " \
            f"grant_id='{self.grant_id}', " \
            f"grant_expire_time='{self.grant_expire_time}', " \
            f"transmit_expire_time='{self.transmit_expire_time}', " \
            f"heartbeat_interval='{self.heartbeat_interval}', " \
            f"last_heartbeat_request_time='{self.last_heartbeat_request_time}', " \
            f"channel_type='{self.channel_type}', " \
            f"created_date='{self.created_date}' " \
            f"updated_date='{self.updated_date}')>"


class DBCbsdState(Base):
    """
    SAS DB CBSD registered state class
    """
    __tablename__ = "cbsd_states"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
               f"name='{self.name}'>"


class DBCbsd(Base):
    """
    SAS DB CBSD class
    """
    __tablename__ = "cbsds"
    id = Column(Integer, primary_key=True, autoincrement=True)
    state_id = Column(
        Integer, ForeignKey(
            "cbsd_states.id", ondelete="CASCADE",
        ), nullable=False,
    )
    desired_state_id = Column(
        Integer, ForeignKey(
            "cbsd_states.id", ondelete="CASCADE",
        ), nullable=False,
    )
    cbsd_id = Column(String)
    user_id = Column(String)
    fcc_id = Column(String)
    cbsd_serial_number = Column(String, unique=True, index=True)
    last_seen = Column(DateTime(timezone=True))
    min_power = Column(Float)
    max_power = Column(Float)
    antenna_gain = Column(Float)
    number_of_ports = Column(Integer)
    preferred_bandwidth_mhz = Column(
        Integer, nullable=False, server_default='0',
    )
    preferred_frequencies_mhz = Column(
        JSON, nullable=False, server_default=sa_text("'[]'::json"),
    )
    single_step_enabled = Column(Boolean, nullable=False, server_default='false')
    cbsd_category = Column(String, nullable=False, server_default='b')
    network_id = Column(String)
    latitude_deg = Column(Float)
    longitude_deg = Column(Float)
    height_m = Column(Float)
    height_type = Column(String)
    indoor_deployment = Column(Boolean, nullable=False, server_default='false')
    is_deleted = Column(Boolean, nullable=False, server_default='false')
    should_deregister = Column(Boolean, nullable=False, server_default='false')
    should_relinquish = Column(Boolean, nullable=False, server_default='false')
    carrier_aggregation_enabled = Column(Boolean, nullable=False, server_default='false')
    max_ibw_mhz = Column(Integer, nullable=False, server_default='150')
    grant_redundancy = Column(Boolean, nullable=False, server_default='true')
    available_frequencies = Column(JSON)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )

    state = relationship("DBCbsdState", foreign_keys=[state_id])
    desired_state = relationship("DBCbsdState", foreign_keys=[desired_state_id])
    requests = relationship("DBRequest", back_populates="cbsd")
    grants = relationship("DBGrant", back_populates="cbsd")
    channels = Column(JSON, nullable=False, server_default=sa_text("'[]'::json"))

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
               f"state_id='{self.state_id}', " \
               f"cbsd_id='{self.cbsd_id}', " \
               f"user_id='{self.user_id}', " \
               f"fcc_id='{self.fcc_id}', " \
               f"cbsd_serial_number='{self.cbsd_serial_number}', " \
               f"created_date='{self.created_date}' " \
               f"updated_date='{self.updated_date}')>"
