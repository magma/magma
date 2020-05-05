#!/usr/bin/env python3

from abc import ABC, abstractmethod
from typing import Any, Dict

from unicodecsv import writer
from unicodecsv.py3 import UnicodeWriter


class FailedOperationException(Exception):
    def __init__(
        self,
        reporter: "Reporter",
        err_msg: str,
        err_id: str,
        operation_name: str,
        variables: Dict[str, Any],
    ) -> None:
        self.reporter = reporter
        self.err_msg = err_msg
        self.err_id = err_id
        self.operation_name = operation_name
        self.variables = variables

    def log_failed_operation(self, row_identifier: str, row: Dict[str, Any]) -> None:
        self.reporter.log_failed_operation(row_identifier, row, self)

    def __str__(self) -> str:
        return "{} ({})".format(self.err_msg, self.err_id)


class Reporter(ABC):
    @abstractmethod
    def log_successful_operation(
        self, operation_name: str, variables: Dict[str, Any], operation_time: float
    ) -> None:
        pass

    @abstractmethod
    def log_failed_operation(
        self, row_identifier: str, row: Dict[str, Any], e: FailedOperationException
    ) -> None:
        pass


class InventoryReporter(Reporter):
    def __init__(self, out_file_path: str, err_file_path: str) -> None:

        """Reporting utility for the InventoryClient to report on
            successful and failed operations.
            In order to report on failed operation, user is required to catch
            FailedOperationException and call logFailedOperation with date
            identifier (row number) & full data for easier debugging later.

            Args:
                out_file_path (str): Path to write csv of successful operations.
                err_file_path (str): Path to write csv of failed operations.

            Example:
            ```
            from pyinventory.reporter import InventoryReporter, FailedOperationException
            reporter = InventoryReporter(csvOutPath, csvErrPath)
            client = InventoryClient(email, password, "fb-test", reporter=reporter)
            try:
                location = client.addLocation(..)
            except FailedOperationException as e:
                e.logFailedOperation(data_identifier, data)
            ```

        """

        self.out_file: UnicodeWriter = writer(
            open(out_file_path, "wb"), encoding="utf-8"
        )
        self.out_file.writerow(["operation_name", "variables", "operation_time"])

        self.err_file: UnicodeWriter = writer(
            open(err_file_path, "wb"), encoding="utf-8"
        )
        self.err_file.writerow(
            [
                "operation_name",
                "variables",
                "error_msg",
                "error_id",
                "row_identifier",
                "row",
            ]
        )

    def log_successful_operation(
        self, operation_name: str, variables: Dict[str, Any], operation_time: float
    ) -> None:
        self.out_file.writerow([operation_name, str(variables), operation_time])

    def log_failed_operation(
        self, row_identifier: str, row: Dict[str, Any], e: FailedOperationException
    ) -> None:
        self.err_file.writerow(
            [
                e.operation_name,
                str(e.variables),
                e.err_msg,
                e.err_id,
                row_identifier,
                str(row),
            ]
        )


class DummyReporter(Reporter):
    def __init__(self) -> None:
        pass

    def log_successful_operation(
        self, operation_name: str, variables: Dict[str, Any], operation_time: float
    ) -> None:
        pass

    def log_failed_operation(
        self, row_identifier: str, row: Dict[str, Any], e: FailedOperationException
    ) -> None:
        pass


DUMMY_REPORTER = DummyReporter()
