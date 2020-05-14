import React, { useEffect, useState } from "react";
import { Button, Modal, Form } from "react-bootstrap";
import WorkflowDia from "../../WorkflowExec/DetailsModal/WorkflowDia/WorkflowDia";
import { HttpClient as http } from "../../../../common/HttpClient";
import { conductorApiUrlPrefix } from "../../../../constants";

const SchedulingModal = props => {
  const [schedule, setSchedule] = useState();
  const [fromDate, setFromDate] = useState();
  const [toDate, setToDate] = useState();
  const [enabled, setEnabled] = useState();
  const [cronString, setCronString] = useState();
  const stateSubmit = "Submit";
  const stateSubmitting = "Submitting..."
  const [status, setStatus] = useState(stateSubmit);
  const [error, setError] = useState();


  useEffect(() => {
    http
      .get(conductorApiUrlPrefix + "/schedule/" + props.name)
      .then(res => {
        setSchedule(res);
        setFromDate(res.fromDate);
        setToDate(res.toDate);
        setEnabled(res.enabled);
        setCronString(res.cronString);
      });
  }, []);

  const handleClose = () => {
    props.modalHandler();
  };


  const submitForm = () => {
    setError(null);
    // update schedule object
    schedule.enabled = enabled;
    schedule.cronString = cronString;
    setStatus(stateSubmitting);
    http.put(conductorApiUrlPrefix + "/schedule/" + props.name, schedule).then(res => {
      handleClose();
    }).catch(error => {
      setStatus(stateSubmit);
      setError("Request failed:" + error);
    });
  }

  return (
    <Modal
      size="lg"
      dialogClassName="modal-70w"
      show={props.show}
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
              value={cronString}
            />
          </Form.Group>
          <Form.Group>
            <Form.Label>Enabled</Form.Label>
            <Form.Control
              type="checkbox"
              onChange={e => setEnabled(e.target.checked)}
              checked={enabled}
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
