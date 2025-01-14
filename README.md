# Poster

Poster is a tool for posting content to social networks. Currently, it supports VK and YouTube, but it is designed to be easily extendable to other platforms.

## Features

- Post content to VK
- Post content to YouTube
- Easily extendable to other social networks

## Installation

1. Clone the repository:

```sh
git clone https://github.com/sunchess/poster.git
cd poster
```

### Install dependencies:

```sh
go mod tidy
```

### Set up your environment variables.
 Create a .env file in the root directory and add the following:

```sh
DB_PATH=your_database_path
VK_API_TOKEN=your_vk_api_token
YOUTUBE_API_KEY=your_youtube_api_key
```
## Usage

To use VK Poster, you can run the application with the following command:
