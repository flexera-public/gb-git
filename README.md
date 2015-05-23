gb-git - Git plugin for gb
==========================

This is a plugin to manage vendored dependencies in [gb](http://getgb.io/) that have been fetched via git.

Installation
------------
```
go get -u github.com/rightscale/gb-git
```

Usage
-----
```
gb git [options...]
```

gb-git status
-------------
`gb-git status` finds all git repository workspaces in the gb project's vendor/src subdirectory and lists their git versions.
For each one it lists the branch that is checked-out and the SHA.

The `-v` option further performs a `git fetch` on each repository that sits on a branch and reports how many more recent
commits exist beyond the currently used version.

gb-git clone
------------
`gb-git clone <git-repo-path>` retrieves the named repository into the vendor subdirectory and 

How it works
------------
`gb-git` distinguishes two aspects of vendoring: locking a dependency to a specific version and creating a copy of
the dependency for checking into the parent project repository.
Locking-in the version of a dependency ensures that builds are repeatable and changes in the dependency don't create
havoc in the project. Making a copy of the dependency ensures that the dependency doesn't suddenly disappear or get
mutilated. In general everyone wants the loging-in but not everyone wants or needs the copying. In particular, for
dependencies that I own I don't need the copying and for dependencies that are so popular that they won't disappear
from the face of the planet I don't need it either, but everyone has slightly different criteria.

`gb-git` supports both forms of dependencies: locked-in ones and copied-ones. The mechanism it uses are as follows:
- For locked-in dependencies it records the command line to retrieve the dependency in a file with the same name as
  the dependency plus a `.fetch` extension.
- For copied dependencies it renames the `.git` subdirectory of the dependency to `_git` such that git checks the whole
  subtree into the project and provides helpers to make it easy to manipulate the dependency. Basically, to use git in the
  dependency subtree, for example to fetch a new version, one has to use `GIT_DIR=_git git ...`
