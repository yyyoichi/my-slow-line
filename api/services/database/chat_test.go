package database

import (
	"testing"
	"time"
)

func TestChatSession(t *testing.T) {

	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	defer close()

	// Create the ChatSessionRepository instance
	csr := &ChatSessionRepository{}

	// Test Create method
	_, err := csr.Create(testUser.Id, "public_key", "Test Chat Session")
	if err != nil {
		t.Errorf("Error creating chat session: %s", err.Error())
	}

	// Test QueryByUserID method
	sessions, err := csr.QueryByUserID(testUser.Id)
	if err != nil {
		t.Errorf("Error querying chat sessions: %s", err.Error())
	}

	if len(sessions) != 1 {
		t.Error("Expected 1 chat session, but got", len(sessions))
	}

	session := sessions[0]
	if session.Name != "Test Chat Session" {
		t.Errorf("Expected chat session name 'Test Chat Session', but got '%s'", session.Name)
	}

	// Test Query method
	sessions, err = csr.Query(session.ID)
	if err != nil {
		t.Errorf("Error querying chat sessions: %s", err.Error())
	}

	if len(sessions) != 1 {
		t.Error("Expected 1 chat session, but got", len(sessions))
	}

	session = sessions[0]
	if session.Name != "Test Chat Session" {
		t.Errorf("Expected chat session name 'Test Chat Session', but got '%s'", session.Name)
	}

	// Test UpdateName method
	err = csr.UpdateName(session.ID, "Updated Chat Session")
	if err != nil {
		t.Errorf("Error updating chat session name: %s", err.Error())
	}

	// Verify updated name
	sessions, err = csr.QueryByUserID(testUser.Id)
	if err != nil {
		t.Errorf("Error querying chat sessions: %s", err.Error())
	}

	if len(sessions) != 1 {
		t.Error("Expected 1 chat session, but got", len(sessions))
	}

	session = sessions[0]
	if session.Name != "Updated Chat Session" {
		t.Errorf("Expected chat session name 'Updated Chat Session', but got '%s'", session.Name)
	}

	// Test Delete method
	err = csr.Delete(session.ID)
	if err != nil {
		t.Errorf("Error deleting chat sessions: %s", err.Error())
	}

	// Verify that all chat sessions have been deleted
	sessions, err = csr.QueryByUserID(testUser.Id)
	if err != nil {
		t.Errorf("Error querying chat sessions: %s", err.Error())
	}

	if len(sessions) != 0 {
		t.Error("Expected 0 chat sessions, but got", len(sessions))
	}
}

func TestChatSessionParticipant(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	defer close()

	// Create a mock chat session
	csr := &ChatSessionRepository{}
	sessionID, err := csr.Create(testUser.Id, "public_key", "Test Chat Session")
	if err != nil {
		t.Fatalf("Error creating chat session: %s", err.Error())
	}
	defer csr.Delete(sessionID)

	// Create the ChatSessionParticipantRepository instance
	cspr := &ChatSessionParticipantRepository{}

	// Test Create method
	err = cspr.Create(sessionID, testUser.Id, Joined)
	if err != nil {
		t.Errorf("Error creating chat session participant: %s", err.Error())
	}

	// Test QueryBySessionID method
	participants, err := cspr.QueryBySessionID(sessionID)
	if err != nil {
		t.Errorf("Error querying chat session participants: %s", err.Error())
	}

	if len(participants) != 1 {
		t.Error("Expected 1 chat session participant, but got", len(participants))
	}

	participant := participants[0]
	if participant.UserID != testUser.Id {
		t.Errorf("Expected chat session participant user ID '%d', but got '%d'", testUser.Id, participant.UserID)
	}

	// Test QueryBySessionAndUser method
	participantByID, err := cspr.QueryBySessionAndUser(sessionID, testUser.Id)
	if err != nil {
		t.Errorf("Error querying chat session participant by session and user: %s", err.Error())
	}

	if participantByID == nil {
		t.Error("Expected chat session participant, but got nil")
	} else {
		if participantByID.UserID != testUser.Id {
			t.Errorf("Expected chat session participant user ID '%d', but got '%d'", testUser.Id, participantByID.UserID)
		}
	}

	// Test UpdateStatus method
	err = cspr.UpdateStatus(participant.ID, Rejected)
	if err != nil {
		t.Errorf("Error updating chat session participant status: %s", err.Error())
	}

	// Verify updated status
	participants, err = cspr.QueryBySessionID(sessionID)
	if err != nil {
		t.Errorf("Error querying chat session participants: %s", err.Error())
	}

	if len(participants) != 1 {
		t.Error("Expected 1 chat session participant, but got", len(participants))
	}

	participant = participants[0]
	if participant.Status != Rejected {
		t.Errorf("Expected chat session participant status '%s', but got '%s'", Rejected, participant.Status)
	}

	// Test Delete method
	err = cspr.Delete(sessionID)
	if err != nil {
		t.Errorf("Error deleting chat session participants: %s", err.Error())
	}

	// Verify that all chat session participants have been deleted
	participants, err = cspr.QueryBySessionID(sessionID)
	if err != nil {
		t.Errorf("Error querying chat session participants: %s", err.Error())
	}

	if len(participants) != 0 {
		t.Error("Expected 0 chat session participants, but got", len(participants))
	}
}

func TestChat(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	defer close()

	// Create a mock chat session
	csr := &ChatSessionRepository{}
	_, err := csr.Create(testUser.Id, "public_key", "Test Chat Session")
	if err != nil {
		t.Fatalf("Error creating chat session: %s", err.Error())
	}
	sessions, err := csr.QueryByUserID(testUser.Id)
	if err != nil {
		t.Fatalf("Error querying chat sessions: %s", err.Error())
	}
	if len(sessions) != 1 {
		t.Fatal("Expected 1 chat session, but got", len(sessions))
	}
	session := sessions[0]
	defer csr.Delete(session.ID)

	// Create the ChatRepository instance
	cr := &ChatRepository{}

	// Test Create method
	err = cr.Create(session.ID, testUser.Id, "Test message")
	if err != nil {
		t.Errorf("Error creating chat message: %s", err.Error())
	}

	// Test QueryBySessionID method
	messages, err := cr.QueryBySessionID(session.ID)
	if err != nil {
		t.Errorf("Error querying chat messages: %s", err.Error())
	}

	if len(messages) != 1 {
		t.Error("Expected 1 chat message, but got", len(messages))
	}

	message := messages[0]
	if message.SessionID != session.ID {
		t.Errorf("Expected chat message session ID '1', but got '%d'", message.SessionID)
	}

	// Test QueryByUserIDAndTimeRange method
	endTime := time.Now()
	messages, err = cr.QueryByUserIDAndTimeRange(testUser.Id, endTime)
	if err != nil {
		t.Errorf("Error querying chat messages: %s", err.Error())
	}

	if len(messages) != 1 {
		t.Error("Expected 1 chat message, but got", len(messages))
	}

	message = messages[0]
	if message.UserID != testUser.Id {
		t.Errorf("Expected chat message user ID '%d', but got '%d'", testUser.Id, message.UserID)
	}

	// Test Delete method
	err = cr.Delete(session.ID)
	if err != nil {
		t.Errorf("Error deleting chat messages: %s", err.Error())
	}

	// Verify that all chat messages have been deleted
	messages, err = cr.QueryBySessionID(session.ID)
	if err != nil {
		t.Errorf("Error querying chat messages: %s", err.Error())
	}

	if len(messages) != 0 {
		t.Error("Expected 0 chat messages, but got", len(messages))
	}
}
