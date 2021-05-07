# Goblinator

## Description
Goblinator is a tool for World of Warcraft auction data crawling.

## Plans
1. Prepare Go package to crawl data
2. Use package in a crawler worker
    - Configure delay between crawls (with validation from Blizzard 
      restrictions)
    - Write logs
3. Allow user to set his own storage for data
    - Google Spreadsheet?
    - PostgreSQL DB
    - Excel File

## Instructions
(Are going to be filled while development process GOes on)

1. Create .env file in root
2. Add BLIZZARD_CLIENT_ID
3. Add BLIZZARD_CLIENT_SECRET
4. Run within root directory: `go run ./cmd/crawler`