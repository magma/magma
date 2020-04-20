import React, { useEffect, useState } from "react";
import { Modal, Button } from "react-bootstrap";
import Highlight from "react-highlight.js";
import axios from 'axios'

const DefinitionModal = props => {
  const [definition, setDefinition] = useState("");

  useEffect(() => {
    const name = props.wf.split(" / ")[0];
    const version = props.wf.split(" / ")[1];
    axios
      .get("/workflows/metadata/workflow/" + name + "/" + version)
      .then(res => {
        setDefinition(JSON.stringify(res.result, null, 2));
      });
  }, []);

  const handleClose = () => {
    props.modalHandler();
  };

  return (
    <Modal size="xl" show={props.show} onHide={handleClose}>
      <Modal.Header>
        <Modal.Title>{props.wf}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <code style={{ fontSize: "17px" }}>
          <pre style={{ maxHeight: "600px" }}>
            <Highlight language="json">{definition}</Highlight>
          </pre>
        </code>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default DefinitionModal;
