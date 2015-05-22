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

gb-git info
-----------
`gb-git info` finds all git repository workspaces in the gb project's vendor/src subdirectory and lists their git versions.
For each one it lists the branch that is checked-out and the SHA.

The `-v` option further performs a `git fetch` on each repository that sits on a branch and reports how many more recent
commits exist beyond the currently used version.

