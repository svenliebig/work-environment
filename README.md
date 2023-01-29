# work environment

work environment is a CLI tool to help you to simplify your daily workflow by:

- removing unnecessary clutter actions that you perform all day
- connecting the things that belong together in one place
- creating an extendable environment of tools

## Installing

The current install is only available by build this project yourself.

Requirements:

- go >= 1.18

### Unix Systems

```sh
git clone git@github.com:svenliebig/work-environment.git
cd work-environment
./build.sh
echo 'export PATH="$PATH:$(pwd)/bin"' >> ~/.bash_profile
```

Thoughts, there would be the possibility to have multiple work-environment. Does this make sense? Or do we want a global bash variable that defines that localtion of that work environment.

## Set up your environment

You can easily initialize your instance by navigating to your default workspace (the root of the place where all your repositories and projects are) and executing this command:

```sh
we init
```

If there is already a `.work-environment.json`, the command will ask you if you want to override the existing configuration or not.

This will search for project's in the directories and create the foundation for your environment.

## Update your environment

If you already have a `.work-environment.json`, but you want to update the configuration with new projects then you can use:

```sh
we update
```

## The Interface

### CI (`we ci`)

Commands that will help you to manage the continuous integration.

Currently supported CI's:

- bamboo

#### `we ci create`

This command will `create` a global CI to your work environment.

Parameters:

- `name` - the identifier for the CI, this must be unique over the whole work-environment
- `type` - the CI type, can be `bamboo`
- `auth` - the basic auth token (base64 encoded)
- `url` - the root url of your CI

Example:

```sh
we ci create --name 'company-ci' --type 'bamboo' --auth 'jfgasdijaskosdf*13asdka)1231' --url 'https://bamboo.mycompany.com'
```

#### `we ci add`

This command will `add` an CI environment to your project. It need to be executed in a directory of a project or you need to provide the project identifier.

(!) Currently there is only `we ci add --suggest` supported, the parameters are not fully implemented.

Parameters:

- `project` - the project identifier in your work-environment. **optional** if you are in the directory of the project.
- `name` - the identifier for the CI. **optional** when there is only one CI setup in your project.
- `projectKey` - the project key on you CI, **optional** if suggest is used
- `suggest (optional: boolean)` - try to use the CI api to find you project build

Example:

```sh
we ci add --name 'company-ci' --projectKey 'PRO-JECT' --project "my-project"
```

#### `we ci open`

Open's the related CI of the project in your browser. If branch builds are supported by that CI, it will be opened on the latest or running branch build.

```sh
we ci open
```

#### `we ci remove`

Removes the CI configuration from a project. It needs to be executed in a project directory.

```sh
we ci remove
```

#### `we ci info`

Prints out the CI informations for the current project, will also print the last build state and a tail of the build log if the last build failed. It needs to be executed in a project directory. If branch builds are supported by that CI, the displayed build will be the branch build.

```sh
we ci info
```

### CD `we cd`

Commands that will help you to manage your continouos delivery.

Currently supported CD's are:

- bamboo

#### `we cd add`

This command will `add` an CD environment to your project. It needs to be executed in a project directory. Currently there is only bamboo supported, so the command will use the bamboo API to get the deployment project by the CI project key. You have to have a CI key configured.

#### `we cd open`

Open's the related CD of the project in your browser. Needs to be executed in a project context.

```sh
we cd open
```

#### `we cd info`

Prints out the CD informations for the current project. It needs to be executed in a project directory.

```sh
we cd info
```

#### `we cd remove`

Removes the CD configuration from a project. It needs to be executed in a project directory.

```sh
we cd remove
```
