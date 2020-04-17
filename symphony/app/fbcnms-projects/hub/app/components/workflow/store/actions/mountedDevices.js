import axios from "axios";
export const MOUNTED_DEVICES = "MOUNTED_DEVICES";

const getTopologyDevices = async topologyType => {
  const devices = await axios.get("/api/odl/oper/all/status/" + topologyType);
  if (
    typeof devices === "object" &&
    Array.isArray(Object.keys(devices[Object.keys(devices)]))
  ) {
    const topologies = Object.keys(devices);
    const topology = Object.keys(devices[Object.keys(devices)]);
    if (devices[topologies][topology]["node"]) {
      return devices[topologies][topology]["node"].map(node => {
        return node["node-id"];
      });
    }
  }
  return [];
};

export const getMountedDevices = () => {
  return async dispatch => {
    const allCliDevices = await getTopologyDevices("cli");
    const allNetconfDevices = await getTopologyDevices("topology-netconf");
    dispatch(updateDevices(allCliDevices.concat(allNetconfDevices)));
  };
};

export const updateDevices = devices => {
  return { type: MOUNTED_DEVICES, devices };
};
