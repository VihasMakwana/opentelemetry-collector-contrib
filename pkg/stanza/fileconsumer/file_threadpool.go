package fileconsumer // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
)

type readerWrapper struct {
	reader *Reader
	fp     *Fingerprint
}

// poll checks all the watched paths for new entries
func (m *Manager) pollConcurrent(ctx context.Context) {
	// Increment the generation on all known readers
	// This is done here because the next generation is about to start

	m.knownFilesLock.Lock()
	for i := 0; i < len(m.knownFiles); i++ {
		m.knownFiles[i].generation++
	}
	m.knownFilesLock.Unlock()

	// Get the list of paths on disk
	matches := m.finder.FindFiles()
	m.consumeConcurrent(ctx, matches)
	m.clearCurrentFingerprints()

	// Any new files that appear should be consumed entirely
	m.readerFactory.fromBeginning = true
	m.syncLastPollFilesConcurrent(ctx)
}

func (m *Manager) worker(ctx context.Context) {
	defer m.workerWg.Done()
	for {
		chanData, ok := <-m.readerChan

		if !ok {
			return
		}
		r, fp := chanData.reader, chanData.fp

		if !m.readToEnd(r, ctx) {
			// Save off any files that were not fully read or if deleteAfterRead is disabled
			m.saveReaders <- readerWrapper{reader: r, fp: fp}
		}
		m.removePath(fp)
	}
}

func (m *Manager) makeReaderConcurrent(filePath string) (*Reader, *Fingerprint) {
	fp, file := m.makeFingerprint(filePath)
	if fp == nil {
		return nil, nil
	}

	// check if the current file is already being consumed
	if m.isCurrentlyConsuming(fp) {
		if err := file.Close(); err != nil {
			m.Errorf("problem closing file", "file", file.Name())
		}
		return nil, nil
	}

	// Exclude any empty fingerprints or duplicate fingerprints to avoid doubling up on copy-truncate files
	if m.checkDuplicates(fp) {
		if err := file.Close(); err != nil {
			m.Errorf("problem closing file", "file", file.Name())
		}
		return nil, nil
	}
	m.currentFps = append(m.currentFps, fp)

	reader, err := m.newReaderConcurrent(file, fp)
	if err != nil {
		m.Errorw("Failed to create reader", zap.Error(err))
		return nil, nil
	}
	return reader, fp
}

func (m *Manager) consumeConcurrent(ctx context.Context, paths []string) {
	m.Debug("Consuming files")
	m.clearOldReadersConcurrent(ctx)
	for _, path := range paths {
		reader, fp := m.makeReaderConcurrent(path)
		if reader != nil {
			// add path and fingerprint as it's not consuming
			m.trieLock.Lock()
			m.trie.Put(fp.FirstBytes, true)
			m.trieLock.Unlock()
			m.readerChan <- readerWrapper{reader: reader, fp: fp}
		}
	}
}

func (m *Manager) isCurrentlyConsuming(fp *Fingerprint) bool {
	m.trieLock.RLock()
	defer m.trieLock.RUnlock()
	return m.trie.Get(fp.FirstBytes) != nil
}

func (m *Manager) removePath(fp *Fingerprint) {
	m.trieLock.Lock()
	defer m.trieLock.Unlock()
	m.trie.Delete(fp.FirstBytes)
}

// saveReadersConcurrent adds the readers from this polling interval to this list of
// known files and removes the fingerprint from the TRIE
func (m *Manager) saveReadersConcurrent(ctx context.Context) {
	defer m._workerWg.Done()
	// Add readers from the current, completed poll interval to the list of known files
	for {
		select {
		case readerWrapper, ok := <-m.saveReaders:
			if !ok {
				return
			}
			m.knownFilesLock.Lock()
			m.knownFiles = append(m.knownFiles, readerWrapper.reader)
			m.knownFilesLock.Unlock()
			m.removePath(readerWrapper.fp)
		}
	}
}

func (m *Manager) clearOldReadersConcurrent(ctx context.Context) {
	m.knownFilesLock.Lock()
	defer m.knownFilesLock.Unlock()
	// Clear out old readers. They are sorted such that they are oldest first,
	// so we can just find the first reader whose poll cycle is less than our
	// limit i.e. last 3 cycles, and keep every reader after that
	oldReaders := make([]*Reader, 0)
	for i := 0; i < len(m.knownFiles); i++ {
		reader := m.knownFiles[i]
		if reader.generation <= 2 {
			oldReaders = m.knownFiles[:i]
			m.knownFiles = m.knownFiles[i:]
			break
		}
	}
	var lostWG sync.WaitGroup
	for _, reader := range oldReaders {
		lostWG.Add(1)
		go func(r *Reader) {
			defer lostWG.Done()
			r.ReadToEnd(ctx)
			r.Close()
		}(reader)
	}
	lostWG.Wait()
}

func (m *Manager) newReaderConcurrent(file *os.File, fp *Fingerprint) (*Reader, error) {
	// Check if the new path has the same fingerprint as an old path
	if oldReader, ok := m.findFingerprintMatchConcurrent(fp); ok {
		return m.readerFactory.copy(oldReader, file)
	}

	// If we don't match any previously known files, create a new reader from scratch
	return m.readerFactory.newReader(file, fp.Copy())
}

func (m *Manager) findFingerprintMatchConcurrent(fp *Fingerprint) (*Reader, bool) {
	// Iterate backwards to match newest first
	m.knownFilesLock.Lock()
	defer m.knownFilesLock.Unlock()

	return m.findFingerprintMatch(fp)
}

// syncLastPollFiles syncs the most recent set of files to the database
func (m *Manager) syncLastPollFilesConcurrent(ctx context.Context) {
	m.knownFilesLock.RLock()
	defer m.knownFilesLock.RUnlock()

	m.syncLastPollFiles(ctx)
}
