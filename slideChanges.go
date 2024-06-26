package main

import (
	"fmt"

	"github.com/tris203/lookatmego/parse"
)

func (m *model) ResetSection() {
	m.CurrenSlideSection = 0
}

func (m *model) NextSection() {
	m.CurrenSlideSection++
}

func (m *model) PrevSection() {
	m.CurrenSlideSection--
}

func (m *model) NextSlide() {
	m.ResetSection()
	m.CurrentSlide++
}

func (m *model) PrevSlide() {
	m.ResetSection()
	m.CurrentSlide--
}

func (m *model) Next() error {
	sectionCount := len(parse.SplitSections(m.presentation.Slides[m.CurrentSlide].Content))
	if m.CurrenSlideSection < sectionCount-1 {
		m.NextSection()
		return nil
	}
	if m.CurrentSlide < len(m.presentation.Slides)-1 {
		m.NextSlide()
		return nil
	}
	return fmt.Errorf("no next slide")
}

func (m *model) Prev() error {
	if m.CurrenSlideSection > 0 {
		m.PrevSection()
		return nil
	}
	if m.CurrentSlide > 0 {
		m.PrevSlide()
		return nil
	}
	return fmt.Errorf("no previous slide")
}
