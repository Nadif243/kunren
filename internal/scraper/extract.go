package scraper

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

// ExtractSentences takes a URL, downloads the HTML, and isolates Japanese sentences.
func ExtractSentences(url string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request: %w", err)
	}
	req.Header.Set("User-Agent", "Renbun/1.0 (Japanese Language Immersion Bot)")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network timeout or failure: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("website rejected request: %d %s", res.StatusCode, res.Status)
	}

	// The Universal Translator
	// This automatically detects Shift_JIS or EUC-JP and streams it as clean UTF-8
	utf8Body, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to convert character encoding: %w", err)
	}

	// Load the sanitized UTF-8 stream into goquery
	doc, err := goquery.NewDocumentFromReader(utf8Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML structure: %w", err)
	}

	// Destroy all furigana reading tags to keep the kanji clean
	doc.Find("rt, rp").Remove()

	hiraganaRegex := regexp.MustCompile(`\p{Hiragana}`)
	var sentences []string

	doc.Find("p, .main_text").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			parts := strings.Split(text, "。")
			for _, part := range parts {
				// Clean up internal newlines caused by Aozora's <br> tags
				cleanPart := strings.ReplaceAll(part, "\n", "")
				cleanPart = strings.ReplaceAll(cleanPart, "\r", "")
				cleanPart = strings.ReplaceAll(cleanPart, " ", "") // Remove full-width Japanese spaces
				cleanPart = strings.TrimSpace(cleanPart)

				if len([]rune(cleanPart)) > 5 && hiraganaRegex.MatchString(cleanPart) {
					sentences = append(sentences, cleanPart+"。")
				}
			}
		}
	})

	return sentences, nil
}
