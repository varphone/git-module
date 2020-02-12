// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// DefaultHooksDir is the default directory for Git hooks.
const DefaultHooksDir = "hooks"

// NewHook creates and returns a new hook with given name. Update method must be called
// to actually save the hook to disk.
func (r *Repository) NewHook(name HookName) *Hook {
	return &Hook{
		name: name,
		path: filepath.Join(r.path, DefaultHooksDir, string(name)),
	}
}

// Hook returns a Git hook by given name in the repository. It returns an os.ErrNotExist
// if both active and sample hook do not exist.
func (r *Repository) Hook(name HookName) (*Hook, error) {
	// 1. Check if there is an active hook.
	fpath := filepath.Join(r.path, DefaultHooksDir, string(name))
	if isFile(fpath) {
		p, err := ioutil.ReadFile(fpath)
		if err != nil {
			return nil, err
		}
		return &Hook{
			name:    name,
			path:    fpath,
			content: string(p),
		}, nil
	}

	// 2. Check if a sample file exists.
	spath := fpath + ".sample"
	if isFile(spath) {
		p, err := ioutil.ReadFile(spath)
		if err != nil {
			return nil, err
		}
		return &Hook{
			name:     name,
			path:     fpath,
			isSample: true,
			content:  string(p),
		}, nil
	}

	return nil, os.ErrNotExist
}

// Hooks returns a list of Git hooks found in the repository. It may return an empty slice
// when no hooks found.
func (r *Repository) Hooks() ([]*Hook, error) {
	if !isDir(path.Join(r.path, DefaultHooksDir)) {
		return []*Hook{}, nil
	}

	hooks := make([]*Hook, 0, len(ServerSideHooks))
	for _, name := range ServerSideHooks {
		h, err := r.Hook(name)
		if err != nil {
			if err == os.ErrNotExist {
				continue
			}
			return nil, err
		}
		hooks = append(hooks, h)
	}
	return hooks, nil
}
