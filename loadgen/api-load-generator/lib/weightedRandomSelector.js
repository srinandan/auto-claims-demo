// Copyright © 2013, 2014 Apigee Corp; 2025 Google LLC
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

// weightedRandomSelector.js
//
// simple weighted random selector.
//
// All rights reserved.
//
// created: Fri, 26 Jul 2013  11:14
// last saved: <2025-August-26 17:39:17>
//
// --------------------------------------------------------

export default class WeightedRandomSelector {
  constructor(a) {
    // This implements a weighted random selector, over an array. The
    // argument a must be an array of arrays. Each inner array is 2
    // elements, with the item value as the first element and the weight
    // for that value as the second. The items need not be in any
    // particular order.  Example usage:
    //
    // var fish = [
    // //  name      weight
    //   ["Shark",      3],
    //   ["Shrimp",     50],
    //   ["Sardine",    10],
    //   ["Herring",    20],
    //   ["Anchovies",  10],
    //   ["Mackerel",   50],
    //   ["Tuna",       8]
    // ];
    //
    // var wrs = new WeightedRandomSelector(fish);
    // var selected = wrs.select();
    // var fishType = selected[0];
    //
    this.totalWeight = 0;
    this.a = a;
    this.selectionCounts = [];
    this.weightThreshold = [];
    // initialize
    for (let i = 0, L = a.length; i < L; i++) {
      this.totalWeight += a[i][1];
      this.weightThreshold[i] = this.totalWeight;
      this.selectionCounts[i] = 0;
    }
  }

  select() {
    // select a random value
    const R = Math.floor(Math.random() * this.totalWeight);

    // now find the bucket that R value falls into.
    for (let i = 0, L = this.a.length; i < L; i++) {
      if (R < this.weightThreshold[i]) {
        this.selectionCounts[i]++;
        return this.a[i];
      }
    }
    return this.a[this.a.length - 1];
  }
}
