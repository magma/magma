import React, { Component } from "react";
import Pagination from "react-bootstrap/Pagination";

type Props = {
  viewedPage: number,
  defaultPages: number,
  count: number,
  indent: number,
  handler: (defaultPageCount: number, pageCount?: number) => void,
  dataSize: number
};

type StateType = {};

class PageSelect extends Component<Props, StateType> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  setPages() {
    let output = [];
    let viewedPage = this.props.viewedPage;
    let pagesCount = this.props.count;
    let indent = this.props.indent ? this.props.indent : 2;
    output.push(
      <Pagination.Prev
        key={"prev"}
        disabled={viewedPage === 1 || pagesCount === 0}
        onClick={e => {
          if (viewedPage !== 1 && pagesCount !== 0)
            this.props.handler(viewedPage - 1);
        }}
      />
    );
    for (let i = 1; i <= pagesCount; i++) {
      if (i >= viewedPage - indent && i <= viewedPage + indent) {
        output.push(
          <Pagination.Item
            key={i}
            active={viewedPage === i}
            onClick={() => this.props.handler(i)}
          >
            {" "}
            {i}{" "}
          </Pagination.Item>
        );
      }
    }
    output.push(
      <Pagination.Next
        key={"next"}
        disabled={viewedPage === pagesCount || pagesCount === 0}
        onClick={() => {
          if (viewedPage !== pagesCount && pagesCount !== 0) {
            this.props.handler(viewedPage + 1);
          }
        }}
      />
    );
    return output;
  }
  render() {
    return (
      <Pagination style={{ float: "right" }}>{this.setPages()}</Pagination>
    );
  }
}

export default PageSelect;
