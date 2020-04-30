import defaultTo from "lodash/fp/defaultTo";
import React, { Component } from "react";
import { Col, Row } from "react-bootstrap";
import Grapher from "../../../../../common/Grapher";
import Workflow2Graph from "../../../../../common/wfegraph";

class WorkflowDia extends Component {
  constructor(props) {
    super(props);

    this.state = WorkflowDia.getGraphState(props);
  }

  static getGraphState(props) {
    const wfe2graph = new Workflow2Graph();
    const subwfs = defaultTo({})(props.subworkflows);
    const wfe = defaultTo({ tasks: [] })(props.wfe);
    const { edges, vertices } = wfe2graph.convert(wfe, props.meta);
    const subworkflows = {};

    for (const refname in subwfs) {
      let submeta = subwfs[refname].meta;
      let subwfe = subwfs[refname].wfe;
      subworkflows[refname] = wfe2graph.convert(subwfe, submeta);
    }

    return { edges, vertices, subworkflows };
  }

  componentWillReceiveProps(nextProps) {
    this.setState(WorkflowDia.getGraphState(nextProps));
  }

  render() {
    const { edges, vertices, subworkflows } = this.state;

    return (
      <div style={{ overflow: "scroll" }}>
        {!this.props.def ? (
          <div>
            <Row style={{ textAlign: "center" }}>
              <Col>
                <h2>Execution Flow</h2>
              </Col>
            </Row>
            <hr />
          </div>
        ) : null}

        <Grapher
          def={this.props.def}
          edges={edges}
          vertices={vertices}
          layout="TD-auto"
          innerGraph={subworkflows}
        />
      </div>
    );
  }
}

export default WorkflowDia;
