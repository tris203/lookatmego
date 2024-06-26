package parse

import (
	"bufio"
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

type Presentation struct {
	Metadata Metadata
	Slides   []*Slide
}

type Metadata struct {
	Title  string `yaml:"title"`
	Author string `yaml:"author"`
	Date   string `yaml:"date"`
	// extensions []string `yaml:"extensions"`
	// styles     []string `yaml:"styles"`
}

type Slide struct {
	Content string
}

func (p *Presentation) AddSlide(content string) {
	slide := &Slide{Content: content}
	p.Slides = append(p.Slides, slide)
}

func LoadFromFile(content []byte) (*Presentation, error) {
	var p = new(Presentation)
	slides, metadata, err := parseSlides(content)

	if err != nil {
		return nil, err
	}

	p.Metadata = *metadata

	for _, slide := range slides {
		p.AddSlide(slide)
	}
	return p, nil
}

func parseSlides(content []byte) ([]string, *Metadata, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var slides []string
	var metadata Metadata
	var slide bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if slide.Len() > 0 {
				if slides == nil {
					err := yaml.Unmarshal(slide.Bytes(), &metadata)
					if err == nil {
						slide.Reset()
						continue
					}

				}
				slides = append(slides, slide.String())
				slide.Reset()
			}
		} else {
			slide.WriteString(line)
			slide.WriteString("\n")
		}
	}

	// Add the last slide if it exists
	if slide.Len() > 0 {
		slides = append(slides, slide.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return slides, &metadata, nil
}

func SplitSections(content string) []string {
	sections := strings.Split(content, "<!-- stop -->\n")

	return sections
}
