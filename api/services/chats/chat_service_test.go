package chat_services

import "himakiwa/services/database"

func NewMock(userID int) *ChatService {
	return &ChatService{
		ChatSessionRepo:            &database.MockChatSessionRepository{},
		ChatSessionParticipantRepo: &database.MockChatSessionParticipantRepository{},
		ChatRepo:                   &database.MockChatRepository{},
		UserID:                     userID,
	}
}
