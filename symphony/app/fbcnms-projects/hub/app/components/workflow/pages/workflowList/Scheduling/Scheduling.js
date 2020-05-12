import React, { Component } from "react";
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
import WfLabels from "../../../common/WfLabels";
import { HttpClient as http } from "../../../common/HttpClient";
import { conductorApiUrlPrefix, frontendUrlPrefix } from "../../../constants";

class Scheduling extends Component {
  constructor(props) {
    super(props);
    this.state = {
      keywords: "",
      labels: [],
      data: [],
      table: [],
      activeRow: null,
      activeWf: null,
      defaultPages: 20,
      pagesCount: 1,
      viewedPage: 1,
      allLabels: []
    };
    this.table = React.createRef();
  }

  componentWillMount() {

  }

  componentDidMount() {
    http.get(conductorApiUrlPrefix + "/schedule").then(res => {
      if (res) {
        let size = ~~(res.length / this.state.defaultPages);
        let dataset =
          res.sort((a, b) =>
            a.name > b.name ? 1 : b.name > a.name ? -1 : 0
          ) || [];
        this.setState({
          data: dataset,
          pagesCount:
            res.length % this.state.defaultPages ? ++size : size
        });
      }
    });
  }

  changeActiveRow(i) {
    let dataset = this.state.data;
    this.setState({
      activeRow: this.state.activeRow === i ? null : i,
      activeWf: dataset[i]["name"] + " / " + dataset[i]["cronString"]
    });
  }


  setCountPages(defaultPages, pagesCount) {
    this.setState({
      defaultPages: defaultPages,
      pagesCount: pagesCount,
      viewedPage: 1
    });
  }

  setViewPage(page) {
    this.setState({
      viewedPage: page
    });
  }

  delete(schedulingEntry) {
    console.log("Deleting", schedulingEntry.name);
    http
      .delete(conductorApiUrlPrefix + "/schedule/" + schedulingEntry.name)
      .then(() => {
        this.componentDidMount();
        let table = this.state.table;
        // if (table.length) {
        //   table.splice(table.findIndex(wf => wf.name === name), 1);
        // }
        this.setState({
          activeRow: null,
          table: table
        });
      });
  }

  repeat() {
    let output = [];
    let defaultPages = this.state.defaultPages;
    let viewedPage = this.state.viewedPage;
    let dataset = this.state.data;
    for (let i = 0; i < dataset.length; i++) {
      if (
        i >= (viewedPage - 1) * defaultPages &&
        i < viewedPage * defaultPages
      ) {
        output.push(
          <div className="wfRow" key={i}>
            <Accordion.Toggle
              id={`wf${i}`}
              onClick={this.changeActiveRow.bind(this, i)}
              className="clickable wfDef"
              as={Card.Header}
              variant="link"
              eventKey={i}
            >
              <b>{dataset[i]["name"]}</b>
              <br />
              <div className="description">
                { dataset[i]["cronString"] }
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
                    onClick={() => alert('Not implemented')}
                  >
                    Edit
                  </Button>

                  <Button
                    variant="outline-danger noshadow"
                    style={{ float: "right" }}
                    onClick={this.delete.bind(this, dataset[i])}
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
  }


  render() {
    return (
      <div>
        <div className="scrollWrapper" style={{ maxHeight: "650px" }}>
          <Table ref={this.table}>
            <thead>
              <tr>
                <th>Name/Cron</th>
              </tr>
            </thead>
            <tbody>
              <Accordion activeKey={this.state.activeRow}>
                {this.repeat()}
              </Accordion>
            </tbody>
          </Table>
        </div>
        <Container style={{ marginTop: "5px" }}>
          <Row>
            <Col sm={2}>
              <PageCount
                dataSize={
                  this.state.data.length
                }
                defaultPages={this.state.defaultPages}
                handler={this.setCountPages.bind(this)}
              />
            </Col>
            <Col sm={8} />
            <Col sm={2}>
              <PageSelect
                viewedPage={this.state.viewedPage}
                count={this.state.pagesCount}
                handler={this.setViewPage.bind(this)}
              />
            </Col>
          </Row>
        </Container>
      </div>
    );
  }
}

export default withRouter(Scheduling);
