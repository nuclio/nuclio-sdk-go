/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nuclio

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	message := "the world is square"
	err := NewErrConflict(message)
	if err.Error() != message {
		t.Fatalf("Bad message: %q != %q", message, err.Error())
	}

	wst, ok := err.(WithStatusCode)
	if !ok {
		t.Fatalf("Not a WithStatusCode error")
	}

	if wst.StatusCode() != http.StatusConflict {
		t.Fatalf("Bad status: %d != %d", wst.StatusCode(), http.StatusConflict)
	}
}

func TestGetByStatusCode(t *testing.T) {
	err := GetByStatusCode(http.StatusCreated)("something")
	if err == nil {
		t.Fatalf("nil error")
	}
	if err.(WithStatusCode).StatusCode() != http.StatusCreated {
		t.Fatalf("unexpected status code ")
	}
}

func TestGetWrapByStatusCode(t *testing.T) {
	err := GetWrapByStatusCode(http.StatusCreated)(errors.New("something bad happened"))
	if err == nil {
		t.Fatalf("nil error")
	}
	if err.(WithStatusCode).StatusCode() != http.StatusCreated {
		t.Fatalf("unexpected status code ")
	}
}
func TestInterface(t *testing.T) {

	// Check that WithStatusCode implements error interface
	var err error = NewErrNotFound("missing page")

	if err == nil {
		t.Fatal("nil error")
	}
}

func ExampleErrNotFound() {
	fmt.Print(ErrNotFound.Error())

	// Output:
	// Not Found
}
