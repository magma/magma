import React, { useState } from "react";
import { Button, Modal, Form } from "react-bootstrap";
import { conductorApiUrlPrefix } from "../../../../constants";
import superagent from "superagent";

const stateSubmit = "Submit";
const stateSubmitting = "Submitting..."

const SchedulingModal = props => {
  const [schedule, setSchedule] = useState();
  const [status, setStatus] = useState();
  const [error, setError] = useState();

  const handleClose = () => {
    props.onClose();
  };

  const handleShow = () => {
    setSchedule(null);
    setStatus(stateSubmit);
    setError(null);
    const path = conductorApiUrlPrefix + "/schedule/" + props.name;
    const req = superagent.get(path).accept("application/json");
    req.end((err, res) => {
      if (res && res.ok) {
        // found in db
        setSchedule(res.body);
      } else {
        // not found, prepare new object to be created
        setSchedule({
          name: props.name,
          workflowName: props.workflowName,
          // workflowVersion must be string
          workflowVersion: props.workflowVersion + '',
          // new schedule is created with enabled: true due to
          // https://github.com/flaviostutz/schellar/issues/5
          enabled: true,
          cronString: '0 * * ? * *'
        })
      }
    });
  }

  const submitForm = () => {
    setError(null);
    setStatus(stateSubmitting);
    const path = conductorApiUrlPrefix + "/schedule/" + props.name;
    const req = superagent.put(path, schedule).set("Content-Type", "application/json");
    req.end((err, res) => {
      if (res && res.ok) {
        handleClose();
      } else {
        setStatus(stateSubmit);
        setError("Request failed:" + err);
      }
    });
  }

  const setCronString = (str) => {
    let mySchedule = {
      ...schedule,
      cronString: str
    };
    setSchedule(mySchedule);
  }

  const setEnabled = (enabled) => {
    let mySchedule = {
      ...schedule,
      enabled: enabled
    };
    setSchedule(mySchedule);
  }

  const getCronString = () => {
    if (schedule != null) {
      if (schedule.cronString != null) {
        return schedule.cronString;
      }
    }
    return "";
  }

  const getEnabled = () => {
    if (schedule != null) {
      if (typeof schedule.enabled === "boolean") {
        return schedule.enabled;
      } // backend does not send this property when disabled
    }
    return false;
  }

  return (
    <Modal
      size="lg"
      dialogClassName="modal-70w"
      show={props.show}
      onHide={handleClose}
      onShow={handleShow}
    >
      <Modal.Header>
        <Modal.Title>Schedule Details - {props.name}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form onSubmit={submitForm}>
          <Form.Group>
            <Form.Label>Cron</Form.Label>
            <Form.Control
              type="input"
              onChange={e => setCronString(e.target.value)}
              placeholder="Enter cron pattern"
              value={getCronString()}
            />
          </Form.Group>
          <Form.Group>
            <Form.Label>Enabled</Form.Label>
            <Form.Control
              type="checkbox"
              onChange={e => setEnabled(e.target.checked)}
              checked={getEnabled()}
            />
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <pre>
        {error}
        </pre>

        <Button
          variant={status === stateSubmit ? "primary" : "info" }
          onClick={submitForm}
          disabled={status === stateSubmitting}
        >
          {status === stateSubmit ? <i className="fas fa-play" /> : null}
          {status === stateSubmitting ?
          (<i className="fas fa-spinner fa-spin" />) : null}
          &nbsp;&nbsp;{status}
        </Button>
        <Button variant="secondary" onClick={handleClose}>
          Close
        </Button>

      </Modal.Footer>
    </Modal>
  );
};

export default SchedulingModal;
