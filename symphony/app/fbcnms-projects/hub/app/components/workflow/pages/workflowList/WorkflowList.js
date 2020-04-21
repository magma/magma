import { saveAs } from "file-saver";
import React from "react";
import { Button, Container, Tab, Tabs } from "react-bootstrap";
import { withRouter } from "react-router-dom";
import { HttpClient as http } from "../../common/HttpClient";
import WorkflowDefs from "./WorkflowDefs/WorkflowDefs";
import WorkflowExec from "./WorkflowExec/WorkflowExec";
import { conductorApiUrlPrefix } from "../../constants";

const JSZip = require("jszip");

const WorkflowList = (props) => {
  const changeUrl = (e) => {
    props.history.push("/workflows/" + e);
  };

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
          http.put(conductorApiUrlPrefix + "/metadata", fileList).then(() => {
            window.location.reload();
          });
        }
      };
      reader.readAsBinaryString(file);
    }
  };

  const exportFile = () => {
    http.get(conductorApiUrlPrefix + "/metadata/workflow").then((res) => {
      const zip = new JSZip();
      let workflows = res.result || [];

      workflows.forEach((wf) => {
        zip.file(wf.name + ".json", JSON.stringify(wf, null, 2));
      });

      zip.generateAsync({ type: "blob" }).then(function(content) {
        saveAs(content, "workflows.zip");
      });
    });
  };

  let query = props.match.params.wfid ? props.match.params.wfid : null;

  const openFileUpload = () => {
    document.getElementById("upload-files").click();
    document
      .getElementById("upload-files")
      .addEventListener("change", importFiles);
  };

  return (
    <Container style={{ textAlign: "left", marginTop: "20px" }}>
      <h1 style={{ marginBottom: "20px" }}>
        <i style={{ color: "grey" }} className="fas fa-cogs" />
        &nbsp;&nbsp;Workflows
        <Button
          variant="outline-primary"
          style={{ marginLeft: "30px" }}
          onClick={() => props.history.push("/workflows/builder")}
        >
          <i className="fas fa-plus" />
          &nbsp;&nbsp;New
        </Button>
        <Button
          variant="outline-primary"
          style={{ marginLeft: "5px" }}
          onClick={openFileUpload}
        >
          <i className="fas fa-file-import" />
          &nbsp;&nbsp;Import
        </Button>
        <Button
          variant="outline-primary"
          style={{ marginLeft: "5px" }}
          onClick={exportFile}
        >
          <i className="fas fa-file-export" />
          &nbsp;&nbsp;Export
        </Button>
      </h1>
      <input id="upload-files" multiple type="file" hidden />
      <Tabs
        onSelect={(e) => changeUrl(e)}
        defaultActiveKey={props.match.params.type || "defs"}
        style={{ marginBottom: "20px" }}
      >
        <Tab eventKey="defs" title="Definitions">
          <WorkflowDefs />
        </Tab>
        <Tab mountOnEnter unmountOnExit eventKey="exec" title="Executed">
          <WorkflowExec query={query} />
        </Tab>
        <Tab eventKey="scheduled" title="Scheduled" disabled></Tab>
      </Tabs>
    </Container>
  );
};

export default withRouter(WorkflowList);
