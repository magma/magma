import React from "react";
import { Container, Tab, Tabs } from "react-bootstrap";
import { withRouter } from "react-router-dom";
import WorkflowDefs from "./WorkflowDefs/WorkflowDefs";
import WorkflowExec from "./WorkflowExec/WorkflowExec";
import {changeUrl, exportButton} from './workflowUtils'

const WorkflowListReadOnly = (props) => {
  let urlUpdater = changeUrl(props.history);
  let query = props.match.params.wfid ? props.match.params.wfid : null;

  return (
      <Container style={{ textAlign: "left", marginTop: "20px" }}>
        <h1 style={{ marginBottom: "20px" }}>
          <i style={{ color: "grey" }} className="fas fa-cogs" />
          &nbsp;&nbsp;Workflows
          {exportButton()}
        </h1>
        <input id="upload-files" multiple type="file" hidden />
        <Tabs
            onSelect={(e) => urlUpdater(e)}
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

export default withRouter(WorkflowListReadOnly);
