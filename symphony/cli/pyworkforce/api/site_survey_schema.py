#!/usr/bin/env python3

import copy
import json
import os
from typing import Any, Dict, Tuple

import pkg_resources
from jsonschema import validate

from ..common.constant import SCHEMA_FILE_NAME


def validate_json(path: str) -> None:
    resource_package = __name__
    schema = pkg_resources.resource_string(resource_package, SCHEMA_FILE_NAME).decode()
    loaded_schema = json.loads(schema)
    with open(path, "rb") as f:
        content = json.loads(f.read().decode())
        validate(content, schema=loaded_schema)


def add_dependencies_to_question(
    question_to_replace: Dict[str, Any], depends_on: Dict[str, Any]
) -> Dict[str, Any]:
    replacable_question = question_to_replace
    while "dependsOn" in replacable_question:
        replacable_question = replacable_question["dependsOn"]
    replacable_question.update({"dependsOn": depends_on})
    return question_to_replace


def set_templates_with_content(
    content: Dict[str, Any],
    template_name_to_questions: Dict[str, Any],
    template_name_to_forms: Dict[str, Any],
) -> None:
    categories = content["categories"]

    changed = True
    while changed:
        changed = False
        for category in categories:
            forms = category["forms"]
            indexes = range(len(forms))
            for i in reversed(indexes):
                if "templateName" in forms[i]:
                    changed = True
                    forms = (
                        forms[:i]
                        + [
                            copy.deepcopy(form)
                            for form in template_name_to_forms[forms[i]["templateName"]]
                        ]
                        + forms[i + 1 :]
                    )
                    category["forms"] = forms
            for form in forms:
                questions = form["questions"]
                indexes = range(len(questions))
                for i in reversed(indexes):
                    if "templateName" in questions[i]:
                        changed = True
                        questions_to_replace = template_name_to_questions[
                            questions[i]["templateName"]
                        ]
                        if "dependsOn" in questions[i]:
                            questions_to_replace = [
                                add_dependencies_to_question(
                                    copy.deepcopy(question_to_replace),
                                    copy.deepcopy(questions[i]["dependsOn"]),
                                )
                                for question_to_replace in questions_to_replace
                            ]
                        questions = (
                            questions[:i] + questions_to_replace + questions[i + 1 :]
                        )
                        form["questions"] = questions


def get_templates_from_content(
    content: Dict[str, Any]
) -> Tuple[Dict[str, Any], Dict[str, Any]]:
    if "templates" not in content:
        return {}, {}
    templates = content["templates"]
    template_name_to_questions = {
        template["templateName"]: template["questions"]
        for template in templates
        if "questions" in template
    }
    template_name_to_forms = {
        template["templateName"]: template["forms"]
        for template in templates
        if "forms" in template
    }
    return template_name_to_questions, template_name_to_forms


def retrieve_tamplates_and_set_them(path: str) -> Dict[str, Any]:
    validate_json(path)
    with open(path, "rb") as f:
        content = json.loads(f.read().decode())
    template_name_to_questions, template_name_to_forms = get_templates_from_content(
        content
    )

    if "imports" in content:
        import_files = content["imports"]
        current_dir_path = os.path.dirname(path)
        for import_file_path in import_files:
            import_file_full_path = os.path.join(current_dir_path, import_file_path)
            validate_json(import_file_full_path)
            import_content = open(import_file_full_path, "rb").read().decode()
            import_content = json.loads(import_content)
            (
                import_question_templates,
                import_form_templates,
            ) = get_templates_from_content(import_content)
            template_name_to_questions.update(import_question_templates)
            template_name_to_forms.update(import_form_templates)
    set_templates_with_content(
        content, template_name_to_questions, template_name_to_forms
    )
    return content


def compile_site_survey_schema(
    input_json_file_path: str, output_json_file_path: str
) -> None:
    content = retrieve_tamplates_and_set_them(input_json_file_path)
    try:
        content.pop("imports")
    except KeyError:
        pass
    try:
        content.pop("templates")
    except KeyError:
        pass
    with open(output_json_file_path, "w") as df:
        json.dump(content, df, indent=4)
