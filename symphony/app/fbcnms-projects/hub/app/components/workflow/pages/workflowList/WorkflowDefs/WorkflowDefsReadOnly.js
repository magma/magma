import React, { Component } from "react";
import {
  Accordion,
  Button,
  Card,
  Col,
  Container,
  Form,
  Row,
  Table,
} from "react-bootstrap";
import { WorkflowDefs } from "./WorkflowDefs";
import { Typeahead } from "react-bootstrap-typeahead";
import "react-bootstrap-typeahead/css/Typeahead.css";
import { withRouter } from "react-router-dom";
import PageCount from "../../../common/PageCount";
import PageSelect from "../../../common/PageSelect";
import WfLabels from "../../../common/WfLabels";
import DefinitionModal from "./DefinitonModal/DefinitionModal";
import DiagramModal from "./DiagramModal/DiagramModal";
import InputModal from "./InputModal/InputModal";
import { HttpClient as http } from "../../../common/HttpClient";
import { conductorApiUrlPrefix, frontendUrlPrefix } from "../../../constants";

class WorkflowDefsReadOnly extends WorkflowDefs {
  constructor(props) {
    super(props);
  }

  repeatDeleteButton() {
    return (null);
  }

  repeatFavouriteButton(dataset) {
    return (null);
  }
  
  repeatScheduleButton(dataset) {
    return (null);
  }

  repeatEditButton() {
    return (null);
  }

  render() {
    return (
      <div>
        {this.renderDefinitionModal()}
        {this.renderInputModal()}
        {this.renderDiagramModal()}
        <Row>
          {this.renderSearchByLabel()}
          {this.renderSearchByKeyword()}
        </Row>
        {this.renderWorkflowTable()}
        {this.renderWorkflowTableFooter()}
      </div>
    );
  }
}

export default withRouter(WorkflowDefsReadOnly);
