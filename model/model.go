// Copyright 2016 The Citar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type ClosedClassSet map[string]interface{}

var _ gob.GobEncoder = Model{}
var _ gob.GobDecoder = &Model{}

// Model stores a model of the training data.
type Model struct {
	tagNumberer  *StringNumberer
	wordTagFreqs map[string]map[Tag]int
	unigramFreqs map[Unigram]int
	bigramFreqs  map[Bigram]int
	trigramFreqs map[Trigram]int
	closedClass  ClosedClassSet
}

type encodedModel struct {
	TagNumberer  *StringNumberer
	WordTagFreqs map[string]map[Tag]int
	UnigramFreqs map[Unigram]int
	BigramFreqs  map[Bigram]int
	TrigramFreqs map[Trigram]int
	ClosedClass  ClosedClassSet
}

func newModel(tagNumberer *StringNumberer, wordTagFreqs map[string]map[Tag]int,
	unigramFreqs map[Unigram]int, bigramFreqs map[Bigram]int,
	trigramFreqs map[Trigram]int, closedClass ClosedClassSet) Model {
	return Model{
		tagNumberer:  tagNumberer,
		wordTagFreqs: wordTagFreqs,
		unigramFreqs: unigramFreqs,
		bigramFreqs:  bigramFreqs,
		trigramFreqs: trigramFreqs,
		closedClass:  closedClass,
	}
}

func (m Model) ClosedClassTags() ClosedClassSet {
	return m.closedClass
}

// WordTagFreqs returns the word-tag frequencies in the training data.
func (m Model) WordTagFreqs() map[string]map[Tag]int {
	return m.wordTagFreqs
}

// UnigramFreqs returns the tag unigram frequencies in the training data.
func (m Model) UnigramFreqs() map[Unigram]int {
	return m.unigramFreqs
}

// BigramFreqs returns the tag bigram frequencies in the training data.
func (m Model) BigramFreqs() map[Bigram]int {
	return m.bigramFreqs
}

// TrigramFreqs returns the tag trigram frequencies in the training data.
func (m Model) TrigramFreqs() map[Trigram]int {
	return m.trigramFreqs
}

// TagNumberer returns the tag <-> number bijection.
func (m Model) TagNumberer() *StringNumberer {
	return m.tagNumberer
}

// String returns a summary of the model as a string.
func (m Model) String() string {
	return fmt.Sprintf("%d words, %d unigrams, %d bigrams, %d trigrams", len(m.wordTagFreqs),
		len(m.unigramFreqs), len(m.bigramFreqs), len(m.trigramFreqs))
}

// GobDecode decodes a Model from a gob.
func (m *Model) GobDecode(data []byte) error {
	var em encodedModel
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&em); err != nil {
		return err
	}

	m.tagNumberer = em.TagNumberer
	m.wordTagFreqs = em.WordTagFreqs
	m.unigramFreqs = em.UnigramFreqs
	m.bigramFreqs = em.BigramFreqs
	m.trigramFreqs = em.TrigramFreqs
	m.closedClass = em.ClosedClass

	return nil
}

// GobEncode encodes a Model as a gob.
func (m Model) GobEncode() ([]byte, error) {
	em := encodedModel{
		TagNumberer:  m.tagNumberer,
		WordTagFreqs: m.wordTagFreqs,
		UnigramFreqs: m.unigramFreqs,
		BigramFreqs:  m.bigramFreqs,
		TrigramFreqs: m.trigramFreqs,
		ClosedClass:  m.closedClass,
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(em); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
