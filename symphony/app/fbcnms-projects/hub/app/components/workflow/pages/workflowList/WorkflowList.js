import axios from "axios";
import { saveAs } from "file-saver";
import React from "react";
import { withRouter } from "react-router-dom";
import { Button, Container, Icon, Tab } from "semantic-ui-react";
import WorkflowDefs from "./WorkflowDefs/WorkflowDefs";
import WorkflowExec from "./WorkflowExec/WorkflowExec";

const JSZip = require("jszip");

const containerStyle = { textAlign: "left", marginTop: "20px" };

const listHeader = {
  display: "flex",
  justifyContent: "space-between",
  h: {
    font: '600 30px/45px "Poppins", sans-serif',
    color: "#282835",
  },
};

const tabPanes = [
  {
    menuItem: "Definitions",
    render: () => (
      <Tab.Pane attached={false}>
        <WorkflowDefs />
      </Tab.Pane>
    ),
  },
  {
    menuItem: "Executed",
    render: () => (
      <Tab.Pane attached={false}>
        <WorkflowExec />
      </Tab.Pane>
    ),
  },
];

const WorkflowList = (props) => {
  const changeUrl = (path) => {
    props.history.push("/workflows/" + path.toLowerCase());
  };

  const importFiles = (e) => {
    const files = e.currentTarget.files;
    const fileList = [];

    Object.keys(files).forEach((i) => {
      readFile(files[i]);
    });

    function readFile(file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        let definition = JSON.parse(e.target.result);
        fileList.push(definition);
        if (!--files.length) {
          axios.put("/workflows/metadata", fileList).then(() => {
            window.location.reload();
          });
        }
      };
      reader.readAsBinaryString(file);
    }
  };

  const exportFile = () => {
    axios.get("/workflows/metadata/workflow").then((res) => {
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

  const openFileUpload = () => {
    document.getElementById("upload-files").click();
    document
      .getElementById("upload-files")
      .addEventListener("change", importFiles);
  };

  return (
    <Container style={containerStyle}>
      <div style={listHeader}>
        <h1 style={listHeader.h}>Workflows</h1>
        <div className="actions">
          <Button
            primary
            onClick={() => props.history.push("/hub/workflows/builder")}
          >
            <Icon name="plus" />
            New
          </Button>
          <Button onClick={openFileUpload}>
            <Icon name="upload" />
            Import
          </Button>
          <Button onClick={exportFile}>
            <Icon name="download" />
            Export
          </Button>
        </div>
      </div>
      <input id="upload-files" multiple type="file" hidden />
      <Tab
        menu={{ pointing: true }}
        onTabChange={(e, data) =>
          changeUrl(data.panes[data.activeIndex].menuItem)
        }
        panes={tabPanes}
      />
    </Container>
  );
};

export default withRouter(WorkflowList);
