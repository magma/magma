import * as React from "react";
import { DefaultLinkFactory } from "@projectstorm/react-diagrams";
import { LinkWithContextWidget } from "./LinkWithContextWidget";

export class LinkWithContextFactory extends DefaultLinkFactory {
  generateReactWidget(diagramEngine, link) {
    return React.createElement(LinkWithContextWidget, {
      diagramEngine: diagramEngine,
      link: link,
    });
  }
}
