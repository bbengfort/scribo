# Scribo

[![Build Status](https://travis-ci.org/bbengfort/scribo.svg?branch=master)](https://travis-ci.org/bbengfort/scribo)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbengfort/scribo)](https://goreportcard.com/report/github.com/bbengfort/scribo)
[![Stories in Ready](https://badge.waffle.io/bbengfort/scribo.png?label=ready&title=Ready)](https://waffle.io/bbengfort/scribo)

**The web API to record data for the Mora ping measurement app.**

![Mora Architecture Diagram](http://bbengfort.github.io/assets/images/2016-05-10-mora-architecture.png)

This application is an uptime collection mechanism associated with the Mora project. Mora contains three pieces: Oro and Scio which ping each other to measure latency inside of the network, then report those pings to Scribo, which is simply a RESTful API designed to record experimental data.

Scribo is intended to be a lightweight, fast microservice -- it is written in Go and deployed and scaled by Heroku (currently).

## Getting Started

In order to get started with Scribo, you can first `go get` or clone the repository:

    $ go get github.com/bbengfort/scribo/...  

This will fetch the Scribo packages and build them along with the `scribo` command that is in `cmd/scribo` as per &ldquo;[Structuring Applications in Go](https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091#.fn6lyl49z)&rdquo;. If you're using `godep` you can restore the dependencies to start development:

    $ godep restore

The web server can then be run as follows:

    $ go run cmd/scribo/main.go

And the tests can be run as follows:

    $ ginkgo -r -v ./scribo/...

Hopefully that's enough to get you up and running!

## About

Mora (delay, waiting) observes ping latencies between nodes in a wide area, heterogenous, user-oriented network by running a local service that pings other nodes in the network. Oro (speak) is the name of the mobile application, and Scio (understand) is the name of the desktop client. The ping data is collected by a centralized RESTful microservice called Scribo (record). This data will be used for scientific research concerning distributed systems.

### Contributing

Scribo is open source, and I'd love your help, particularly if you are a student at the University of Maryland and are interested in studying distributed systems. If you would like to contribute, you can do so in the following ways:

1. Add issues or bugs to the bug tracker: [https://github.com/bbengfort/scribo/issues](https://github.com/bbengfort/scribo/issues)
2. Work on a card on the dev board: [https://waffle.io/bbengfort/scribo](https://waffle.io/bbengfort/scribo)
3. Create a pull request in Github: [https://github.com/bbengfort/scribo/pulls](https://github.com/bbengfort/scribo/pulls)

Note that labels in the Github issues are defined in the blog post: [How we use labels on GitHub Issues at Mediocre Laboratories](https://mediocre.com/forum/topics/how-we-use-labels-on-github-issues-at-mediocre-laboratories).

When doing a pull request, keep in mind that the project is set up in a typical production/release/development cycle as described in _[A Successful Git Branching Model](http://nvie.com/posts/a-successful-git-branching-model/)_. A typical workflow is as follows:

1. Select a card from the [dev board](https://waffle.io/bbengfort/scribo) - preferably one that is "ready" then move it to "in-progress".

2. Create a branch off of develop called "feature-[feature name]", work and commit into that branch.

        ~$ git checkout -b feature-myfeature develop

3. Once you are done working (and everything is tested) merge your feature into develop.

        ~$ git checkout develop
        ~$ git merge --no-ff feature-myfeature
        ~$ git branch -d feature-myfeature
        ~$ git push origin develop

4. Repeat. Releases will be routinely pushed into master via release branches, then deployed to the server.

Note that no pull requests into master will be considered; only those that pull into develop.

### Throughput

[![Throughput Graph](https://graphs.waffle.io/bbengfort/scribo/throughput.svg)](https://waffle.io/bbengfort/scribo/metrics/throughput)

## Contributors

Thank you for all your help contributing to make Scribo a great project!

### Maintainers

- Benjamin Bengfort: [@bbengfort](https://github.com/bbengfort/)

### Contributors

- Your name here!

## Changelog

The release versions that are tagged in Git. You can see the tags through the GitHub web application and download the tarball of the version you'd like.

The versioning uses a three part version system, "a.b.c" - "a" represents a major release that may not be backwards compatible. "b" is incremented on minor releases that may contain extra features, but are backwards compatible. "c" releases are bug fixes or other micro changes that developers should feel free to immediately update to.

## Version 0.1

* **tag**: [v0.1](https://github.com/bbengfort/scribo/releases/tag/v0.1)
* **deployment**: Tuesday, May 10, 2016
* **commit**: [3ad53af](https://github.com/bbengfort/scribo/commit/3ad53af3b06f05016477efe49a0977e8f19b6d7d)

This is just a pre-release and represents the first (working) push of the simple code base to Heroku for testing. There is still a long way to go before we're ready for an official release. However, this is my first Go microservice, so I'm excited to release it into the wild!
