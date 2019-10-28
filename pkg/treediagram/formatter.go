package treediagram

import (
	"github.com/jukeizu/contract"
	"github.com/machinebox/sdk-go/textbox"
)

func FormatSentimentReaction(request contract.Request, analysis *textbox.Analysis) *contract.Reaction {
	averageSentiment := computeAverageSentiment(analysis)

	emoji := lookupEmojiForSentiment(averageSentiment)
	if emoji == "" {
		return nil
	}

	reaction := contract.Reaction{
		MessageId: request.Id,
		ChannelId: request.ChannelId,
		EmojiId:   emoji,
	}

	return &reaction
}

func lookupEmojiForSentiment(sentiment float64) string {
	switch s := sentiment; {
	case s <= 0.4:
		return "😒"
	case s >= 0.6 && s < 0.9:
		return "😄"
	case s >= 0.9:
		return "😊"
	}

	return ""
}

func computeAverageSentiment(analysis *textbox.Analysis) float64 {
	sum := 0.0
	total := 0.0

	for _, sentence := range analysis.Sentences {
		sum += sentence.Sentiment
		total++
	}

	return sum / total
}
