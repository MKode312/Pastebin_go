# Pastebin

This is a pet-project Pastebin API, that features writing blocks of text, saving them and getting a link for them as well. When you follow the link, you see the textblock itself.

## Table of contents

* [Introduction](#introduction)
* [Installation](#installation)
* [Quick start](#quick-start)
* [Usage](#usage)
* [Known issues and limitations](#known-issues-and-limitations)


## Introduction

The main goal of this project was to practice using some new techologies for me, such as MinIO and Redis cache. I believe, that I succeeded in this. The project isn't meant to solve real problems, but it might be useful in some ways.


## Installation

### Prerequisites

* [Go](https://go.dev/doc/install) v1.24.2
* [Docker Engine](https://docs.docker.com/engine/install/)
* [Postman](https://www.postman.com/downloads/) (optionally)

### Installing the project

To install the project, simply use this command:

* Clone the repository:
```bash
git clone https://github.com/MKode312/Pastebin_go.git
```

## Quick start

To start the system, use the following commands:

* Change directory:
```bash
cd Pastebin_go
```
* Run the redis-cache docker-container:
```bash
docker run --name redis-server -p 6379:6379 redis:latest
```

* Run the main docker-container:
```bash
docker-compose up
```

* Run the MinIO client:
```bash
go run cmd/minio/main.go
```

* Run the API:
```bash
go run cmd/pastebin/main.go
```

1. To create a profile, you need to send a POST request to this URL:
[http://localhost:8082/pastebin/register](http://localhost:8082/pastebin/register) 

Example request:



## Usage

The [Usage](#usage) section would explain in more detail how to run the software, what kind of output or behavior to expect, and so on. It would cover basic operations as well as more advanced uses.

Some of the information in this section will repeat what is in the [Quick start](#quick-start) section. This repetition is unavoidable, but also, not entirely undesirable: the more detailed explanations in this [Usage](#usage) section can help provide more context as well as clarify possible ambiguities that may exist in the more concise [Quick start](#quick-start) section.

If your software is complex and has many features, it may be better to create a dedicated website for your documentation (e.g., in [GitHub Pages](https://pages.github.com), [Read the Docs](https://about.readthedocs.com), or similar) rather than to cram everything into a single linear README file. In that case, the [Usage](#usage) section can be shortened to just a sentence or two pointing people to your documentation site.


### Basic operation

When learning how to use anything but the simplest software, new users may appreciate beginning with basic features and modes of operation. If your software has a help system of some kind (e.g., in the form of a command-line flag such as `--help`, or a menu item in a GUI), explaining it is an excellent starting point for this section.

The basic approach for using this README file is as follows:

1. Copy the [README source file](https://raw.githubusercontent.com/mhucka/readmine/main/README.md) to your repository
2. Delete the body text but keep the section headings
3. Replace the title heading (the first line of the file) with the name of your software
4. Save the resulting skeleton file in your version control system
5. Continue by writing your real README content in the file

The first paragraph in the README file (under the title at the top) should summarize your software in a concise fashion, preferably using no more than one or two sentences as illustrated by the circled text in the figure below.

<p align="center">
<img alt="Screenshot showing the top portion of this file on the web." width="80%" src="https://raw.githubusercontent.com/mhucka/readmine/main/.graphics/screenshot-top-paragraph.png"><br>
<em>Figure: Screenshot showing elements of the top portion of this file.</em>
</p>

The space under the first paragraph and _before_ the [Table of Contents](#table-of-contents) is a good location for optional [badges](https://github.com/badges/shields), which are small visual tokens commonly used on GitHub repositories to communicate project status, dependencies, versions, DOIs, and other information. (Two example badges are shown in the figure above, under the circled text.) The particular badges and colors you use depend on your project and personal tastes.


### More options

Some projects need to communicate additional information to users and can benefit from additional sections in the README file. It's difficult to give specific instructions here – a lot depends on your software, your intended audience, etc. Use your judgment and ask for feedback from users or colleagues to help figure out what else is worth explaining.


## Known issues and limitations

In this section, summarize any notable issues and/or limitations of your software. If none are known yet, this section can be omitted (and don't forget to remove the corresponding entry in the [Table of Contents](#table-of-contents) too); alternatively, you can leave this section in place and write something along the lines of "none are known at this time".


## Getting help

Inform readers how they can contact you, or at least how they can report problems they may encounter. This could take the form of a request to use the issue tracker on your repository. Some projects have associated discussion forums or mailing lists, and this section is a good place to mention those.


## Contributing

If your project accepts open-source contributions, this is where you can welcome contributions and explain to readers how they can go about it. Mention the [`CONTRIBUTING.md`](CONTRIBUTING.md) file in your repository, if you have one.


## License

This section should state any copyright asserted on the project materials as well as the terms of use for the software, files and other materials found in the project repository.

_This_ README file is itself distributed under the terms of the [Creative Commons 1.0 Universal license (CC0)](https://creativecommons.org/publicdomain/zero/1.0/). The license applies to this file and other files in the [GitHub repository](http://github.com/mhucka/readmine) hosting this file. This does _not_ mean that you, as a user of this README file in your software project, must also use CC0 license!  You may use whatever license for your work you prefer, or whatever you are required to use by your employer or sponsor.


## Acknowledgments

This final section is where you should acknowledge funding and/or institutional support, prior work that influenced or inspired your project, resources that you used (such as other people's software), important contributions from other people, and anything else that deserves mention. After all, nothing is truly done in isolation; everything is built on top of something, and we all owe debts to other projects and people who helped us, supported us, and influenced us.

For example, in the process of developing this file, I used not only my own ideas: I read many (sometimes contradictory) recommendations for README files, examined real READMEs in actual use, and tried to distill the best ideas into the result you see here. Sources included the following:

* [Readme Driven Development](http://tom.preston-werner.com/2010/08/23/readme-driven-development.html)
* [How to Write a Good README](https://dev.to/merlos/how-to-write-a-good-readme-bog)
* [How To Write A Great README](https://thoughtbot.com/blog/how-to-write-a-great-readme)
* [How to Write an Awesome Readme](https://dev.to/documatic/how-to-write-an-awesome-readme-cfl)
* [Readme Best Practices](https://github.com/jehna/readme-best-practices)
* [Art of README](https://github.com/noffle/art-of-readme)
* [Making a useful README file for research projects](http://jonathanpeelle.net/making-a-readme-file)
* [Tips for Making your GitHub Profile Page Accessible](https://github.com/orgs/community/discussions/64778)
* [Make a README](https://www.makeareadme.com)
* [ReadMe.so](https://readme.so)
* [common readme](https://github.com/noffle/common-readme)
* [Standard Readme](https://github.com/RichardLitt/standard-readme)
* [CFPB Open Source Project Template Instructions](https://github.com/cfpb/open-source-project-template)
* [README-Template.md](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
* [open-source-template](https://github.com/davidbgk/open-source-template/)
* [Feedmereadmes: A README Help Exchange](https://github.com/lappleapple/feedmereadmes)
* [Awesome README List](https://github.com/matiassingers/awesome-readme)
* [Top ten reasons why I won't use your open source project](https://changelog.com/posts/top-ten-reasons-why-i-wont-use-your-open-source-project)
*


```bash
docker run -d --name redis-server -p 6379:6379 redis:latest
```

```bash
docker-compose up
```

```bash
go run cmd/minio/main.go
```

```bash
go run cmd/pastebin/main.go
```
