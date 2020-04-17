import * as _ from "lodash";

export const getWfInputsRegex = wf => {
  let def = JSON.stringify(wf);
  let inputsArray = [...new Set(def.match(/(?<=workflow\.input\.)([a-zA-Z0-9-_]+)/gim))];
  let inputParameters = {};

  inputsArray.forEach(el => {
    inputParameters[el] = "${workflow.input." + el + "}";
  });

  return inputParameters;
};

export const getTaskInputsRegex = t => {
  let inputParameters = {};
  if (t.inputKeys) {
    t.inputKeys.forEach(el => {
      inputParameters[el] = "${workflow.input." + el + "}";
    });
  }

  return inputParameters;
};

export const hash = () =>
  Math.random()
    .toString(36)
    .toUpperCase()
    .substr(2, 4);

export const encode = s => {
  let out = [];
  for (let i = 0; i < s.length; i++) {
    out[i] = s.charCodeAt(i);
  }
  return new Uint8Array(out);
};

export const getLabelsFromString = str => {
  let labelsString = str
    .split("-")
    .pop()
    .replace(/ /g, "");
  return labelsString === "" ? [] : labelsString.split(",");
};

export const getWfInputs = wf => {
  let taskArray = wf.tasks;
  let inputParams = [];
  let inputParameters = {};

  taskArray.forEach(task => {
    if (task !== undefined) {
      let nonSystemTask = fn(task, "inputParameters");

      if (_.isArray(nonSystemTask)) {
        nonSystemTask.forEach(el => {
          if (el.inputParameters) {
            inputParams.push(el.inputParameters);
          }
        });
      } else if (nonSystemTask.inputParameters) {
        inputParams.push(task.inputParameters);
      }
    }
  });

  for (let i = 0; i < inputParams.length; i++) {
    inputParameters = { ...inputParameters, ...inputParams[i] };
  }

  return inputParameters;
};

// function to get nested key (inputParameters) from system tasks
export const fn = (obj, key) => {
  if (_.has(obj, key)) return obj;

  return _.flatten(
    _.map(obj, function(v) {
      return typeof v == "object" ? fn(v, key) : [];
    }),
    true
  );
};

export const getLinksArray = (type, node) => {
  let linksArray = [];
  _.values(node.ports).forEach(port => {
    if (type === "in" || type === "inputPort") {
      if (port.in || port.name === "left") {
        linksArray = _.values(port.links);
      }
    } else if (type === "out") {
      if (!port.in || port.name === "right") {
        linksArray = _.values(port.links);
      }
    }
  });
  return linksArray;
};

export const getStartNode = links => {
  for (let i = 0; i < _.values(links).length; i++) {
    let link = _.values(links)[i];
    if (link.sourcePort.type === "start") {
      return link.sourcePort.parent;
    }
  }
};

export const getEndNode = links => {
  for (let i = 0; i < _.values(links).length; i++) {
    let link = _.values(links)[i];
    if (link.targetPort.type === "end") {
      return link.targetPort.parent;
    }
  }
};

export const handleForkNode = forkNode => {
  let joinNode = null;
  let forkTasks = [];
  let joinOn = [];
  let forkBranches = forkNode.ports.right.links;

  //for each branch chain tasks
  _.values(forkBranches).forEach(link => {
    let tmpBranch = [];
    let parent = link.targetPort.getNode();
    let current = link.targetPort.getNode();

    //iterate trough tasks in each branch till join node
    while (current) {
      let outputLinks = getLinksArray("out", current);
      switch (current.type) {
        case "join":
          joinOn.push(parent.extras.inputs.taskReferenceName);
          joinNode = current;
          current = null;
          break;
        case "fork":
          let innerForkNode = handleForkNode(current).forkNode;
          let innerJoinNode = handleForkNode(current).joinNode;
          let innerJoinOutLinks = getLinksArray("out", innerJoinNode);
          tmpBranch.push(
            innerForkNode.extras.inputs,
            innerJoinNode.extras.inputs
          );
          parent = innerJoinNode;
          current = innerJoinOutLinks[0].targetPort.getNode();
          break;
        case "decision":
          let { decideNode, firstNeutralNode } = handleDecideNode(current);
          tmpBranch.push(decideNode.extras.inputs);
          current = firstNeutralNode;
          break;
        case "default":
          tmpBranch.push(current.extras.inputs);
          parent = current;
          if (outputLinks.length > 0) {
            current = outputLinks[0].targetPort.getNode();
          } else {
            current = null;
          }
          break;
        default:
          current = null;
      }
    }
    forkTasks.push(tmpBranch);
  });

  forkNode.extras.inputs.forkTasks = forkTasks;
  joinNode.extras.inputs.joinOn = joinOn;

  return { forkNode, joinNode };
};

export const handleDecideNode = decideNode => {
  let failBranchLink = _.values(decideNode.ports.failPort.links)[0];
  let neutralBranchLink = _.values(decideNode.ports.neutralPort.links)[0];
  let firstNeutralNode = null;

  [failBranchLink, neutralBranchLink].forEach((branch, i) => {
    let branchArray = [];

    if (branch) {
      let currentNode = branch.targetPort.getNode();
      let inputLinks = getLinksArray("in", currentNode);
      let outputLink = getLinksArray("out", currentNode)[0];

      while (
        (inputLinks.length === 1 ||
          currentNode.type === "join" ||
          currentNode.type === "fork") &&
        outputLink
      ) {
        switch (currentNode.type) {
          case "fork":
            let { forkNode, joinNode } = handleForkNode(currentNode);
            branchArray.push(forkNode.extras.inputs, joinNode.extras.inputs);
            currentNode = getLinksArray(
              "out",
              joinNode
            )[0].targetPort.getNode();
            break;
          case "decision":
            let innerDecideNode = handleDecideNode(currentNode).decideNode;
            let innerFirstNeutralNode = handleDecideNode(currentNode)
              .firstNeutralNode;
            branchArray.push(innerDecideNode.extras.inputs);
            if (innerFirstNeutralNode && innerFirstNeutralNode.extras.inputs) {
              branchArray.push(innerFirstNeutralNode.extras.inputs);
              currentNode = getLinksArray(
                "out",
                innerFirstNeutralNode
              )[0].targetPort.getNode();
            } else {
              currentNode = innerFirstNeutralNode;
            }
            break;
          default:
            branchArray.push(currentNode.extras.inputs);
            currentNode = outputLink.targetPort.getNode();
            break;
        }
        inputLinks = getLinksArray("in", currentNode);
        outputLink = getLinksArray("out", currentNode)[0];
      }

      firstNeutralNode = currentNode;
    }

    let casesValues = Object.keys(decideNode.extras.inputs.decisionCases);

    switch (i) {
      case 0:
        decideNode.extras.inputs.decisionCases[casesValues[0]] = branchArray;
        break;
      case 1:
        decideNode.extras.inputs.defaultCase = branchArray;
        break;
      default:
        break;
    }
  });

  return { decideNode, firstNeutralNode };
};

export const linkNodes = (node1, node2, whichPort) => {
  if (
    node1.type === "fork" ||
    node1.type === "join" ||
    node1.type === "start"
  ) {
    const fork_join_start_outPort = node1.getPort("right");

    if (node2.type === "default") {
      return fork_join_start_outPort.link(node2.getInPorts()[0]);
    }
    if (node2.type === "fork") {
      return fork_join_start_outPort.link(node2.getPort("left"));
    }
    if (node2.type === "join") {
      return fork_join_start_outPort.link(node2.getPort("left"));
    }
    if (node2.type === "decision") {
      return fork_join_start_outPort.link(node2.getPort("inputPort"));
    }
    if (node2.type === "end") {
      return fork_join_start_outPort.link(node2.getPort("left"));
    }
  } else if (node1.type === "default") {
    const defaultOutPort = node1.getOutPorts()[0];

    if (node2.type === "default") {
      return defaultOutPort.link(node2.getInPorts()[0]);
    }
    if (node2.type === "fork") {
      return defaultOutPort.link(node2.getPort("left"));
    }
    if (node2.type === "join") {
      return defaultOutPort.link(node2.getPort("left"));
    }
    if (node2.type === "decision") {
      return defaultOutPort.link(node2.getPort("inputPort"));
    }
    if (node2.type === "end") {
      return defaultOutPort.link(node2.getPort("left"));
    }
  } else if (node1.type === "decision") {
    const currentPort = node1.getPort(whichPort);

    if (node2.type === "default") {
      return currentPort.link(node2.getInPorts()[0]);
    }
    if (node2.type === "fork") {
      return currentPort.link(node2.getPort("left"));
    }
    if (node2.type === "join") {
      return currentPort.link(node2.getPort("left"));
    }
    if (node2.type === "decision") {
      return currentPort.link(node2.getPort("inputPort"));
    }
    if (node2.type === "end") {
      return currentPort.link(node2.getPort("left"));
    }
  }
};
