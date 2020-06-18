import React, { useEffect, useState } from "react";
import { Button, Modal } from "react-bootstrap";
import WorkflowDia from "../../WorkflowExec/DetailsModal/WorkflowDia/WorkflowDia";
import { HttpClient as http } from "../../../../common/HttpClient";
import { conductorApiUrlPrefix } from "../../../../constants";

const DiagramModal = props => {
  const [meta, setMeta] = useState([]);
  const backendApiUrlPrefix = props.backendApiUrlPrefix ?? conductorApiUrlPrefix;

  useEffect(() => {
    const name = props.wf.split(" / ")[0];
    const version = props.wf.split(" / ")[1];
    http
      .get(backendApiUrlPrefix + "/metadata/workflow/" + name + "/" + version)
      .then(res => {
        setMeta(res.result);
      });
  }, []);

  const handleClose = () => {
    props.modalHandler();
  };

  return (
    <Modal
      size="lg"
      dialogClassName="modal-70w"
      show={props.show}
      onHide={handleClose}
    >
      <Modal.Header>
        <Modal.Title>Workflow Diagram</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <WorkflowDia meta={meta} tasks={[]} def={true} />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default DiagramModal;
