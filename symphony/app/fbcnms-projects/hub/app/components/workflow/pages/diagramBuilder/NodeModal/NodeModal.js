import React, { useEffect, useState } from "react";
import { Button, Modal, Tab, Tabs } from "react-bootstrap";
import { hash } from "../builder-utils";
import GeneralTab from "./GeneralTab";
import InputsTab from "./InputsTab";
import { HttpClient as http } from "../../../common/HttpClient";
import { conductorApiUrlPrefix } from "../../../constants";

const OBJECT_KEYWORDS = ["template", "body"];

const renameObjKey = (oldObj, oldKey, newKey) => {
  return Object.keys(oldObj).reduce((acc, val) => {
    if (val === oldKey) acc[newKey] = oldObj[oldKey];
    else acc[val] = oldObj[val];
    return acc;
  }, {});
};

function NodeModal(props) {
  const [inputs, setInputs] = useState([]);
  const [name, setName] = useState();
  const [version, setVersion] = useState();
  const [inputParameters, setInputParameters] = useState([]);

  useEffect(() => {
    setName(props.inputs.inputs.name);
    setInputs(props.inputs.inputs);

    const { subWorkflowParam } = props.inputs.inputs;

    if (subWorkflowParam) {
      const { name, version } = subWorkflowParam;
      setName(name);
      setVersion(version);

      http
        .get(conductorApiUrlPrefix + "/metadata/workflow/" + name + "/" + version)
        .then(res => {
          setInputParameters(res.result.inputParameters);
        });
    }
  }, [props.inputs]);

  function handleSave() {
    props.saveInputs(inputs, props.inputs.id);
    props.modalHandler();
  }

  function addNewInputParam(param) {
    let inputParameters = { ...inputs.inputParameters };

    let updatedInputs = {
      ...inputs,
      inputParameters: {
        ...inputParameters,
        [param]: "${workflow.input." + param + "}"
      }
    };

    setInputs(updatedInputs);
  }

  function addRemoveHeader(handleOperation, i) {
    let updatedInputs = { ...inputs };
    let headers = updatedInputs["inputParameters"]["http_request"]["headers"];

    if (handleOperation) {
      let key = "key_" + hash();
      headers[key] = "value_" + hash();
    } else {
      delete headers[Object.keys(headers)[i]];
    }
    updatedInputs["inputParameters"]["http_request"]["headers"] = headers;
    setInputs(updatedInputs);
  }

  function updateInputParams(value, key, entry) {
    let updatedInputs = { ...inputs };
    let inputParameters = updatedInputs.inputParameters;

    if (typeof key[1] === "object") {
      if (OBJECT_KEYWORDS.find(e => entry[0].includes(e))) {
        try {
          value = JSON.parse(value);
        } catch (e) {
          console.log(e);
        }
      }
    }

    updatedInputs = {
      ...updatedInputs,
      inputParameters: {
        ...inputParameters,
        [entry[0]]: value
      }
    };

    setInputs(updatedInputs);
  }

  function updateHTTPHeader(value, i, headerKey) {
    let copiedInputs = { ...inputs };
    let headers = copiedInputs["inputParameters"]["http_request"]["headers"];
    let header = Object.keys(headers)[i];

    if (headerKey) {
      copiedInputs["inputParameters"]["http_request"]["headers"] = renameObjKey(
        headers,
        header,
        value
      );
    } else {
      copiedInputs["inputParameters"]["http_request"]["headers"][
        header
      ] = value;
    }

    return copiedInputs["inputParameters"]["http_request"]["headers"];
  }

  function updateHTTPRequest(value, key, entry, i, headerKey) {
    let copiedInputs = { ...inputs };
    let http_request = copiedInputs.inputParameters.http_request;
    let inputParameters = copiedInputs.inputParameters;

    if (typeof key[1] === "object") {
      if (entry[0] === "headers") {
        value = updateHTTPHeader(value, i, headerKey);
      } else if (OBJECT_KEYWORDS.find(e => entry[0].includes(e))) {
        try {
          value = JSON.parse(value);
        } catch (e) {
          console.log(e);
        }
      }
    }

    if (entry[0] === "method") {
      if (value === "PUT" || value === "POST")
        http_request = { ...http_request, body: "${workflow.input.body}" };
      else delete http_request["body"];
    }

    copiedInputs = {
      ...copiedInputs,
      inputParameters: {
        ...inputParameters,
        http_request: {
          ...http_request,
          [entry[0]]: value
        }
      }
    };

    setInputs(copiedInputs);
  }

  function updateDecisionCase(value) {
    let copiedInputs = { ...inputs };
    let decisionCases = { ...copiedInputs.decisionCases };
    let keyNames = Object.keys(decisionCases);
    let falseCase = decisionCases[keyNames[0]] || [];

    decisionCases = {
      [value]: falseCase
    };

    copiedInputs.decisionCases = decisionCases;
    setInputs(copiedInputs);
  }

  function handleInput(value, key, entry, i, headerKey) {
    switch (key[0]) {
      case "inputParameters":
        updateInputParams(value, key, entry);
        break;
      case "headers":
      case "http_request":
        updateHTTPRequest(value, key, entry, i, headerKey);
        break;
      case "decisionCases":
        updateDecisionCase(value);
        break;
      default: {
        let copiedInputs = { ...inputs };
        copiedInputs = {
          ...copiedInputs,
          [key]: value
        };
        setInputs(copiedInputs);
        break;
      }
    }
  }

  return (
    <Modal size="lg" show={props.show} onHide={props.modalHandler}>
      <Modal.Header>
        <Modal.Title style={{ fontSize: "20px" }}>
          {name} / {version}
        </Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ padding: "30px" }}>
        <Tabs style={{ marginBottom: "20px" }}>
          <Tab eventKey={1} title="General">
            <GeneralTab inputs={inputs} handleInput={handleInput} />
          </Tab>
          <Tab eventKey={2} title="Input parameters">
            <InputsTab
              inputs={inputs}
              handleInput={handleInput}
              addNewInputParam={addNewInputParam}
              inputParameters={inputParameters}
              addRemoveHeader={addRemoveHeader}
            />
          </Tab>
        </Tabs>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="primary" onClick={() => handleSave()}>
          Save
        </Button>
        <Button variant="secondary" onClick={props.modalHandler}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

export default NodeModal;
