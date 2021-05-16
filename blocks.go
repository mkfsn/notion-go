package notion

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mkfsn/notion-go/typed"
)

type Block interface {
	isBlock()
}

// FIXME: reduce function length
// nolint:funlen
func newBlock(data []byte) (Block, error) {
	var base BlockBase

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case typed.BlockTypeParagraph:
		var block ParagraphBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeHeading1:
		var block Heading1Block

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeHeading2:
		var block Heading2Block

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeHeading3:
		var block Heading3Block

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeBulletedListItem:
		var block BulletedListItemBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeNumberedListItem:
		var block NumberedListItemBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeToDo:
		var block ToDoBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeToggle:
		var block ToggleBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeChildPage:
		var block ChildPageBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil

	case typed.BlockTypeUnsupported:
		var block UnsupportedBlock

		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}

		return block, nil
	}

	return nil, ErrUnknown
}

type BlockBase struct {
	// Always "block".
	Object typed.ObjectType `json:"object"`
	// Identifier for the block.
	ID string `json:"id,omitempty"`
	// Type of block.
	Type typed.BlockType `json:"type"`
	// Date and time when this block was created. Formatted as an ISO 8601 date time string.
	CreatedTime *time.Time `json:"created_time,omitempty"`
	// Date and time when this block was last updated. Formatted as an ISO 8601 date time string.
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	// Whether or not the block has children blocks nested within it.
	HasChildren *bool `json:"has_children,omitempty"`
}

func (b BlockBase) isBlock() {}

type ParagraphBlock struct {
	BlockBase
	Paragraph RichTextBlock `json:"paragraph"`
}

type HeadingBlock struct {
	Text []RichText `json:"text"`
}

type Heading1Block struct {
	BlockBase
	Heading1 HeadingBlock `json:"heading_1"`
}

type Heading2Block struct {
	BlockBase
	Heading2 HeadingBlock `json:"heading_2"`
}

type Heading3Block struct {
	BlockBase
	Heading3 HeadingBlock `json:"heading_3"`
}

type RichTextBlock struct {
	Text     []RichText  `json:"text"`
	Children []BlockBase `json:"children,omitempty"`
}

type BulletedListItemBlock struct {
	BlockBase
	BulletedListItem RichTextBlock `json:"bulleted_list_item"`
}

type NumberedListItemBlock struct {
	BlockBase
	NumberedListItem RichTextBlock `json:"numbered_list_item"`
}

type RichTextWithCheckBlock struct {
	Text     []RichText  `json:"text"`
	Checked  bool        `json:"checked"`
	Children []BlockBase `json:"children"`
}

type ToDoBlock struct {
	BlockBase
	ToDo RichTextWithCheckBlock `json:"todo"`
}

type ToggleBlock struct {
	BlockBase
	Toggle RichTextBlock `json:"toggle"`
}

type TitleBlock struct {
	Title string `json:"title"`
}

type ChildPageBlock struct {
	BlockBase
	ChildPage TitleBlock `json:"child_page"`
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
	BlockID string `json:"-" url:"-"`
	// Child content to append to a container block as an array of block objects
	Children []Block `json:"children"  url:"-"`
}

type BlocksChildrenAppendResponse struct {
	Block
}

func (b *BlocksChildrenAppendResponse) UnmarshalJSON(data []byte) error {
	block, err := newBlock(data)
	if err != nil {
		return err
	}

	b.Block = block

	return nil
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
	endpoint := strings.Replace(APIBlocksAppendChildrenEndpoint, "{block_id}", params.BlockID, 1)

	data, err := b.client.Request(ctx, http.MethodPatch, endpoint, params)
	if err != nil {
		return nil, err
	}

	var response BlocksChildrenAppendResponse

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
