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

import {
  logWrite,
  trackFailure,
  getSleepTimeInMs,
  resolveExpression,
  getLoadGeneratorSource,
  getNow,
  getState,
  getLogLevel,
} from "./utils.js";
import {
  expandEmbeddedTemplates,
  evalJavascript,
  resolveExtract,
} from "./context-builder.js";
import ipGenerator from "./ipGenerator.js";
import Handlebars from "handlebars";

const MAX_ITERATIONS = 50;
const PROGRESS_LOG_INTERVAL = 60 * 1000;

async function invokeOneRequest(linkUrl, method, payload, headers, job) {
  method = method.toLowerCase();
  const options = {
    method,
    headers,
    redirect: "manual",
  };

  if (method === "post" || method === "put") {
    logWrite(6, "expanding templates...");
    payload = expandEmbeddedTemplates(payload);
    logWrite(6, "payload: " + JSON.stringify(payload));
    const t2 = Object.prototype.toString.call(payload);
    if (t2 === "[object String]") {
      options.body = payload;
      if (!options.headers["content-type"]) {
        options.headers["content-type"] = "application/x-www-form-urlencoded";
      }
    } else {
      options.body = JSON.stringify(payload);
      if (!options.headers["content-type"]) {
        options.headers["content-type"] = "application/json";
      }
    }
  }

  options.headers["x-runload-source"] = await getLoadGeneratorSource();

  if (job.simulateGeoDistribution) {
    const selected = ipGenerator.generateIp();
    logWrite(5, "selected client IP = %s", JSON.stringify(selected));
    options.headers["x-forwarded-for"] = selected.ip;
  } else {
    logWrite(5, "not contriving an IP");
  }
  logWrite(8, "headers " + JSON.stringify(options.headers));

  try {
    const state = getState();
    const response = await fetch(linkUrl, options);
    state.g.status.nRequests++;
    logWrite(
      4,
      "%s %s  ==> %d",
      options.method.toUpperCase(),
      linkUrl,
      response.status,
    );
    const aIndex = response.status + "";
    if (state.g.status.responseCounts.hasOwnProperty(aIndex)) {
      state.g.status.responseCounts[aIndex]++;
    } else {
      state.g.status.responseCounts[aIndex] = 1;
    }
    state.g.status.responseCounts.total++;

    const responseHeaders = {};
    response.headers.forEach((value, name) => {
      responseHeaders[name] = value;
    });

    let body = await response.text();
    try {
      body = JSON.parse(body);
    } catch (e) {
      // ignore
    }
    return { body, headers: responseHeaders, status: response.status };
  } catch (e) {
    if (linkUrl.startsWith("http:")) {
      logWrite(
        4,
        "%s %s  ==> FAIL (as expected)",
        options.method.toUpperCase(),
        linkUrl,
      );
    } else {
      logWrite(2, "%s %s  ==> FAIL", options.method.toUpperCase(), linkUrl);
      console.error(e);
    }
    return {};
  }
}

async function invokeOneBatchOfRequests(r, job) {
  const state = getState();
  if (r.imports) {
    Object.keys(r.imports).forEach((name) => {
      const s = r.imports[name];
      try {
        if (s.startsWith("{{") && s.endsWith("}}")) {
          // handlebars
          const tmpl = Handlebars.compile(s);
          const result = tmpl(state.g.context);
          logWrite(4, "import: context[%s] = %s", name, result);
          state.g.context[name] = result;
        } else if (s.startsWith("{") && s.endsWith("}")) {
          // eval, because sometimes I don't want a string.
          const result = evalJavascript(s.slice(1, -1), state.g.context);
          logWrite(4, "import: context[%s] = %s", name, JSON.stringify(result));
          state.g.context[name] = result;
        }
      } catch (e) {
        logWrite(
          2,
          "exception while applying import template: " + r.imports[name],
        );
        logWrite(2, "exception: " + e);
      }
    });
  }

  let linkUrl = r.url || r.endpoint;
  // faulty input from user can cause exception in handlebars
  try {
    const template = Handlebars.compile(linkUrl);
    linkUrl = template(state.g.context);
  } catch (e) {
    logWrite(2, "exception while applying URL template: " + linkUrl);
    logWrite(2, "exception : " + e);
  }

  let method = r.method || "GET";
  if (method.toUpperCase() !== "GET") {
    try {
      const template = Handlebars.compile(method);
      method = template(state.g.context);
    } catch (e) {
      logWrite(2, "exception while applying method template: " + method);
      logWrite(2, "exception : " + e);
    }
  }

  const headers = {};
  if (r.headers) {
    logWrite(4, "applying headers");
    Object.keys(r.headers).forEach((name) => {
      try {
        //logWrite(7, 'header[%s]...', name);
        const nTemplate = Handlebars.compile(name, { noEscape: true });
        name = nTemplate(state.g.context);
        let value = r.headers[name];
        logWrite(8, "header[%s] raw value %s", name, value);
        const vTemplate = Handlebars.compile(value, { noEscape: true });
        //logWrite(4, 'gContext %s', JSON.stringify(g.context));
        value = vTemplate(state.g.context);
        logWrite(4, "header[%s]= %s", name, value);
        headers[name] = value;
      } catch (e) {
        logWrite(
          2,
          "exception while applying header template: " + r.headers[name],
        );
        logWrite(2, "exception : " + e);
      }
    });
  }

  const batchsize = Math.min(
    MAX_ITERATIONS,
    resolveExpression(r.iterations, 1),
  );
  const results = await Promise.all(
    Array.from({ length: batchsize }, () =>
      invokeOneRequest(linkUrl, method, r.payload, headers, job),
    ),
  );

  // perform extracts
  const lastResult = results[results.length - 1];
  if (
    r.extracts &&
    Object.prototype.toString.call(r.extracts) === "[object Object]" &&
    Object.keys(r.extracts).length > 0
  ) {
    logWrite(4, "extracts: " + JSON.stringify(Object.keys(r.extracts)));
    Object.keys(r.extracts).forEach(
      resolveExtract(
        lastResult.body,
        lastResult.headers,
        lastResult.status,
        r.extracts,
      ),
    );
  }
}

async function invokeSequence(obj) {
  const state = getState();
  const sequenceStartTime = getNow();
  const job = obj.job;
  if (!job.hasOwnProperty("simulateGeoDistribution")) {
    job.simulateGeoDistribution = state.defaults.simulateGeoDistribution;
  }

  const startingCount = state.g.status.responseCounts.total;
  state.g.context = { ...state.defaults.initialContext, ...job.initialContext };

  // the batches of requests must be serialized
  for (const r of job.requests) {
    await invokeOneBatchOfRequests(r, job);
  }

  const now = getNow();
  let tps =
    (state.g.status.responseCounts.total - startingCount) /
    ((now - sequenceStartTime) / 1000);

  state.g.status.mostRecentSequenceTps =
    Math.round(Math.round(1000 * tps) / 10) / 100;
  tps =
    state.g.status.responseCounts.total /
    ((now - new Date(state.g.status.times.start)) / 1000);

  state.g.status.netTps = Math.round(Math.round(1000 * tps) / 10) / 100;
  state.g.status.nCycles++;
  state.g.status.times.lastRun = new Date(now).toISOString();
  const sleepTime = getSleepTimeInMs(sequenceStartTime);
  state.g.status.times.wake = new Date(now + sleepTime).toISOString();
  setTimeout(wakeup, sleepTime);
}

export function setInitialContext(model) {
  return { job: model };
}

function periodicProgressLog() {
  const ofInterest = { ...getState()?.g?.status };
  delete ofInterest.version;
  delete ofInterest.pid;
  delete ofInterest.jobId;
  delete ofInterest.description;
  console.log(`progress: ` + JSON.stringify(ofInterest));
}

export async function initializeJobRun(context) {
  const state = getState();
  state.g.status.jobId = context.job.id || "-none-";
  state.g.status.description = context.job.description || "-none-";

  // Only during initial startup, set loglevel, and set a logging interval.
  if (context.initializing) {
    state.g.status.loglevel = getLogLevel(context);
    setInterval(periodicProgressLog, PROGRESS_LOG_INTERVAL);
  }

  // launch the loop
  state.g.status.runState = "running";
  try {
    await invokeSequence(context);
  } catch (e) {
    trackFailure(e);
  }
}

export function wakeup() {
  const state = getState();
  const wakeTime = getNow();
  logWrite(3, "awake");
  delete state.g.status.times.wake;
  delete state.g.status.sleepTimeInMs;

  function maybeRun() {
    logWrite(4, `runstate: ${state.g.status.runState}`);
    if (state.g.status.runState === "running") {
      const context = setInitialContext(state.g.model);
      initializeJobRun(context);
    } else {
      state.g.status.nCycles++;
      state.g.status.times.lastRun = new Date(wakeTime).toISOString();
      const sleepTime = getSleepTimeInMs(wakeTime);
      state.g.status.times.wake = new Date(getNow() + sleepTime).toISOString();
      setTimeout(wakeup, sleepTime);
    }
  }

  let value = state.g.status.loglevel;
  logWrite(4, `wakeup() retrieved loglevel: ${value}`);
  maybeRun();
}
