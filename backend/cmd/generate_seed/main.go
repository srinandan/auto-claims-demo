// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"example.com/claims-app/pkg/seeder"
)

func main() {
	policies, err := seeder.FetchPolicyHolders()
	if err != nil {
		log.Fatalf("Failed to fetch policies: %v", err)
	}

	fmt.Println("-- Seed data for policy_holders")

	now := time.Now().Format("2006-01-02 15:04:05")

	for _, p := range policies {
		// Escape strings for SQL
		escape := func(s string) string {
			return strings.ReplaceAll(s, "'", "''")
		}

		// Format time
		date := p.PolicyBindDate.Format("2006-01-02 15:04:05")

		fmt.Printf("INSERT INTO policy_holders (created_at, updated_at, policy_number, first_name, last_name, months_as_customer, age, policy_bind_date, policy_state, policy_csl, policy_deductible, policy_annual_premium, umbrella_limit, insured_zip, insured_sex, insured_education_level, insured_occupation, insured_hobbies, insured_relationship, capital_gains, capital_loss, auto_make, auto_model, auto_year) VALUES ('%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', %d, %f, %d, %d, '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', %d) ON CONFLICT (policy_number) DO NOTHING;\n",
			now, now,
			escape(p.PolicyNumber),
			escape(p.FirstName),
			escape(p.LastName),
			p.MonthsAsCustomer,
			p.Age,
			date,
			escape(p.PolicyState),
			escape(p.PolicyCSL),
			p.PolicyDeductible,
			p.PolicyAnnualPremium,
			p.UmbrellaLimit,
			p.InsuredZip,
			escape(p.InsuredSex),
			escape(p.InsuredEducationLevel),
			escape(p.InsuredOccupation),
			escape(p.InsuredHobbies),
			escape(p.InsuredRelationship),
			p.CapitalGains,
			p.CapitalLoss,
			escape(p.AutoMake),
			escape(p.AutoModel),
			p.AutoYear,
		)
	}
}
