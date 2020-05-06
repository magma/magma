import * as React from "react";
import { DefaultLinkWidget } from "@projectstorm/react-diagrams";
import { LinkContextMenu, LinkMenuProvider } from "./ContextMenu";

export class LinkWithContextWidget extends DefaultLinkWidget {
  render() {
    return (
      <g>
        <LinkMenuProvider link={this.props.link}>{super.render()}</LinkMenuProvider>
        <LinkContextMenu
          link={this.props.link}
          diagramEngine={this.props.diagramEngine}
        />
      </g>
    );
  }
}
