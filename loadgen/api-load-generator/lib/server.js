// Copyright © 2025 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

import express from "express";
import bodyParser from "body-parser";
import { getState, isNumber, logWrite, getLogLevel } from "./utils.js";

const app = express();
const state = getState();
let logResetTimeout = null;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

function setGlobals(request) {
  ["org", "env", "proxy"].forEach((item) => {
    const headerName = "x-apigee-" + item;
    if (request.header(headerName)) {
      state.globals[item] = request.header(headerName);
    }
  });
}

app.get("/status", (request, response) => {
  setGlobals(request);
  response.header({ "Content-Type": "application/json" });
  if (!(state.g.model && state.g.model.id)) {
    response.status(200).send('{ "status" : "starting" }\n');
    return;
  }
  state.g.status.times.current = new Date().toISOString();
  const payload = { status: { ...state.g.status } };
  response.status(200).send(JSON.stringify(payload, null, 2) + "\n");
});

app.post("/control", (request, response) => {
  try {
    setGlobals(request);
    if (!(state.g.model && state.g.model.id)) {
      response.status(503).send('{ "status" : "starting" }\n');
      return;
    }
    const payload = { status: { ...state.g.status } };
    payload.status.times.current = new Date().toISOString();

    const action = request.body.action || request.query.action;
    const requestedLoglevel = request.body.loglevel || request.query.loglevel;

    response.header({ "Content-Type": "application/json" });

    if (action !== "stop" && action !== "start" && action !== "setlog") {
      payload.error = "unsupported request (action)";
      response.status(400).send(JSON.stringify(payload, null, 2) + "\n");
      return;
    }

    if (action === "setlog") {
      if (!isNumber(requestedLoglevel)) {
        payload.error = "must pass loglevel";
        response.status(400).send(JSON.stringify(payload, null, 2) + "\n");
        return;
      }
      if (logResetTimeout) {
        clearTimeout(logResetTimeout);
      }
      // coerce
      const newLevel = Math.max(
        0,
        Math.min(10, parseInt(requestedLoglevel, 10)),
      );
      state.g.status.loglevel = newLevel;
      let value = state.g.status.loglevel;
      logWrite(2, `setlog loglevel: ${value}`);
      payload.status.loglevel = newLevel;
      response.status(200).send(JSON.stringify(payload, null, 2) + "\n");

      // setTimeout to restore loglevel reset
      let requestedTimeout =
        request.body.timeout == undefined
          ? state.globals.LOG_RESET_TIMEOUT
          : Math.max(0, Math.min(120000, parseInt(request.body.timeout, 10)));
      if (requestedTimeout > 0) {
        logWrite(
          6,
          `log will reset to ${getLogLevel()} after ${requestedTimeout}ms`,
        );
        logResetTimeout = setTimeout(() => {
          const value = getLogLevel();
          state.g.status.loglevel = value;
          logWrite(2, `log reset timer, loglevel now ${value}`);
          logResetTimeout = null;
        }, requestedTimeout);
      } else {
        logWrite(2, `no reset log timeout`);
      }
      return;
    }

    const currentRunState = state.g.status.runState;
    const noChange =
      (currentRunState === "stopped" && action === "stop") ||
      (currentRunState === "running" && action === "start");
    if (noChange) {
      payload.error = `already ${currentRunState}`;
      response.status(400).send(JSON.stringify(payload, null, 2) + "\n");
      return;
    }
    state.g.status.runState = action === "stop" ? "stopped" : "running";
    payload.status.runState = state.g.status.runState;
    payload.status.priorRunState = currentRunState;
    response.status(200).send(JSON.stringify(payload, null, 2) + "\n");
  } catch (e) {
    console.log(`error: ` + e);
    throw e;
  }
});

// default behavior
app.all("*", (_request, response) => {
  response
    .header({ "Content-Type": "application/json" })
    .status(404)
    .send('{ "message" : "not found. That\'s all we know." }\n');
});

// This must come AFTER all other app.use() and routes.
app.use((err, _req, res, next) => {
  // Check if the error is a syntax error from JSON parsing
  if (err instanceof SyntaxError && err.status === 400 && "body" in err) {
    console.error("Bad JSON:", err.message);
    return res.status(400).send({ message: "The JSON payload is malformed." });
  }
  // If it's another kind of error, pass it to the default Express error handler
  next();
});

export default app;
