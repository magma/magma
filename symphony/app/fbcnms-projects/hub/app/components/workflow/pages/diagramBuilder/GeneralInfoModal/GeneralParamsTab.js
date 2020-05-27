import React from "react";
import {
  Button,
  ButtonGroup,
  Col,
  Form,
  InputGroup,
  Row
} from "react-bootstrap";
import { Typeahead } from "react-bootstrap-typeahead";
import { workflowDescriptions } from "../../../constants";
import { getLabelsFromString } from "../builder-utils";

const GeneralParamsTab = props => {
  const { isWfNameLocked, isWfNameValid } = props;
  const hiddenParams = [
    "name",
    "description",
    "schemaVersion",
    "workflowStatusListenerEnabled",
    "tasks",
    "outputParameters",
    "inputParameters",
    "updateTime"
  ];

  const handleSubmit = e => {
    if (e.key === "Enter" || e === "Enter") {
      props.handleSubmit(e);
    }
  };

  const lockedNameField = () => (
    <Form.Group>
      <InputGroup size="lg">
        <InputGroup.Prepend>
          <InputGroup.Text>
            <i className="fas fa-lock" />
            &nbsp;&nbsp;name
          </InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control
          disabled
          type="input"
          onChange={e => props.handleInput(e.target.value, "name")}
          value={props.finalWf["name"]}
        />
      </InputGroup>
    </Form.Group>
  );

  const unlockedNameField = () => (
    <Form.Group>
      <InputGroup size="lg">
        <InputGroup.Prepend>
          <InputGroup.Text>name:</InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control
          isValid={isWfNameValid}
          isInvalid={!isWfNameValid}
          type="input"
          onChange={e => props.handleInput(e.target.value, "name")}
          value={props.finalWf["name"]}
        />
        <Form.Control.Feedback type={isWfNameValid ? "valid" : "invalid"}>
          {isWfNameValid
            ? "unique name of workflow"
            : props.finalWf["description"].length < 1
            ? "unique name of workflow"
            : "workflow with this name already exits"}
        </Form.Control.Feedback>
      </InputGroup>
    </Form.Group>
  );

  const description = () => {
    let desc = "";
    let labels = [];
    let existingLabels = [];

    if (props.finalWf["description"]) {
      desc = props.finalWf["description"].split(" - ")[0];
      labels = getLabelsFromString(props.finalWf["description"]);
      existingLabels = Array.from(props.getExistingLabels());
    }

    return (
      <Form.Group>
        <InputGroup style={{ marginBottom: "8px" }}>
          <InputGroup.Prepend>
            <InputGroup.Text>description:</InputGroup.Text>
          </InputGroup.Prepend>
          <Form.Control
            type="input"
            onChange={e =>
              props.handleInput(e.target.value + " - " + labels, "description")
            }
            value={desc}
          />
        </InputGroup>
        <Typeahead
          allowNew
          multiple
          clearButton
          newSelectionPrefix="Add a new label: "
          defaultSelected={labels}
          value={labels}
          onChange={e =>
            props.handleInput(
              desc +
                " - " +
                e.map(item => (item.label ? item.label.toUpperCase() : item)),
              "description"
            )
          }
          options={existingLabels}
          placeholder="Add labels..."
        />
        <Form.Text className="text-muted">
          {workflowDescriptions["description"]}
        </Form.Text>
      </Form.Group>
    );
  };

  const buttonWrappedField = (item, left, right) => (
    <Form.Group>
      <InputGroup>
        <InputGroup.Prepend>
          <InputGroup.Text>{item[0]}:</InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control value={item[1]} onChange={() => {}} />
        <InputGroup.Append>
          <ButtonGroup>
            <Button
              variant="outline-primary"
              onClick={() => props.handleInput(left[1], item[0])}
            >
              {left[0]}
            </Button>
            <Button
              variant="outline-primary"
              onClick={() => props.handleInput(right[1], item[0])}
            >
              {right[0]}
            </Button>
          </ButtonGroup>
        </InputGroup.Append>
      </InputGroup>
      <Form.Text className="text-muted">
        {workflowDescriptions[item[0]]}
      </Form.Text>
    </Form.Group>
  );

  return (
    <Form onKeyPress={handleSubmit}>
      {isWfNameLocked ? lockedNameField() : unlockedNameField()}
      {description()}
      <Row>
        {Object.entries(props.finalWf).map((item, i) => {
          if (!hiddenParams.includes(item[0])) {
            if (item[0] === "version") {
              return (
                <Col sm={6} key={`col2-${i}`}>
                  {buttonWrappedField(
                    item,
                    ["-", item[1] - 1],
                    ["+", item[1] + 1]
                  )}
                </Col>
              );
            }
            if (item[0] === "restartable") {
              return (
                <Col sm={6} key={`col2-${i}`}>
                  {buttonWrappedField(item, ["<", !item[1]], [">", !item[1]])}
                </Col>
              );
            } else {
              return (
                <Col sm={6} key={`col3-${i}`}>
                  <Form.Group>
                    <InputGroup>
                      <InputGroup.Prepend>
                        <InputGroup.Text>{item[0]}:</InputGroup.Text>
                      </InputGroup.Prepend>
                      <Form.Control
                        type="input"
                        onChange={e =>
                          props.handleInput(e.target.value, item[0])
                        }
                        value={item[1]}
                      />
                      <Form.Text className="text-muted">
                        {workflowDescriptions[item[0]]}
                      </Form.Text>
                    </InputGroup>
                  </Form.Group>
                </Col>
              );
            }
          }
          return null;
        })}
      </Row>
    </Form>
  );
};

export default GeneralParamsTab;
