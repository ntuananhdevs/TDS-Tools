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
│   └── Cookie_FB.txt           # Lưu trữ cookie của các tài khoản Facebook
└── go.mod                      # Cấu hình module Go


.env                            # Lưu trữ biến môi trường, như access_token, cookie FB
.env.example                    # Mẫu cấu hình biến môi trường

.gitignore                       # Các file và thư mục cần bỏ qua khi commit
go.mod                           # Quản lý phụ thuộc trong Go
go.sum                           # Đảm bảo tính toàn vẹn của các phụ thuộc Go


