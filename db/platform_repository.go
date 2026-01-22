package database

import (
	"context"
	"fmt"
	"time"
)

// PlatformAccountRepository handles all platform account-related database operations
type PlatformAccountRepository struct {
	conn *SupabaseConnection
}

// NewPlatformAccountRepository creates a new PlatformAccountRepository
func NewPlatformAccountRepository(conn *SupabaseConnection) *PlatformAccountRepository {
	return &PlatformAccountRepository{conn: conn}
}

// GetByID retrieves a platform account by its UUID
func (r *PlatformAccountRepository) GetByID(ctx context.Context, id string) (*PlatformAccount, error) {
	var account PlatformAccount
	query := `
		SELECT id, platform, account_handle, display_name, is_active,
		       posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
		       credential_key, created_at, updated_at, metadata
		FROM platform_accounts
		WHERE id = $1`

	err := r.conn.DB.GetContext(ctx, &account, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform account by ID: %w", err)
	}
	return &account, nil
}

// GetByPlatformAndHandle retrieves a platform account by platform and handle
func (r *PlatformAccountRepository) GetByPlatformAndHandle(ctx context.Context, platform PlatformType, handle string) (*PlatformAccount, error) {
	var account PlatformAccount
	query := `
		SELECT id, platform, account_handle, display_name, is_active,
		       posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
		       credential_key, created_at, updated_at, metadata
		FROM platform_accounts
		WHERE platform = $1 AND account_handle = $2`

	err := r.conn.DB.GetContext(ctx, &account, query, platform, handle)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform account by platform and handle: %w", err)
	}
	return &account, nil
}

// ListAll retrieves all platform accounts
func (r *PlatformAccountRepository) ListAll(ctx context.Context) ([]PlatformAccount, error) {
	var accounts []PlatformAccount
	query := `
		SELECT id, platform, account_handle, display_name, is_active,
		       posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
		       credential_key, created_at, updated_at, metadata
		FROM platform_accounts
		ORDER BY platform, account_handle`

	err := r.conn.DB.SelectContext(ctx, &accounts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all platform accounts: %w", err)
	}
	return accounts, nil
}

// ListActive retrieves all active platform accounts
func (r *PlatformAccountRepository) ListActive(ctx context.Context) ([]PlatformAccount, error) {
	var accounts []PlatformAccount
	query := `
		SELECT id, platform, account_handle, display_name, is_active,
		       posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
		       credential_key, created_at, updated_at, metadata
		FROM platform_accounts
		WHERE is_active = true
		ORDER BY platform, account_handle`

	err := r.conn.DB.SelectContext(ctx, &accounts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active platform accounts: %w", err)
	}
	return accounts, nil
}

// ListByPlatform retrieves all accounts for a specific platform
func (r *PlatformAccountRepository) ListByPlatform(ctx context.Context, platform PlatformType) ([]PlatformAccount, error) {
	var accounts []PlatformAccount
	query := `
		SELECT id, platform, account_handle, display_name, is_active,
		       posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
		       credential_key, created_at, updated_at, metadata
		FROM platform_accounts
		WHERE platform = $1
		ORDER BY account_handle`

	err := r.conn.DB.SelectContext(ctx, &accounts, query, platform)
	if err != nil {
		return nil, fmt.Errorf("failed to list platform accounts by platform: %w", err)
	}
	return accounts, nil
}

// GetAccountsDueForPosting retrieves accounts that are due for posting
// An account is due if:
// 1. It's active
// 2. next_post_due_at is NULL or in the past
// OR
// 3. last_posted_at is NULL or (now - last_posted_at) >= posting_schedule_days
func (r *PlatformAccountRepository) GetAccountsDueForPosting(ctx context.Context) ([]PlatformAccount, error) {
	var accounts []PlatformAccount
	query := `
		SELECT id, platform, account_handle, display_name, is_active,
		       posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
		       credential_key, created_at, updated_at, metadata
		FROM platform_accounts
		WHERE is_active = true
		  AND (
		    next_post_due_at IS NULL
		    OR next_post_due_at <= NOW()
		    OR last_posted_at IS NULL
		    OR last_posted_at <= NOW() - (posting_schedule_days || ' days')::interval
		  )
		ORDER BY
		  COALESCE(last_posted_at, '1970-01-01'::timestamp) ASC,
		  platform, account_handle`

	err := r.conn.DB.SelectContext(ctx, &accounts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts due for posting: %w", err)
	}
	return accounts, nil
}

// GetAccountStatus retrieves the posting status for all accounts
func (r *PlatformAccountRepository) GetAccountStatus(ctx context.Context) ([]PlatformAccountStatus, error) {
	var statuses []PlatformAccountStatus
	query := `SELECT * FROM v_platform_account_status ORDER BY platform, account_handle`

	err := r.conn.DB.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get account status: %w", err)
	}
	return statuses, nil
}

// Create inserts a new platform account
func (r *PlatformAccountRepository) Create(ctx context.Context, account *PlatformAccount) error {
	query := `
		INSERT INTO platform_accounts (
			platform, account_handle, display_name, is_active,
			posting_schedule_days, cooldown_days, last_posted_at, next_post_due_at,
			credential_key, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id, created_at, updated_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		account.Platform,
		account.AccountHandle,
		account.DisplayName,
		account.IsActive,
		account.PostingScheduleDays,
		account.CooldownDays,
		account.LastPostedAt,
		account.NextPostDueAt,
		account.CredentialKey,
		account.Metadata,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create platform account: %w", err)
	}
	return nil
}

// Update updates an existing platform account
func (r *PlatformAccountRepository) Update(ctx context.Context, account *PlatformAccount) error {
	query := `
		UPDATE platform_accounts SET
			display_name = $1,
			is_active = $2,
			posting_schedule_days = $3,
			cooldown_days = $4,
			last_posted_at = $5,
			next_post_due_at = $6,
			credential_key = $7,
			metadata = $8,
			updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		account.DisplayName,
		account.IsActive,
		account.PostingScheduleDays,
		account.CooldownDays,
		account.LastPostedAt,
		account.NextPostDueAt,
		account.CredentialKey,
		account.Metadata,
		account.ID,
	).Scan(&account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update platform account: %w", err)
	}
	return nil
}

// UpdateLastPostedAt updates the last_posted_at timestamp and calculates next_post_due_at
func (r *PlatformAccountRepository) UpdateLastPostedAt(ctx context.Context, accountID string, postedAt time.Time) error {
	query := `
		UPDATE platform_accounts SET
			last_posted_at = $1,
			next_post_due_at = $1 + (posting_schedule_days || ' days')::interval,
			updated_at = NOW()
		WHERE id = $2`

	_, err := r.conn.DB.ExecContext(ctx, query, postedAt, accountID)
	if err != nil {
		return fmt.Errorf("failed to update last_posted_at: %w", err)
	}
	return nil
}

// SetActive enables or disables a platform account
func (r *PlatformAccountRepository) SetActive(ctx context.Context, accountID string, active bool) error {
	query := `
		UPDATE platform_accounts SET
			is_active = $1,
			updated_at = NOW()
		WHERE id = $2`

	_, err := r.conn.DB.ExecContext(ctx, query, active, accountID)
	if err != nil {
		return fmt.Errorf("failed to set account active status: %w", err)
	}
	return nil
}

// Count returns the total number of platform accounts
func (r *PlatformAccountRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM platform_accounts`

	err := r.conn.DB.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count platform accounts: %w", err)
	}
	return count, nil
}
