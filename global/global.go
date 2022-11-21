/*
 * global.go, part of miria-cli (https://github.com/aportelli/miria-cli)
 * Copyright (C) 2022 Antonin Portelli
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free Software
 * Foundation, either version 3 of the License, or (at your option) any later
 * version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of  MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package global

import (
	"encoding/json"
	"fmt"
	"os"
)

func AppCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return cacheDir + "/miria-cli", nil
}

func JsonMapToStruct(pt any, dict map[string]any) error {
	jbuf, err := json.Marshal(dict)
	if err != nil {
		return err
	}
	json.Unmarshal(jbuf, pt)
	return nil
}

func PrettyPrintResponse(dict map[string]any) {
	jbuf, _ := json.MarshalIndent(dict, "", "  ")
	fmt.Println(string(jbuf))
}
