import React, { useState } from "react";
import { Button, Col, Form, InputGroup, Row } from "react-bootstrap";

const OutputParamsTab = props => {
  const [customParam, setCustomParam] = useState("");
  const outputParameters = props.finalWf.outputParameters || [];

  const handleCustomParam = e => {
    e.preventDefault();
    e.stopPropagation();
    props.handleCustomParam(customParam);
    setCustomParam("");
  };

  return (
    <div>
      <Row>
        <Form onSubmit={handleCustomParam}>
          <InputGroup style={{ padding: "10px 215px 10px" }}>
            <Form.Control
              value={customParam}
              onChange={e => setCustomParam(e.target.value)}
              placeholder="Add new output parameter name"
            />
            <InputGroup.Append>
              <Button variant="outline-primary" type="submit">
                Add
              </Button>
            </InputGroup.Append>
          </InputGroup>
        </Form>
      </Row>
      <hr className="hr-text" data-content="existing output parameters" />
      <Form onSubmit={props.handleSubmit}>
        <Row>
          {Object.keys(outputParameters).length > 0
            ? Object.entries(outputParameters).map((entry, i) => {
                return (
                  <Col sm={6} key={`col4-${i}`}>
                    <Form.Group>
                      <Form.Label>
                        {entry[0]}&nbsp;&nbsp;
                        <i
                          className="fas fa-times clickable"
                          style={{ color: "red" }}
                          onClick={() => props.deleteOutputParam(entry[0])}
                        />
                      </Form.Label>
                      <Form.Control
                        type="input"
                        onChange={e =>
                          props.handleOutputParam(entry[0], e.target.value)
                        }
                        value={entry[1]}
                      />
                    </Form.Group>
                  </Col>
                );
              })
            : null}
        </Row>
      </Form>
    </div>
  );
};

export default OutputParamsTab;
