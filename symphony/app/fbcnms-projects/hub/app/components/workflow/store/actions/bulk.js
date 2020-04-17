import axios from "axios";
import { round } from "lodash/math";
import { fetchNewData, fetchParentWorkflows } from "./searchExecs";

export const IS_FLAT = "IS_FLAT";
export const REQUEST_BULK_OPERATION = "REQUEST_BULK_OPERATION";
export const RECEIVE_BULK_OPERATION_RESPONSE =
  "RECEIVE_BULK_OPERATION_RESPONSE";
export const FAIL_BULK_OPERATION = "FAIL_BULK_OPERATION";
export const RESET_BULK_OPERATION_RESULT = "RESET_BULK_OPERATION_RESULT";
export const UPDATE_LOADING_BAR = "UPDATE_LOADING_BAR";

export const requestBulkOperation = () => {
  return { type: REQUEST_BULK_OPERATION };
};

export const receiveBulkOperationResponse = (
  successfulResults,
  errorResults,
  defaultPages
) => {
  return (dispatch, getState) => {
    dispatch(storeResponse(successfulResults, errorResults));
    const { isFlat } = getState().bulkReducer;
    isFlat
      ? dispatch(fetchNewData(1, defaultPages))
      : dispatch(fetchParentWorkflows(1, defaultPages));
    setTimeout(() => dispatch(resetBulkOperationResult()), 2000);
  };
};

export const storeResponse = (successfulResults, errorResults) => {
  return {
    type: RECEIVE_BULK_OPERATION_RESPONSE,
    successfulResults,
    errorResults
  };
};

export const failBulkOperation = error => {
  return { type: FAIL_BULK_OPERATION, error };
};

export const resetBulkOperationResult = () => {
  return { type: RESET_BULK_OPERATION_RESULT };
};

export const updateLoadingBar = percentage => {
  return { type: UPDATE_LOADING_BAR, percentage };
};

export const checkDeleted = (deletedWfs, workflows, defaultPages) => {
  return dispatch => {
    if (deletedWfs.length === workflows.length) {
      dispatch(receiveBulkOperationResponse(deletedWfs, {}, defaultPages));
    } else {
      setTimeout(
        () => dispatch(checkDeleted(deletedWfs, workflows, defaultPages)),
        200
      );
    }
  };
};

export const performBulkOperation = (operation, workflows, defaultPages) => {
  const url = `/api/conductor/bulk/${operation}`;
  let deletedWfs = [];

  return dispatch => {
    dispatch(requestBulkOperation());
    try {
      switch (operation) {
        case "retry":
        case "restart":
          axios.post(url, workflows).then(res => {
            const { bulkSuccessfulResults, bulkErrorResults } = res.body.text
              ? JSON.parse(res.body.text)
              : [];
            dispatch(
              receiveBulkOperationResponse(
                bulkSuccessfulResults,
                bulkErrorResults,
                defaultPages
              )
            );
          });
          break;
        case "pause":
        case "resume":
          axios.put(url, workflows).then(res => {
            const { bulkSuccessfulResults, bulkErrorResults } = res.body.text
              ? JSON.parse(res.body.text)
              : [];
            dispatch(
              receiveBulkOperationResponse(
                bulkSuccessfulResults,
                bulkErrorResults,
                defaultPages
              )
            );
          });
          break;
        case "terminate":
          axios.delete(url, workflows).then(res => {
            const { bulkSuccessfulResults, bulkErrorResults } = res.body.text
              ? JSON.parse(res.body.text)
              : [];
            dispatch(
              receiveBulkOperationResponse(
                bulkSuccessfulResults,
                bulkErrorResults,
                defaultPages
              )
            );
          });
          break;
        case "delete":
          workflows.map(wf => {
            axios.delete("/api/conductor/workflow/" + wf).then(() => {
              deletedWfs.push(wf);
              dispatch(
                updateLoadingBar(
                  round((deletedWfs.length / workflows.length) * 100)
                )
              );
            });
            return null;
          });
          dispatch(checkDeleted(deletedWfs, workflows, defaultPages));
          break;
        default:
          dispatch(failBulkOperation("Invalid operation requested."));
      }
    } catch (e) {
      dispatch(failBulkOperation(e.message));
    }
  };
};

export const setView = isFlat => {
  return { type: IS_FLAT, isFlat };
};
