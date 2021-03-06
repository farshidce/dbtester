// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	gax "github.com/googleapis/gax-go"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
)

// runWithRetry calls the function until it returns nil or a non-retryable error, or
// the context is done.
func runWithRetry(ctx context.Context, call func() error) error {
	var backoff gax.Backoff // use defaults for gax exponential backoff
	for {
		err := call()
		if err == nil {
			return nil
		}
		e, ok := err.(*googleapi.Error)
		if !ok {
			return err
		}
		// Retry on 429 and 5xx, according to
		// https://cloud.google.com/storage/docs/exponential-backoff.
		if e.Code == 429 || (e.Code >= 500 && e.Code < 600) {
			if err := gax.Sleep(ctx, backoff.Pause()); err != nil {
				return err
			}
			continue
		}
		return err
	}
}
