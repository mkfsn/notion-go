package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Block interface {
	isBlock()
}

func newBlock(data []byte) (Block, error) {
	var base BlockBase

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case BlockTypeParagraph:
		var block ParagraphBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeHeading1:
		var block Heading3Block

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeHeading2:
		var block Heading3Block

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeHeading3:
		var block Heading3Block

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeBulletedListItem:
		var block BulletedListItemBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeNumberedListItem:
		var block NumberedListItemBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeToDo:
		var block ToDoBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeToggle:
		var block ToggleBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeChildPage:
		var block ChildPageBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case BlockTypeUnsupported:
		var block UnsupportedBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil
	}

	return nil, ErrUnknown
}

type BlockType string

const (
	BlockTypeParagraph        BlockType = "paragraph"
	BlockTypeHeading1         BlockType = "heading_1"
	BlockTypeHeading2         BlockType = "heading_2"
	BlockTypeHeading3         BlockType = "heading_3"
	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"
	BlockTypeToDo             BlockType = "to_do"
	BlockTypeToggle           BlockType = "toggle"
	BlockTypeChildPage        BlockType = "child_page"
	BlockTypeUnsupported      BlockType = "unsupported"
)

type BlockBase struct {
	// Always "block".
	Object string `json:"object"`
	// Identifier for the block.
	ID string `json:"id"`
	// Type of block.
	Type BlockType `json:"type"`
	// Date and time when this block was created. Formatted as an ISO 8601 date time string.
	CreatedTime time.Time `json:"created_time"`
	// Date and time when this block was last updated. Formatted as an ISO 8601 date time string.
	LastEditedTime time.Time `json:"last_edited_time"`
	// Whether or not the block has children blocks nested within it.
	HasChildren bool `json:"has_children"`
}

func (b BlockBase) isBlock() {}

type ParagraphBlock struct {
	BlockBase
	Text     []RichText  `json:"text"`
	Children []BlockBase `json:"children"`
}

type Heading3Block struct {
	BlockBase
	Text []RichText `json:"text"`
}

type HeadingTwoBlock struct {
	BlockBase
	Text []RichText `json:"text"`
}

type HeadingThreeBlock struct {
	BlockBase
	Text []RichText `json:"text"`
}

type BulletedListItemBlock struct {
	BlockBase
	Text     []RichText  `json:"text"`
	Children []BlockBase `json:"children"`
}

type NumberedListItemBlock struct {
	BlockBase
	Text     []RichText  `json:"text"`
	Children []BlockBase `json:"children"`
}

type ToDoBlock struct {
	BlockBase
	Text     []RichText  `json:"text"`
	Checked  bool        `json:"checked"`
	Children []BlockBase `json:"children"`
}

type ToggleBlock struct {
	BlockBase
	Text     []RichText  `json:"text"`
	Children []BlockBase `json:"children"`
}

type ChildPageBlock struct {
	BlockBase
	Title string `json:"title"`
}

type UnsupportedBlock struct {
	BlockBase
}

type BlocksInterface interface {
	Children() BlocksChildrenInterface
}

type blocksClient struct {
	childrenClient *blocksChildrenClient
}

func newBlocksClient(client client) *blocksClient {
	return &blocksClient{
		childrenClient: newBlocksChildrenClient(client),
	}
}

func (b *blocksClient) Children() BlocksChildrenInterface {
	if b == nil {
		return nil
	}
	return b.childrenClient
}

type BlocksChildrenListParameters struct {
	PaginationParameters

	// Identifier for a block
	BlockID string `json:"-"`
}

type BlocksChildrenListResponse struct {
	PaginatedList
	Results []Block `json:"results"`
}

func (b *BlocksChildrenListResponse) UnmarshalJSON(data []byte) error {
	type Alias BlocksChildrenListResponse

	alias := struct {
		*Alias
		Results []json.RawMessage `json:"results"`
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	b.Results = make([]Block, 0, len(alias.Results))

	for _, value := range alias.Results {
		block, err := newBlock(value)
		if err != nil {
			return err
		}
		b.Results = append(b.Results, block)
	}

	return nil
}

type BlocksChildrenAppendParameters struct {
	// Identifier for a block
	BlockID string
	// Child content to append to a container block as an array of block objects
	Children []Block `json:"children"`
}

type BlocksChildrenAppendResponse struct {
	// TODO: check if this is correct
	BlockBase
}

type BlocksChildrenInterface interface {
	List(ctx context.Context, params BlocksChildrenListParameters) (*BlocksChildrenListResponse, error)
	Append(ctx context.Context, params BlocksChildrenAppendParameters) (*BlocksChildrenAppendResponse, error)
}

type blocksChildrenClient struct {
	client client
}

func newBlocksChildrenClient(client client) *blocksChildrenClient {
	return &blocksChildrenClient{
		client: client,
	}
}

func (b *blocksChildrenClient) List(ctx context.Context, params BlocksChildrenListParameters) (*BlocksChildrenListResponse, error) {
	endpoint := strings.Replace(APIBlocksListChildrenEndpoint, "{block_id}", params.BlockID, 1)

	fmt.Printf("endpoint: %s\n", endpoint)

	data, err := b.client.Request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response BlocksChildrenListResponse

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (b *blocksChildrenClient) Append(ctx context.Context, params BlocksChildrenAppendParameters) (*BlocksChildrenAppendResponse, error) {
	return nil, ErrUnimplemented
}
