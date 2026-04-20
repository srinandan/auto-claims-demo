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

import util from "node:util";
import { execSync } from "node:child_process";

const state = {
  globals: {
    LOG_RESET_TIMEOUT: 10000,
  },
  defaults: {
    initialLogLevel: 2,
  },
  g: {
    model: null,
    context: null,
    status: {
      version: "20250826-1642",
      pid: process.pid,
      nRequests: 0,
      nCycles: 0,
      times: {
        start: new Date().toISOString(),
        lastRun: new Date().toISOString(),
      },
      jobId: "",
      description: "",
      runState: "none",
      responseCounts: { total: 0 },
    },
  },
  cache: {
    _data: {},
    get: (key) => state.cache._data[key],
    put: (key, value) => {
      state.cache._data[key] = value;
    },
  },
};

export const logWrite = function (level, ...args) {
  if (state.g.status.loglevel >= level) {
    const time = new Date().toString();
    const tstr =
      "[" +
      time.slice(11, 15) +
      "-" +
      time.slice(4, 7) +
      "-" +
      time.slice(8, 10) +
      " " +
      time.slice(16, 24) +
      "] ";
    console.log(tstr + util.format.apply(util, args));
  }
};

export function isNumber(n) {
  if (typeof n === "undefined") {
    return false;
  }
  return !isNaN(parseFloat(n)) && isFinite(n);
}

export function trackFailure(reason) {
  if (reason) {
    logWrite(0, "failure: " + reason);
    logWrite(1, reason.stack);
    state.g.status.lastError = {
      message: reason.stack.toString(),
      time: new Date().toISOString(),
    };
  } else {
    logWrite(0, "unknown failure?");
  }
}

export function cacheKey(tag) {
  return "runload-" + tag + "-" + (state.g.model.id || "xxx");
}

export function resolveExpression(input, defaultValue) {
  let I = input;
  if (typeof input === "undefined") {
    I = defaultValue;
  } else if (typeof input === "string") {
    const trimmedInput = input.trim();
    if (trimmedInput === "") {
      I = input;
    } else {
      try {
        I = JSON.parse(input);
      } catch (e) {
        const n = Number(input);
        if (!isNaN(n)) {
          I = n;
        } else {
          I = input;
        }
      }
    }
  }
  return I;
}

export async function getLoadGeneratorSource() {
  const service = process.env.K_SERVICE;
  if (service) {
    try {
      const response = await fetch(
        "http://metadata.google.internal/computeMetadata/v1/project/project-id",
        {
          headers: { "Metadata-Flavor": "Google" },
        },
      );
      if (response.ok) {
        const projectId = await response.text();
        return `cloudrun::project:${projectId},service:${service}`;
      }
    } catch (e) {
      // fall through to default
    }
  }

  let hostname = "unknown";
  try {
    hostname = execSync("hostname").toString().trim();
  } catch (e) {
    // ignore
  }

  return `host:${hostname}`;
}

export function capitalizeOne(str) {
  return str.charAt(0).toUpperCase().concat(str.slice(1).toLowerCase());
}

export const dayNames = [
  "sunday",
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
];

export const ONE_HOUR_IN_MS = 60 * 60 * 1000;
export const MIN_SLEEP_TIME_MS = 120;

export function getVariationByDayOfWeek(currentDayOfWeek) {
  const dayName = dayNames[currentDayOfWeek % 7];

  function trySource(v) {
    const vtype = Object.prototype.toString.call(v);

    if (vtype === "[object Array]") {
      if (
        v.length &&
        v.length === 7 &&
        v[currentDayOfWeek] &&
        v[currentDayOfWeek] > 0 &&
        v[currentDayOfWeek] <= 10
      ) {
        return v[currentDayOfWeek];
      } else {
        logWrite(2, "variationByDayOfWeek seems of wrong length, or value");
      }
    } else if (vtype === "[object Object]") {
      if (dayName) {
        return v[dayName] || v[capitalizeOne(dayName)];
      } else {
        logWrite(2, "variationByDayOfWeek seems wrong: " + dayName);
      }
    }
    return undefined;
  }

  let v =
    trySource(state.g.model.variationByDayOfWeek) ||
    trySource(state.defaults.variationByDayOfWeek) ||
    1;
  if (!(v > 0 && v <= 10)) {
    logWrite(2, "variationByDayOfWeek seems wrong: " + dayName);
    v = 1;
  }
  logWrite(5, "day variation: " + v);
  return v;
}

class Gaussian {
  constructor(mean, stddev) {
    this.mean = mean;
    this.stddev = stddev || mean * 0.1;
  }

  #normal() {
    let u1 = 0,
      u2 = 0;
    while (u1 * u2 === 0) {
      u1 = Math.random();
      u2 = Math.random();
    }
    return Math.sqrt(-2 * Math.log(u1)) * Math.cos(2 * Math.PI * u2);
  }

  next() {
    return this.stddev * this.#normal() + 1 * this.mean;
  }
}

export function getTargetRunsPerHour(nowMs) {
  const now = new Date(nowMs);
  let currentHour = now.getHours();
  const currentMinute = now.getMinutes();
  const currentDayOfWeek = now.getDay();
  const speedfactor = 3.2;
  function getInvocationsPerHour(hour) {
    const r =
      state.g.model.invocationsPerHour &&
      state.g.model.invocationsPerHour[hour];
    if (!r) {
      return state.defaults.invocationsPerHour[hour];
    }
    return r;
  }
  if (currentHour < 0 || currentHour > 23) {
    currentHour = 0;
  }
  const nextHour = currentHour === 23 ? 0 : currentHour + 1;
  const hourFraction = (currentMinute + 1) / 60;
  const baseRunsPerHour = getInvocationsPerHour(currentHour);
  const smoothedRunsPerHour =
    speedfactor *
    (baseRunsPerHour +
      hourFraction * (getInvocationsPerHour(nextHour) - baseRunsPerHour));
  const v = getVariationByDayOfWeek(currentDayOfWeek);
  const scaledRunsPerHour = v
    ? Math.floor(smoothedRunsPerHour * v)
    : smoothedRunsPerHour;
  const gaussian = new Gaussian(scaledRunsPerHour, 0.05 * scaledRunsPerHour); // fuzz
  const fuzzedRunsPerHour = Math.floor(gaussian.next());
  logWrite(
    4,
    `runsPerHour base(${baseRunsPerHour}) smoothed(${smoothedRunsPerHour}) scaled(${scaledRunsPerHour}) fuzzed(${fuzzedRunsPerHour})`,
  );
  return fuzzedRunsPerHour;
}

export function getNow() {
  return Date.now();
}

export function getState() {
  return state;
}

export function getLogLevel(context) {
  if (!context) {
    context = { job: state?.g?.model };
  }
  return (
    context?.job?.initialLogLevel ||
    context?.job?.loglevel ||
    state.defaults?.initialLogLevel ||
    2
  );
}

export function getSleepTimeInMs(startOfMostRecentRequest) {
  const nowMs = getNow();
  const runsPerHour = getTargetRunsPerHour(nowMs);
  const delayBtwnRuns = Math.floor(ONE_HOUR_IN_MS / runsPerHour);
  const durationOfLastRun = nowMs - startOfMostRecentRequest;
  const calculatedSleepTime =
    Math.floor(ONE_HOUR_IN_MS / runsPerHour) - durationOfLastRun;

  const sleepTimeInMs =
    calculatedSleepTime < MIN_SLEEP_TIME_MS
      ? MIN_SLEEP_TIME_MS
      : calculatedSleepTime;

  logWrite(
    5,
    `dRPH(${runsPerHour}) interval(${delayBtwnRuns}) durationOfLast(${durationOfLastRun}) calcSleep(${calculatedSleepTime}) sleep(${sleepTimeInMs}) wake(${new Date(nowMs + sleepTimeInMs).toString().slice(16, 24)})`,
  );

  return sleepTimeInMs;
}
