import { saveAs } from "file-saver";
import {conductorApiUrlPrefix, frontendUrlPrefix} from "../../constants";
import {HttpClient as http} from "../../common/HttpClient";
import { Button } from "react-bootstrap";
import React from "react";

const JSZip = require("jszip");

export const changeUrl = (history) => {
 return function(e) {
    history.push(frontendUrlPrefix + "/" + e);
  };
};

const exportFile = () => {
  http.get(conductorApiUrlPrefix + '/metadata/workflow').then((res) => {
    const zip = new JSZip();
    let workflows = res.result || [];

    workflows.forEach((wf) => {
      zip.file(wf.name + ".json", JSON.stringify(wf, null, 2));
    });

    zip.generateAsync({ type: "blob" }).then(function(content) {
      saveAs(content, "workflows.zip");
    });
  });
};

export const exportButton = () => {
  return (
      <Button
          variant="outline-primary"
          style={{ marginLeft: "5px" }}
          onClick={exportFile}>
        <i className="fas fa-file-export" />
        &nbsp;&nbsp;Export
      </Button>
  );
}