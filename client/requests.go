/*
Copyright © 2022 Antonin Portelli <antonin.portelli@me.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package client

type AuthRequest struct {
	Db        string `json:"dbName"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	SuperUser bool   `json:"superUser"`
}

type FindInstanceRequest struct {
	RootObjectPath string `json:"rootObjectPath"`
	ResultType     string `json:"resultType"`
	PageSize       int    `json:"pageSize"`
	Criteria       struct {
		Condition string     `json:"condition"`
		Rules     []FindRule `json:"rules"`
	} `json:"criteria"`
}

type FindRule struct {
	Type     string `json:"type"`
	Value    any    `json:"value"`
	Value2   any    `json:"value2"`
	Operator string `json:"operator"`
}
