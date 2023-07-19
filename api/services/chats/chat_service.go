package chat_services

import (
	"database/sql"
	"errors"
	"himakiwa/services/database"
	"time"
)

type ChatService struct {
	ChatSessionRepo            database.ChatSessionRepositoryInterface
	ChatSessionParticipantRepo database.ChatSessionParticipantRepositoryInterface
	ChatRepo                   database.ChatRepositoryInterface
	UserID                     int
}

var (
	ErrCannotAccess = errors.New("cannot access session")
)

func NewChatService(loginUserID int) *ChatService {
	return &ChatService{
		ChatSessionRepo:            &database.ChatSessionRepository{},
		ChatSessionParticipantRepo: &database.ChatSessionParticipantRepository{},
		ChatRepo:                   &database.ChatRepository{},
		UserID:                     loginUserID,
	}
}

// CreateChatSession creates a new chat session.
// [invitedUsers] is users invited to the session by cs.user.
func (cs *ChatService) CreateChatSession(publicKey, name string, invitedUsers []int) error {
	return database.WithTransaction(func(tx *sql.Tx) error {
		sessionID, err := cs.ChatSessionRepo.Create(tx, cs.UserID, publicKey, name)
		if err != nil {
			return err
		}
		return cs.InviteParticipant(sessionID, invitedUsers)
	})
}

// UpdateChatSessionName updates the name of a chat session.
func (cs *ChatService) UpdateChatSessionName(sessionID int, name string) error {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return ErrCannotAccess
	}
	return database.WithTransaction(func(tx *sql.Tx) error {
		err := cs.ChatSessionRepo.UpdateName(tx, sessionID, name)
		if err != nil {
			return err
		}
		return nil
	})
}

// DeleteChatSession deletes a chat session and all related data.
func (cs *ChatService) DeleteChatSession(sessionID int) error {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return ErrCannotAccess
	}
	return database.WithTransaction(func(tx *sql.Tx) error {
		err := cs.ChatSessionParticipantRepo.Delete(tx, sessionID)
		if err != nil {
			return err
		}

		err = cs.ChatRepo.Delete(tx, sessionID)
		if err != nil {
			return err
		}

		err = cs.ChatSessionRepo.Delete(tx, sessionID)
		if err != nil {
			return err
		}
		return nil
	})
}

// InviteParticipant invites a user to join a chat session.
func (cs *ChatService) InviteParticipant(sessionID int, invitedUsers []int) error {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return ErrCannotAccess
	}
	return database.WithTransaction(func(tx *sql.Tx) error {
		for _, invited := range invitedUsers {
			participant, err := cs.ChatSessionParticipantRepo.QueryBySessionAndUser(tx, sessionID, invited)
			if err != nil {
				return err
			}
			if participant != nil {
				// Participant already exists, update the status to "invited"
				err = cs.ChatSessionParticipantRepo.UpdateStatus(tx, participant.ID, database.Invited)
				if err != nil {
					return err
				}
			} else {
				// Participant does not exist, create a new participant with status "invited"
				err = cs.ChatSessionParticipantRepo.Create(tx, sessionID, invited, cs.UserID, database.Invited)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// AcceptInvitation accepts an invitation to join a chat session.
func (cs *ChatService) AcceptInvitation(sessionID int) error {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return ErrCannotAccess
	}
	userID := cs.UserID
	return database.WithTransaction(func(tx *sql.Tx) error {
		participant, err := cs.ChatSessionParticipantRepo.QueryBySessionAndUser(tx, sessionID, userID)
		if err != nil {
			return err
		}

		if participant != nil {
			// Update the participant's status to "joined"
			err = cs.ChatSessionParticipantRepo.UpdateStatus(tx, participant.ID, database.Joined)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// RejectInvitation rejects an invitation to join a chat session.
func (cs *ChatService) RejectInvitation(sessionID int) error {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return ErrCannotAccess
	}
	userID := cs.UserID
	return database.WithTransaction(func(tx *sql.Tx) error {

		participant, err := cs.ChatSessionParticipantRepo.QueryBySessionAndUser(tx, sessionID, userID)
		if err != nil {
			return err
		}
		if participant != nil {
			// Update the participant's status to "rejected"
			err = cs.ChatSessionParticipantRepo.UpdateStatus(tx, participant.ID, database.Rejected)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// SendMessage sends a chat message in a chat session.
func (cs *ChatService) SendMessage(sessionID int, content string) error {
	userID := cs.UserID
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return ErrCannotAccess
	}
	return database.WithTransaction(func(tx *sql.Tx) error {
		err := cs.ChatRepo.Create(tx, sessionID, userID, content)
		if err != nil {
			return err
		}
		return nil
	})
}

// GetChatSessionByID retrieves a chat session by its ID.
func (cs *ChatService) GetChatSessionByID(sessionID int) (*database.ChatSession, error) {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return nil, ErrCannotAccess
	}
	var chatSession *database.ChatSession
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		chatSession, err = cs.ChatSessionRepo.Query(tx, sessionID)
		return err
	})
	return chatSession, err
}

// GetChatSessionParticipantBySessionID retrieves a chat session participant by its sessionID.
func (cs *ChatService) GetChatSessionParticipantBySessionID(sessionID int) ([]database.ChatSessionParticipant, error) {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return nil, ErrCannotAccess
	}
	var chatSessionParticipant []database.ChatSessionParticipant
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		chatSessionParticipant, err = cs.ChatSessionParticipantRepo.QueryBySessionID(tx, sessionID)
		return err
	})
	return chatSessionParticipant, err
}

// GetChatMessagesBySessionID retrieves chat messages for a specific chat session.
func (cs *ChatService) GetChatMessagesBySessionID(sessionID int) ([]database.Chat, error) {
	if enable, err := cs.enableSessionAccess(sessionID); !enable || err != nil {
		return nil, ErrCannotAccess
	}
	var chat []database.Chat
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		chat, err = cs.ChatRepo.QueryBySessionID(tx, sessionID)
		return err
	})
	return chat, err
}

// GetChatMessagesByUserIDWithTimeRange retrieves chat messages for a specific range time.
func (cs *ChatService) GetChatMessagesByTimeRange(endTime time.Time) ([]database.Chat, error) {
	userID := cs.UserID
	var chat []database.Chat
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		chat, err = cs.ChatRepo.QueryByUserIDAndTimeRange(tx, userID, endTime)
		return err
	})
	return chat, err
}

// check access right
// true if user is join the session of [sessionID].
func (cs *ChatService) enableSessionAccess(sessionID int) (bool, error) {
	include := false
	participants, err := cs.GetChatSessionParticipantBySessionID(sessionID)
	if err != nil {
		return false, err
	}
	for _, pt := range participants {
		if pt.UserID == cs.UserID && pt.Status == database.Joined {
			include = true
		}
	}
	return include, nil
}
