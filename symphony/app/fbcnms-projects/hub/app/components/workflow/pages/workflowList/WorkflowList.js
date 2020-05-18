import React from "react";
import { Button, Container, Tab, Tabs } from "react-bootstrap";
import { withRouter } from "react-router-dom";
import { HttpClient as http } from "../../common/HttpClient";
import WorkflowDefs from "./WorkflowDefs/WorkflowDefs";
import WorkflowExec from "./WorkflowExec/WorkflowExec";
import Scheduling from "./Scheduling/Scheduling";
import { conductorApiUrlPrefix, frontendUrlPrefix } from "../../constants";
import {changeUrl, exportButton} from './workflowUtils'

const workflowModifyButtons = (openFileUpload, history) => {
  return [
      <Button
          variant="outline-primary"
          style={{ marginLeft: "30px" }}
          onClick={() => history.push(frontendUrlPrefix + "/builder")}
      >
        <i className="fas fa-plus" />
        &nbsp;&nbsp;New
      </Button>,
      <Button
      variant="outline-primary"
      style={{ marginLeft: "5px" }}
      onClick={openFileUpload}
  >
    <i className="fas fa-file-import" />
    &nbsp;&nbsp;Import
  </Button>
  ];
}

const upperMenu = (history, openFileUpload) => {
  return(
    <h1 style={{ marginBottom: "20px" }}>
    <i style={{ color: "grey" }} className="fas fa-cogs" />
    &nbsp;&nbsp;Workflows
    { workflowModifyButtons(openFileUpload, history) }
    { exportButton() }
  </h1>);
}

const WorkflowList = (props) => {
  let urlUpdater = changeUrl(props.history);
  let query = props.match.params.wfid ? props.match.params.wfid : null;

  const importFiles = (e) => {
    const files = e.currentTarget.files;
    const fileList = [];
    let count = files.length;

    Object.keys(files).forEach((i) => {
      readFile(files[i]);
    });

    function readFile(file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        let definition = JSON.parse(e.target.result);
        fileList.push(definition);
        if (!--count) {
          http.put(conductorApiUrlPrefix + '/metadata', fileList).then(() => {
            window.location.reload();
          });
        }
      };
      reader.readAsBinaryString(file);
    }
  };

  const openFileUpload = () => {
    document.getElementById("upload-files").click();
    document
      .getElementById("upload-files")
      .addEventListener("change", importFiles);
  };

  let menu = upperMenu(props.history, openFileUpload);

  return (
    <Container style={{ textAlign: "left", marginTop: "20px" }}>
      {menu}
      <input id="upload-files" multiple type="file" hidden />
      <Tabs
        onSelect={(e) => urlUpdater(e)}
        defaultActiveKey={props.match.params.type || "defs"}
        style={{ marginBottom: "20px" }}
      >
        <Tab mountOnEnter unmountOnExit eventKey="defs" title="Definitions">
          <WorkflowDefs />
        </Tab>
        <Tab mountOnEnter unmountOnExit eventKey="exec" title="Executed">
          <WorkflowExec query={query} />
        </Tab>
        <Tab mountOnEnter unmountOnExit eventKey="scheduled" title="Scheduled">
          <Scheduling />
        </Tab>
      </Tabs>
    </Container>
  );
};

export default withRouter(WorkflowList);
