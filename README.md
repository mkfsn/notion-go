# notion-go

A go client for the [Notion API](https://developers.notion.com/)

## Description

This aims to be a Go version of the [official SDK](https://github.com/makenotion/notion-sdk-js)
which is written in JavaScript.

## Installation

```
go get -u github.com/mkfsn/notion-go
```

## TODOs

- [x] Users ✅
   * [x] Retrieve ✅
   * [x] List ✅
- [x] Databases ✅
  * [x] Retrieve ✅
  * [x] List ✅
  * [x] Query ✅
- [ ] Pages ⚠️
  * [x] Retrieve ✅
  * [ ] Create ❌
  * [ ] Update ❌
- [ ] Blocks ❌
  * [ ] Children ❌
    - [ ] Retrieve ❌
    - [ ] Append ❌
- [x] Search ✅
