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


class DBRequestState(Base):
    """
    SAS DB Request state class
    """
    __tablename__ = "request_states"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)

    requests = relationship("DBRequest", back_populates="state")

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
    type_id = Column(
        Integer, ForeignKey(
            "request_types.id", ondelete="CASCADE",
        ),
    )
    state_id = Column(
        Integer, ForeignKey(
            "request_states.id", ondelete="CASCADE",
        ),
    )
    cbsd_id = Column(Integer, ForeignKey("cbsds.id", ondelete="CASCADE"))
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )
    payload = Column(JSON)

    state = relationship("DBRequestState", back_populates="requests")
    type = relationship("DBRequestType", back_populates="requests")
    response = relationship("DBResponse", back_populates="request")
    cbsd = relationship("DBCbsd", back_populates="requests")

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        type_name = self.type.name
        state_name = self.state.name
        return f"<{class_name}(id='{self.id}', " \
            f"type='{type_name}', " \
            f"state='{state_name}', " \
            f"cbsd_id='{self.cbsd_id}' " \
            f"created_date='{self.created_date}' " \
            f"updated_date='{self.updated_date}' " \
            f"payload='{self.payload}')>"


class DBResponse(Base):
    """
    SAS DB Response class
    """
    __tablename__ = "responses"
    id = Column(Integer, primary_key=True, autoincrement=True)
    request_id = Column(Integer, ForeignKey("requests.id", ondelete="CASCADE"))
    grant_id = Column(
        Integer, ForeignKey(
            "grants.id", ondelete="CASCADE",
        ), nullable=True,
    )
    response_code = Column(Integer, nullable=False)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    payload = Column(JSON)

    request = relationship("DBRequest", back_populates="response")
    grant = relationship("DBGrant", back_populates="responses")

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
            f"request_id='{self.request_id}', " \
            f"response_code='{self.response_code}', " \
            f"created_date='{self.created_date}' " \
            f"payload='{self.payload}')>"


class DBGrantState(Base):
    """
    SAS DB Grant state class
    """
    __tablename__ = "grant_states"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False, unique=True)

    grants = relationship(
        "DBGrant", back_populates="state", cascade="all, delete",
        passive_deletes=True,
    )

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
    cbsd_id = Column(Integer, ForeignKey("cbsds.id", ondelete="CASCADE"))
    channel_id = Column(Integer, ForeignKey("channels.id", ondelete="CASCADE"))
    grant_id = Column(String, nullable=False)
    grant_expire_time = Column(DateTime(timezone=True))
    transmit_expire_time = Column(DateTime(timezone=True))
    heartbeat_interval = Column(Integer)
    last_heartbeat_request_time = Column(DateTime(timezone=True))
    channel_type = Column(String)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )

    state = relationship(
        "DBGrantState", back_populates="grants", cascade="all, delete",
        passive_deletes=True,
    )
    responses = relationship(
        "DBResponse", back_populates="grant", cascade="all, delete",
        passive_deletes=True,
    )
    cbsd = relationship(
        "DBCbsd", back_populates="grants", cascade="all, delete",
        passive_deletes=True,
    )
    channel = relationship(
        "DBChannel", back_populates="grants", cascade="all, delete",
        passive_deletes=True,
    )

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

    cbsds = relationship(
        "DBCbsd", back_populates="state", cascade="all, delete",
        passive_deletes=True,
    )
    active_mode_configs = relationship(
        "DBActiveModeConfig", back_populates="desired_state",
        cascade="all, delete", passive_deletes=True,
    )

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
    cbsd_id = Column(String)
    user_id = Column(String)
    fcc_id = Column(String)
    cbsd_serial_number = Column(String)
    last_seen = Column(DateTime(timezone=True))
    min_power = Column(Float)
    max_power = Column(Float)
    antenna_gain = Column(Float)
    number_of_ports = Column(Integer)
    network_id = Column(String)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )

    state = relationship(
        "DBCbsdState", back_populates="cbsds", cascade="all, delete",
        passive_deletes=True,
    )
    requests = relationship(
        "DBRequest", back_populates="cbsd", cascade="all, delete",
        passive_deletes=True,
    )
    grants = relationship(
        "DBGrant", back_populates="cbsd", cascade="all, delete",
        passive_deletes=True,
    )
    channels = relationship(
        "DBChannel", back_populates="cbsd", cascade="all, delete",
        passive_deletes=True,
    )
    active_mode_config = relationship(
        "DBActiveModeConfig", back_populates="cbsd",
        cascade="all, delete", passive_deletes=True,
    )

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        state_name = self.state.name
        return f"<{class_name}(id='{self.id}', " \
               f"state='{state_name}', " \
               f"cbsd_id='{self.cbsd_id}', " \
               f"user_id='{self.user_id}', " \
               f"fcc_id='{self.fcc_id}', " \
               f"cbsd_serial_number='{self.cbsd_serial_number}', " \
               f"created_date='{self.created_date}' " \
               f"updated_date='{self.updated_date}')>"


class DBChannel(Base):
    """
    SAS DB Channel class
    """
    __tablename__ = "channels"
    id = Column(Integer, primary_key=True, autoincrement=True)
    cbsd_id = Column(Integer, ForeignKey("cbsds.id", ondelete="CASCADE"))
    low_frequency = Column(BigInteger, nullable=False)
    high_frequency = Column(BigInteger, nullable=False)
    channel_type = Column(String, nullable=False)
    rule_applied = Column(String, nullable=False)
    max_eirp = Column(Float)
    last_used_max_eirp = Column(Float)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )

    cbsd = relationship(
        "DBCbsd", back_populates="channels", cascade="all, delete",
        passive_deletes=True,
    )
    grants = relationship(
        "DBGrant", back_populates="channel", cascade="all, delete",
        passive_deletes=True,
    )

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', cbsd_id='{self.cbsd_id}')>"


class DBActiveModeConfig(Base):
    """
    DB CBSD Active Mode Configuration class
    """
    __tablename__ = "active_mode_configs"
    id = Column(Integer, primary_key=True, autoincrement=True)
    cbsd_id = Column(
        Integer, ForeignKey("cbsds.id", ondelete="CASCADE"),
        nullable=False, unique=True,
    )
    desired_state_id = Column(
        Integer, ForeignKey(
            "cbsd_states.id",
            ondelete="CASCADE",
        ), nullable=False,
    )
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )
    updated_date = Column(
        DateTime(timezone=True),
        server_default=now(), onupdate=now(),
    )

    cbsd = relationship(
        "DBCbsd", back_populates="active_mode_config", cascade="all, delete",
        passive_deletes=True,
    )
    desired_state = relationship(
        "DBCbsdState", back_populates="active_mode_configs",
    )

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
               f"cbsd_id='{self.cbsd_id}', " \
               f"desired_state='{self.desired_state}', " \
               f"created_date='{self.created_date}', " \
               f"updated_date='{self.updated_date}')>"


class DBLog(Base):
    """
    Domain Proxy DB request/response log class
    """
    __tablename__ = "domain_proxy_logs"
    id = Column(Integer, primary_key=True, autoincrement=True)
    log_from = Column(String)
    log_to = Column(String)
    log_name = Column(String)
    log_message = Column(String)
    cbsd_serial_number = Column(String)
    network_id = Column(String)
    fcc_id = Column(String)
    response_code = Column(Integer)
    created_date = Column(
        DateTime(timezone=True),
        nullable=False, server_default=now(),
    )

    def __repr__(self):
        """
        Return string representation of DB object
        """
        class_name = self.__class__.__name__
        return f"<{class_name}(id='{self.id}', " \
               f"log_from='{self.log_from}', " \
               f"log_to='{self.log_to}', " \
               f"log_name='{self.log_name}', " \
               f"log_message='{self.log_message}', " \
               f"cbsd_serial_number='{self.cbsd_serial_number}', " \
               f"network_id='{self.network_id}', " \
               f"fcc_id='{self.fcc_id}', " \
               f"created_date='{self.created_date}')>"
