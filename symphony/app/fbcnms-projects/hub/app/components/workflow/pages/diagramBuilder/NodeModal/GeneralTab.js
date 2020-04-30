import React from "react";
import { Form, Col, InputGroup, ButtonGroup, Button } from "react-bootstrap";
import { taskDescriptions } from "../../../constants";

const NOT_GENERAL = [
  "type",
  "subWorkflowParam",
  "joinOn",
  "name",
  "taskReferenceName",
  "forkTasks",
  "inputParameters",
  "defaultCase"
];

const GeneralTab = props => {
  const taskName = props.inputs?.name || "";
  const taskRefName = props?.inputs?.taskReferenceName || "";
  const decisionCases = [];
  const caseValueParam = [];

  const renderTaskName = item => (
    <Form.Group>
      <InputGroup size="lg">
        <InputGroup.Prepend>
          <InputGroup.Text>name:</InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control
          type="input"
          disabled={props.inputs?.type === "SIMPLE"}
          onChange={e => props.handleInput(e.target.value, "name")}
          value={item}
        />
      </InputGroup>
      <Form.Text className="text-muted">{taskDescriptions["name"]}</Form.Text>
    </Form.Group>
  );

  const renderTaskRefName = item => (
    <Form.Group>
      <InputGroup size="lg">
        <InputGroup.Prepend>
          <InputGroup.Text>taskReferenceName:</InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control
          type="input"
          onChange={e => props.handleInput(e.target.value, "taskReferenceName")}
          value={item}
        />
      </InputGroup>
      <Form.Text className="text-muted">
        {taskDescriptions["taskReferenceName"]}
      </Form.Text>
    </Form.Group>
  );

  const buttonWrappedField = (item, left, right) => (
    <Form.Group>
      <InputGroup>
        <InputGroup.Prepend>
          <InputGroup.Text>{item[0]}:</InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control type="input" value={item[1]} onChange={() => {}} />
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
      <Form.Text className="text-muted">{taskDescriptions[item[0]]}</Form.Text>
    </Form.Group>
  );

  return (
    <Form>
      {renderTaskName(taskName)}
      {renderTaskRefName(taskRefName)}

      <Form.Row>
        {caseValueParam}
        {decisionCases}
      </Form.Row>
      <hr />

      <Form.Row>
        {Object.entries(props.inputs).map((item, i) => {
          if (!NOT_GENERAL.includes(item[0])) {
            if (item[0] === "decisionCases") {
              return Object.entries(item[1]).forEach((entry, i) => {
                decisionCases.push(
                  <Col sm={6} key={`colGeneral-${i}`}>
                    <Form.Group>
                      <InputGroup>
                        <InputGroup.Prepend>
                          <InputGroup.Text>is equal to</InputGroup.Text>
                        </InputGroup.Prepend>
                        <Form.Control
                          type="input"
                          onChange={e =>
                            props.handleInput(e.target.value, item)
                          }
                          value={entry[0]}
                        />
                      </InputGroup>
                      <Form.Text className="text-muted">
                        {taskDescriptions[item[0]]}
                      </Form.Text>
                    </Form.Group>
                  </Col>
                );
              });
            } else {
              if (item[0] === "optional") {
                return (
                  <Col sm={6} key={`colGeneral-${i}`}>
                    {buttonWrappedField(item, ["<", !item[1]], [">", !item[1]])}
                  </Col>
                );
              } else if (item[0] === "caseValueParam") {
                caseValueParam.push(
                  <Col sm={6} key={`colGeneral-${i}`}>
                    <Form.Group>
                      <InputGroup>
                        <InputGroup.Prepend>
                          <InputGroup.Text>if</InputGroup.Text>
                        </InputGroup.Prepend>
                        <Form.Control
                          type="input"
                          onChange={e =>
                            props.handleInput(e.target.value, item[0])
                          }
                          value={item[1]}
                        />
                      </InputGroup>
                      <Form.Text className="text-muted">
                        {taskDescriptions[item[0]]}
                      </Form.Text>
                    </Form.Group>
                  </Col>
                );
              } else {
                return (
                  <Col sm={6} key={`colGeneral-${i}`}>
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
                      </InputGroup>
                      <Form.Text className="text-muted">
                        {taskDescriptions[item[0]]}
                      </Form.Text>
                    </Form.Group>
                  </Col>
                );
              }
            }
          }
          return null;
        })}
      </Form.Row>
    </Form>
  );
};

export default GeneralTab;
