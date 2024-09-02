# simple-todo

Simple TODO is a cli is test CLI application that is built to practice GoLang skills
In this branch, we have replaced SQLite DB with simple CSV file as SQLite on windows needs additional configurations.

### Tasks:

1. [x] ~~Connect with SQLite DB for saving TODOs~~
1. [x] Use CSV file as DB (Windows need some additional configuration to work with SQLite)
1. [x] Use Cobra to build TODO application
       a. [x] Create TODO
       b. [x] List all TODOs
       b. [x] Update TODO
       c. [x] Set status of TODO
       d. [x] Delete TODO
1. [x] Update the UX by using lipgloss by charm
       a. [x] Add table view
       b. [x] Add Form for creating TODO
       c. [x] Add form for updating TODO
       d. [x] Use List for selecting TODO for update and status-update
       e. [x] Give option to set-status as closed true or false if no flags are passed
