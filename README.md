
```bash
App
├── cmd/
│   └── main.go                 # Điểm khởi chạy chính của ứng dụng
├── config/
│   └── config.go               # Cấu hình chung cho tool (token TDS, delay, ...)
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
├── assets/
│   ├── configtds.txt           # Lưu trữ token TDS
│   └── Cookie_FB.txt  
├── .env      # Lưu trữ token TDS
├── cookie.json                  # Lưu trữ cookie Facebook
├── go.sum                      # Cấu hình module Go
└── go.mod                      # Cấu hình module Go
```


