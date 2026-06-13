package migrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	notionAPIVersion = "2022-06-28"
	notionMaxChunk   = 1990
	notionMaxBlocks  = 100
)

type notionRichText struct {
	Type string         `json:"type"`
	Text notionTextBody `json:"text"`
}

type notionTextBody struct {
	Content string `json:"content"`
}

type notionParagraphBlock struct {
	Object    string              `json:"object"`
	Type      string              `json:"type"`
	Paragraph notionRichTextBlock `json:"paragraph"`
}

type notionRichTextBlock struct {
	RichText []notionRichText `json:"rich_text"`
}

type notionCreatePageRequest struct {
	Parent     notionParent           `json:"parent"`
	Properties notionPageProperties   `json:"properties"`
	Children   []notionParagraphBlock `json:"children"`
}

type notionParent struct {
	PageID string `json:"page_id"`
}

type notionPageProperties struct {
	Title notionTitleProperty `json:"title"`
}

type notionTitleProperty struct {
	Title []notionRichText `json:"title"`
}

type notionPageResponse struct {
	ID string `json:"id"`
}

type notionAppendBlocksRequest struct {
	Children []notionParagraphBlock `json:"children"`
}

func splitChunks(s string, size int) []string {
	var chunks []string
	runes := []rune(s)
	for len(runes) > 0 {
		if len(runes) <= size {
			chunks = append(chunks, string(runes))
			break
		}
		chunks = append(chunks, string(runes[:size]))
		runes = runes[size:]
	}
	return chunks
}

func buildParagraphBlocks(content string) []notionParagraphBlock {
	chunks := splitChunks(content, notionMaxChunk)
	blocks := make([]notionParagraphBlock, 0, len(chunks))
	for _, chunk := range chunks {
		blocks = append(blocks, notionParagraphBlock{
			Object: "block",
			Type:   "paragraph",
			Paragraph: notionRichTextBlock{
				RichText: []notionRichText{{Type: "text", Text: notionTextBody{Content: chunk}}},
			},
		})
	}
	return blocks
}

func createPage(token, parentPageID, title, content string) (string, error) {
	blocks := buildParagraphBlocks(content)

	first := blocks
	var rest []notionParagraphBlock
	if len(blocks) > notionMaxBlocks {
		first = blocks[:notionMaxBlocks]
		rest = blocks[notionMaxBlocks:]
	}

	reqBody := notionCreatePageRequest{
		Parent: notionParent{PageID: parentPageID},
		Properties: notionPageProperties{
			Title: notionTitleProperty{
				Title: []notionRichText{{Type: "text", Text: notionTextBody{Content: title}}},
			},
		},
		Children: first,
	}

	pageID, err := doCreatePage(token, reqBody)
	if err != nil {
		return "", err
	}

	for len(rest) > 0 {
		batch := rest
		if len(rest) > notionMaxBlocks {
			batch = rest[:notionMaxBlocks]
			rest = rest[notionMaxBlocks:]
		} else {
			rest = nil
		}
		if err := appendBlocks(token, pageID, batch); err != nil {
			return pageID, err
		}
	}

	return pageID, nil
}

func doCreatePage(token string, body notionCreatePageRequest) (string, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Notion-Version", notionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("notion API %d: %s", resp.StatusCode, respBody)
	}

	var pageResp notionPageResponse
	if err := json.Unmarshal(respBody, &pageResp); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	return pageResp.ID, nil
}

func appendBlocks(token, blockID string, blocks []notionParagraphBlock) error {
	body := notionAppendBlocksRequest{Children: blocks}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", blockID), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Notion-Version", notionAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("notion API %d: %s", resp.StatusCode, b)
	}

	return nil
}
