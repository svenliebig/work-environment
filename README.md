# dev organizer

dev organizer is a CLI tool to help you to simplify your daily workflow by:

* removing unnecessary clutter actions that you perform all day
* connecting the things that belong together in one place
* creating an extendable environment of tools

## Installing

```sh
```

## Set up your environment

You can easily initialize your instance by navigating to your default workspace (the root of the place where all your
repositories and projects are) and executing this command:

```sh
do init
```

Thoughts: `do` needs to be renamed, it's reserved basically everywhere.

This will search for project's in the directories and create the foundation for your environment.

## The Interface

### Open CI

```sh
do open ci
```

or

```sh
do ci open
```

Thoughts: `do open ci` is more naturally, but `do ci open`, is more friendly to be extended. Maybe it's just a rephrasing
in the CLI, but in the code you could have an interface that implements the `browserOpen` method, which can be inherited by
`ticket` and by `ci`, but a CI can do more than a ticket. Seems to be a better approach than writing an open interface which
has all the different entities it can use. Anyways, it hasn't much todo with the CLI interface that is implemented.

### Open Ticket

```sh
do open ticket
```