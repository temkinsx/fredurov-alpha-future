package llm

/*

func (s *Service) Reply(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, userText string, scenarioCode *string) (*domain.Message, error) {
	if userText == "" {
		return nil, nil
	}

	chat, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat.UserID != userID {
		return nil, errors.New("wrong userID for this chat")
	}

	rawMsgHistory, err := s.msgRepo.GetLastN(ctx, chatID, 10)
	if err != nil {
		return nil, err
	}

	var msgHistory strings.Builder
	for _, msg := range rawMsgHistory {
		msgHistory.WriteString(msg.String())
		msgHistory.WriteString("\n")
	}

}


*/
