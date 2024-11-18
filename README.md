# Go-SQLite-Blog

Go-SQLite-Blog is a lightweight, full-featured blogging platform written in Go (Golang). It uses SQLite for data storage, making it easy to set up without the need for a separate database server. Ideal for personal blogs, Go-SQLite-Blog offers essential features for content creation and management.

## Table of Contents

- [Features](#features)
- [Demo](#demo)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [Contributing](#contributing)
- [License](#license)

## Features

- Admin Friendly: Easy to use, manage and backup - there is just the executable and the database file
- User Management: Support for multiple user roles (admin, editor, author)
- Article Management: Create, edit, and delete blog posts with Markdown support
- Comments: Enable readers to leave comments on posts (with moderation options)
- Categories and Tags: Organize posts using categories and tags
- Media Uploads: Upload and manage images and other media files
- Static Pages: Create static pages like "About Us" or "Contact"
- Search Functionality: Full-text search on articles
- Pagination: Navigate through posts with ease

## Demo

Check out a live demo of Go-SQLite-Blog [here](#). - _Coming Soon_

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/AndreHeber/go-sqlite-blog.git
   cd go-sqlite-blog
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o go-sqlite-blog
   ```

## Usage

### Configuration

Before running the application, you can configure settings in the `config.yaml` file:

```yaml
app:
  port: 8080
  templates_dir: templates
  static_dir: static
database:
  path: ./data/Go-SQLite-Blog.db
```

### Running the Application

1. Start the server

```bash
./go-sqlite-blog
```

2. Access the application

    Open your web browser and navigate to http://localhost:8080.

## Project Structure

```
Go-SQLite-Blog/
├── cmd/
│   └── Go-SQLite-Blog/
│       └── main.go          # Entry point of the application
├── internal/
│   ├── db/
│   │   └── sqlite.go        # Database connection and queries
│   ├── handlers/
│   │   ├── article.go       # Handlers for article routes
│   │   ├── user.go          # Handlers for user routes
│   │   └── comment.go       # Handlers for comment routes
│   ├── models/
│   │   ├── article.go       # Article model
│   │   ├── user.go          # User model
│   │   └── comment.go       # Comment model
│   └── templates/           # HTML templates
├── static/                  # Static files (CSS, JS, images)
├── config.yaml              # Configuration file
├── schema.sql               # Database schema
├── go.mod                   # Go module file
├── LICENSE
└── README.md
```

## Database Schema

Go-SQLite-Blog uses SQLite with the following tables:

```sql
users: Manage user accounts and roles.
articles: Store blog posts with metadata.
comments: Allow users to comment on articles.
categories: Organize articles into categories.
article_categories: Link articles to categories (many-to-many relationship).
tags: Tag articles for detailed classification.
article_tags: Link articles to tags (many-to-many relationship).
media: Manage uploaded media files.
pages: Create static pages.
settings: Store application settings.
```

Refer to the schema.sql file for detailed definitions.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository

```bash
git clone https://github.com/yourusername/go-sqlite-blog.git
cd go-sqlite-blog
```

2. Create a feature branch

```bash
git checkout -b feature/YourFeature
```

3. Make your changes

4. Run the tests

```bash
go test ./...
```

5. Commit your changes

```bash
git commit -m 'Add YourFeature'
```

6. Push to the branch

```bash
git push origin feature/YourFeature
```

7. Open a pull request

## License

This project is licensed under the MIT License. See the LICENSE file for details.