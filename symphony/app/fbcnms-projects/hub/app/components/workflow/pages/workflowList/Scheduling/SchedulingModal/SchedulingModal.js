import React, { useEffect, useState } from "react";
import { Button, Modal, Form } from "react-bootstrap";
import WorkflowDia from "../../WorkflowExec/DetailsModal/WorkflowDia/WorkflowDia";
import { HttpClient as http } from "../../../../common/HttpClient";
import { conductorApiUrlPrefix } from "../../../../constants";

const stateSubmit = "Submit";
const stateSubmitting = "Submitting..."

const SchedulingModal = props => {
  const [schedule, setSchedule] = useState(()=>{
    http
      .get(conductorApiUrlPrefix + "/schedule/" + props.name)
      .then(res => {
        setSchedule(res);
      });
    return null;
  });
  const [status, setStatus] = useState(stateSubmit);
  const [error, setError] = useState();

  const handleClose = () => {
    props.onClose();
  };

  const submitForm = () => {
    setError(null);
    setStatus(stateSubmitting);
    http.put(conductorApiUrlPrefix + "/schedule/" + schedule.name, schedule).then(res => {
      handleClose();
    }).catch(error => {
      setStatus(stateSubmit);
      setError("Request failed:" + error);
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

  return (
    <Modal
      size="lg"
      dialogClassName="modal-70w"
      show="true"
      onHide={handleClose}
    >
      <Modal.Header>
        <Modal.Title>Schedule Details</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form onSubmit={submitForm}>
          <Form.Group>
            <Form.Label>Cron</Form.Label>
            <Form.Control
              type="input"
              onChange={e => setCronString(e.target.value)}
              placeholder="Enter cron pattern"
              value={schedule?.cronString}
            />
          </Form.Group>
          <Form.Group>
            <Form.Label>Enabled</Form.Label>
            <Form.Control
              type="checkbox"
              onChange={e => setEnabled(e.target.checked)}
              checked={schedule?.enabled}
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
