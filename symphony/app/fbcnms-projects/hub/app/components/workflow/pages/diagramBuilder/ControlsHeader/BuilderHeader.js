import React, { useEffect } from "react";
import logo from "./logo-min.png";
import x from "./X_icon_RGB-min.png";
import { Navbar } from "react-bootstrap";
import { Button, Dropdown, Icon, Popup } from "semantic-ui-react";
import { NavLink, withRouter } from "react-router-dom";
import { motion } from "framer-motion";
import "./BuilderHeader.css";

const XLetter = () => (
  <motion.img
    key={x}
    src={x}
    initial={{ opacity: 1, x: 82, scale: 0.9 }}
    animate={{
      x: 0,
      scale: 1.1,
      rotate: 180
    }}
    transition={{ duration: 0.5, delay: 1.1 }}
  />
);

const Logo = () => (
  <motion.img
    key={logo}
    src={logo}
    initial={{ opacity: 1, x: -30 }}
    animate={{
      opacity: 0
    }}
    transition={{ duration: 0.2, delay: 1 }}
  />
);

const Title = () => (
  <motion.div
    initial={{ opacity: 0, x: -120 }}
    animate={{
      opacity: 1,
      x: -80
    }}
    transition={{ duration: 0.5, delay: 1.3 }}
  >
    <h3>Workflow Builder</h3>
  </motion.div>
);

const ControlsButton = props => (
  <motion.div
    initial={{ opacity: 0 }}
    animate={{
      opacity: 1
    }}
    transition={{ duration: 1, delay: 1 }}
  >
    <Button.Group basic inverted size="small" style={{ marginRight: "20px" }}>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{
          opacity: 1
        }}
        transition={{ duration: 0.8, delay: 1 }}
      >
        <Popup
          trigger={<Button onClick={props.showNewModal} icon="file" />}
          header="New"
          content="Create new workflow."
          basic
        />
        <Popup
          trigger={<Button onClick={props.saveWorkflow} icon="save" />}
          header={
            <h4>
              Save <kbd>Ctrl</kbd>+<kbd>S</kbd>
            </h4>
          }
          content="Save the workflow into database."
          basic
        />
        <Popup
          trigger={<Button onClick={props.openFileUpload} icon="download" />}
          header="Import"
          content="Import workflow on canvas (must be in valid JSON format)."
          basic
        />
        <Popup
          trigger={<Button onClick={props.saveFile} icon="upload" />}
          header="Export"
          content="Export and download workflow in JSON format."
          basic
        />
      </motion.div>
    </Button.Group>

    <Button.Group basic inverted size="small" style={{ marginRight: "20px" }}>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{
          opacity: 1
        }}
        transition={{ duration: 0.8, delay: 1.1 }}
      >
        <Button
          title="Zoom out"
          onClick={() =>
            props.setZoomLevel(props.workflowDiagram.getZoomLevel() - 10)
          }
          icon="zoom-out"
        />
        <Dropdown
          text={props.workflowDiagram.getZoomLevel().toFixed(1) + "%"}
          button
          style={{ paddingLeft: "10px", paddingRight: "10px" }}
        >
          <Dropdown.Menu>
            {[25, 50, 75, 100, 125].map(level => {
              return (
                <Dropdown.Item onClick={() => props.setZoomLevel(level)}>
                  <span className="text">{level}</span>
                </Dropdown.Item>
              );
            })}
          </Dropdown.Menu>
        </Dropdown>
        <Button
          title="Zoom in"
          onClick={() =>
            props.setZoomLevel(props.workflowDiagram.getZoomLevel() + 10)
          }
          icon="zoom-in"
        />
      </motion.div>
    </Button.Group>

    <Button.Group basic inverted size="small" style={{ marginRight: "20px" }}>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{
          opacity: 1
        }}
        transition={{ duration: 0.8, delay: 1.2 }}
      >
        <Popup
          trigger={
            <Button
              id="expand"
              onClick={props.expandNodeToWorkflow}
              icon="expand"
            />
          }
          header={
            <h4>
              Expand <kbd>Ctrl</kbd>+<kbd>X</kbd>
            </h4>
          }
          content="Expand selected sub-workflows."
          basic
        />
        <Popup
          trigger={
            <Button
              id="delete"
              onClick={() => props.workflowDiagram.deleteSelected()}
              icon="trash"
            />
          }
          header={
            <h4>
              Delete <kbd>LMB</kbd>+<kbd>Delete</kbd>
            </h4>
          }
          content="Delete selected nodes."
          basic
        />
        <Popup
          trigger={
            <Button
              onClick={props.setLocked}
              icon={props.workflowDiagram.isLocked() ? <Icon className="lock" style={{color: "red"}}/>  : "unlock"}
            />
          }
          header={
            props.workflowDiagram.isLocked() ? (
              <h4>
                Unlock diagram <kbd>Ctrl</kbd>+<kbd>L</kbd>
              </h4>
            ) : (
              <h4>
                Lock diagram <kbd>Ctrl</kbd>+<kbd>L</kbd>
              </h4>
            )
          }
          content={
            props.workflowDiagram.isLocked()
              ? "Unlock diagram state (currently locked)"
              : "Lock diagram state (currently unlocked)"
          }
          basic
        />
      </motion.div>
    </Button.Group>

    <Button.Group basic inverted size="small" style={{ marginRight: "20px" }}>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{
          opacity: 1
        }}
        transition={{ duration: 0.8, delay: 1.3 }}
      >
        <Popup
          trigger={<Button onClick={props.showGeneralInfoModal} icon="edit" />}
          header="Edit general"
          content="Edit workflow general information."
          basic
        />
        <Popup
          trigger={
            <Button onClick={props.showDefinitionModal} icon="file code" />
          }
          header="Show Definition"
          content="Show workflow current definition (JSON)."
          basic
        />
      </motion.div>
    </Button.Group>

    <Button.Group size="small" style={{ marginRight: "20px" }}>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{
          opacity: 1
        }}
        transition={{ duration: 0.8, delay: 1.4 }}
      >
        <Popup
          trigger={<Button icon="play" onClick={props.saveAndExecute} />}
          header={
            <h4>
              Save & Execute <kbd>Alt</kbd>+<kbd>Enter</kbd>
            </h4>
          }
          content="Save and execute current workflow (opens input dialog)."
          basic
        />
      </motion.div>
    </Button.Group>
  </motion.div>
);

const BuilderHeader = props => {
  useEffect(() => {
    document.addEventListener("click", handleClickInside, true);
    return () => {
      document.removeEventListener("click", handleClickInside, true);
    };
  }, []);

  const handleClickInside = event => {
    const headerEl = document.getElementById("builder-header");
    const sideMenu = document.getElementById("sidebar-secondary");
    const expandBtn = document.getElementById("expand");
    const deleteBtn = document.getElementById("delete");

    // workaround to prevent deleting nodes while typing (e.g. pressing delete)
    // focus on node is lost when sidebar or header is clicked
    if (
      headerEl &&
      sideMenu &&
      (headerEl.contains(event.target) || sideMenu.contains(event.target)) &&
      (!expandBtn.contains(event.target) && !deleteBtn.contains(event.target))
    ) {
      props.workflowDiagram.getDiagramModel().clearSelection();
      props.workflowDiagram.renderDiagram();
    }
  };

  const openFileUpload = () => {
    document.getElementById("upload-file").click();
    document
      .getElementById("upload-file")
      .addEventListener("change", props.submitFile);
  };

  return (
    <Navbar id="builder-header" className="builder-header">
      <Navbar.Brand>
        <NavLink to="/">
          <XLetter />
          <Logo />
        </NavLink>
      </Navbar.Brand>
      <Title />
      <ControlsButton
        openFileUpload={openFileUpload}
        saveWorkflow={props.saveWorkflow}
        saveFile={props.saveFile}
        showExitModal={props.showExitModal}
        expandNodeToWorkflow={props.expandNodeToWorkflow}
        showGeneralInfoModal={props.showGeneralInfoModal}
        showDefinitionModal={props.showDefinitionModal}
        saveAndExecute={props.saveAndExecute}
        showNewModal={props.showNewModal}
        workflowDiagram={props.workflowDiagram}
        setZoomLevel={props.setZoomLevel}
        setLocked={props.setLocked}
      />
      <input id="upload-file" type="file" hidden />
      <Navbar.Collapse className="justify-content-end">
        <Button
          basic
          inverted
          animated="vertical"
          onClick={props.showExitModal}
        >
          <Button.Content hidden>Exit</Button.Content>
          <Button.Content visible>
            <Icon name="x" />
          </Button.Content>
        </Button>
      </Navbar.Collapse>
    </Navbar>
  );
};

export default withRouter(BuilderHeader);
