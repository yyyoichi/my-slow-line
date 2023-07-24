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
		users := []participantTUser{
			{
				userID: cs.UserID,
				status: database.Joined,
			},
		}
		for _, userID := range invitedUsers {
			users = append(users, participantTUser{userID, database.Invited})
		}
		return cs.addParticipant(tx, sessionID, users)
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
		users := []participantTUser{}
		for _, invited := range invitedUsers {
			user := participantTUser{
				userID: invited,
				status: database.Invited,
			}
			users = append(users, user)
		}
		cs.addParticipant(tx, sessionID, users)
		return nil
	})
}

type participantTUser struct {
	userID int
	status database.ParticipantStatus
}

func (cs *ChatService) addParticipant(tx *sql.Tx, sessionID int, users []participantTUser) error {
	for _, user := range users {
		participant, err := cs.ChatSessionParticipantRepo.QueryBySessionAndUser(tx, sessionID, user.userID)
		if err != nil {
			return err
		}
		if participant != nil {
			// Participant already exists, update the status to ""
			err = cs.ChatSessionParticipantRepo.UpdateStatus(tx, participant.ID, user.status)
			if err != nil {
				return err
			}
		} else {
			// Participant does not exist, create a new participant with status ""
			err = cs.ChatSessionParticipantRepo.Create(tx, sessionID, user.userID, cs.UserID, user.status)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

// GetInvitedByUserID retrieves a chat session by its ID.
func (cs *ChatService) GetInvitedJoinedByUserID(userID int) ([]database.ChatSessionParticipant, error) {
	var chatSessionParticipant []database.ChatSessionParticipant
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		chatSessionParticipant, err = cs.ChatSessionParticipantRepo.QueryInvitedJoinedByUserID(tx, userID)
		return err
	})
	return chatSessionParticipant, err
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
	var chatSessionParticipant []database.ChatSessionParticipant
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		chatSessionParticipant, err = cs.ChatSessionParticipantRepo.QueryBySessionID(tx, sessionID)
		if err != nil {
			return err
		}
		if ok := cs.incluedPartipants(chatSessionParticipant); !ok {
			return ErrCannotAccess
		}
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
	participants, err := cs.GetChatSessionParticipantBySessionID(sessionID)
	if err != nil {
		return false, err
	}
	return cs.incluedPartipants(participants), nil
}

func (cs *ChatService) incluedPartipants(participants []database.ChatSessionParticipant) bool {
	include := false
	for _, pt := range participants {
		if pt.UserID == cs.UserID && pt.Status == database.Joined {
			include = true
		}
	}
	return include
}
