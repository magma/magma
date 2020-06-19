#!/usr/bin/env python3

from datetime import datetime
from typing import Any, Dict, NamedTuple, Optional


class SiteSurvey(NamedTuple):
    """
    Attributes:
        name (str): name
        survey_id (str): ID
        completionTime (datetime): complition time
        sourceFileId (Optional[str]): source file ID
        sourceFileName (Optional[str]): source file name
        sourceFileKey (Optional[str]): source file key
        forms (Dict[str, Dict[str, Any]]): forms
    """

    name: str
    survey_id: str
    completionTime: datetime
    sourceFileId: Optional[str]
    sourceFileName: Optional[str]
    sourceFileKey: Optional[str]
    forms: Dict[str, Dict[str, Any]]


class WorkOrderType(NamedTuple):
    """
    Attributes:
        id (str): ID
    """

    id: str


class WorkOrder(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
    """

    id: str
    name: str
