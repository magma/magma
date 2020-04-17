
export const wfLabelsColor = [
  "#7D6608",
  "#43ABC9",
  "#EBC944",
  "#CD6155",
  "#F4D03F",
  "#808B96",
  "#212F3D",
  "#4A340C",
  "#00cd00",
  "#18b5b5",
  "#3A48EC",
  "#EA9D16",
  "#7D3C98",
  "#A6ACAF",
  "#F1948A",
  "#02d500",
  "#AF4141",
  "#EA7616",
  "#A569BD",
  "#68386C",
  "#5A5144",
  "#6F927D",
  "#3AEC60",
  "#EDB152",
  "#C52F38",
  "#A3A042",
  "#249D83",
  "#0DAA79",
  "#3A96EC",
  "#3ADFEC",
  "#5D6D7E",
  "#000080",
  "#229954",
  "#117864",
  "#16A085",
  "#107896"
];

export const workflowDescriptions = {
  name: "name of the workflow",
  description: "description of the workflow (optional)",
  version:
    "numeric field used to identify the version of the schema (use incrementing numbers)",
  tasks: [],
  outputParameters: {},
  schemaVersion:
    "current Conductor Schema version, schemaVersion 1 is discontinued",
  restartable: "boolean flag to allow workflow restarts",
  workflowStatusListenerEnabled:
    "ff true, every workflow that gets terminated or completed will send a notification"
};

export const taskDescriptions = {
  name: "name of the task",
  taskReferenceName:
    "alias used to refer the task within the workflow (MUST be unique within workflow)",
  optional: "when set to true - workflow continues even if the task fails.",
  startDelay: "time period before task executes"
};