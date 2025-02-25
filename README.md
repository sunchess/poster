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
YOUTUBE_API_KEY=your_youtube_api_key.json
```
## Usage

To use VK Poster, you can run the application with the following command:

```sh
go run cmd/app/main.go --scope vk --limit 2 --post_gap 1800
```

## Command Line Arguments

- `--scope:` The platform to post to (`vk` or `youtube`). Default is vk.
- `--limit:` The number of posts to process in one run. Default is 2.
- `--post_gap:` The gap between posts in seconds. Default is 1800 (30 minutes).

## Extending Functionality
VK Poster is designed to be easily extendable. To add support for a new social network:

1. Implement a new service that satisfies the `PostingServiceInterface`.
2. Add the new service to the `GetGateway` function in `posting.go`.

## Example
Here is an example of how to run the application:

```sh
./poster --scope vk --limit 5 --post_gap 3600
```

This command will post to VK, processing up to 5 posts with a 1-hour gap between each post.

## Credentials
To set up posting videos to YouTube Shorts via the API, you need to use the YouTube Data API v3. This API allows you to upload videos, specify their settings, and metadata.

Setup steps:

### Create a project in `Google Cloud Console`:

- Go to Google Cloud Console.
- Create a new project.
- Enable YouTube Data API v3 for this project.
- Create credentials (OAuth 2.0 Client ID) to access the API.

During the execution of the program, a link for authorization will appear in the console. The key obtained as a result of authorization needs to be inserted into the console.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
This project is licensed under the MIT License.