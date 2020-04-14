#!/usr/bin/env python3

import math
from datetime import datetime
from itertools import combinations
from typing import Any, Dict, List, Optional, Sequence, Tuple, Union

import pandas as pd
from dacite import Config, from_dict
from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException
from xlsxwriter.format import Format
from xlsxwriter.utility import xl_col_to_name
from xlsxwriter.workbook import Workbook
from xlsxwriter.worksheet import Worksheet

from .api.file import add_site_survey_image, delete_site_survey_image
from .client import SymphonyClient
from .consts import Entity, Location, SiteSurvey
from .exceptions import EntityNotFoundError
from .graphql.create_survey_mutation import CreateSurveyMutation
from .graphql.location_surveys_query import LocationSurveysQuery
from .graphql.remove_site_survey_mutation import RemoveSiteSurveyMutation
from .graphql.survey_create_data_input import SurveyCreateData
from .graphql.survey_fragment import SurveyFragment
from .graphql.survey_question_fragment import SurveyQuestionFragment
from .graphql.survey_question_type_enum import SurveyQuestionType
from .site_survey_schema import retrieve_tamplates_and_set_them


CREATE_SURVEY_MUTATION_NAME = "createSurvey"
SIMPLE_QUESTION_TYPE_TO_REQUIRED_PROPERTY_NAME = {
    "DATE": "dateData",
    "BOOL": "boolData",
    "EMAIL": "emailData",
    "TEXT": "textData",
    "FLOAT": "floatData",
    "INTEGER": "intData",
    "PHONE": "phoneData",
}


def _get_dependencies(question: Dict[str, Any]) -> Tuple[List[str], List[str]]:
    if "dependsOn" not in question:
        return [], []

    question_names = []
    dependency_names = []

    q = question
    while "dependsOn" in q:
        question_names.insert(0, q["dependsOn"]["questionName"])
        dependency_names.insert(0, q["dependsOn"]["name"])
        q = q["dependsOn"]

    return question_names, dependency_names


def _extract_questions(
    content: Dict[str, Any]
) -> Tuple[List[Tuple[str, ...]], List[Dict[str, Any]], int]:
    form_to_questions = [
        (form["formTitle"], list(form["questions"])) for form in content["forms"]
    ]
    full_question_paths = []
    questions = []
    max_deps = 0
    for form_to_question in form_to_questions:
        form_name = form_to_question[0]
        form_questions = form_to_question[1]
        for question in form_questions:
            question_names, question_deps = _get_dependencies(question)
            max_deps = max(max_deps, len(question_deps))
            full_question_paths.append(
                (
                    form_name,
                    *zip(question_deps, question_names),
                    question["questionName"],
                )
            )
            questions.append(question)

    return full_question_paths, questions, max_deps


def _extract_form_descriptions(content: Dict[str, Any]) -> Dict[str, str]:
    form_name_to_description = {}
    for form in content["forms"]:
        if "formDescription" in form:
            form_name_to_description[form["formTitle"]] = form["formDescription"]

    return form_name_to_description


def rpad(l: List[str], char: str, size: int) -> List[str]:
    return l + [char] * (size - len(l))


def adjust_column_and_row_sizes(
    worksheet: Worksheet, rows: List[Tuple[str, ...]], num_columns: int
) -> None:
    """
    Rows height:
    - Title row is thicker than data row
    - blank separator row is thinner than data row
    Columns width:
    - Each column width is wider enough all the text in cells it contains

    :param worksheet: worksheet we operate on
    :param rows: list of tuples where each tuple represents a row in worksheet
    :param num_columns: number of columns in worksheet
    :return:
    """
    num_rows = len(rows)

    # adjust row sizes
    column_max_sizes = [0] * num_columns
    worksheet.set_row(0, 30)
    for i in range(num_rows):
        if rows[i][0] != " ":
            worksheet.set_row(i + 1, 20)
        else:
            worksheet.set_row(i + 1, 10)
        for j in range(num_columns):
            column_max_sizes[j] = max(column_max_sizes[j], len(rows[i][j]) + 5)

    # adjust column sizes
    for j in range(num_columns):
        worksheet.set_column(j, j, column_max_sizes[j])
    worksheet.set_column(num_columns, num_columns, 20)


def write_column_titles(
    worksheet: Worksheet,
    num_columns: int,
    title_format: Format,
    translations: Dict[str, str],
) -> None:
    worksheet.write("A1", translations.get("Form", "Form"), title_format)
    for i in range(1, num_columns - 1):
        worksheet.write_blank(f"{xl_col_to_name(i)}1", "", title_format)
    worksheet.write(
        f"{xl_col_to_name(num_columns - 1)}1",
        translations.get("Question", "Question"),
        title_format,
    )
    worksheet.write(
        f"{xl_col_to_name(num_columns)}1",
        translations.get("Answer", "Answer"),
        title_format,
    )


def write_conditional_blank_color_if_dependency_fails(
    worksheet: Worksheet,
    row_index: int,
    column_index: int,
    question: Dict[str, Any],
    full_question_path: Tuple[str, ...],
    question_to_cell: Dict[Tuple[str, ...], str],
    bool_options: Dict[str, str],
    black_format: Format,
) -> None:
    index = -1
    criterias = []
    q = question
    while "dependsOn" in q:
        compare, value = q["dependsOn"]["compare"], q["dependsOn"]["value"]
        new_question_path = full_question_path[:index]
        dependent_cell_name = question_to_cell[
            (
                new_question_path[0],
                *list(map(lambda x: x[0], list(new_question_path[1:-1]))),
                new_question_path[-1][1],
            )
        ]
        comparison_value = (
            f'"{bool_options["Yes"]}"'
            if value is True
            else '"{boolOptions["No"]}"'
            if value is False
            else f"{value}"
            if isinstance(value, int) or isinstance(value, float)
            else f'"{value}"'
        )
        opposite_compare = (
            "<>"
            if compare == "="
            else "="
            if compare == "<>"
            else ">"
            if compare == "<="
            else "<="
            if compare == ">"
            else ">="
            if compare == "<"
            else "<"
        )
        criterias.append(
            'AND(${0}{1}{2},${0}<>"")'.format(
                dependent_cell_name, opposite_compare, comparison_value
            )
        )
        index = index - 1
        q = q["dependsOn"]
    if len(criterias) > 1:
        final_criteria = "=OR("
        for i, criteria in enumerate(criterias):
            if i == 0:
                final_criteria = final_criteria + criteria
            else:
                final_criteria = final_criteria + "," + criteria
        final_criteria = final_criteria + ")"
    else:
        final_criteria = "=" + criterias[0]
    worksheet.conditional_format(
        row_index,
        column_index,
        row_index,
        column_index,
        {"type": "formula", "criteria": final_criteria, "format": black_format},
    )


def write_form_description_as_cell_comment(
    worksheet: Worksheet,
    form_names: List[str],
    question_index_to_row_index: Dict[int, int],
    form_name_to_description: Dict[str, str],
) -> None:
    # add form description as comment to column and make every other form
    # appear in different color
    last_form_name = ""
    i = 1
    for form_name in form_names:
        if form_name != last_form_name:
            last_form_name = form_name
            if form_name in form_name_to_description:
                worksheet.write_comment(
                    question_index_to_row_index[i + 1],
                    0,
                    form_name_to_description[form_name],
                )
        i = i + 1


def add_workbook_formats(workbook: Workbook) -> Tuple[Format, Format, Format, Format]:
    title_format = workbook.add_format(
        {"bg_color": "#DCDCDC", "bold": True, "align": "center", "valign": "vcenter"}
    )
    date_format = workbook.add_format(
        {
            "num_format": "dd/mm/yyyy",
            "locked": 0,
            "align": "center",
            "valign": "vcenter",
        }
    )
    black_format = workbook.add_format({"bg_color": "#000000"})
    cell_format = workbook.add_format(
        {"locked": 0, "align": "center", "valign": "vcenter"}
    )

    return title_format, date_format, black_format, cell_format


def add_cell_validation(
    worksheet: Worksheet,
    row_index: int,
    column_index: int,
    full_question_path: Tuple[str, ...],
    question: Dict[str, Any],
    question_to_cell: Dict[Tuple[str, ...], str],
    bool_options: Dict[str, str],
    date_format: Format,
    black_format: Format,
) -> None:
    type = question["questionType"]
    has_options = "options" in question
    if type == "BOOL":
        worksheet.data_validation(
            row_index,
            column_index,
            row_index,
            column_index,
            {"validate": "list", "source": list(bool_options.values())},
        )
    elif has_options:
        options = question["options"]
        if "multiOptions" in question and question["multiOptions"]:
            s = set(options)
            multi_options = sum(
                map(lambda r: sorted(list(combinations(s, r))), range(1, len(s) + 1)),
                [],
            )
            options = [";".join(multi_option) for multi_option in multi_options]
        worksheet.data_validation(
            row_index,
            column_index,
            row_index,
            column_index,
            {"validate": "list", "source": options},
        )
    elif type == "INTEGER":
        worksheet.data_validation(
            row_index,
            column_index,
            row_index,
            column_index,
            {"validate": "integer", "criteria": ">", "value": -(2 ** 32)},
        )
    elif type == "FLOAT":
        worksheet.data_validation(
            row_index,
            column_index,
            row_index,
            column_index,
            {"validate": "decimal", "criteria": ">", "value": -(2 ** 32)},
        )
    elif type == "DATE":
        worksheet.write_datetime(
            row_index,
            column_index,
            datetime.strptime("01/01/1970", "%d/%m/%Y"),
            date_format,
        )
    if "questionDescription" in question:
        worksheet.write_comment(
            row_index, column_index, question["questionDescription"]
        )
    if "dependsOn" in question:
        write_conditional_blank_color_if_dependency_fails(
            worksheet,
            row_index,
            column_index,
            question,
            full_question_path,
            question_to_cell,
            bool_options,
            black_format,
        )


def export_to_excel(json_file_path: str, excel_file_path: str) -> None:
    content = retrieve_tamplates_and_set_them(json_file_path)

    # pyre-fixme[45]: Cannot instantiate abstract class `ExcelWriter`.
    writer = pd.ExcelWriter(excel_file_path, engine="xlsxwriter")
    workbook = writer.book
    assert workbook is not None  # pyre-ignore

    title_format, date_format, black_format, cell_format = add_workbook_formats(
        workbook
    )

    translations = content.get("translations", {})
    bool_options = {
        "Yes": translations.get("Yes", "Yes"),
        "No": translations.get("No", "No"),
    }
    for category in content["categories"]:
        full_question_paths, questions, max_deps = _extract_questions(category)
        form_name_to_description = _extract_form_descriptions(category)

        rows = [
            (
                full_question_path[0],
                *rpad(
                    list(map(lambda x: x[0], list(full_question_path[1:-1]))),
                    " ",
                    max_deps,
                ),
                full_question_path[-1],
            )
            for full_question_path in full_question_paths
        ]
        df = pd.DataFrame({row: [] for row in rows})

        for i in range(len(rows) - 1, -1, -1):
            if i != 0 and rows[i - 1][0] != rows[i][0]:
                df.insert(
                    i,
                    column=tuple([" "] * (max_deps + 2)),
                    value="",
                    allow_duplicates=True,
                )
        df = df.transpose()
        df.to_excel(writer, sheet_name=category["categoryName"])

        rows = [record for record in df.to_records()]
        row_indexes = [i + 2 for i, row in enumerate(rows) if row[0] != " "]
        question_index_to_row_index = dict(
            zip(range(2, len(row_indexes) + 2), row_indexes)
        )

        num_columns = max_deps + 2
        worksheet = writer.sheets[category["categoryName"]]
        write_column_titles(worksheet, num_columns, title_format, translations)
        adjust_column_and_row_sizes(worksheet, rows, num_columns)

        question_to_cell = {}

        # add validations and formatting to columns
        for i, (full_question_path, question) in enumerate(
            zip(full_question_paths, questions)
        ):
            cell_name = "{}{}".format(
                xl_col_to_name(num_columns), question_index_to_row_index[i + 2]
            )
            question_to_cell[
                (
                    full_question_path[0],
                    *list(map(lambda x: x[0], list(full_question_path[1:-1]))),
                    full_question_path[-1],
                )
            ] = cell_name

            worksheet.write_blank(cell_name, "", cell_format)
            add_cell_validation(
                worksheet,
                question_index_to_row_index[i + 2],
                num_columns,
                question_to_cell,
                full_question_path,
                question,
                question_to_cell,
                bool_options,
                date_format,
                black_format,
            )

        form_names = [
            full_question_path[0] for full_question_path in full_question_paths
        ]
        write_form_description_as_cell_comment(
            worksheet, form_names, question_index_to_row_index, form_name_to_description
        )

        worksheet.protect()

    writer.save()


def _nan_to_empty_string(value: Union[float, int, str]) -> str:
    return "" if isinstance(value, float) and math.isnan(value) else str(value)


def dependency_check(
    compare: str,
    value: Union[float, int, str, bool],
    result: str,
    bool_options: Dict[str, str],
) -> bool:
    if isinstance(value, bool):
        assert compare in ["=", "<>"]
        if compare == "=":
            return (result == bool_options["Yes"] and value is True) or (
                result == bool_options["No"] and value is False
            )
        else:
            return (result == bool_options["Yes"] and value is False) or (
                result == bool_options["No"] and value is True
            )
    elif isinstance(value, int) or isinstance(value, float):
        float_result = float(result)
        if compare == ">=":
            return float_result >= value
        elif compare == ">":
            return float_result > value
        elif compare == "<":
            return float_result < value
        elif compare == "<=":
            return float_result <= value
        elif compare == "<>":
            return float_result != value
        else:
            return float_result == value
    else:
        assert compare in ["=", "<>"]
        if compare == "=":
            return result == value
        else:
            return result != value


def break_and_validate_coordinates(
    value: str, invalid_type_err_msg: str
) -> Tuple[float, float]:
    value_list = value.strip('"').split(",")
    assert len(value_list) == 2, invalid_type_err_msg
    try:
        new_value = (float(value_list[0]), float(value_list[1]))
    except ValueError:
        raise AssertionError(invalid_type_err_msg)
    assert (
        new_value[0] >= -90 and new_value[0] <= 90
    ), "Latitude is not between -90 and 90 ({})".format(new_value[0])
    assert (
        new_value[1] >= -90 and new_value[1] <= 90
    ), "Longtitude is not between -90 and 90 ({})".format(new_value[1])
    return new_value


def get_response_value(
    full_representation_path: Tuple[str, ...],
    question_type: str,
    value: str,
    bool_options: Dict[str, str],
) -> Dict[str, Any]:
    invalid_type_err_msg = "Value of {} ({}) is not of type {}".format(
        full_representation_path, value, question_type
    )

    if question_type == "COORDS":
        lat, long = break_and_validate_coordinates(value, invalid_type_err_msg)
        return {"latitude": lat, "longitude": long}
    elif question_type == "BOOL":
        assert value in list(bool_options.values()), invalid_type_err_msg
        return {"boolData": True if value == bool_options["Yes"] else False}
    elif question_type == "EMAIL":
        return {"emailData": str(value)}
    elif question_type == "TEXT":
        return {"textData": str(value)}
    elif question_type == "FLOAT":
        try:
            return {"floatData": float(value)}
        except ValueError:
            raise AssertionError(invalid_type_err_msg)
    elif question_type == "INTEGER":
        try:
            return {"intData": int(value)}
        except ValueError:
            raise AssertionError(invalid_type_err_msg)
    elif question_type == "DATE":
        return {
            "dateData": int(datetime.strptime(value, "%Y-%m-%d %H:%M:%S").timestamp())
        }
    elif question_type == "PHONE":
        return {"phoneData": str(value)}
    raise AssertionError(f"question type {question_type} not found")


def get_survey_response(
    category_name: str,
    form_name: str,
    form_index: int,
    full_representation_path: Tuple[str, ...],
    question_index: int,
    question_type: str,
    value: str,
    bool_options: Dict[str, str],
    form_description: Optional[str] = None,
) -> Dict[str, Any]:
    response = {
        "formName": category_name + " - " + form_name,
        "formIndex": form_index,
        "questionText": " - ".join(full_representation_path[1:]),
        "questionFormat": SurveyQuestionType(question_type),
        "questionIndex": question_index,
        "wifiData": [],
        "cellData": [],
        "imagesData": [],
    }

    if form_description:
        response.update({"formDescription": form_description})

    response_value = get_response_value(
        full_representation_path, question_type, value, bool_options
    )
    response.update(response_value)

    return response


def _get_survey_reponses(
    excel_file_path: str, json_file_path: str
) -> List[Dict[str, Any]]:
    content = retrieve_tamplates_and_set_them(json_file_path)

    responses = []

    form_index = -1

    translations = content.get("translations", {})
    bool_options = {
        "Yes": translations.get("Yes", "Yes"),
        "No": translations.get("No", "No"),
    }
    for category in content["categories"]:
        results = pd.read_excel(
            excel_file_path, sheet_name=category["categoryName"]
        ).values.tolist()
        results = [result for result in results if result[0] != " "]
        results = [_nan_to_empty_string(result[-1]) for result in results]
        full_question_paths, questions, max_deps = _extract_questions(category)

        results = {
            (
                full_question_path[0],
                *list(map(lambda x: x[0], list(full_question_path[1:-1]))),
                full_question_path[-1],
            ): results[i]
            for i, full_question_path in enumerate(full_question_paths)
        }
        form_name_to_description = _extract_form_descriptions(category)

        full_representation_path_to_question = {
            (
                full_question_path[0],
                *list(map(lambda x: x[0], list(full_question_path[1:-1]))),
                full_question_path[-1],
            ): question
            for (full_question_path, question) in zip(full_question_paths, questions)
        }

        last_form_name = ""
        question_index = -1
        for full_question_path in full_question_paths:
            full_representation_path = (
                full_question_path[0],
                *list(map(lambda x: x[0], list(full_question_path[1:-1]))),
                full_question_path[-1],
            )
            form_name = full_question_path[0]
            if form_name != last_form_name:
                last_form_name = form_name
                form_index = form_index + 1
                question_index = 0
            else:
                question_index = question_index + 1

            if len(full_question_path) > 2:
                dependencies_checks = []
                depends_on = full_representation_path_to_question[
                    full_representation_path
                ]
                q_path = full_question_path
                while len(q_path) > 2:
                    depends_on = depends_on["dependsOn"]
                    compare, value = depends_on["compare"], depends_on["value"]
                    result = results[
                        (
                            q_path[:-1][0],
                            *list(map(lambda x: x[0], list(q_path[:-1][1:-1]))),
                            q_path[:-1][-1][1],
                        )
                    ]
                    dependencies_checks.insert(0, (compare, value, result))
                    q_path = q_path[:-1]

                dependency_failed = False
                for compare, value, result in dependencies_checks:
                    if not dependency_check(compare, value, result, bool_options):
                        dependency_failed = True
                        break
                if dependency_failed:
                    continue

            question = full_representation_path_to_question[full_representation_path]

            value = results[full_representation_path]
            if "options" in question:
                options = question["options"]
                if "multi_options" in question:
                    s = set(options)
                    multi_options = sum(
                        map(
                            lambda r: sorted(list(combinations(s, r))),
                            range(1, len(s) + 1),
                        ),
                        [],
                    )
                    options = [";".join(multi_option) for multi_option in multi_options]
                assert value in options, "Value of {} ({}) not in options: {}".format(
                    full_representation_path, value, options
                )

            question_type = question["questionType"]

            response = get_survey_response(
                category["categoryName"],
                form_name,
                form_index,
                full_representation_path,
                question_index,
                question_type,
                value,
                bool_options,
                form_name_to_description.get(form_name, None),
            )

            responses.append(response)

    return responses


def upload_site_survey(
    client: SymphonyClient,
    location: Location,
    name: str,
    completion_date: datetime,
    excel_file_path: str,
    json_file_path: str,
) -> None:
    """Upload the site survey to the given completion with the data in the
        given excel file. We use the schema file to validate the input in the
        excel is as needed for upload.

        Args:
            location ( `pyinventory.consts.Location` ): could be retrieved from getLocation or addLocation api
            name (str): name of the site survey
            completion_date (datetime.datetime object): the time the site survey was completed
            excel_file_path (str): the path for the excel with the site survey information
                                The format of this excel should be created by calling site_survey.exportToExcel
                                with the paremeter jsonFilePath (the next parameter of this function)
            json_file_path(str): the path for the json file of the schema of the site survey
                               the json file should comply to the schema found in survey_schema.json
                               Example of the format:
            ```
            {
                "forms": [
                {
                    "formTitle": "Site Management - General Information",
                    "questions": [
                    {
                        "questionName": "Exact address",
                        "questionType": "TEXT"
                    },
                    {
                        "questionName": "Reference for address",
                        "questionType": "TEXT"
                    },
                    {
                        "questionName": "Ubigeo",
                        "questionType": "TEXT"
                    }
                    ]
                }
                ]
            }
            ```

        Raises:
            AssertionException: if input values in the excel are incorrect
            FailedOperationException: internal inventory error
    """
    data_variables = {
        "name": name,
        "completionTimestamp": int(datetime.timestamp(completion_date)),
        "locationID": location.id,
        "surveyResponses": _get_survey_reponses(excel_file_path, json_file_path),
    }
    create_survey_variables = {
        "data": from_dict(
            data_class=SurveyCreateData, data=data_variables, config=Config(strict=True)
        )
    }

    try:
        site_survey_id = CreateSurveyMutation.execute(
            client, **create_survey_variables
        ).__dict__[CREATE_SURVEY_MUTATION_NAME]
        client.reporter.log_successful_operation(
            CREATE_SURVEY_MUTATION_NAME, create_survey_variables
        )
        add_site_survey_image(client, excel_file_path, site_survey_id)
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            CREATE_SURVEY_MUTATION_NAME,
            create_survey_variables,
        )


def _survey_responses_to_forms(
    responses: Sequence[SurveyQuestionFragment]
) -> Dict[str, Dict[str, Any]]:
    forms = {}
    value_does_not_match_type_error_msg = (
        "{} question type doesn't have the correct response type {}"
    )
    for response in responses:
        form_name = response.formName
        if form_name not in forms:
            forms[form_name] = {}
        form = forms[form_name]
        question_text = response.questionText
        question_type = response.questionFormat
        if question_type == "COORDS":
            assert (
                response.latitude is not None and response.longitude is not None
            ), value_does_not_match_type_error_msg.format(
                "COORDS", "[latitude, longitude]"
            )
            form[question_text] = (response.latitude, response.longitude)
        elif question_type == "DATE":
            date_data = response.dateData
            assert date_data is not None, value_does_not_match_type_error_msg.format(
                "DATE", "dateData"
            )
            form[question_text] = datetime.fromtimestamp(date_data)
        else:
            for (
                simple_question_type,
                required_property,
            ) in SIMPLE_QUESTION_TYPE_TO_REQUIRED_PROPERTY_NAME.items():
                if question_type == simple_question_type:
                    assert (
                        response.__dict__[required_property] is not None
                    ), value_does_not_match_type_error_msg.format(
                        simple_question_type, required_property
                    )
                    form[question_text] = response.__dict__[required_property]
                    if question_type == "DATE":
                        form[question_text] = datetime.fromtimestamp(
                            form[question_text]
                        )
                    break

    return forms


def build_site_survey_from_survey_response(survey: SurveyFragment) -> SiteSurvey:
    id = survey.id
    name = survey.name
    completion_time = datetime.fromtimestamp(survey.completionTimestamp)
    source_file = survey.sourceFile
    source_file_id = source_file.id if source_file else None
    source_file_name = source_file.fileName if source_file else None
    source_file_key = source_file.storeKey if source_file else None
    forms = _survey_responses_to_forms(survey.surveyResponses)

    return SiteSurvey(
        id=id,
        name=name,
        completionTime=completion_time,
        sourceFileId=source_file_id,
        sourceFileName=source_file_name,
        sourceFileKey=source_file_key,
        forms=forms,
    )


def get_site_surveys(client: SymphonyClient, location: Location) -> List[SiteSurvey]:
    """Retrieve all site survey completed in the location.

        Args:
            location ( `pyinventory.consts.Location` ): could be retrieved from getLocation or addLocation api

        Returns:
            List[ `pyinventory.consts.SiteSurvey` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: location does not exist
    """

    location_with_surveys = LocationSurveysQuery.execute(
        client, id=location.id
    ).location
    if not location_with_surveys:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location.id)
    surveys = location_with_surveys.surveys
    return [build_site_survey_from_survey_response(survey) for survey in surveys]


def delete_site_survey(client: SymphonyClient, site_survey: SiteSurvey) -> None:
    delete_site_survey_image(client, site_survey)
    RemoveSiteSurveyMutation.execute(client, id=site_survey.id)
