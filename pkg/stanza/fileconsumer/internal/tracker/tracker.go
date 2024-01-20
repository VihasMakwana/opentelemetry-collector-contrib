// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracker // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/tracker"

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/checkpoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/fileset"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/fingerprint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/reader"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"go.uber.org/zap"
)

var errTooManyActiveFiles = errors.New("number of actively read files exceeds max_concurrent_files")

type Tracker struct {
	*zap.SugaredLogger
	activeFiles        *fileset.Fileset[*reader.Reader]
	openFiles          *fileset.Fileset[*reader.Reader]
	closedFiles        *fileset.Fileset[*reader.Metadata]
	ReaderFactory      reader.Factory
	maxConcurrentFiles int

	// This value approximates the expected number of files which we will find in a single poll cycle.
	// It is updated each poll cycle using a simple moving average calculation which assigns 20% weight
	// to the most recent poll cycle.
	// It is used to regulate the size of knownFiles. The goal is to allow knownFiles
	// to contain checkpoints from a few previous poll cycles, but not grow unbounded.
	MovingAverageMatches int
}

func New(logger *zap.SugaredLogger, maxConcurrentFiles int, readerFactory reader.Factory) *Tracker {
	return &Tracker{
		SugaredLogger:      logger,
		ReaderFactory:      readerFactory,
		openFiles:          fileset.New[*reader.Reader](100),
		activeFiles:        fileset.New[*reader.Reader](maxConcurrentFiles),
		closedFiles:        fileset.New[*reader.Metadata](100),
		maxConcurrentFiles: maxConcurrentFiles,
	}
}

func (t *Tracker) ReadFile(path string) {
	if t.activeFiles.Len()+t.openFiles.Len() > t.maxConcurrentFiles {
		// pop the oldest open files and add them to closed filelist
		if r, err := t.openFiles.PopN(1); err == nil {
			t.closedFiles.Add(r[0].Close())
		} else {
			t.Errorw("cannot open file", zap.Error(errTooManyActiveFiles))
			return
		}
	}
	fp, file := t.makeFingerprint(path)
	if fp == nil {
		return
	}
	// Exclude duplicate paths with the same content. This can happen when files are
	// being rotated with copy/truncate strategy. (After copy, prior to truncate.)
	for _, r := range t.activeFiles.Get() {
		if fp.Equal(r.Fingerprint) {
			if err := file.Close(); err != nil {
				t.Debugw("problem closing file", zap.Error(err))
			}
			return
		}
	}
	r, err := t.newReader(file, fp)
	if err != nil {
		t.Errorw("Failed to create reader", zap.Error(err))
		return
	}
	t.activeFiles.Add(r)
}

func (t *Tracker) ActiveFiles() []*reader.Reader {
	return t.activeFiles.Get()
}

func (t *Tracker) OpenFiles() []*reader.Reader {
	return t.openFiles.Get()
}

func (t *Tracker) ClosedFiles() []*reader.Metadata {
	return t.closedFiles.Get()
}

func (t *Tracker) FromBeginning() {
	t.ReaderFactory.FromBeginning = true
}

func (t *Tracker) newReader(file *os.File, fp *fingerprint.Fingerprint) (*reader.Reader, error) {
	// Find a prefix match in previous poll's open fileset
	if reader := t.openFiles.Match(fp); reader != nil {
		return t.ReaderFactory.NewReaderFromMetadata(file, reader.Close())
	}
	// Find a prefix match in previous known files
	if metadata := t.closedFiles.Match(fp); metadata != nil {
		r, err := t.ReaderFactory.NewReaderFromMetadata(file, metadata)
		return r, err
	}
	// If we don't match any previously known files, create a new reader from scratch
	t.Infow("Started watching file", "path", file.Name())
	return t.ReaderFactory.NewReader(file, fp)
}

func (t *Tracker) makeFingerprint(path string) (*fingerprint.Fingerprint, *os.File) {
	file, err := os.Open(path) // #nosec - operator must read in files defined by user
	if err != nil {
		t.Errorw("Failed to open file", zap.Error(err))
		return nil, nil
	}

	fp, err := t.ReaderFactory.NewFingerprint(file)
	if err != nil {
		if err = file.Close(); err != nil {
			t.Debugw("problem closing file", zap.Error(err))
		}
		return nil, nil
	}

	if len(fp.FirstBytes) == 0 {
		// Empty file, don't read it until we can compare its fingerprint
		if err = file.Close(); err != nil {
			t.Debugw("problem closing file", zap.Error(err))
		}
		return nil, nil
	}
	return fp, file
}

func (t *Tracker) Persist(persister operator.Persister) {
	allCheckpoints := make([]*reader.Metadata, 0, t.closedFiles.Len()+t.openFiles.Len())
	for _, r := range t.openFiles.Get() {
		allCheckpoints = append(allCheckpoints, r.Metadata)
	}
	for _, r := range t.closedFiles.Get() {
		allCheckpoints = append(allCheckpoints, r)
	}
	if err := checkpoint.Save(context.Background(), persister, allCheckpoints); err != nil {
		t.Errorw("save offsets", zap.Error(err))
	}
}

func (t *Tracker) Load(persister operator.Persister) error {
	offsets, err := checkpoint.Load(context.Background(), persister)
	if err != nil {
		return fmt.Errorf("read known files from database: %w", err)
	}
	if len(offsets) > 0 {
		t.ReaderFactory.FromBeginning = true
		t.Infow("Resuming from previously known offset(s). 'start_at' setting is not applicable.")
	}
	t.closedFiles.Add(offsets...)
	return nil
}

func (t *Tracker) closePreviousFiles() {
	if t.closedFiles.Len() > 4*t.MovingAverageMatches {
		t.closedFiles.PopN(t.MovingAverageMatches)
	}
	readers := t.openFiles.Reset()
	for _, r := range readers {
		t.closedFiles.Add(r.Close())
	}
}
