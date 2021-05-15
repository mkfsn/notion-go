package notion

import (
	"context"
	"time"
)

type Block interface {
	isBlock()
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

type HeadingOneBlock struct {
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
	BlockID string
}

type BlocksChildrenListResponse struct {
	PaginatedList
	Results []BlockBase `json:"results"`
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
	return nil, ErrUnimplemented
}

func (b *blocksChildrenClient) Append(ctx context.Context, params BlocksChildrenAppendParameters) (*BlocksChildrenAppendResponse, error) {
	return nil, ErrUnimplemented
}
