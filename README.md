Tools Directory Structure
==========================
```bash
├── cmd/
│   └── main.go                 # Main entry point of the application
├── config/
│   └── config.go               # General configuration for the tool (TDS token, delay, etc.), loaded from `config.
├── internal/
│   ├── facebook/
│   │   └── facebook.go         # Facebook task: Like, Follow, Comment, Share
│   ├── traodoisub/
│   │   └── tds.go              # Interact with Traodoisub API: Get job, confirm, get xu
│   └── job/
│       └── executor.go         # Handle job for each Facebook account, using concurrency
├── models/
│   └── types.go                # Data structure definitions: Job, CookieUser, TDSProfile
├── utils/
│   └── logger.go               # Log with color, timestamp, easier to debug
├── .env                        # Store TDS token
├── cookie.json                 # Store Facebook cookie
├── go.sum                      # Go module configuration
└── go.mod                      # Go module configuration
```


