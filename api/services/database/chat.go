package database

import (
	"database/sql"
	"time"
)

type ChatSession struct {
	ID        int
	UserID    int
	PublicKey string
	Name      string
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type ChatSessionRepository struct{}

func (csr *ChatSessionRepository) Query(tx *sql.Tx, sessionID int) (*ChatSession, error) {
	// query
	query := `SELECT id, user_id, public_key, name, create_at, update_at, deleted FROM chat_sessions WHERE id = ?`
	row := tx.QueryRow(query, sessionID)

	// result
	result := &ChatSession{}
	err := row.Scan(&result.ID, &result.UserID, &result.PublicKey, &result.Name, &result.CreateAt, &result.UpdateAt, &result.Deleted)
	if err != nil {
			return nil, err
		}
		results = append(results, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (csr *ChatSessionRepository) QueryByUserID(tx *sql.Tx, userID int) ([]ChatSession, error) {
	// query
	query := `SELECT id, user_id, public_key, name, create_at, update_at, deleted FROM chat_sessions WHERE user_id = ?`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []ChatSession
	for rows.Next() {
		cs := ChatSession{}
		err := rows.Scan(&cs.ID, &cs.UserID, &cs.PublicKey, &cs.Name, &cs.CreateAt, &cs.UpdateAt, &cs.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (csr *ChatSessionRepository) UpdateName(tx *sql.Tx, id int, name string) error {
	// query
	query := `UPDATE chat_sessions SET name = ?, update_at = NOW() WHERE id = ?`
	_, err := DB.Exec(query, name, id)
	if err != nil {
		return err
	}
	return nil
}

func (csr *ChatSessionRepository) Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error) {
	// query
	query := `INSERT INTO chat_sessions (user_id, public_key, name) VALUES (?, ?, ?)`
	result, err := DB.Exec(query, userID, publicKey, name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (csr *ChatSessionRepository) Delete(tx *sql.Tx, sessionID int) error {
	// query
	query := `DELETE FROM chat_sessions WHERE id = ?`
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////
//chat session participants//
/////////////////////////////

type ChatSessionParticipant struct {
	ID        int
	SessionID int
	UserID    int
	Status    ParticipantStatus
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type ParticipantStatus string

const (
	Invited  ParticipantStatus = "invited"
	Joined   ParticipantStatus = "joined"
	Rejected ParticipantStatus = "rejected"
)

type ChatSessionParticipantRepository struct{}

func (cspr *ChatSessionParticipantRepository) QueryBySessionID(tx *sql.Tx, sessionID int) ([]ChatSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted FROM chat_session_participants WHERE chat_session_id = ?`
	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []ChatSessionParticipant
	for rows.Next() {
		csp := ChatSessionParticipant{}
		err := rows.Scan(&csp.ID, &csp.SessionID, &csp.UserID, &csp.Status, &csp.CreateAt, &csp.UpdateAt, &csp.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, csp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (cspr *ChatSessionParticipantRepository) QueryBySessionAndUser(tx *sql.Tx, sessionID int, userID int) (*ChatSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted FROM chat_session_participants WHERE chat_session_id = ? AND user_id = ?`
	row := DB.QueryRow(query, sessionID, userID)

	// result
	csp := ChatSessionParticipant{}
	err := row.Scan(&csp.ID, &csp.SessionID, &csp.UserID, &csp.Status, &csp.CreateAt, &csp.UpdateAt, &csp.Deleted)
	if err != nil {
		return nil, err
	}

	return &csp, nil
}

func (cspr *ChatSessionParticipantRepository) QueryInvitedByUserID(tx *sql.Tx, userID int) ([]ChatSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted 
	FROM chat_session_participants 
	WHERE user_id = ?
		AND status = 'invited'`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	// result
	var results []ChatSessionParticipant
	for rows.Next() {
		csp := ChatSessionParticipant{}
		err := rows.Scan(&csp.ID, &csp.SessionID, &csp.UserID, &csp.Status, &csp.CreateAt, &csp.UpdateAt, &csp.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, csp)
	}
	return results, nil
}

// [invite_user_id] is ID of the user invited to the session
func (cspr *ChatSessionParticipantRepository) Create(tx *sql.Tx, sessionID int, userID, inviteUserID int, status ParticipantStatus) error {
	// query
	query := `INSERT INTO chat_session_participants (chat_session_id, user_id, status) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, sessionID, userID, status)
	if err != nil {
		return err
	}
	return nil
}

func (cspr *ChatSessionParticipantRepository) Delete(tx *sql.Tx, sessionID int) error {
	// query
	query := `DELETE FROM chat_session_participants WHERE chat_session_id = ?`
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (cspr *ChatSessionParticipantRepository) UpdateStatus(tx *sql.Tx, id int, status ParticipantStatus) error {
	// query
	query := `UPDATE chat_session_participants SET status = ? WHERE id = ?`
	_, err := DB.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////
//         chats           //
/////////////////////////////

type Chat struct {
	ID        int
	SessionID int
	UserID    int
	Content   string
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type ChatRepository struct{}

func (cr *ChatRepository) QueryBySessionID(tx *sql.Tx, sessionID int) ([]Chat, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, content, create_at, update_at, deleted FROM chats WHERE chat_session_id = ?`
	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []Chat
	for rows.Next() {
		c := Chat{}
		err := rows.Scan(&c.ID, &c.SessionID, &c.UserID, &c.Content, &c.CreateAt, &c.UpdateAt, &c.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// QueryByUserIDAndTimeRange retrieves chat messages created within a specified time range for a given user ID.
// It takes the userID as the ID of the user to query, endTime as the upper limit of the time range.
// The function returns a slice of Chat structs representing the retrieved messages, or an error if the query fails.
func (cr *ChatRepository) QueryByUserIDAndTimeRange(tx *sql.Tx, userID int, endTime time.Time) ([]Chat, error) {
	startTime := time.Now().Add(-24 * time.Hour)

	// query
	query := `
	SELECT id, chat_session_id, user_id, content, create_at, update_at, deleted 
	FROM chats 
	WHERE user_id = ?
		AND create_at >= ? 
		AND create_at <= ?`
	rows, err := DB.Query(query, userID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []Chat
	for rows.Next() {
		c := Chat{}
		err := rows.Scan(&c.ID, &c.SessionID, &c.UserID, &c.Content, &c.CreateAt, &c.UpdateAt, &c.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (cr *ChatRepository) Create(tx *sql.Tx, sessionID int, userID int, content string) error {
	// query
	query := `INSERT INTO chats (chat_session_id, user_id, content) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, sessionID, userID, content)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ChatRepository) Delete(tx *sql.Tx, sessionID int) error {
	// query
	query := `DELETE FROM chats WHERE chat_session_id = ?`
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ChatRepository) DeleteBySessionAndUser(tx *sql.Tx, sessionID int, userID int) error {
	// query
	query := `DELETE FROM chats WHERE chat_session_id = ? AND user_id = ?`
	_, err := DB.Exec(query, sessionID, userID)
	if err != nil {
		return err
	}
	return nil
}
