
```bash
App
├── cmd/
│   └── main.go                 # Main entry point of the application
├── config/
│   └── config.go               # Cấu hình chung cho tool (token TDS, delay, ...), được load từ file `config.json`
├── internal/
│   ├── facebook/
│   │   └── facebook.go         # Tác vụ xử lý Facebook: Like, Follow, Comment, Share
│   ├── traodoisub/
│   │   └── tds.go              # Tương tác với API Traodoisub: Lấy job, xác nhận, lấy xu
│   └── job/
│       └── executor.go         # Xử lý nhiệm vụ cho từng tài khoản Facebook, sử dụng concurrency
├── models/
│   └── types.go                # Các định nghĩa cấu trúc dữ liệu: Job, CookieUser, TDSProfile
├── utils/
│   └── logger.go               # In log với màu sắc, timestamp, giúp debug dễ dàng hơn
├── .env      # Lưu trữ token TDS
├── cookie.json                  # Lưu trữ cookie Facebook
├── go.sum                      # Cấu hình module Go
└── go.mod                      # Cấu hình module Go
```


