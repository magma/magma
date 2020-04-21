import { saveAs } from "file-saver";
import * as _ from "lodash";
import React, { Component } from "react";
import { Button, Modal } from "react-bootstrap";
import { HotKeys } from "react-hotkeys";
import { connect } from "react-redux";
import "semantic-ui-css/semantic.min.css";
import { DiagramWidget, Toolkit } from "storm-react-diagrams";
import * as builderActions from "../../store/actions/builder";
import InputModal from "../workflowList/WorkflowDefs/InputModal/InputModal";
import DetailsModal from "../workflowList/WorkflowExec/DetailsModal/DetailsModal";
import { Application } from "./Application";
import { encode } from "./builder-utils";
import BuilderHeader from "./ControlsHeader/BuilderHeader";
import CustomAlert from "./CustomAlert";
import "./DiagramBuilder.css";
import GeneralInfoModal from "./GeneralInfoModal/GeneralInfoModal";
import NodeModal from "./NodeModal/NodeModal";
import Sidemenu from "./Sidemenu/Sidemenu";
import SidemenuRight from "./Sidemenu/SidemenuRight";
import WorkflowDefModal from "./WorkflowDefModal/WorkflowDefModal";
import { WorkflowDiagram } from "./WorkflowDiagram";
import { HttpClient as http } from "../../common/HttpClient";
import { conductorApiUrlPrefix, frontendUrlPrefix } from "../../constants";

class DiagramBuilder extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showNodeModal: false,
      showDefinitionModal: false,
      showGeneralInfoModal: false,
      showInputModal: false,
      showDetailsModal: false,
      showExitModal: false,
      showNewModal: false,
      modalInputs: null,
      zoomLevel: 100,
      isLocked: false,
      workflowDiagram: new WorkflowDiagram(
        new Application(),
        this.props.finalWorkflow,
        { x: 600, y: 300 }
      ),
    };

    this.setLocked = this.setLocked.bind(this);
    this.onNodeDrop = this.onNodeDrop.bind(this);
    this.importFile = this.importFile.bind(this);
    this.exportFile = this.exportFile.bind(this);
    this.setZoomLevel = this.setZoomLevel.bind(this);
    this.showNewModal = this.showNewModal.bind(this);
    this.saveWorkflow = this.saveWorkflow.bind(this);
    this.showExitModal = this.showExitModal.bind(this);
    this.showNodeModal = this.showNodeModal.bind(this);
    this.redirectOnNew = this.redirectOnNew.bind(this);
    this.showInputModal = this.showInputModal.bind(this);
    this.saveAndExecute = this.saveAndExecute.bind(this);
    this.redirectOnExit = this.redirectOnExit.bind(this);
    this.closeInputModal = this.closeInputModal.bind(this);
    this.showDetailsModal = this.showDetailsModal.bind(this);
    this.parseDiagramToJSON = this.parseDiagramToJSON.bind(this);
    this.showDefinitionModal = this.showDefinitionModal.bind(this);
    this.expandNodeToWorkflow = this.expandNodeToWorkflow.bind(this);
    this.showGeneralInfoModal = this.showGeneralInfoModal.bind(this);
    this.saveNodeInputsHandler = this.saveNodeInputsHandler.bind(this);
    this.createDiagramByDefinition = this.createDiagramByDefinition.bind(this);
  }

  componentDidMount() {
    document.addEventListener("dblclick", this.doubleClickListener.bind(this));

    console.log(this.props);
    this.props.hideHeader();

    http.get(conductorApiUrlPrefix + "/metadata/workflow").then((res) => {
      this.props.storeWorkflows(
        res.result?.sort((a, b) => a.name.localeCompare(b.name)) || []
      );
    });

    http.get(conductorApiUrlPrefix + "/metadata/taskdefs").then((res) => {
      this.props.storeTasks(
        res.result?.sort((a, b) => a.name.localeCompare(b.name)) || []
      );
    });

    if (!_.isEmpty(this.props.match.params)) {
      this.createExistingWorkflow();
    } else {
      this.createNewWorkflow();
    }
  }

  componentWillUnmount() {
    this.props.resetToDefaultWorkflow();
  }

  createNewWorkflow() {
    this.setState({ showGeneralInfoModal: true });
    this.state.workflowDiagram.placeDefaultNodes();
    this.props.showCustomAlert(
      true,
      "primary",
      "Start to drag & drop tasks from left menu on canvas."
    );
  }

  createExistingWorkflow() {
    const { name, version } = this.props.match.params;
    http
      .get(conductorApiUrlPrefix + "/metadata/workflow/" + name + "/" + version)
      .then((res) => {
        this.createDiagramByDefinition(res.result);
      })
      .catch(() => {
        return this.props.showCustomAlert(
          true,
          "danger",
          `Cannot find selected sub-workflow: ${name}.`
        );
      });
  }

  createDiagramByDefinition(definition) {
    this.props.updateFinalWorkflow(definition);
    this.props.showCustomAlert(
      true,
      "info",
      `Editing workflow ${definition.name} / ${definition.version}.`
    );
    this.props.lockWorkflowName();

    this.state.workflowDiagram
      .setDefinition(definition)
      .createDiagram()
      .withStartEnd()
      .renderDiagram();
  }

  onNodeDrop(e) {
    this.props.showCustomAlert(false);
    this.state.workflowDiagram.dropNewNode(e);
  }

  doubleClickListener(event) {
    let diagramModel = this.state.workflowDiagram.getDiagramModel();
    let element = Toolkit.closest(event.target, ".node[data-nodeid]");
    let node = null;

    if (element) {
      node = diagramModel.getNode(element.getAttribute("data-nodeid"));
      if (node && node.type !== "start" && node.type !== "end") {
        node.setSelected(false);
        this.setState({
          showNodeModal: true,
          modalInputs: { inputs: node.extras.inputs, id: node.id },
        });
      }
    }
  }

  parseDiagramToJSON() {
    try {
      this.props.showCustomAlert(false);
      const finalWf = this.state.workflowDiagram.parseDiagramToJSON(
        this.props.finalWorkflow
      );
      this.props.updateFinalWorkflow(finalWf);
      return finalWf;
    } catch (e) {
      this.props.showCustomAlert(true, "danger", e.message);
    }
  }

  expandNodeToWorkflow(e) {
    e.preventDefault();
    try {
      this.props.showCustomAlert(false);
      this.state.workflowDiagram.expandSelectedNodes();
    } catch (e) {
      this.props.showCustomAlert(true, "danger", e.message);
    }
  }

  saveWorkflow(e) {
    e.preventDefault();
    this.state.workflowDiagram
      .saveWorkflow(this.props.finalWorkflow)
      .then((res) => {
        this.props.showCustomAlert(
          true,
          "info",
          `Workflow ${res.name} saved successfully.`
        );
      })
      .catch((e) => {
        this.props.showCustomAlert(
          true,
          "danger",
          e.path + ":\xa0\xa0\xa0" + e.message
        );
      });
  }

  saveAndExecute(e) {
    e.preventDefault();
    this.state.workflowDiagram
      .saveWorkflow(this.props.finalWorkflow)
      .then(() => {
        this.showInputModal();
      })
      .catch((e) => {
        this.props.showCustomAlert(
          true,
          "danger",
          e.path + ":\xa0\xa0\xa0" + e.message
        );
      });
  }

  /*** MODAL HANDLERS ***/

  // NODE MODAL
  showNodeModal() {
    this.setState({
      showNodeModal: !this.state.showNodeModal,
    });
  }

  // DEFINITION MODAL
  showDefinitionModal() {
    this.parseDiagramToJSON();
    this.setState({
      showDefinitionModal: !this.state.showDefinitionModal,
    });
  }

  // GENERAL INFO MODAL
  showGeneralInfoModal() {
    this.setState({
      showGeneralInfoModal: !this.state.showGeneralInfoModal,
    });
  }

  // WORKFLOW EXECUTION INPUT MODAL
  showInputModal() {
    this.setState({
      showInputModal: true,
    });
  }

  closeInputModal() {
    this.setState({
      showInputModal: false,
    });
    this.showDetailsModal();
  }

  // WORKFLOW EXECUTION DETAILS MODAL
  showDetailsModal() {
    this.setState({
      showDetailsModal: !this.state.showDetailsModal,
    });
  }

  // EXIT MODAL
  showExitModal() {
    this.setState({
      showExitModal: !this.state.showExitModal,
    });
  }

  showNewModal(e) {
    e.preventDefault();
    this.setState({
      showNewModal: !this.state.showNewModal,
    });
  }
  /*************** ***************/

  saveNodeInputsHandler(savedInputs, id) {
    let nodes = this.state.workflowDiagram.getNodes();

    nodes.forEach((node) => {
      if (node.id === id) {
        node.extras.inputs = savedInputs;
      }
    });
    this.parseDiagramToJSON();
  }

  importFile() {
    const fileToLoad = document.getElementById("upload-file").files[0];
    const fileReader = new FileReader();

    fileReader.onload = (() => {
      return function(e) {
        try {
          let jsonObj = JSON.parse(e.target.result);
          this.createDiagramByDefinition(jsonObj);
        } catch (err) {
          alert("Error when trying to parse json." + err);
        }
      };
    })(fileToLoad).bind(this);
    fileReader.readAsText(fileToLoad);
  }

  exportFile() {
    const definition = this.parseDiagramToJSON();

    if (!definition) {
      return null;
    }

    const data = encode(JSON.stringify(definition, null, 2));
    const file = new Blob([data], {
      type: "application/octet-stream",
    });

    saveAs(file, definition.name + ".json");
  }

  redirectOnExit() {
    this.props.history.push(frontendUrlPrefix + "/defs");
    window.location.reload();
  }

  redirectOnNew() {
    this.props.history.push(frontendUrlPrefix + "/builder");
    window.location.reload();
  }

  setZoomLevel(percentage, e) {
    if (e) {
      e.preventDefault();
    }
    this.setState({
      zoomLevel: percentage,
    });
    this.state.workflowDiagram.setZoomLevel(percentage);
  }

  setLocked(e) {
    if (e) {
      e.preventDefault();
    }
    this.setState({
      isLocked: !this.state.workflowDiagram.isLocked(),
    });
    this.state.workflowDiagram.setLocked(
      !this.state.workflowDiagram.isLocked()
    );
  }

  render() {
    let inputsModal = this.state.showInputModal ? (
      <InputModal
        wf={
          this.props.finalWorkflow.name +
          " / " +
          this.props.finalWorkflow.version
        }
        modalHandler={this.closeInputModal}
        fromBuilder
        show={this.state.showInputModal}
      />
    ) : null;

    let detailsModal = this.state.showDetailsModal ? (
      <DetailsModal
        wfId={this.props.workflowId}
        modalHandler={this.showDetailsModal}
      />
    ) : null;

    let nodeModal = this.state.showNodeModal ? (
      <NodeModal
        modalHandler={this.showNodeModal}
        inputs={this.state.modalInputs}
        saveInputs={this.saveNodeInputsHandler}
        show={this.state.showNodeModal}
      />
    ) : null;

    let generalInfoModal = this.state.showGeneralInfoModal ? (
      <GeneralInfoModal
        finalWorkflow={this.props.finalWorkflow}
        workflows={this.props.workflows}
        closeModal={this.showGeneralInfoModal}
        saveInputs={this.props.updateFinalWorkflow}
        show={this.state.showGeneralInfoModal}
        lockWorkflowName={this.props.lockWorkflowName}
        isWfNameLocked={this.props.isWfNameLocked}
        redirectOnExit={this.redirectOnExit}
      />
    ) : null;

    let workflowDefModal = this.state.showDefinitionModal ? (
      <WorkflowDefModal
        definition={this.props.finalWorkflow}
        closeModal={this.showDefinitionModal}
        show={this.state.showDefinitionModal}
      />
    ) : null;

    let exitModal = this.state.showExitModal ? (
      <Modal show={this.state.showExitModal}>
        <Modal.Header>
          <Modal.Title>Do you want to exit builder?</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          All changes since last <b>Save</b> or <b>Execute</b> operation will be
          lost
        </Modal.Body>
        <Modal.Footer>
          <Button variant="outline-primary" onClick={this.showExitModal}>
            Cancel
          </Button>
          <Button variant="danger" onClick={this.redirectOnExit}>
            Exit
          </Button>
        </Modal.Footer>
      </Modal>
    ) : null;

    let newModal = this.state.showNewModal ? (
      <Modal show={this.state.showNewModal}>
        <Modal.Header>
          <Modal.Title>Create new workflow</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          All changes since last <b>Save</b> or <b>Execute</b> operation will be
          lost
        </Modal.Body>
        <Modal.Footer>
          <Button variant="outline-primary" onClick={this.showNewModal}>
            Cancel
          </Button>
          <Button variant="primary" onClick={this.redirectOnNew}>
            New
          </Button>
        </Modal.Footer>
      </Modal>
    ) : null;

    const keyMap = {
      ZOOM_IN: ["ctrl++"],
      ZOOM_OUT: ["ctrl+-"],
      LOCK: ["ctrl+l"],
      SAVE: ["ctrl+s"],
      EXECUTE: ["alt+enter"],
      EXPAND: ["ctrl+x"],
    };

    const handlers = {
      ZOOM_IN: (e) => this.setZoomLevel(this.state.zoomLevel + 10, e),
      ZOOM_OUT: (e) => this.setZoomLevel(this.state.zoomLevel - 10, e),
      LOCK: (e) => this.setLocked(e),
      SAVE: (e) => this.saveWorkflow(e),
      EXECUTE: (e) => this.saveAndExecute(e),
      EXPAND: (e) => this.expandNodeToWorkflow(e),
    };

    return (
      <HotKeys keyMap={keyMap}>
        <HotKeys handlers={handlers}>
          <div style={{ position: "relative", height: "100vh" }}>
            {workflowDefModal}
            {nodeModal}
            {inputsModal}
            {detailsModal}
            {generalInfoModal}
            {exitModal}
            {newModal}

            <BuilderHeader
              showDefinitionModal={this.showDefinitionModal}
              showGeneralInfoModal={this.showGeneralInfoModal}
              showExitModal={this.showExitModal}
              showNewModal={this.showNewModal}
              saveAndExecute={this.saveAndExecute}
              saveWorkflow={this.saveWorkflow}
              expandNodeToWorkflow={this.expandNodeToWorkflow}
              updateQuery={this.props.updateQuery}
              submitFile={this.importFile}
              saveFile={this.exportFile}
              workflowDiagram={this.state.workflowDiagram}
              setZoomLevel={this.setZoomLevel}
              setLocked={this.setLocked}
            />

            <Sidemenu
              workflows={this.props.workflows}
              tasks={this.props.tasks}
              updateQuery={this.props.updateQuery}
              openCard={this.props.openCard}
            />
            <SidemenuRight functional={this.props.functional} />

            <CustomAlert
              showCustomAlert={this.props.showCustomAlert}
              show={this.props.customAlert.show}
              msg={this.props.customAlert.msg}
              alertVariant={this.props.customAlert.variant}
            />

            <div
              style={{ height: "calc(100% - 50px)" }}
              onDrop={(e) => this.onNodeDrop(e)}
              onDragOver={(event) => {
                event.preventDefault();
              }}
            >
              <DiagramWidget
                className="srd-demo-canvas"
                diagramEngine={this.state.workflowDiagram.getDiagramEngine()}
              />
            </div>
          </div>
        </HotKeys>
      </HotKeys>
    );
  }
}

const mapStateToProps = (state) => {
  return {
    workflows: state.buildReducer.workflows,
    tasks: state.buildReducer.tasks,
    functional: state.buildReducer.functional,
    finalWorkflow: state.buildReducer.finalWorkflow,
    customAlert: state.buildReducer.customAlert,
    isWfNameLocked: state.buildReducer.workflowNameLock,
    workflowId: state.buildReducer.executedWfId,
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    storeWorkflows: (wfList) => dispatch(builderActions.storeWorkflows(wfList)),
    storeTasks: (taskList) => dispatch(builderActions.storeTasks(taskList)),
    updateFinalWorkflow: (finalWorkflow) =>
      dispatch(builderActions.updateFinalWorkflow(finalWorkflow)),
    resetToDefaultWorkflow: () =>
      dispatch(builderActions.resetToDefaultWorkflow()),
    updateQuery: (query, labels) =>
      dispatch(builderActions.requestUpdateByQuery(query, labels)),
    openCard: (which) => dispatch(builderActions.openCard(which)),
    showCustomAlert: (show, variant, msg) =>
      dispatch(builderActions.showCustomAlert(show, variant, msg)),
    lockWorkflowName: () => dispatch(builderActions.lockWorkflowName()),
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(DiagramBuilder);
