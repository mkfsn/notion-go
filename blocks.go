package notion

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/mkfsn/notion-go/rest"
)

type Block interface {
	isBlock()
}

type BlockBase struct {
	// Always "block".
	Object ObjectType `json:"object"`
	// Identifier for the block.
	ID string `json:"id,omitempty"`
	// Type of block.
	Type BlockType `json:"type"`
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

func (h *HeadingBlock) UnmarshalJSON(data []byte) error {
	var text []richTextDecoder

	if err := json.Unmarshal(data, &text); err != nil {
		return nil
	}

	h.Text = make([]RichText, 0, len(text))

	for _, decoder := range text {
		h.Text = append(h.Text, decoder.RichText)
	}

	return nil
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
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

func (r *RichTextBlock) UnmarshalJSON(data []byte) error {
	var alias struct {
		Text     []richTextDecoder `json:"text"`
		Children []blockDecoder    `json:"children"`
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return nil
	}

	r.Text = make([]RichText, 0, len(r.Text))

	for _, decoder := range alias.Text {
		r.Text = append(r.Text, decoder.RichText)
	}

	r.Children = make([]Block, 0, len(r.Children))

	for _, decoder := range alias.Children {
		r.Children = append(r.Children, decoder.Block)
	}

	return nil
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

func newBlocksClient(restClient rest.Interface) *blocksClient {
	return &blocksClient{
		childrenClient: newBlocksChildrenClient(restClient),
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
		Results []blockDecoder `json:"results"`
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	b.Results = make([]Block, 0, len(alias.Results))

	for _, decoder := range alias.Results {
		b.Results = append(b.Results, decoder.Block)
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
	var decoder blockDecoder

	if err := json.Unmarshal(data, &decoder); err != nil {
		return err
	}

	b.Block = decoder.Block

	return nil
}

type BlocksChildrenInterface interface {
	List(ctx context.Context, params BlocksChildrenListParameters) (*BlocksChildrenListResponse, error)
	Append(ctx context.Context, params BlocksChildrenAppendParameters) (*BlocksChildrenAppendResponse, error)
}

type blocksChildrenClient struct {
	restClient rest.Interface
}

func newBlocksChildrenClient(restClient rest.Interface) *blocksChildrenClient {
	return &blocksChildrenClient{
		restClient: restClient,
	}
}

func (b *blocksChildrenClient) List(ctx context.Context, params BlocksChildrenListParameters) (*BlocksChildrenListResponse, error) {
	var result BlocksChildrenListResponse
	var failure HTTPError

	err := b.restClient.New().Get().
		Endpoint(strings.Replace(APIBlocksListChildrenEndpoint, "{block_id}", params.BlockID, 1)).
		QueryStruct(params).
		BodyJSON(nil).
		Receive(ctx, &result, &failure)

	return &result, err
}

func (b *blocksChildrenClient) Append(ctx context.Context, params BlocksChildrenAppendParameters) (*BlocksChildrenAppendResponse, error) {
	var result BlocksChildrenAppendResponse
	var failure HTTPError

	err := b.restClient.New().Patch().
		Endpoint(strings.Replace(APIBlocksAppendChildrenEndpoint, "{block_id}", params.BlockID, 1)).
		QueryStruct(params).
		BodyJSON(params).
		Receive(ctx, &result, &failure)

	return &result, err
}

type blockDecoder struct {
	Block
}

func (b *blockDecoder) UnmarshalJSON(data []byte) error {
	var decoder struct {
		Type BlockType `json:"type"`
	}

	if err := json.Unmarshal(data, &decoder); err != nil {
		return err
	}

	switch decoder.Type {
	case BlockTypeParagraph:
		b.Block = &ParagraphBlock{}

	case BlockTypeHeading1:
		b.Block = &Heading1Block{}

	case BlockTypeHeading2:
		b.Block = &Heading2Block{}

	case BlockTypeHeading3:
		b.Block = &Heading3Block{}

	case BlockTypeBulletedListItem:
		b.Block = &BulletedListItemBlock{}

	case BlockTypeNumberedListItem:
		b.Block = &NumberedListItemBlock{}

	case BlockTypeToDo:
		b.Block = &ToDoBlock{}

	case BlockTypeToggle:
		b.Block = &ToggleBlock{}

	case BlockTypeChildPage:
		b.Block = &ChildPageBlock{}

	case BlockTypeUnsupported:
		b.Block = &UnsupportedBlock{}
	}

	return json.Unmarshal(data, &b.Block)
}
