const request = require("superagent");

const HttpClient = {
  get: (path: string, token: string) =>
    new Promise((resolve, reject) => {
      const req = request.get(path).accept("application/json");
      if (token) {
        req.set("Authorization", token);
      }
      req.end((err, res) => {
        if (err) {
          if (res && res.error) {
            resolve(res.error.status);
            console.log(res.error.message);
          } else {
            reject(err);
          }
        } else {
          resolve(res.body);
        }
      });
    }),

  delete: (path: string, data: [{}], token: string) =>
    new Promise((resolve, reject) => {
      const req = request
        .delete(path, data)
        .accept("application/json")
        .query("archiveWorkflow=false");
      if (token) {
        req.set("Authorization", token);
      }
      req.end((err, res) => {
        if (err) {
          resolve(err);
          reject(err);
        } else {
          if (res) {
            resolve(res);
          }
        }
      });
    }),

  post: (path: string, data: [{}], token: string) =>
    new Promise((resolve, reject) => {
      const req = request
        .post(path, data)
        .set("Content-Type", "application/json");
      if (token) {
        req.set("Authorization", token);
      }
      req.end((err, res) => {
        if (err || !res.ok) {
          console.error("Error on post! " + res);
          reject(err);
        } else {
          if (res) {
            resolve(res);
          }
        }
      });
    }),

  put: (path: string, data: [{}], token: string) =>
    new Promise((resolve, reject) => {
      const req = request.put(path, data).set("Accept", "application/json");

      if (token) {
        req.set("Authorization", token);
      }

      req
        .then(res => {
          resolve(res);
        })
        .catch(error => {
          reject(error);
        });
    })
};

exports.HttpClient = HttpClient;
