# work environment

work environment is a CLI tool to help you to simplify your daily workflow by:

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
we init
```

Errors:

* there is already a work environment file

Thoughts: `do` needs to be renamed, it's reserved basically everywhere. Also important is the position of the keys, you 
will not have fun pressing `qp` all the time. `cd` is a good example, you can enter it by just rolling with two fingers
over your keyboard. `df`

This will search for project's in the directories and create the foundation for your environment.

## The Interface

Whereever you are in your work environment structure, you can always use all commands of your work environment. When you are
located in a specific project, the command will be execute on that project. If you are in an directory that contains multiple
projects, the CLI will provide either ask you to select one of the projects, or execute the command on all project when you
provided the `--all` flag.

### CI (`we ci`)

Commands that will help you to manage the continuous integration.

#### `we ci create`

This command will `create` a global CI to your work environment.

Parameters:

* `name` - the identifier for the CI, this must be unique over the whole work-environment
* `type` - the CI type, can be `bamboo`
* `auth` - the basic auth token (base64 encoded)
* `url`  - the root url of your CI

Example:

```sh
we ci create --name 'company-ci' --type 'bamboo' --auth 'jfgasdijaskosdf*13asdka)1231' --url 'https://bamboo.mycompany.com'
```

#### `we ci add`

This command will `add` an CI environment to your project. It need to be executed in a directory of a project or you need to provide the project identifier.

Parameters:

* `project` - the project identifier in your work-environment. __optional__ if you are in the directory of the project.
* `name` - the identifier for the CI. __optional__ when there is only one CI setup in your project.
* `projectKey` - the project key on you CI

Example:

```sh
we ci add --name 'company-ci' --projectKey 'PRO-JECT' --project "my-project"
```

### `open`

Open's the related CI of the project in your browser. If branch builds are supported by that CI, it will be opened on
the latest or running branch build.

```sh
we open ci
```

or

```sh
we ci open
```

Thoughts: `we open ci` is more naturally, but `we ci open`, is more friendly to be extended. Maybe it's just a rephrasing
in the CLI, but in the code you could have an interface that implements the `browserOpen` method, which can be inherited by
`ticket` and by `ci`, but a CI can do more than a ticket. Seems to be a better approach than writing an open interface which
has all the different entities it can use. Anyways, it hasn't much todo with the CLI interface that is implemented.

Thoughts 2: Why not both? Support both, dev's should not have the problem of figuring out "how do I do things", it should
come naturally to them and there should be the possibility to give both wordings.

### Open Ticket

```sh
we open ticket
```