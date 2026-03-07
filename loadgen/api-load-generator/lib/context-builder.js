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

import Handlebars from "handlebars";
import jsonpath from "jsonpath";
import { getState, logWrite } from "./utils.js";
import registerHelpers from "./handlebars-helpers.js";

const state = getState();

registerHelpers();

export function expandEmbeddedTemplates(obj) {
  let newObj, tmpl;
  const type = Object.prototype.toString.call(obj);

  if (type === "[object String]") {
    tmpl = Handlebars.compile(obj);
    newObj = tmpl(state.g.context);
  } else if (type === "[object Array]") {
    // iterate
    newObj = [];
    for (let i = 0; i < obj.length; i++) {
      newObj.push(expandEmbeddedTemplates(obj[i]));
    }
  } else if (type === "[object Object]") {
    newObj = {};
    Object.keys(obj).forEach((prop) => {
      const type = Object.prototype.toString.call(obj[prop]);
      if (type === "[object String]") {
        // replace all templates in a string
        tmpl = Handlebars.compile(obj[prop]);
        newObj[prop] = tmpl(state.g.context);
      } else if (type === "[object Object]" || type === "[object Array]") {
        // recurse
        newObj[prop] = expandEmbeddedTemplates(obj[prop]);
      } else {
        // no replacement
        newObj[prop] = obj[prop];
      }
    });
  }
  return newObj;
}

export function evalJavascript(code, ctx) {
  const values = [],
    names = [];
  let result = "";

  // create the fn signature
  Object.keys(ctx).forEach((prop) => {
    names.push(prop);
    values.push(ctx[prop]);
  });

  code = code.trim();
  let src = code.endsWith(";") ? code : "return (" + code + ");";
  logWrite(9, "evalJavascript: " + src);
  try {
    let f = new Function(names.join(","), src);
    logWrite(9, "fn: " + f.toString());
    // call the function with all its arguments
    result = f.apply(null, values);
  } catch (exc1) {
    logWrite(3, "evalJavascript, exception: " + exc1.toString());
    result = "";
  }
  logWrite(7, "evalJavascript, result: " + result);
  return result;
}

export const resolveExtract =
  (payload, headers, status, extracts) => (name) => {
    let value = extracts[name];
    if (name && value) {
      if (value.startsWith("{{") && value.endsWith("}}")) {
        logWrite(4, "extract: handlebars...");
        // handlebars
        let template = Handlebars.compile(value, { noEscape: true });
        let result = template(state.g.context);
        logWrite(4, "extract: context[%s] = %s", name, result);
        state.g.context[name] = result;
      } else if (value.startsWith("{") && value.endsWith("}")) {
        // eval, because sometimes I don't want a string.
        logWrite(4, "extract: evalJavascript...");
        let result = evalJavascript(
          value.slice(1, -1),
          Object.assign(state.g.context, { payload, headers, status }),
        );
        logWrite(4, "extract: context[%s] = %s", name, JSON.stringify(result));
        state.g.context[name] = result;
      } else {
        logWrite(4, "extract: jsonpath...");
        // jsonpath
        try {
          // https://www.npmjs.com/package/jsonpath
          let result = jsonpath.query(payload, value);
          logWrite(
            4,
            "extract: context[%s] = %s",
            name,
            JSON.stringify(result),
          );
          state.g.context[name] =
            Array.isArray(result) && result.length === 1 ? result[0] : result;
        } catch (e) {
          // gulp
        }
      }
    }
  };
