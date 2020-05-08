import React from "react";
import { Alert } from "react-bootstrap";

const CustomAlert = props => {
  return (
    <div className={props.show ? "custom-alert-in" : "custom-alert-out"}>
      <Alert
        show={props.show}
        variant={props.alertVariant}
        dismissible
        onClose={() => props.showCustomAlert(false)}
        style={{ textAlign: "right", opacity: "70%" }}
      >
        {props.alertVariant === "danger" ? (
          <i className="fas fa-exclamation-triangle" />
        ) : (
          <i className="fas fa-info-circle" />
        )}
        &nbsp;&nbsp;{props.msg}
      </Alert>
    </div>
  );
};

export default CustomAlert;
