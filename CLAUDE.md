# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
```bash
# Development mode
go run main.go

# Build and run
go build -o linebot main.go
./linebot

# With Docker
docker build -t chinese-learning-linebot .
docker run -p 8080:8080 chinese-learning-linebot
```

### Cloudflare Tunnel (for LINE Bot webhook)
```bash
# Expose local server to public (port 8081)
cloudflared tunnel --url http://localhost:8081
```

### Dependency Management
```bash
# Install/update dependencies
go mod tidy

# Download dependencies
go mod download
```

### Environment Setup
The application requires environment variables. Create a `.env` file with:
```env
# LINE Bot Configuration
LINE_CHANNEL_SECRET=your_line_channel_secret_here
LINE_CHANNEL_ACCESS_TOKEN=your_line_channel_access_token_here

# Firebase Configuration
FIREBASE_PROJECT_ID=your_firebase_project_id
GOOGLE_APPLICATION_CREDENTIALS=path/to/your/firebase-service-account-key.json

# Server Configuration
PORT=8080
GIN_MODE=release
```

## Architecture Overview

### Core Architecture Pattern
This is a **layered architecture** with clear separation of concerns:

- **Handlers** (`handlers/`) - HTTP request routing and LINE Bot event handling
- **Services** (`services/`) - Business logic and data operations
- **Models** (`models/`) - Data structures and domain objects
- **Config** (`config/`) - Firebase and LINE Bot initialization
- **Utils** (`utils/`) - Shared utility functions

### Key Dependencies
- **Gin** - Web framework for HTTP routing
- **LINE Bot SDK v7** - LINE messaging platform integration
- **Firebase SDK** - Firestore database operations
- **godotenv** - Environment variable management

### Data Flow Architecture
1. LINE Platform → Webhook (`/webhook`) → `handlers/webhook.go`
2. Event routing → `handlers/message.go` for text messages
3. Business logic → `services/` layer (character.go, lesson.go, practice.go)
4. Data persistence → Firebase Firestore via `config/firebase.go`
5. Response formatting → `utils/response.go` → LINE Platform

## Key Service Layer Functions

### CharacterService (`services/character.go`)
- `LookupCharacter(char string)` - Get detailed character information
- `SearchCharacters(keyword, limit)` - Search for characters
- `GetRandomCharacters(count)` - Generate random characters for practice

### LessonService (`services/lesson.go`)
- `GetLearningProgress(publisher, grade, semester)` - Track learning progress
- `GetLessonsByGrade(publisher, grade)` - Retrieve grade-specific lessons
- `GetCharactersFromLessons(publisher, grade, semester)` - Extract lesson characters

### PracticeService (`services/practice.go`)
- `GeneratePhoneticQuestion()` - Create pronunciation practice
- `GenerateStrokeQuestion()` - Create stroke count practice
- `GenerateSentenceQuestion()` - Create sentence completion practice
- `CheckAnswer(questionID, answer)` - Validate practice answers

## Data Models Structure

### Core Entities
- **CharacterInfo** - Chinese character with phonetic, stroke, radical, meaning, examples
- **LessonInfo** - Curriculum lesson with characters, publisher, grade, semester
- **PracticeQuestion** - Interactive questions with options and explanations
- **PracticeSession** - User practice sessions with scoring and progress

### Database Collections (Firestore)
- `characters` - Character information and metadata
- `lessons` - Lesson content organized by publisher/grade/semester
- `cumulative_characters` - Aggregated character learning statistics
- Practice sessions are cached in-memory with expiration

## Important Implementation Notes

### Event Handling Pattern
The webhook handler (`handlers/webhook.go`) uses event type switching:
- `EventTypeMessage` → `handleMessage()` (implemented in `message.go`)
- `EventTypePostback` → `handlePostback()` (currently placeholder)
- `EventTypeFollow`/`EventTypeUnfollow` → User lifecycle events

### Error Handling Strategy
- Firebase initialization failures are logged as warnings, app continues without Firebase
- LINE Bot initialization failures are logged as warnings, app continues without LINE Bot
- Individual operation errors are logged but don't crash the application

### Missing Components
- `handlers/postback.go` - Referenced in README but not implemented
- `.env.example` - Template for environment variables
- Test files - No test coverage currently exists
- No build scripts or CI/CD configuration

## Firebase Integration

### Configuration
Firebase client is initialized in `config/firebase.go` and passed as dependency to handlers and services. The app uses service account authentication via `GOOGLE_APPLICATION_CREDENTIALS`.

### Collections Structure
- Publisher-based organization (康軒, 翰林, 南一)
- Grade-level curriculum (1-6)
- Semester-based progression (上學期, 下學期)

## Development Workflow

When adding new features:
1. Define data structures in `models/`
2. Implement business logic in `services/`
3. Add HTTP handlers in `handlers/`
4. Update response formatting in `utils/response.go`
5. Test via LINE Bot webhook endpoint

The application follows Go conventions with proper error handling, dependency injection, and modular design patterns suitable for educational chatbot functionality.