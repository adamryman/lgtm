# lgtm; slackbot for #pr-party [![Build Status](https://travis-ci.org/StudentRND/lgtm.svg?branch=master)](https://travis-ci.org/StudentRND/lgtm)

lgtm lets users add github integrations to the slack channel #pr-party for
posting to the channel when a new pull request opens.

## Usage

See `lgtm.go` for environment variables to configure.

[Register a github oauth application](https://github.com/settings/applications/new) for your host and path `lgtm/authorize/callback`.

```
e.x.
https://srnd.org/lgtm/authorize/callback
```

## (UN)LICENSE

[This is free and unencumbered software released into the public domain.](./UNLICENSE)
