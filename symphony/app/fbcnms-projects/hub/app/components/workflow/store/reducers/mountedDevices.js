import { MOUNTED_DEVICES } from "../actions/mountedDevices";

const initialState = {
  devices: []
};

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case MOUNTED_DEVICES: {
      let { devices } = action;
      devices = devices ? devices : [];
      return { ...state, devices };
    }
    default:
      break;
  }
  return state;
};

export default reducer;
