# resume-generator

This is a toy project of mine that allows one to generate custom resume PDFs from subsets of their job history.

## Use

This project uses [Devbox][devbox] to manage the development environment.

With Devbox installed:

1. Run a dev shell.
1. Compile the generator: `go build`
1. Copy the example resume data: `cp data.example.json data.json`
1. Edit the resume data and add your own info.
1. Run the generator and make your custom resume.

## Issues

You're free to add issues into Github but it is **highly unlikely** that I respond to them.
I do not intend to maintain this app long-term. Sorry 💛.

## Development

1. Install [Devbox][devbox].
1. Use your favorite editor to make changes.
1. Compile and test: `go build; ./resume-generator`

There's no tests here. Like I said, it's a toy project.

[devbox]: https://www.jetify.com/devbox/docs/
