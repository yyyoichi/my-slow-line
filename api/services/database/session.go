package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type TSessionStatus string

const (
	TActiveSession   TSessionStatus = "active"
	TArchivedSession TSessionStatus = "archived"
	TBreakupSession  TSessionStatus = "breakup"
)

type TParticipantStatus string

const (
	TInvitedParty  TParticipantStatus = "invited"
	TJoinedParty   TParticipantStatus = "joined"
	TRejectedParty TParticipantStatus = "rejected"
)

type SessionRepositories struct {
	SessionRepository            SessionRepositoryInterface
	SessionParticipantRepository SessionParticipantRepositoryInterface
	SessionChatRepository        SessionChatRepositoryInterface
}

func NewSessionRepositories() *SessionRepositories {
	return &SessionRepositories{
		&SessionRepository{},
		&SessionParticipantRepository{},
		&SessionChatRepository{},
	}
}

type SessionRepository struct{}
type SessionRepositoryInterface interface {
	QueryByUserID(tx *sql.Tx, userID int, options TQuerySessionsOptions) ([]*TQuerySessions, error)
	QueryBySessionUserID(tx *sql.Tx, sessionID, userID int) (*TQuerySessions, error)
	HasStatusAt(tx *sql.Tx, sessionID, userID int, inStatus []TParticipantStatus) (bool, error)
	Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error)
	UpdateName(tx *sql.Tx, id int, name string) error
	UpdateStatus(tx *sql.Tx, id int, status TSessionStatus) error
	HardDelete(tx *sql.Tx, sessionID int) error
	SoftDelete(tx *sql.Tx, sessionID int) error
	HardDeleteAll(tx *sql.Tx, sessionID int) error
}

// query sessions by userID

type TQuerySessions struct {
	ID            int
	Name          string
	PublicKey     string
	SessionStatus TSessionStatus
	Status        TParticipantStatus
	CreateAt      time.Time
	UpdateAt      time.Time
	Deleted       bool
}
type TQuerySessionsOptions struct {
	InPartyStatus   []TParticipantStatus
	InSessionStatus []TSessionStatus
}

// query sessions
func (sr *SessionRepository) QueryByUserID(tx *sql.Tx, userID int, options TQuerySessionsOptions) ([]*TQuerySessions, error) {
	// query
	query := `
	SELECT s.id, s.public_key, s.name, s.status, s.create_at, s.update_at, s.deleted, p.status 
	FROM chat_sessions AS s 
		RIGHT JOIN chat_session_participants AS p ON s.id = p.chat_session_id
	WHERE p.user_id = ? 
		AND p.status IN (?)
		AND s.status IN (?)`
	query, params, err := sqlx.In(query, userID, options.InPartyStatus, options.InSessionStatus)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*TQuerySessions
	for rows.Next() {
		s := &TQuerySessions{}
		err := rows.Scan(&s.ID, &s.PublicKey, &s.Name, &s.SessionStatus, &s.CreateAt, &s.UpdateAt, &s.Deleted, &s.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// query a session
func (sr *SessionRepository) QueryBySessionUserID(tx *sql.Tx, sessionID, userID int) (*TQuerySessions, error) {
	// query
	query := `
	SELECT s.id, s.public_key, s.name, s.status, s.create_at, s.update_at, s.deleted, p.status 
	FROM chat_sessions AS s 
		RIGHT JOIN chat_session_participants AS p ON s.id = p.chat_session_id
	WHERE s.id = ? AND p.user_id = ?`
	row := tx.QueryRow(query, sessionID, userID)

	// result
	s := &TQuerySessions{}
	err := row.Scan(&s.ID, &s.PublicKey, &s.Name, &s.SessionStatus, &s.CreateAt, &s.UpdateAt, &s.Deleted, &s.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return nil without an error.
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

// has status at session for user
func (sr *SessionRepository) HasStatusAt(tx *sql.Tx, sessionID, userID int, inStatus []TParticipantStatus) (bool, error) {
	// query
	query := `
	SELECT s.id
	FROM chat_sessions AS s 
		RIGHT JOIN chat_session_participants AS p ON s.id = p.chat_session_id
	WHERE s.id = ? AND p.user_id = ? AND p.status IN (?)`
	query, params, err := sqlx.In(query, sessionID, userID, inStatus)
	if err != nil {
		return false, err
	}
	row := tx.QueryRow(query, params...)

	// result
	id := 0
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return nil without an error.
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// create
func (sr *SessionRepository) Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error) {
	// query
	query := `INSERT INTO chat_sessions (user_id, public_key, name) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, userID, publicKey, name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// update
func (sr *SessionRepository) UpdateName(tx *sql.Tx, id int, name string) error {
	// query
	query := `UPDATE chat_sessions SET name = ? WHERE id = ?`
	_, err := tx.Exec(query, name, id)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionRepository) UpdateStatus(tx *sql.Tx, id int, status TSessionStatus) error {
	// query
	query := `UPDATE chat_sessions SET status = ? WHERE id = ?`
	_, err := tx.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}

// delete row
func (sr *SessionRepository) HardDelete(tx *sql.Tx, sessionID int) error {
	// query
	query := `DELETE FROM chat_sessions WHERE id = ?`
	_, err := tx.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

// up deleted flag
func (sr *SessionRepository) SoftDelete(tx *sql.Tx, sessionID int) error {
	// query
	query := `UPDATE chat_sessions SET deleted = 1 WHERE id = ?`
	_, err := tx.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

// delete hard all session contents
func (sr *SessionRepository) HardDeleteAll(tx *sql.Tx, sessionID int) error {
	// query
	querys := []string{
		`DELETE FROM chats WHERE chat_session_id = ?`,
		`DELETE FROM chat_session_participants WHERE chat_session_id = ?`,
		`DELETE FROM chat_sessions WHERE id = ?`,
	}
	for _, query := range querys {
		_, err := tx.Exec(query, sessionID)
		if err != nil {
			return err
		}
	}
	return nil
}

// query session participants

type SessionParticipantRepository struct{}
type SessionParticipantRepositoryInterface interface {
	QueryBySessionID(tx *sql.Tx, sessionID int) ([]*TSessionParticipant, error)
	Create(tx *sql.Tx, sessionID, userID, inviteUserID int, status TParticipantStatus) (int, error)
	UpdateStatus(tx *sql.Tx, id int, status TParticipantStatus) error
	HardDelete(tx *sql.Tx, sessionID int) error
}

type TSessionParticipant struct {
	ID        int
	SessionID int
	UserID    int
	Status    TParticipantStatus
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

func (spr *SessionParticipantRepository) QueryBySessionID(tx *sql.Tx, sessionID int) ([]*TSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted 
	FROM chat_session_participants 
	WHERE chat_session_id = ?`
	rows, err := tx.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []*TSessionParticipant
	for rows.Next() {
		sp := &TSessionParticipant{}
		err := rows.Scan(&sp.ID, &sp.SessionID, &sp.UserID, &sp.Status, &sp.CreateAt, &sp.UpdateAt, &sp.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, sp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// [invite_user_id] is ID of the user exec to the session
func (spr *SessionParticipantRepository) Create(tx *sql.Tx, sessionID, userID, inviteUserID int, status TParticipantStatus) (int, error) {
	// query
	query := `INSERT INTO chat_session_participants (chat_session_id, user_id, invite_user_id, status) VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(query, sessionID, userID, inviteUserID, status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// update
func (spr *SessionParticipantRepository) UpdateStatus(tx *sql.Tx, id int, status TParticipantStatus) error {
	// query
	query := `UPDATE chat_session_participants SET status = ? WHERE id = ?`
	_, err := tx.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}

// delete
func (spr *SessionParticipantRepository) HardDelete(tx *sql.Tx, participantID int) error {
	// query
	query := `DELETE FROM chat_session_participants WHERE id = ?`
	_, err := tx.Exec(query, participantID)
	if err != nil {
		return err
	}
	return nil
}

type SessionChatRepository struct{}
type SessionChatRepositoryInterface interface {
	Create(tx *sql.Tx, sessionID, userID int, content string) (int, error)
	QueryByUserIDInRange(tx *sql.Tx, userID int, inRange TQuerySessionChatInRange) ([]*TQuerySessionChat, error)
	QueryBySessionIDInRange(tx *sql.Tx, sessionID int, inRange TQuerySessionChatInRange) ([]*TQuerySessionChat, error)
	QueryLastChatInActiveSessions(tx *sql.Tx, userID int) ([]*TQueryLastChat, error)
	HardDelete(tx *sql.Tx, chatID int) error
}

type TQuerySessionChat struct {
	ID        int
	SessionID int
	UserID    int
	Content   string
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type TQuerySessionChatInRange struct {
	StartDate, EndDate time.Time
}

func (cr *SessionChatRepository) QueryByUserIDInRange(tx *sql.Tx, userID int, inRange TQuerySessionChatInRange) ([]*TQuerySessionChat, error) {
	// query
	query := `
	SELECT c.id, c.chat_session_id, c.user_id, c.content, c.create_at, c.update_at, c.deleted 
	FROM chats AS c
	LEFT JOIN chat_session_participants AS p 
		ON c.chat_session_id = p.chat_session_id AND c.user_id = p.user_id
	WHERE c.chat_session_id IN (
		SELECT chat_session_id FROM chat_session_participants AS p 
		WHERE p.user_id = ? AND p.status = 'joined' 
	)
		AND c.create_at BETWEEN ? AND ?`
	rows, err := tx.Query(query, userID, inRange.StartDate, inRange.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []*TQuerySessionChat
	for rows.Next() {
		c := &TQuerySessionChat{}
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

func (cr *SessionChatRepository) QueryBySessionIDInRange(tx *sql.Tx, sessionID int, inRange TQuerySessionChatInRange) ([]*TQuerySessionChat, error) {
	// query
	query := `
	SELECT id, chat_session_id, user_id, content, create_at, update_at, deleted 
	FROM chats 
	WHERE chat_session_id = ?
		AND create_at BETWEEN ? AND ?`
	rows, err := tx.Query(query, sessionID, inRange.StartDate, inRange.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []*TQuerySessionChat
	for rows.Next() {
		c := &TQuerySessionChat{}
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

type TQueryLastChat struct {
	SessionName string
	SessionID   int
	UserID      int
	ID          int
	Content     string
	CreateAt    time.Time
	UpdateAt    time.Time
	Deleted     bool
}

func (cr *SessionChatRepository) QueryLastChatInActiveSessions(tx *sql.Tx, userID int) ([]*TQueryLastChat, error) {
	query := `
	SELECT name, id, user_id, chat_id, content, create_at, update_at, deleted
	FROM (
		SELECT s.name, s.id, c.user_id, c.id AS chat_id, c.content, c.create_at, c.update_at, c.deleted,
					ROW_NUMBER() OVER (PARTITION BY s.id ORDER BY c.id DESC) AS rn
		FROM chats AS c
		LEFT JOIN chat_sessions AS s ON c.chat_session_id = s.id
		LEFT JOIN chat_session_participants AS p ON p.chat_session_id = s.id
		WHERE p.user_id = ?
			AND p.status = 'joined'
			AND s.status = 'active'
	) AS temp
	WHERE rn = 1
	ORDER BY create_at DESC
	`
	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []*TQueryLastChat
	for rows.Next() {
		c := &TQueryLastChat{}
		err := rows.Scan(&c.SessionName, &c.SessionID, &c.UserID, &c.ID, &c.Content, &c.CreateAt, &c.UpdateAt, &c.Deleted)
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

// create
func (scr *SessionChatRepository) Create(tx *sql.Tx, sessionID, userID int, content string) (int, error) {
	// query
	query := `INSERT INTO chats (chat_session_id, user_id, content) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, sessionID, userID, content)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// delete
func (cr *SessionChatRepository) HardDelete(tx *sql.Tx, chatID int) error {
	//query
	query := `DELETE FROM chats WHERE id = ?`
	_, err := tx.Exec(query, chatID)
	if err != nil {
		return err
	}
	return nil
}
