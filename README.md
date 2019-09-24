# GitTeams

[![Build status](https://img.shields.io/travis/NoUseFreak/gitteams/master?style=flat-square)](https://travis-ci.org/NoUseFreak/gitteams)
[![Release](https://img.shields.io/github/v/release/NoUseFreak/gitteams?style=flat-square)](https://github.com/NoUseFreak/gitteams/releases)
[![Maintained](https://img.shields.io/maintenance/yes/2019?style=flat-square)](https://github.com/NoUseFreak/gitteams)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/nousefreak/gitteams?style=flat-square)](https://hub.docker.com/r/nousefreak/gitteams)
[![License](https://img.shields.io/github/license/NoUseFreak/gitteams?style=flat-square)](https://github.com/NoUseFreak/gitteams/blob/master/LICENSE)
[![Coffee](https://img.shields.io/badge/☕️-Buy%20me%20a%20coffee-blue?style=flat-square&color=blueviolet)](https://www.buymeacoffee.com/driesdepeuter)

GitTeams gives you insight into multiple repositories at once.

## Example

```
$ gitteams --github-token=<token> --github-username=<username>  stats --sort=branches 
INFO[0000] Collecting repos                             
INFO[0000] Processing                                   
INFO[0001] Report                                       
┌─────────────────────────────────────────────────────────────────────────────────────────────────┐
│ REPO                                     BRANCH COUNT  LANGUAGE           LINES OF CODE  MERGED │
├─────────────────────────────────────────────────────────────────────────────────────────────────┤
│ gh:NoUseFreak/Cron                                  3  PHP (76%)                     33       0 │
│ gh:NoUseFreak/docker-multi-cache                    1  Go (22%)                       6       0 │
│ gh:NoUseFreak/cicd                                  1  Go (60%)                      33       0 │
...
└─────────────────────────────────────────────────────────────────────────────────────────────────┘
```

## More info

To see all options available, check the help function.

```
$ gitteams help
Git Teams helps you manage all project at once.

Usage:
  gitteams [command]

Available Commands:
  author      Count number of authors
  branch      Count number of branches
  commits     Count commits in repository
  help        Help about any command
  language    Show main language in repository
  loc         Get LOC count in repositories
  merged      Count merged branches
...
```

## Platforms

 - Github
 - Gitlab
 - Bitbucket

## Output
|                 | Description                                           |
| --------------- | ----------------------------------------------------- |
| Author count    | Count authors in each repository.                     |
| Branch count    | Count branches in each repository.                    |
| Commit count    | Count commits in each repository.                     |
| Lines of Code   | Count the lines of code in each repository.           |
| Language        | Find the most common language in each repository.     |
| Merged branches | Count all branched fully merged into the main branch. |
| Size            | Calculate the size of the repository in kb.           |
| Tag count       | Count all tags into each repository.                  |

## Install

__Binary__

```
curl -sL http://bit.ly/gh-get | PROJECT=NoUseFreak/gitteams bash
```

__Docker__

```
docker pull nousefreak/gitteams
```
