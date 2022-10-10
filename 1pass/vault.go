// onepass configures the vault in-memory key store and pulls the relevant secrets from a 1password connect server.
package onepass

import (
	"fmt"

	"github.com/1Password/connect-sdk-go/connect"
)

//Vault stores secrets and the connection client for getting secrets
type Vault struct {
	// ophost  string // for docker compose this is http://localhost:8080
	// optoken string // this is set by the cli or ui interface
	client connect.Client
	// Secrets contains information for connecting to discord, twitter, and the cockroachDB
	Secrets map[string]interface{}
}

func NewVault(ophost, optoken string) *Vault {
	client := connect.NewClient(ophost, optoken)
	return &Vault{client: client}
}

// GetSecrets retreive secrets for all vaults associated to a client. 1password gives us a payload of this format
// type Item struct {
// 	ID    string `json:"id"`
// 	Title string `json:"title"`

// 	URLs     []ItemURL `json:"urls,omitempty"`
// 	Favorite bool      `json:"favorite,omitempty"`
// 	Tags     []string  `json:"tags,omitempty"`
// 	Version  int       `json:"version,omitempty"`

// 	Vault    ItemVault    `json:"vault"`
// 	Category ItemCategory `json:"category,omitempty"` // TODO: switch this to `category`

// 	Sections []*ItemSection `json:"sections,omitempty"`
// 	Fields   []*ItemField   `json:"fields,omitempty"`
// 	Files    []*File        `json:"files,omitempty"`

// 	LastEditedBy string    `json:"lastEditedBy,omitempty"`
// 	CreatedAt    time.Time `json:"createdAt,omitempty"`
// 	UpdatedAt    time.Time `json:"updatedAt,omitempty"`

// 	// Deprecated: Connect does not return trashed items.
// 	Trashed bool `json:"trashed,omitempty"`
// }
func (v *Vault) GetSecrets() error {
	vaults, err := v.client.GetVaults()
	if err != nil {
		return fmt.Errorf("cannot retrieve 1password secrets. %w", err)
	}

	for _, vault := range vaults {
		items, err := v.client.GetItems(vault.ID)
		if err != nil {
			return fmt.Errorf("cannot retrieve 1password secrets. %w", err)
		}
		for _, item := range items {
			v.Secrets[item.Title] = item.GetValue("api_credential")
		}
	}
	return nil
}
