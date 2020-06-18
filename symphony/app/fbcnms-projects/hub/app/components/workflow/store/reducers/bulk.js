import {
  IS_FLAT,
  REQUEST_BULK_OPERATION,
  RECEIVE_BULK_OPERATION_RESPONSE,
  FAIL_BULK_OPERATION,
  RESET_BULK_OPERATION_RESULT,
  UPDATE_LOADING_BAR
} from "../actions/bulk";

const initialState = {
  isFetching: false,
  isFlat: true,
  error: null,
  successfulResults: [],
  errorResults: {},
  data: [],
  table: [],
  query: "",
  label: [],
  loading: 0
};

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case IS_FLAT: {
      const { isFlat } = action;
      return { ...state, isFlat };
    }
    case REQUEST_BULK_OPERATION: {
      return { ...state, isFetching: true, error: null };
    }
    case RECEIVE_BULK_OPERATION_RESPONSE: {
      const { successfulResults = [], errorResults = {} } = action;

      return {
        ...state,
        isFetching: false,
        error: null,
        successfulResults,
        errorResults
      };
    }
    case FAIL_BULK_OPERATION: {
      const { error } = action;

      return { ...state, isFetching: false, error };
    }
    case RESET_BULK_OPERATION_RESULT: {
      return { ...state, successfulResults: [], errorResults: [], loading: 0 };
    }
    case UPDATE_LOADING_BAR: {
      const { percentage } = action;
      return { ...state, loading: percentage };
    }
    default:
      break;
  }
  return state;
};

export default reducer;
