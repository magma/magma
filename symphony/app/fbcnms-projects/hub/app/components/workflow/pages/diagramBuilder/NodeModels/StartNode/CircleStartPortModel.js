import * as _ from "lodash";
import {
  LinkModel,
  DiagramEngine,
  DefaultLinkModel,
  PortModel
} from "storm-react-diagrams";

export class CircleStartPortModel extends PortModel {
  position: string | "top" | "bottom" | "left" | "right";

  constructor(isInput: boolean, pos: string = "right") {
    super(pos, "start");
    this.in = isInput;
    this.position = pos;
  }

  serialize() {
    return _.merge(super.serialize(), {
      position: this.position
    });
  }

  link(port: PortModel): LinkModel {
    let link = this.createLinkModel();
    link.setSourcePort(this);
    link.setTargetPort(port);
    return link;
  }

  deSerialize(data: any, engine: DiagramEngine) {
    super.deSerialize(data, engine);
    this.position = data.position;
  }

  canLinkToPort(port: PortModel): boolean {
    return !this.in && port.in;
  }

  createLinkModel(): LinkModel {
    return new DefaultLinkModel();
  }
}
