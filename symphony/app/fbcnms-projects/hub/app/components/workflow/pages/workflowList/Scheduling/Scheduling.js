import React, { useEffect, useState } from "react";
import {
  Accordion,
  Button,
  Card,
  Col,
  Container,
  Form,
  Row,
  Table
} from "react-bootstrap";
import { Typeahead } from "react-bootstrap-typeahead";
import "react-bootstrap-typeahead/css/Typeahead.css";
import { withRouter } from "react-router-dom";
import PageCount from "../../../common/PageCount";
import PageSelect from "../../../common/PageSelect";
import { HttpClient as http } from "../../../common/HttpClient";
import { conductorApiUrlPrefix, frontendUrlPrefix } from "../../../constants";
import SchedulingModal from "./SchedulingModal/SchedulingModal";
import DiagramModal from "../WorkflowDefs/DiagramModal/DiagramModal";

const Scheduling = props => {
  const [showSchedulingModal, setShowSchedulingModal] = useState(false);
  const [activeScheduleName, setActiveScheduleName] = useState();
  const [activeRow, setActiveRow] = useState();

  const refresh = () => {
    http.get(conductorApiUrlPrefix + "/schedule").then(res => {
      let size = Math.floor(res.length / defaultPages);
      let dataset = res.sort((a, b) =>
          a.name > b.name ? 1 : b.name > a.name ? -1 : 0
        ) || [];
      setData(dataset);
      setPagesCount(res.length % defaultPages ? ++size : size);
      deselectActiveRow();
      return dataset;
    });
    return [];
  };

  const [data, setData] = useState(refresh);
  const [pagesCount, setPagesCount] = useState(1);
  const [defaultPages, setDefaultPages] = useState(20);
  const [viewedPage, setViewedPage] = useState(1);


  const deselectActiveRow = () => {
    setActiveRow(null);
    setActiveScheduleName(null);
  }

  const changeActiveRow = (i) => {
    const deselectingCurrentRow = activeRow === i;
    if (deselectingCurrentRow) {
      deselectActiveRow();
    } else {
      setActiveRow(i);
      setActiveScheduleName(data[i]["name"]);
    }
  };

  const setCountPages = (defaultPages, pagesCount) => {
    setDefaultPages(defaultPages);
    setPagesCount(pagesCount);
    setViewedPage(1);
  };

  const deleteEntry = (schedulingEntry) => {
    console.log("Deleting", schedulingEntry.name);
    http
      .delete(conductorApiUrlPrefix + "/schedule/" + schedulingEntry.name)
      .then(() => {
        deselectActiveRow();
      });
  };

  const repeat = () => {
    let output = [];
    for (let i = 0; i < data.length; i++) {
      if (
        i >= (viewedPage - 1) * defaultPages &&
        i < viewedPage * defaultPages
      ) {
        output.push(
          <div className="wfRow" key={i}>
            <Accordion.Toggle
              id={`wf${i}`}
              onClick={changeActiveRow.bind(this, i)}
              className="clickable wfDef"
              as={Card.Header}
              variant="link"
              eventKey={i}
            >
              <b>{data[i]["workflowName"]}</b> v.{data[i]["workflowVersion"]}
              <br />
              <div className="description">
                { data[i]["cronString"] }
              </div>
            </Accordion.Toggle>
            <Accordion.Collapse eventKey={i}>
              <Card.Body style={{ padding: "0px" }}>
                <div
                  style={{
                    background:
                      "linear-gradient(-120deg, rgb(0, 147, 255) 0%, rgb(0, 118, 203) 100%)",
                    padding: "15px",
                    marginBottom: "10px"
                  }}
                >
                  <Button
                    variant="outline-light noshadow"
                    onClick={flipShowSchedulingModal}
                  >
                    Edit
                  </Button>

                  <Button
                    variant="outline-danger noshadow"
                    style={{ float: "right" }}
                    onClick={deleteEntry.bind(this, data[i])}
                  >
                    <i className="fas fa-trash-alt" />
                  </Button>
                </div>
              </Card.Body>
            </Accordion.Collapse>
          </div>
        );
      }
    }
    return output;
  };

  const flipShowSchedulingModal = () => {
    setShowSchedulingModal(!showSchedulingModal);
  };

  const onModalClose = () => {
    flipShowSchedulingModal();
    refresh();
  }

  return (
    <div>
      {showSchedulingModal ? (
      <SchedulingModal
        name={activeScheduleName}
        onClose={onModalClose}
      />
      ) : null}

      <Button variant="outline-primary" style={{ marginLeft: "30px" }}
      onClick={() => refresh()}>
        <i className="fas fa-sync" />&nbsp;&nbsp;Refresh
      </Button>

      <div className="scrollWrapper" style={{ maxHeight: "650px" }}>
        <Table>
          <thead>
            <tr>
              <th>Name/Cron</th>
            </tr>
          </thead>
          <tbody>
            <Accordion activeKey={activeRow}>
              {repeat()}
            </Accordion>
          </tbody>
        </Table>
      </div>
      <Container style={{ marginTop: "5px" }}>
        <Row>
          <Col sm={2}>
            <PageCount
              dataSize={
                data.length
              }
              defaultPages={defaultPages}
              handler={setCountPages.bind(this)}
            />
          </Col>
          <Col sm={8} />
          <Col sm={2}>
            <PageSelect
              viewedPage={viewedPage}
              count={pagesCount}
              handler={setViewedPage}
            />
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default withRouter(Scheduling);
