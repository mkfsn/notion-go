# notion-go

[![Go Report Card](https://goreportcard.com/badge/github.com/mkfsn/notion-go)](https://goreportcard.com/report/github.com/mkfsn/notion-go)
[![Actions Status](https://github.com/mkfsn/notion-go/actions/workflows/develop.yaml/badge.svg)](https://github.com/mkfsn/notion-go/actions)
[![codecov](https://codecov.io/gh/mkfsn/notion-go/branch/develop/graph/badge.svg?token=NA64P6EPQ0)](https://codecov.io/gh/mkfsn/notion-go)


A go client for the [Notion API](https://developers.notion.com/)

## Description

This aims to be an unofficial Go version of [the official SDK](https://github.com/makenotion/notion-sdk-js)
which is written in JavaScript.

## Installation

```
go get -u github.com/mkfsn/notion-go
```

## Usage

```go
c := notion.New("<NOTION_AUTH_TOKEN>")

// Retrieve block children
c.Blocks().Children().List(context.Background(), notion.BlocksChildrenListParameters{...})

// Append block children
c.Blocks().Children().Append(context.Background(), notion.BlocksChildrenAppendParameters{...})

// List databases
c.Databases().List(context.Background(), notion.DatabasesListParameters{...})

// Query a database
c.Databases().Query(context.Background(), notion.DatabasesQueryParameters{...})

// Retrieve a database
c.Databases().Retrieve(context.Background(), notion.DatabasesRetrieveParameters{...})

// Create a page
c.Pages().Create(context.Background(), notion.PagesCreateParameters{...})

// Retrieve a page
c.Pages().Retreive(context.Background(), notion.PagesRetrieveParameters{...})

// Update page properties
c.Pages().Update(context.Background(), notion.PagesUpdateParameters{...})

// List all users
c.Users().List(context.Background(), notion.UsersListParameters{...})

// Retrieve a user
c.Users().Retrieve(context.Background(), notion.UsersRetrieveParameters{...})

// Search
c.Search(context.Background(), notion.SearchParameters{...})
```

For more information, please see [examples](./examples).

## Supported Features

This client supports all endpoints in the [Notion API](https://developers.notion.com/reference/intro).

- [x] Users ✅
   * [x] [Retrieve](https://developers.notion.com/reference/get-user) ✅
   * [x] [List](https://developers.notion.com/reference/get-users) ✅
- [x] Databases ✅
  * [x] [Retrieve](https://developers.notion.com/reference/get-database) ✅
  * [x] [List](https://developers.notion.com/reference/get-databases) ✅
  * [x] [Query](https://developers.notion.com/reference/post-database-query) ✅
- [x] Pages ✅
  * [x] [Retrieve](https://developers.notion.com/reference/get-page) ✅
  * [x] [Create](https://developers.notion.com/reference/post-page) ✅️
  * [x] [Update](https://developers.notion.com/reference/patch-page) ✅️
- [x] Blocks ✅️
  * [x] Children ✅
    - [x] [Retrieve](https://developers.notion.com/reference/get-block-children) ✅
    - [x] [Append](https://developers.notion.com/reference/patch-block-children) ✅
- [x] [Search](https://developers.notion.com/reference/post-search) ✅
