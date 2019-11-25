#!/usr/bin/env python3
# pyre-strict

from abc import ABC, abstractmethod
from typing import Any, Dict

import unicodecsv as csv


class FailedOperationException(Exception):
    def __init__(
        self,
        reporter: "Reporter",
        err_msg: str,
        err_id: str,
        mutation_name: str,
        variables: Dict[str, Any],
    ) -> None:
        self.reporter = reporter
        self.err_msg = err_msg
        self.err_id = err_id
        self.mutation_name = mutation_name
        self.variables = variables

    def log_failed_operation(self, row_identifier: str, row: Dict[str, Any]) -> None:
        self.reporter.log_failed_operation(row_identifier, row, self)

    def __str__(self) -> str:
        return "{} ({})".format(self.err_msg, self.err_id)


class Reporter(ABC):
    @abstractmethod
    def log_successful_operation(
        self, mutation_name: str, variables: Dict[str, Any]
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
            successful and failed mutations.
            In order to report on failed mutation, user is required to catch
            FailedOperationException and call logFailedOperation with date
            identifier (row number) & full data for easier debugging later.

            Args:
                out_file_path (str): Path to write csv of successful mutations.
                    Format: mutationName, variables
                err_file_path (str): Path to write csv of failed mutations.
                    Format:
                        mutationName,
                        variables,
                        error_msg,
                        error_id,
                        row_identifier,
                        row

            Example:
                from pyinventory.reporter import InventoryReporter, FailedOperationException
                reporter = InventoryReporter(csvOutPath, csvErrPath)
                client = InventoryClient(email, password, "fb-test", reporter=reporter)
                try:
                    location = client.addLocation(..)
                except FailedOperationException as e:
                    e.logFailedOperation(data_identifier, data)

        """

        # pyre-fixme[4]: Attribute must be annotated.
        self.outFile = csv.writer(open(out_file_path, "wb"), encoding="utf-8")
        self.outFile.writerow(["mutationName", "variables"])

        # pyre-fixme[4]: Attribute must be annotated.
        self.errFile = csv.writer(open(err_file_path, "wb"), encoding="utf-8")
        self.errFile.writerow(
            [
                "mutationName",
                "variables",
                "error_msg",
                "error_id",
                "row_identifier",
                "row",
            ]
        )

    def log_successful_operation(
        self, mutation_name: str, variables: Dict[str, Any]
    ) -> None:
        self.outFile.writerow([mutation_name, str(variables)])

    def log_failed_operation(
        self, row_identifier: str, row: Dict[str, Any], e: FailedOperationException
    ) -> None:
        self.errFile.writerow(
            [
                e.mutation_name,
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
        self, mutation_name: str, variables: Dict[str, Any]
    ) -> None:
        pass

    def log_failed_operation(
        self, row_identifier: str, row: Dict[str, Any], e: FailedOperationException
    ) -> None:
        pass
