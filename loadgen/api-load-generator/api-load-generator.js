// Copyright © 2019,2025 Google LLC.
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

import fs from "fs";

import { getState, trackFailure } from "./lib/utils.js";
import { setInitialContext, initializeJobRun } from "./lib/load-generator.js";
import app from "./lib/server.js";
import stripJsonComments from "strip-json-comments";
const DEFAULT_CONFIG_FILE = "./config/config.json";

function reportModel(context) {
  console.log("================================================");
  console.log("==         Job Definition Retrieved           ==");
  console.log("================================================");
  console.log(JSON.stringify(context, null, 2));
  return context;
}

async function kickoff(arg) {
  try {
    let configFile = arg;
    if (!configFile || !fs.existsSync(configFile)) {
      configFile = DEFAULT_CONFIG_FILE;
    }
    if (fs.existsSync(configFile)) {
      const state = getState();
      console.log(configFile);
      const cleanJson = stripJsonComments(fs.readFileSync(configFile, "utf8"));
      state.g.model = JSON.parse(cleanJson);
      state.g.status.runState = "starting";
      if (fs.existsSync("config/defaults.json")) {
        state.defaults = JSON.parse(
          fs.readFileSync("config/defaults.json", "utf8"),
        );
      }

      reportModel(state.g.model);
      const context = setInitialContext(state.g.model);
      context.initializing = true;
      await initializeJobRun(context);
    } else {
      console.log(`That file does not exist. (${arg})`);
    }
  } catch (exc1) {
    console.log("Exception:" + exc1);
    trackFailure(exc1);
  }
}

const positionalArgs = process.argv.slice(2);
const modelFilename = positionalArgs[0] || DEFAULT_CONFIG_FILE;
const port = process.env.PORT || 5950;

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
  setTimeout(() => kickoff(modelFilename), 600);
});

process.on("uncaughtException", function (err) {
  console.log("uncaughtException: " + err);
  console.log(err.stack);
});
