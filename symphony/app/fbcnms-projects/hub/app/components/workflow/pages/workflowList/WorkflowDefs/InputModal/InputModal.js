import React, { useEffect, useState } from "react";
import {
  Modal,
  Button,
  Form,
  Row,
  Col,
  ToggleButton,
  ToggleButtonGroup
} from "react-bootstrap";
import { useDispatch, useSelector } from "react-redux";
import { Typeahead } from "react-bootstrap-typeahead";
import { getMountedDevices } from "../../../../store/actions/mountedDevices";
import { storeWorkflowId } from "../../../../store/actions/builder";
import { HttpClient as http } from "../../../../common/HttpClient";
import { conductorApiUrlPrefix, frontendUrlPrefix } from "../../../../constants";


const getInputs = def => {
  let matchArray = def.match(/(?<=workflow\.input\.)([a-zA-Z0-9-_]+)/gim);
  return [...new Set(matchArray)];
};

const getDetails = (def, inputsArray) => {
  let [detailsArray, tmpDesc, tmpValue, descs, values] = [[], [], [], [], []];

  if (inputsArray.length > 0) {
    for (let i = 0; i < inputsArray.length; i++) {
      let RegExp3 = new RegExp(`\\b${inputsArray[i]}\\[.*?]"`, "igm");
      detailsArray[i] = def.match(RegExp3);
    }
  }
  for (let i = 0; i < detailsArray.length; i++) {
    if (detailsArray[i]) {
      tmpDesc[i] = detailsArray[i][0].match(/\[.*?\[/);
      tmpValue[i] = detailsArray[i][0].match(/].*?]/);
      if (tmpDesc[i] == null) {
        tmpDesc[i] = detailsArray[i][0].match(/\[(.*?)]/);
        descs[i] = tmpDesc[i][1];
        values[i] = null;
      } else {
        tmpDesc[i] = tmpDesc[i][0].match(/[^[\]"]+/);
        tmpValue[i] = tmpValue[i][0].match(/[^[\]*]+/);
        descs[i] = tmpDesc[i] ? tmpDesc[i][0] : null;
        values[i] = tmpValue[i] ? tmpValue[i][0].replace(/\\/g, "") : null;
      }
    } else {
      descs[i] = null;
      values[i] = null;
    }
  }
  return { descs, values };
};

function InputModal(props) {
  const dispatch = useDispatch();
  const devices = useSelector(state => state.mountedDeviceReducer.devices);
  const [wfdesc, setWfDescs] = useState();
  const [wfId, setWfId] = useState();
  const [name, setName] = useState();
  const [version, setVersion] = useState();
  const [warning, setWarning] = useState([]);
  const [status, setStatus] = useState("Execute");
  const [workflowForm, setWorkflowForm] = useState({
    labels: [],
    descs: [],
    values: []
  });
  const [waitingWfs, setWaitingWfs] = useState([]);
  const backendApiUrlPrefix = props.backendApiUrlPrefix ?? conductorApiUrlPrefix;

  useEffect(() => {
    let name = props.wf.split(" / ")[0];
    let version = props.wf.split(" / ")[1];
    setName(name);
    setVersion(Number(props.wf.split(" / ")[1]));

    http
      .get(backendApiUrlPrefix + "/metadata/workflow/" + name + "/" + version)
      .then(res => {
        let definition = JSON.stringify(res.result, null, 2);
        let description = res.result?.description?.split("-")[0] || "";
        let labels = getInputs(definition);
        let { descs, values } = getDetails(definition, labels);

        if (definition.match(/\bEVENT_TASK\b/)) {
          getWaitingWorkflows().then(waitingWfs => {
            setWaitingWfs(waitingWfs);
          });
        }

        setWfDescs(description);
        setWorkflowForm({
          labels,
          descs,
          values
        });

        if (descs.some(rx => rx && rx.match(/.*#node_id.*/g))) {
          dispatch(getMountedDevices());
        }
      });
  }, [props]);

  const getWaitingWorkflows = () => {
    return new Promise((resolve, reject) => {
      let waitingWfs = [];
      let q = 'status:"RUNNING"';
      http
        .get(
          backendApiUrlPrefix + "/executions/?q=&h=&freeText=" +
            q +
            "&start=" +
            0 +
            "&size="
        )
        .then(res => {
          let runningWfs = res.result?.hits || [];
          let promises = runningWfs.map(wf => {
            return http.get(backendApiUrlPrefix + "/id/" + wf.workflowId);
          });

          Promise.all(promises).then(results => {
            results.forEach(r => {
              let workflow = r.result;
              const waitTasks = workflow?.tasks
                .filter(task => task.taskType === "WAIT")
                .map(t => t.referenceTaskName);
              if (waitTasks.length > 0) {
                let waitingWf = {
                  id: workflow.workflowId,
                  name: workflow.workflowName,
                  waitingTasks: waitTasks
                };
                waitingWfs.push(waitingWf);
              }
            });
            resolve(waitingWfs);
          });
        });
    });
  };

  const handleClose = () => {
    props.modalHandler();
  };

  const handleInput = (e, i) => {
    const workflowFormCopy = { ...workflowForm };
    const warningCopy = { ...warning };

    workflowFormCopy.values[i] = e.target.value;
    warningCopy[i] = !!(
      workflowFormCopy.values[i].match(/^\s.*$/) ||
      workflowFormCopy.values[i].match(/^.*\s$/)
    );

    setWorkflowForm(workflowFormCopy);
    setWarning(warningCopy);
  };

  const handleTypeahead = (e, i) => {
    const workflowFormCopy = { ...workflowForm };
    workflowFormCopy.values[i] = e.toString();
    setWorkflowForm(workflowFormCopy);
  };

  const handleSwitch = (e, i) => {
    const workflowFormCopy = { ...workflowForm };
    workflowFormCopy.values[i] = e ? "true" : "false";
    setWorkflowForm(workflowFormCopy);
  };

  const executeWorkflow = () => {
    let { labels, values } = { ...workflowForm };
    let input = {};
    let payload = {
      name: name,
      version: version,
      input
    };

    for (let i = 0; i < labels.length; i++) {
      if (values[i]) {
        input[labels[i]] = values[i].startsWith("{")
          ? JSON.parse(values[i])
          : values[i];
      }
    }
    setStatus("Executing...");
    http.post(backendApiUrlPrefix + "/workflow", JSON.stringify(payload)).then(res => {
      setStatus(res.statusText);
      setWfId(res.body.text);

      dispatch(storeWorkflowId(res.body.text));
      timeoutBtn();

      if (props.fromBuilder) {
        handleClose();
      }
    });
  };

  const timeoutBtn = () => {
    setTimeout(() => setStatus("Execute"), 1000);
  };

  const inputModel = (type, i) => {
    switch (true) {
      case waitingWfs.length > 0 && type.toLowerCase().includes("id"):
        return (
          <Typeahead
            id={`input-${type}`}
            onChange={e => handleTypeahead(e, i)}
            placeholder="Enter or select workflow id"
            options={waitingWfs.map(w => w.id)}
            defaultSelected={workflowForm.values[i] || ""}
            onInputChange={e => handleTypeahead(e, i)}
            renderMenuItemChildren={option => (
              <div>
                {option}
                <div>
                  <small>
                    name: {waitingWfs.find(w => w.id === option)?.name}
                  </small>
                </div>
              </div>
            )}
          />
        );
      case waitingWfs.length > 0 && type.toLowerCase().includes("task"):
        return (
          <Typeahead
            id={`input-${type}`}
            onChange={e => handleTypeahead(e, i)}
            placeholder="Enter or select task reference name"
            options={waitingWfs.map(w => w.waitingTasks).flat()}
            onInputChange={e => handleTypeahead(e, i)}
            renderMenuItemChildren={option => (
              <div>
                {option}
                <div>
                  <small>
                    name:{" "}
                    {
                      waitingWfs.find(w => w.waitingTasks.includes(option))
                        ?.name
                    }
                  </small>
                </div>
              </div>
            )}
          />
        );
      case /node_id.*/g.test(type):
        return (
          <Typeahead
            id={`input-${i}`}
            onChange={e => handleTypeahead(e, i)}
            placeholder="Enter the node id"
            multiple={!!type.match(/node_ids/g)}
            options={devices}
            selected={devices.filter(
              device => device === workflowForm.values[i]
            )}
            onInputChange={e => handleTypeahead(e, i)}
          />
        );
      case /template/g.test(type):
        return (
          <Form.Control
            type="input"
            as="textarea"
            rows="2"
            onChange={e => handleInput(e, i)}
            placeholder="Enter the input"
            value={workflowForm.values[i] || ""}
            isInvalid={warning[i]}
          />
        );
      case /bool/g.test(type):
        return (
          <ToggleButtonGroup
            type="radio"
            value={workflowForm.values[i] === "true"}
            name={`switch-${i}`}
            onChange={e => handleSwitch(e, i)}
            style={{
              height: "calc(1.5em + .75rem + 2px)",
              width: "100%",
              paddingTop: ".375rem"
            }}
          >
            <ToggleButton size="sm" variant="outline-primary" value={true}>
              On
            </ToggleButton>
            <ToggleButton size="sm" variant="outline-primary" value={false}>
              Off
            </ToggleButton>
          </ToggleButtonGroup>
        );
      default:
        return (
          <Form.Control
            type="input"
            onChange={e => handleInput(e, i)}
            placeholder="Enter the input"
            value={workflowForm.values[i] || ""}
            isInvalid={warning[i]}
          />
        );
    }
  };

  return (
    <Modal size="lg" show={props.show} onHide={handleClose}>
      <Modal.Body style={{ padding: "30px" }}>
        <h4>
          {name} / {version}
        </h4>
        <p className="text-muted">{wfdesc}</p>
        <hr />
        <Form onSubmit={executeWorkflow}>
          <Row>
            {workflowForm.labels.map((item, i) => {
              return (
                <Col sm={6} key={`col1-${i}`}>
                  <Form.Group>
                    <Form.Label>{item}</Form.Label>
                    {warning[i] ? (
                      <div
                        style={{
                          color: "red",
                          fontSize: "12px",
                          float: "right",
                          marginTop: "5px"
                        }}
                      >
                        Unnecessary space
                      </div>
                    ) : null}
                    {inputModel(
                      workflowForm?.descs[i]?.split("#")[1] || item,
                      i
                    )}
                    <Form.Text className="text-muted">
                      {workflowForm?.descs[i]?.split("#")[0] || []}
                    </Form.Text>
                  </Form.Group>
                </Col>
              );
            })}
          </Row>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <a
          style={{ float: "left", marginRight: "50px" }}
          href={`${frontendUrlPrefix}/exec/${wfId}`}
        >
          {wfId}
        </a>
        <Button
          variant={
            status === "OK"
              ? "success"
              : status === "Executing..."
              ? "info"
              : status === "Execute"
              ? "primary"
              : "danger"
          }
          onClick={executeWorkflow}
        >
          {status === "Execute" ? <i className="fas fa-play" /> : null}
          {status === "Executing..." ? (
            <i className="fas fa-spinner fa-spin" />
          ) : null}
          {status === "OK" ? <i className="fas fa-check-circle" /> : null}
          &nbsp;&nbsp;{status}
        </Button>
        <Button variant="secondary" onClick={handleClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

export default InputModal;
