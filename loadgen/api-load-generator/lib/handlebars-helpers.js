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

// handlebars-helpers.js
// ------------------------------------------------------------------
//
// created: Fri Feb  2 15:36:32 2018
// last saved: <2025-August-26 17:38:14>
/* global Buffer */

import Handlebars from "handlebars";
import WeightedRandomSelector from "./weightedRandomSelector.js";
import ScheduleSelector from "./ScheduleSelector.js";

const helpers = {};
const rStringChars =
  "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";

const isNumber = (n) => typeof n === "number";

helpers.base64 = (s) => Buffer.from(s).toString("base64");

helpers.json = (obj) => JSON.stringify(obj);

helpers.httpbasicauth = (u, p) => {
  const userpass =
    Object.prototype.toString.call(p) === "[object String]" &&
    Object.prototype.toString.call(u) === "[object String]"
      ? u + ":" + p
      : u;
  const a = Buffer.from(userpass || "unknown").toString("base64");
  return "Basic " + a;
};

helpers.random = (min, max) => {
  if (!isNumber(min)) {
    min = 0;
  }
  if (!isNumber(max)) {
    max = 1000000;
  }
  return Math.floor(Math.random() * (max - min)) + min;
};

helpers.randomSelect = (a) => {
  if (!Array.isArray(a)) {
    console.log("randomSelect: ERROR");
    throw new Error("expected a to be an array");
  }
  const L = a.length;
  const s = a[Math.floor(Math.random() * L)];
  console.log("randomSelect: " + JSON.stringify(s));
  if (typeof s === "string") {
    return s;
  }
  return s;
};

helpers.split = (s, separator) => {
  const pieces = s.split(separator);
  return pieces;
};

helpers.jsonprop = (json, prop) => {
  const o = JSON.parse(json);
  const v = o[prop];
  console.log("jsonprop: " + v);
  return v;
};

helpers.index_of = (context, ndx) => {
  if (!Array.isArray(context)) {
    context = context.split(",");
  }
  return context[ndx];
};

helpers.randomString = (length) => {
  let result = "";
  if (Object.prototype.toString.call(length) === "[object Object]") {
    length = 0;
  }
  length = length || Math.ceil(Math.random() * 28) + 12;
  length = Math.abs(Math.min(length, 1024));
  for (let i = length; i > 0; --i) {
    result +=
      rStringChars[Math.round(Math.random() * (rStringChars.length - 1))];
  }
  return result;
};

helpers.weightedRandomSelect = (aa) => {
  const wrs = new WeightedRandomSelector(aa);
  const result = wrs.select()[0];
  if (Object.prototype.toString.call(result) === "[object String]") {
    return new Handlebars.SafeString(result);
  }
  return JSON.stringify(result);
};

helpers.firstString = (s) => {
  const pieces = s.split(" ");
  return pieces[0];
};

helpers.lastString = (s) => {
  const pieces = s.split(" ");
  return pieces[pieces.length - 1];
};

const randomInRange = function (min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
};

helpers.scheduledRandomSelect = (aa) => {
  const srs = new ScheduleSelector(aa);
  const range = srs.select().split("-");
  const result = randomInRange(range[0], range[1]);
  if (Object.prototype.toString.call(result) === "[object String]") {
    return new Handlebars.SafeString(result);
  }
  return JSON.stringify(result);
};

export default function registerHelpers() {
  Object.keys(helpers).forEach((key) =>
    Handlebars.registerHelper(key, helpers[key]),
  );
}
