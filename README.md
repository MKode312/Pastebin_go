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

### Running
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
### Using
1. To create a profile, you need to send a POST request to this URL: [http://localhost:8082/pastebin/register](http://localhost:8082/pastebin/register)

Example request:
<p align="center">
<img alt="Screenshot showing the registering example." src="imgs/registerExample.png"><br>
</p>

2. To login, send a POST request to this URL: [http://localhost:8082/pastebin/login](http://localhost:8082/pastebin/login)

Example request:
<p align="center">
<img alt="Screenshot showing the login example." src="imgs/loginExample.png"><br>
</p>

3. To save a textblock, send a POST request to this URL: [http://localhost:8082/pastebin/write](http://localhost:8082/pastebin/write)

Example request:
<p align="center">
<img alt="Screenshot showing the writing text example." src="imgs/writeExample.png"><br>
</p>

4. To get a textblock, send a GET request to URL, which you recieved in the previous step.

Example request:
<p align="center">
<img alt="Screenshot showing the getting text example." src="imgs/gettingTextExample.png"><br>
</p>

## Usage
The [Usage](#usage) section would explain in more detail how to run the software, what kind of output or behavior to expect, and so on. It would cover basic operations as well as more advanced uses.

Some of the information in this section will repeat what is in the [Quick start](#quick-start) section. This repetition is unavoidable, but also, not entirely undesirable: the more detailed explanations in this [Usage](#usage) section can help provide more context as well as clarify possible ambiguities that may exist in the more concise [Quick start](#quick-start) section.

If your software is complex and has many features, it may be better to create a dedicated website for your documentation (e.g., in [GitHub Pages](https://pages.github.com), [Read the Docs](https://about.readthedocs.com), or similar) rather than to cram everything into a single linear README file. In that case, the [Usage](#usage) section can be shortened to just a sentence or two pointing people to your documentation site.


### More options
Some projects need to communicate additional information to users and can benefit from additional sections in the README file. It's difficult to give specific instructions here – a lot depends on your software, your intended audience, etc. Use your judgment and ask for feedback from users or colleagues to help figure out what else is worth explaining.

## Known issues and limitations
In this section, summarize any notable issues and/or limitations of your software. If none are known yet, this section can be omitted (and don't forget to remove the corresponding entry in the [Table of Contents](#table-of-contents) too); alternatively, you can leave this section in place and write something along the lines of "none are known at this time".
