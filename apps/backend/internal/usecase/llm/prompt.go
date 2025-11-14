package llm

import (
	"strings"
)

const (
	defaultSysPrompt string = "Ты — умный и осторожный ассистент для владельцев микробизнеса." +
		"  \nОтвечай кратко, ясно и структурированно." +
		"  \nЕсли вопрос непонятный — уточняй." +
		"  \nЕсли информации недостаточно — прямо скажи об этом." +
		"  \nНе придумывай факты и не выдумывай данные документов." +
		"  \nЕсли пользователь прикрепил документ — используй только текст, который тебе дали." +
		"  \nДелай ответы простыми и полезными, без воды и лишней формальности." +
		"  \nПиши всегда по-русски."
)

type promptBudget struct {
	MaxTotal int
	Used     int
}

// Take принимает запрос и максимальное количество символов, обрезает его в зависимости от лимитов и возвращает новую строку
func (b *promptBudget) Take(s string, max int) string {
	if s == "" {
		return ""
	}

	remaining := b.MaxTotal - b.Used
	if remaining <= 0 {
		return ""
	}

	if max > remaining {
		max = remaining
	}

	if len(s) > max {
		s = s[:max]
	}

	b.Used += len(s)
	return s
}

// buildPrompt собирает корректный промпт для LLM с системным промптом, историей сообщений, данными с ресурсов и запросом пользователя
func (s *Service) buildPrompt(sysPrompt, msgHistory, documents, userReq string) string {
	if sysPrompt == "" {
		sysPrompt = defaultSysPrompt
	}

	var b strings.Builder

	pBudget := promptBudget{MaxTotal: s.limits.MaxPromptChars}

	sys := pBudget.Take(sysPrompt, 2000)
	b.WriteString("SYSTEM:\n")
	b.WriteString(sys)
	b.WriteString("\n\n")

	if msgHistory != "" {
		msg := pBudget.Take(msgHistory, s.limits.MaxHistoryChars)
		b.WriteString("HISTORY:\n")
		b.WriteString(msg)
		b.WriteString("\n\n")
	}

	if documents != "" {
		doc := pBudget.Take(documents, s.limits.MaxHistoryChars)
		b.WriteString("DOCUMENTS:\n")
		b.WriteString(doc)
		b.WriteString("\n\n")
	}

	if userReq != "" {
		req := pBudget.Take(userReq, s.limits.MaxRequestChars)
		b.WriteString("USER:\n")
		b.WriteString(req)
	}

	prompt := b.String()
	if len(prompt) > s.limits.MaxPromptChars {
		prompt = prompt[:s.limits.MaxPromptChars]
	}

	return prompt
}
