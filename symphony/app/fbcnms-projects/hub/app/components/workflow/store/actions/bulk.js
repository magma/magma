import { HttpClient as http } from "../../common/HttpClient";
import { round } from "lodash/math";
import { fetchNewData, fetchParentWorkflows } from "./searchExecs";
import { conductorApiUrlPrefix } from "../../constants";

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
  defaultPages,
  backendApiUrlPrefix
) => {
  return (dispatch, getState) => {
    dispatch(storeResponse(successfulResults, errorResults));
    const { isFlat } = getState().bulkReducer;
    isFlat
      ? dispatch(fetchNewData(1, defaultPages, backendApiUrlPrefix))
      : dispatch(fetchParentWorkflows(1, defaultPages, backendApiUrlPrefix));
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

export const checkDeleted = (deletedWfs, workflows, defaultPages, backendApiUrlPrefix) => {
  return dispatch => {
    if (deletedWfs.length === workflows.length) {
      dispatch(receiveBulkOperationResponse(deletedWfs, {}, defaultPages, backendApiUrlPrefix));
    } else {
      setTimeout(
        () => dispatch(checkDeleted(deletedWfs, workflows, defaultPages)),
        200
      );
    }
  };
};

export const performBulkOperation = (operation, workflows, defaultPages, backendApiUrlPrefix) => {
  const url = backendApiUrlPrefix + "/bulk/" + operation;
  let deletedWfs = [];

  return dispatch => {
    dispatch(requestBulkOperation());
    try {
      switch (operation) {
        case "retry":
        case "restart":
          http.post(url, workflows).then(res => {
            const { bulkSuccessfulResults, bulkErrorResults } = res.body.text
              ? JSON.parse(res.body.text)
              : [];
            dispatch(
              receiveBulkOperationResponse(
                bulkSuccessfulResults,
                bulkErrorResults,
                defaultPages,
                backendApiUrlPrefix
              )
            );
          });
          break;
        case "pause":
        case "resume":
          http.put(url, workflows).then(res => {
            const { bulkSuccessfulResults, bulkErrorResults } = res.body.text
              ? JSON.parse(res.body.text)
              : [];
            dispatch(
              receiveBulkOperationResponse(
                bulkSuccessfulResults,
                bulkErrorResults,
                defaultPages,
                backendApiUrlPrefix
              )
            );
          });
          break;
        case "terminate":
          http.delete(url, workflows).then(res => {
            const { bulkSuccessfulResults, bulkErrorResults } = res.body.text
              ? JSON.parse(res.body.text)
              : [];
            dispatch(
              receiveBulkOperationResponse(
                bulkSuccessfulResults,
                bulkErrorResults,
                defaultPages,
                backendApiUrlPrefix
              )
            );
          });
          break;
        case "delete":
          workflows.map(wf => {
            http.delete(backendApiUrlPrefix + "/workflow/" + wf).then(() => {
              deletedWfs.push(wf);
              dispatch(
                updateLoadingBar(
                  round((deletedWfs.length / workflows.length) * 100)
                )
              );
            });
            return null;
          });
          dispatch(checkDeleted(deletedWfs, workflows, defaultPages, backendApiUrlPrefix));
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
