package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"ai-clipping-backend/models"
)

// GeminiAnalyzer uses Google Gemini API to find highlights
type GeminiAnalyzer struct {
	apiKey      string
	videoType   models.VideoType
	preferences models.ClipPreferences
}

// NewGeminiAnalyzer creates a new Gemini-based highlight analyzer
func NewGeminiAnalyzer(apiKey string, videoType models.VideoType, prefs models.ClipPreferences) *GeminiAnalyzer {
	return &GeminiAnalyzer{
		apiKey:      apiKey,
		videoType:   videoType,
		preferences: prefs,
	}
}

// FindHighlights analyzes transcript using Gemini API to find interesting moments
func (ga *GeminiAnalyzer) FindHighlights(transcript []models.TranscriptSegment) ([]models.Highlight, error) {
	if ga.apiKey == "" {
		return nil, fmt.Errorf("gemini API key not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Create Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(ga.apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %v", err)
	}
	defer client.Close()

	// Use Gemini 2.5 Flash (latest stable, fast, free)
	// Supports up to 1M tokens - perfect for video transcripts
	model := client.GenerativeModel("gemini-2.5-flash")

	// Build prompt

	prompt := ga.buildPrompt(transcript)

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini api call failed: %v", err)
	}

	// Parse response
	highlights, err := ga.parseResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %v", err)
	}

	log.Printf("Gemini found %d highlights", len(highlights))
	return highlights, nil
}

// buildPrompt creates the prompt for Gemini based on video type
func (ga *GeminiAnalyzer) buildPrompt(transcript []models.TranscriptSegment) string {
	// Convert transcript to readable format
	var transcriptText strings.Builder
	for _, seg := range transcript {
		transcriptText.WriteString(fmt.Sprintf("[%.1f - %.1f] %s\n", seg.Start, seg.End, seg.Text))
	}

	// Video type specific instructions
	typeInstructions := ga.getTypeSpecificInstructions()

	prompt := fmt.Sprintf(`You are an expert video highlight detector. Your task is to analyze a video transcript and identify the most interesting moments that would make great short clips.

VIDEO TYPE: %s
%s

CLIP PREFERENCES:
- Minimum duration: %d seconds
- Maximum duration: %d seconds
- Maximum number of clips: %d
- Target format: Vertical video (9:16) for social media

TRANSCRIPT:
%s

INSTRUCTIONS:
1. Analyze the transcript and find the %d MOST interesting moments
2. Look for:
   - Key information or important points
   - Emotional peaks (excitement, surprise, humor)
   - Actionable advice or tips
   - Dramatic or engaging moments
   - Complete thoughts (not cut mid-sentence)
3. Each highlight should be self-contained and understandable on its own
4. Ensure highlights are between %d and %d seconds long
5. Avoid silence or filler words
6. For each highlight, write a CATCHY, click-worthy TITLE (max 70 characters) ready to be used directly as the video title when uploaded to YouTube Shorts / TikTok / Reels. Write the title in the SAME LANGUAGE as the transcript. Make it engaging but not clickbait-lying.

OUTPUT FORMAT (JSON only, no explanation):
[
  {
    "start": 10.5,
    "end": 25.3,
    "score": 85,
    "title": "Catchy ready-to-upload video title here",
    "reason": "Brief description why this is interesting",
    "keywords": ["key", "words", "found"]
  }
]

Return ONLY the JSON array, nothing else.`,
		ga.videoType,
		typeInstructions,
		ga.preferences.MinDuration,
		ga.preferences.MaxDuration,
		ga.preferences.MaxClips,
		transcriptText.String(),
		ga.preferences.MaxClips,
		ga.preferences.MinDuration,
		ga.preferences.MaxDuration,
	)

	return prompt
}

// getTypeSpecificInstructions returns specific instructions based on video type
func (ga *GeminiAnalyzer) getTypeSpecificInstructions() string {
	switch ga.videoType {
	case models.VideoTypeTutorial:
		return `This is an EDUCATIONAL/TUTORIAL video. Focus on:
- Step-by-step instructions
- Key concepts being explained
- "Aha!" moments where complicated things become clear
- Practical tips and tricks
- Before/after comparisons
- Common mistakes to avoid`

	case models.VideoTypeGaming:
		return `This is a GAMING video. Focus on:
- Epic moments (kills, victories, clutches)
- Funny or unexpected moments
- Skillful plays or combos
- Reactions to exciting events
- Commentary highlights
- "Rage" or intense moments`

	case models.VideoTypePodcast:
		return `This is a PODCAST/INTERVIEW video. Focus on:
- Interesting stories or anecdotes
- Strong opinions or controversial takes
- Funny or emotional moments
- Key insights or advice
- Surprising revelations
- Quotable statements`

	case models.VideoTypeGeneral:
		fallthrough
	default:
		return `This is a GENERAL video. Focus on:
- Most engaging moments
- Information-dense segments
- Emotional or entertaining parts
- Unique or surprising content
- Clear and concise statements`
	}
}

// parseResponse parses Gemini's JSON response into Highlight structs
func (ga *GeminiAnalyzer) parseResponse(resp *genai.GenerateContentResponse) ([]models.Highlight, error) {
	if resp == nil || len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	// Get text from response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Extract JSON from response (sometimes Gemini adds markdown formatting)
	jsonStart := strings.Index(responseText, "[")
	jsonEnd := strings.LastIndex(responseText, "]")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no JSON array found in response: %s", responseText)
	}

	jsonStr := responseText[jsonStart : jsonEnd+1]

	// Parse JSON
	var highlights []models.Highlight
	if err := json.Unmarshal([]byte(jsonStr), &highlights); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v, response: %s", err, jsonStr)
	}

	// Validate and adjust highlights
	for i := range highlights {
		// Ensure duration is within bounds
		duration := highlights[i].End - highlights[i].Start

		if duration < float64(ga.preferences.MinDuration) {
			// Extend to minimum duration
			center := (highlights[i].Start + highlights[i].End) / 2
			halfMin := float64(ga.preferences.MinDuration) / 2
			highlights[i].Start = center - halfMin
			highlights[i].End = center + halfMin
		}

		if duration > float64(ga.preferences.MaxDuration) {
			// Trim to maximum duration
			highlights[i].End = highlights[i].Start + float64(ga.preferences.MaxDuration)
		}

		// Ensure non-negative start time
		if highlights[i].Start < 0 {
			highlights[i].Start = 0
		}

		// Fallback title jika Gemini tidak mengisi
		if strings.TrimSpace(highlights[i].Title) == "" {
			highlights[i].Title = fmt.Sprintf("Highlight %d", i+1)
		}
	}

	return highlights, nil
}
