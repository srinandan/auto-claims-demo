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
// --------------------------------------------------------
import parser from "cron-parser";

export default class ScheduleSelector {
  constructor(a) {
    // This implements a scheduled random selector, over an array. The
    // argument a must be an array of arrays. Each inner array is 2
    // elements, with the value as the first element and the schedule in crontab format
    // as the second element. The value will be selected
    // based on the first matching schedule, so 'anomalies' should be listed first.
    //
    // Example usage:
    /*
     var latency = [
      // value      schedule
      ["1000-1500", "* 25-30 * * * *"], //every second from minutes 25 - 30 of every hour of every day of every month
     ];

     var srs = new ScheduleSelector(latency);
     var selected = srs.select();
     console.log(selected);
    */
    this.a = a;
  }

  select() {
    const now = new Date();
    const parserOptions = {
      currentDate: now,
      endDate: new Date(now.getTime() + 1000), //1 second from now
      iterator: true,
    };
    for (let i = 0, L = this.a.length; i < L; i++) {
      //determine if the next minute falls within crontab schedule
      const schedule = this.a[i][1];
      try {
        const interval = parser.parseExpression(schedule, parserOptions);
        while (true) {
          try {
            interval.next();
            // we are within the scheduled interval, so return random integer from selected range
            const range = this.a[i][0];
            return range;
          } catch (e) {
            break;
          }
        }
      } catch (e) {
        console.log("Error: " + e.message);
      }
    }

    // default to the last range if no matching schedules
    const range = this.a[this.a.length - 1][0];
    return range;
  }
}
