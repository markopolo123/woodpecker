// Copyright 2018 Drone.IO Inc.
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

package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/woodpecker-ci/woodpecker/server/model"
)

func TestRegistryFind(t *testing.T) {
	store, closer := newTestStore(t, new(model.Registry))
	defer closer()

	err := store.RegistryCreate(&model.Registry{
		RepoID:   1,
		Address:  "index.docker.io",
		Username: "foo",
		Password: "bar",
		Email:    "foo@bar.com",
		Token:    "12345",
	})
	if err != nil {
		t.Errorf("Unexpected error: insert registry: %s", err)
		return
	}

	registry, err := store.RegistryFind(&model.Repo{ID: 1}, "index.docker.io")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := registry.RepoID, int64(1); got != want {
		t.Errorf("Want repo id %d, got %d", want, got)
	}
	if got, want := registry.Address, "index.docker.io"; got != want {
		t.Errorf("Want registry address %s, got %s", want, got)
	}
	if got, want := registry.Username, "foo"; got != want {
		t.Errorf("Want registry username %s, got %s", want, got)
	}
	if got, want := registry.Password, "bar"; got != want {
		t.Errorf("Want registry password %s, got %s", want, got)
	}
	if got, want := registry.Email, "foo@bar.com"; got != want {
		t.Errorf("Want registry email %s, got %s", want, got)
	}
	if got, want := registry.Token, "12345"; got != want {
		t.Errorf("Want registry token %s, got %s", want, got)
	}
}

func TestRegistryList(t *testing.T) {
	store, closer := newTestStore(t, new(model.Registry))
	defer closer()

	assert.NoError(t, store.RegistryCreate(&model.Registry{
		RepoID:   1,
		Address:  "index.docker.io",
		Username: "foo",
		Password: "bar",
	}))
	assert.NoError(t, store.RegistryCreate(&model.Registry{
		RepoID:   1,
		Address:  "foo.docker.io",
		Username: "foo",
		Password: "bar",
	}))

	list, err := store.RegistryList(&model.Repo{ID: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := len(list), 2; got != want {
		t.Errorf("Want %d registries, got %d", want, got)
	}
}

func TestRegistryUpdate(t *testing.T) {
	store, closer := newTestStore(t, new(model.Registry))
	defer closer()

	registry := &model.Registry{
		RepoID:   1,
		Address:  "index.docker.io",
		Username: "foo",
		Password: "bar",
	}
	if err := store.RegistryCreate(registry); err != nil {
		t.Errorf("Unexpected error: insert registry: %s", err)
		return
	}
	registry.Password = "qux"
	if err := store.RegistryUpdate(registry); err != nil {
		t.Errorf("Unexpected error: update registry: %s", err)
		return
	}
	updated, err := store.RegistryFind(&model.Repo{ID: 1}, "index.docker.io")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := updated.Password, "qux"; got != want {
		t.Errorf("Want registry password %s, got %s", want, got)
	}
}

func TestRegistryIndexes(t *testing.T) {
	store, closer := newTestStore(t, new(model.Registry))
	defer closer()

	if err := store.RegistryCreate(&model.Registry{
		RepoID:   1,
		Address:  "index.docker.io",
		Username: "foo",
		Password: "bar",
	}); err != nil {
		t.Errorf("Unexpected error: insert registry: %s", err)
		return
	}

	// fail due to duplicate addr
	if err := store.RegistryCreate(&model.Registry{
		RepoID:   1,
		Address:  "index.docker.io",
		Username: "baz",
		Password: "qux",
	}); err == nil {
		t.Errorf("Unexpected error: duplicate address")
	}
}
