import React from "react";
import { Col, Container, Modal, Row, Tab, Tabs } from "react-bootstrap";
import Highlight from "react-highlight.js";
import type {Task} from "./flowtypes";

type Props = {
  task: Task,
  show: boolean,
  handle: () => void
};

const TaskModal = (props:Props) => {
  let task = props.task;
  let show = props.show;
  return (
    <Modal
      style={{ marginTop: "-20px" }}
      size="lg"
      scrollable
      show={show}
      onHide={props.handle}
    >
      <Modal.Header closeButton>
        <Modal.Title>
          {task.taskType} ({task.status})
          <div
            style={{
              color: "#ff0000",
              display: task.status === "FAILED" ? "" : "none"
            }}
          >
            {task.reasonForIncompletion}
          </div>
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Tabs defaultActiveKey={1}>
          <Tab eventKey={1} title="Summary">
            <Container>
              <Row style={{ marginTop: "20px", marginBottom: "20px" }}>
                <Col sm={8}>
                  <b>Task Ref. Name:&nbsp;&nbsp;</b>
                  {task.referenceTaskName}
                </Col>
                <Col>
                  <b>Callback After:&nbsp;&nbsp;</b>
                  {task.callbackAfterSeconds
                    ? task.callbackAfterSeconds
                    : 0}{" "}
                  (second)
                  <br />
                  <b>Poll Count:&nbsp;&nbsp;</b>
                  {task.pollCount}
                </Col>
              </Row>
              <hr />
              <Row style={{ marginBottom: "15px" }}>
                <b>
                  Input
                  <i
                    title="copy to clipboard"
                    className="btn fa fa-clipboard"
                    data-clipboard-target="#t_input"
                  />
                </b>
              </Row>
              <Row>
                <code>
                  <pre style={{ width: "770px" }} id="t_input">
                    <Highlight language="json">
                      {JSON.stringify(task.inputData, null, 3)}
                    </Highlight>
                  </pre>
                </code>
              </Row>
              <Row style={{ marginBottom: "15px" }}>
                <b>
                  Output
                  <i
                    title="copy to clipboard"
                    className="btn fa fa-clipboard"
                    data-clipboard-target="#t_output"
                  />
                </b>
              </Row>
              <Row>
                <code>
                  <pre style={{ width: "770px" }} id="t_output">
                    <Highlight language="json">
                      {JSON.stringify(task.outputData, null, 3)}
                    </Highlight>
                  </pre>
                </code>
              </Row>
            </Container>
          </Tab>
          <Tab eventKey={2} title="JSON">
            <br />
            <b>
              JSON
              <i
                title="copy to clipboard"
                className="btn fa fa-clipboard"
                data-clipboard-target="#t_json"
              />
            </b>

            <code>
              <pre
                style={{
                  maxHeight: "500px",
                  marginTop: "20px",
                  backgroundColor: "#eaeef3"
                }}
                id="t_json"
              >
                {JSON.stringify(task, null, 3)}
              </pre>
            </code>
          </Tab>
          <Tab eventKey={3} title="Logs">
            <br />
            <b>
              Logs
              <i
                title="copy to clipboard"
                className="btn fa fa-clipboard"
                data-clipboard-target="#t_logs"
              />
            </b>
            <code>
              <pre
                style={{ maxHeight: "500px", marginTop: "20px" }}
                id="t_logs"
              >
                <Highlight language="json">
                  {JSON.stringify(task.logs, null, 3)}
                </Highlight>
              </pre>
            </code>
          </Tab>
        </Tabs>
      </Modal.Body>
    </Modal>
  );
};

export default TaskModal;
