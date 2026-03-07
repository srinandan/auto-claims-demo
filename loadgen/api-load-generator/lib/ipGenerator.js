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

import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";
import WeightedRandomSelector from "./weightedRandomSelector.js";
import pop from "./data/population2016.json" with { type: "json" };
import countryAdjustments from "./data/population_weight_adjustments.json" with { type: "json" };

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function getCountrySelector() {
  const countries = pop.records.map((elt) => [
    elt,
    Number(elt.value * (countryAdjustments[elt.name] || 1)),
  ]);
  return new WeightedRandomSelector(countries);
}

function readZone(zoneName) {
  const p = path.resolve(
    __dirname,
    "./data/zones/" + zoneName.toLowerCase() + ".zone",
  );
  // console.log("path for zone: %s", p);
  return fs
    .readFileSync(p, "utf-8")
    .split(/\n/)
    .filter((line) => line && line.trim() !== "")
    .map((line) => [line, Math.pow(2, 32 - Number(line.split(/\//, 2)[1]))]);
}

function randomIpFromCidr(cidr) {
  const [network, mask] = cidr.split(/\//, 2);
  const parsedNetwork = network.split(/\./);
  const mapped = parsedNetwork.map((n) =>
    ("00" + parseInt(n, 10).toString(16)).substr(-2),
  );
  const networkAsHex = parseInt(mapped.join(""), 16);
  const hostMax = Math.pow(2, 32 - mask);
  const chosenHost = Math.floor(Math.random() * hostMax);
  const finalHostAsHex = (
    "00" + (networkAsHex + chosenHost).toString(16)
  ).substr(-8);
  const groups = [];
  for (let i = 0; i < 8; i += 2) {
    groups.push(parseInt(finalHostAsHex.slice(i, i + 2), 16));
  }
  return groups.join("."); // eg, 192.168.1.176
}

const countrySelector = getCountrySelector();

function generateIp() {
  let country;
  do {
    country = countrySelector.select()[0];
  } while (!country.zone);

  //console.log('selected country: %s', country.name);

  const blockSelector = new WeightedRandomSelector(readZone(country.zone));
  const selectedBlock = blockSelector.select()[0];
  //console.log('selected block  : %s', selectedBlock);

  const ip = randomIpFromCidr(selectedBlock);
  //console.log('random IP       : %s', ip);
  return { ip, country: country.name, block: selectedBlock };
}

export default {
  generateIp,
};
