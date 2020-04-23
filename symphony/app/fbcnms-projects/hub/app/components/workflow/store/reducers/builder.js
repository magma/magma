import {
  LOCK_WORKFLOW_NAME, OPEN_CARD,
  RESET_TO_DEFAULT_WORKFLOW,
  SHOW_CUSTOM_ALERT, STORE_TASKS,
  STORE_WORKFLOW_ID,
  STORE_WORKFLOWS,
  SWITCH_SMART_ROUTING, UPDATE_BUILDER_LABELS,
  UPDATE_BUILDER_QUERY,
  UPDATE_FINAL_WORKFLOW, UPDATE_TASKS,
  UPDATE_WORKFLOWS
} from "../actions/builder";

const finalWorkflowTemplate = {
  name: "",
  description: "",
  version: 1,
  tasks: [],
  outputParameters: {},
  inputParameters: [],
  schemaVersion: 2,
  restartable: true,
  workflowStatusListenerEnabled: false
};

const initialState = {
  workflows: [],
  tasks: [],
  originalWorkflows: [],
  originalTasks: [],
  query: "",
  labels: [],
  openCard: null,
  functional: [
    { name: "start", description: "Starting point of every workflow" },
    { name: "end", description: "Successful finish of a workflow" },
    { name: "terminate", description: "Unsuccessful termination of a workflow" },
    { name: "decision", description: "Conditional branching point in a workflow" },
    { name: "fork", description: "Concurrent execution fork in a workflow" },
    { name: "join", description: "Concurrent execution join in a workflow" },
    { name: "http", description: "HTTP execution task" },
    { name: "lambda", description: "Arbitrary javascript code execution task" },
    { name: "wait", description: "Wait for a specific event before continuing" },
    { name: "event", description: "Publish a specific event (other workflow might be waiting for)" },
    { name: "raw", description: "Specify a task in JSON" },
  ],
  workflowNameLock: false,
  switchSmartRouting: false,
  executedWfId: null,
  customAlert: {
    show: false,
    variant: "danger",
    msg: "",
  },
  finalWorkflow: {
    name: "",
    description: "",
    version: 1,
    tasks: [],
    outputParameters: {},
    inputParameters: [],
    schemaVersion: 2,
    restartable: true,
    workflowStatusListenerEnabled: false,
  },
};

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case UPDATE_BUILDER_QUERY: {
      let { query } = action;
      return { ...state, query };
    }
    case UPDATE_BUILDER_LABELS: {
      let { labels } = action;
      return { ...state, labels };
    }
    case STORE_WORKFLOWS: {
      const { originalWorkflows, workflows } = action;
      return { ...state, originalWorkflows, workflows };
    }
    case STORE_TASKS: {
      const { originalTasks, tasks } = action;
      return { ...state, originalTasks, tasks };
    }
    case OPEN_CARD: {
      const { openCard } = action;
      return { ...state, openCard };
    }
    case UPDATE_WORKFLOWS: {
      const { workflows } = action;
      return { ...state, workflows };
    }
    case UPDATE_TASKS: {
      const { tasks } = action;
      return { ...state, tasks };
    }
    case RESET_TO_DEFAULT_WORKFLOW: {
      return {
        ...state,
        finalWorkflow: finalWorkflowTemplate,
        workflowNameLock: false
      };
    }
    case STORE_WORKFLOW_ID: {
      const { executedWfId } = action;
      return { ...state, executedWfId };
    }
    case LOCK_WORKFLOW_NAME: {
      return { ...state, workflowNameLock: true };
    }
    case SWITCH_SMART_ROUTING: {
      const { switchSmartRouting } = state;
      return { ...state, switchSmartRouting: !switchSmartRouting };
    }
    case UPDATE_FINAL_WORKFLOW: {
      let { finalWorkflow } = action;
      return { ...state, finalWorkflow };
    }
    case SHOW_CUSTOM_ALERT: {
      let { show, variant, msg } = action;
      return { ...state, customAlert: { show, variant, msg } };
    }
    default:
      break;
  }
  return state;
};

export default reducer;
