import React, { useState } from "react";
import { Button, Col, Form, InputGroup } from "react-bootstrap";
import { getWfInputsRegex } from "../builder-utils";
import _ from "lodash";

const createInputParamsList = props => {
  const existingInputParameters = props.finalWf.inputParameters || [];
  let inputParametersKeys = Object.keys(getWfInputsRegex(props.finalWf)) || [];

  existingInputParameters.forEach(param => {
    inputParametersKeys.push(param.match(/^(.*?)\[/)[1]);
  });

  inputParametersKeys = _.uniq(inputParametersKeys);

  return inputParametersKeys;
};

const DefaultsDescsTab = props => {
  const inputParamsList = createInputParamsList(props);
  const [selectedParam, setSelectedParam] = useState(inputParamsList[0]);

  const getDescriptionAndDefault = () => {
    let inputParameters = props.finalWf.inputParameters || [];
    let result = [];

    inputParameters.forEach(param => {
      if (param.match(/^(.*?)\[/)[1] === selectedParam) {
        param.match(/\[(.*?)]/g).forEach(group => {
          result.push(group.replace(/[[\]']+/g, ""));
        });
      }
    });
    return result.length > 0 ? result : ["", ""];
  };

  let currentDescription = getDescriptionAndDefault(selectedParam)[0];
  let currentDefault = getDescriptionAndDefault(selectedParam)[1];

  return (
    <div>
      <Form>
        <Form.Group>
          <InputGroup>
            <InputGroup.Prepend>
              <InputGroup.Text>Available input parameters:</InputGroup.Text>
            </InputGroup.Prepend>
            <Form.Control
              disabled={inputParamsList.length === 0}
              onClick={e => setSelectedParam(e.target.value)}
              as="select"
            >
              {inputParamsList.map(param => (
                <option>{param}</option>
              ))}
            </Form.Control>
            <InputGroup.Append>
              <Button
                disabled={inputParamsList.length === 0}
                title="delete parameter's default and description"
                onClick={() => props.deleteDefaultAndDesc(selectedParam)}
                variant="outline-danger"
              >
                <i className="fas fa-times" />
              </Button>
            </InputGroup.Append>
          </InputGroup>
        </Form.Group>
        <Form.Row>
          <Col>
            <Form.Label>default value</Form.Label>
            <Form.Control
              placeholder="default value"
              disabled={inputParamsList.length === 0}
              value={currentDefault}
              onChange={e =>
                props.handleCustomDefaultAndDesc(
                  selectedParam,
                  e.target.value,
                  currentDescription
                )
              }
            />
          </Col>
          <Col>
            <Form.Label>description</Form.Label>
            <Form.Control
              placeholder="description"
              disabled={inputParamsList.length === 0}
              value={currentDescription}
              onChange={e =>
                props.handleCustomDefaultAndDesc(
                  selectedParam,
                  currentDefault,
                  e.target.value
                )
              }
            />
          </Col>
        </Form.Row>
      </Form>
    </div>
  );
};

export default DefaultsDescsTab;
