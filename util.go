/*
 * Copyright (C) 2020 Yunify, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this work except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package edge_driver_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

const (
	device_id = "orgi"
	thing_id  = "thid"
)

func parseToken(tokenString string) (string, string, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if token == nil || token.Claims == nil {
		return "", "", errors.New("parse token fail")
	}
	payload := token.Claims.(jwt.MapClaims)
	if err := payload.Valid(); err != nil {
		return "", "", err
	}
	id, ok := payload[device_id].(string)
	if !ok {
		return "", "", errors.New("device id type error")
	}
	thingId, ok := payload[thing_id].(string)
	if !ok {
		return "", "", errors.New("device id type error")
	}
	return id, thingId, nil
}

func wait(f func() error) <-chan error {
	done := make(chan error)
	go func() {
		done <- f()
	}()
	return done
}
func isUserDevice(id string) bool {
	return id == userThingId
}
