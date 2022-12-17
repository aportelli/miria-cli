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

type AuthToken struct {
	Db      string `json:"dbName"`
	Expire  int    `json:"expire"`
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

type SearchResponse struct {
	Next         string         `json:"next"`
	NextPage     string         `json:"nextPage"`
	Previous     any            `json:"previous"`
	PreviousPage any            `json:"previousPage"`
	Results      []SearchResult `json:"results"`
}

type SearchResult struct {
	InstanceBackupDate string `json:"instanceBackupDate"`
	InstanceId         int    `json:"instanceId"`
	ObjectId           int    `json:"objectId"`
	ObjectName         string `json:"objectName"`
	ObjectPath         string `json:"objectPath"`
	ObjectSize         uint64 `json:"objectSize"`
	ObjectType         string `json:"objectType"`
	RepositoryId       int    `json:"repositoryId"`
}

type ObjectId struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}
