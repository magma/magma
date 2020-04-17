import React, { useState, useEffect } from "react";
import { Modal, Button, Tab, Tabs, ButtonGroup } from "react-bootstrap";
import DefaultsDescsTab from "./DefaultsDescsTab";
import OutputParamsTab from "./OutputParamsTab";
import GeneralParamsTab from "./GeneralParamsTab";
import { getLabelsFromString } from "../builder-utils";

const GeneralInfoModal = props => {
  const [isWfNameValid, setWfNameValid] = useState(false);
  const [finalWorkflow, setFinalWf] = useState(props.finalWorkflow);
  const isNameLocked = props.isWfNameLocked;

  useEffect(() => {
    setFinalWf(props.finalWorkflow);
  }, [props.finalWorkflow]);

  const handleSave = () => {
    props.saveInputs(finalWorkflow);
    props.lockWorkflowName();
    props.closeModal();
  };

  const handleSubmit = e => {
    if (props.isWfNameLocked || isWfNameValid) {
      handleSave();
    } else {
      e.preventDefault();
      e.stopPropagation();
    }
  };

  const handleClose = () => {
    setFinalWf(props.finalWorkflow);
    props.closeModal();
  };

  const handleInput = (value, key) => {
    let finalWf = { ...finalWorkflow };

    if (key === "name") {
      validateWorkflowName(value);
    }

    finalWf = {
      ...finalWf,
      [key]: value
    };

    setFinalWf(finalWf);
  };

  const validateWorkflowName = name => {
    let isValid = name.length >= 1;
    let workflows = props.workflows || [];

    workflows.forEach(wf => {
      if (wf.name === name) {
        isValid = false;
      }
    });
    setWfNameValid(isValid);
  };

  const getExistingLabels = () => {
    let workflows = props.workflows || [];
    let labels = [];
    workflows.forEach(wf => {
      if (wf.description) {
        labels.push(...getLabelsFromString(wf.description));
      }
    });
    return new Set(labels);
  };

  const handleOutputParam = (key, value) => {
    let finalWf = { ...finalWorkflow };
    let outputParameters = finalWf.outputParameters;

    finalWf = {
      ...finalWf,
      outputParameters: {
        ...outputParameters,
        [key]: value
      }
    };

    setFinalWf(finalWf);
  };

  const handleCustomParam = param => {
    let finalWf = { ...finalWorkflow };
    let outputParameters = finalWf.outputParameters;

    finalWf = {
      ...finalWf,
      outputParameters: {
        ...outputParameters,
        [param]: "provide path"
      }
    };

    setFinalWf(finalWf);
  };

  const handleCustomDefaultAndDesc = (param, defaultValue, description) => {
    let finalWf = { ...finalWorkflow };
    let inputParameters = finalWf.inputParameters || [];
    // eslint-disable-next-line no-useless-concat
    let entry = `${param}` + `[${description}]` + `[${defaultValue}]`;
    let isUnique = true;

    if (inputParameters.length > 0) {
      inputParameters.forEach((elem, i) => {
        if (elem.startsWith(param)) {
          inputParameters[i] = entry;
          return (isUnique = false);
        }
      });
    }

    if (isUnique) {
      inputParameters.push(entry);
    }

    finalWf = { ...finalWf, inputParameters };
    setFinalWf(finalWf);
  };

  const deleteDefaultAndDesc = selectedParam => {
    let finalWf = { ...finalWorkflow };
    let inputParameters = finalWf.inputParameters || [];

    inputParameters.forEach((param, i) => {
      if (param.match(/^(.*?)\[/)[1] === selectedParam) {
        inputParameters.splice(i, 1);
      }
    });

    finalWf = { ...finalWf, inputParameters };
    setFinalWf(finalWf);
  };

  const deleteOutputParam = selectedParam => {
    let finalWf = { ...finalWorkflow };
    let outputParameters = finalWf.outputParameters || [];

    delete outputParameters[selectedParam];

    finalWf = { ...finalWf, outputParameters };
    setFinalWf(finalWf);
  };

  return (
    <Modal
      size="lg"
      show={props.show}
      onHide={isNameLocked ? handleClose : () => false}
    >
      <Modal.Header>
        <Modal.Title>
          {isNameLocked ? "Edit general informations" : "Create new workflow"}
        </Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ padding: "30px" }}>
        <Tabs style={{ marginBottom: "20px" }}>
          <Tab eventKey={1} title="General">
            <GeneralParamsTab
              finalWf={finalWorkflow}
              handleInput={handleInput}
              isWfNameValid={isWfNameValid}
              handleSubmit={handleSubmit}
              isWfNameLocked={isNameLocked}
              getExistingLabels={getExistingLabels}
            />
          </Tab>
          <Tab eventKey={2} title="Output parameters">
            <OutputParamsTab
              finalWf={finalWorkflow}
              handleSubmit={handleSubmit}
              handleOutputParam={handleOutputParam}
              handleCustomParam={handleCustomParam}
              deleteOutputParam={deleteOutputParam}
            />
          </Tab>
          <Tab eventKey={3} title="Defaults & description">
            <DefaultsDescsTab
              finalWf={finalWorkflow}
              deleteDefaultAndDesc={deleteDefaultAndDesc}
              handleCustomDefaultAndDesc={handleCustomDefaultAndDesc}
            />
          </Tab>
        </Tabs>
        <ButtonGroup style={{ width: "100%", marginTop: "20px" }}>
          {!isNameLocked ? (
            <Button variant="outline-secondary" onClick={props.redirectOnExit}>
              Cancel
            </Button>
          ) : null}
          <Button onClick={handleSubmit} variant="primary">
            Save
          </Button>
        </ButtonGroup>
      </Modal.Body>
    </Modal>
  );
};

export default GeneralInfoModal;
