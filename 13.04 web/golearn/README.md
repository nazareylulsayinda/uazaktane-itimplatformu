# GoLearn LMS Backend API

GoLearn, öğretmenlerin kurs ve ders (video/içerik) oluşturup öğrencilerin kayıt olduğu, quizler çözüp ilerlemelerini takip ettikleri bir Uzaktan Eğitim Platformu (LMS - Learning Management System) backend hizmetidir.

## 🚀 Kullanılan Teknolojiler
- **Go (Golang) & Gin Framework:** Hızlı ve modern HTTP Router
- **GORM & SQLite:** ORM ve Veritabanı
- **JWT (JSON Web Token):** Güvenli Kimlik Doğrulama
- **Gorilla WebSocket:** Canlı Sınıf İçi Mesajlaşma Odaları
- **Swagger:** API Dokümantasyonu
- **Docker & Docker Compose:** Container Mimari

## 🛠️ Nasıl Çalıştırılır?

Projeyi Docker kullanarak çok basit bir şekilde ayağa kaldırabilirsiniz. Host makinenizde Go yüklü olması şart değildir.

```bash
cd "13.04 web/golearn"
docker-compose up --build
```
> Bu komut, Dockerfile içindeki adımları yürütür (Gereken Go paketlerini kurar, swag init ile swagger dökümanlarını yaratır ve derler).

Uygulama **localhost:8090** portunda çalışacaktır.

## 📚 API Dokümantasyonu (Swagger)

Proje Docker üzerinden başlatıldıktan sonra tarayıcınızda şu adresi ziyaret edebilirsiniz:
=> http://localhost:8090/swagger/index.html

*(Not: Lokal PC'nizde geliştirmek isterseniz `go install github.com/swaggo/swag/cmd/swag@latest` ve `swag init` çalıştırıp ardından `go run main.go` diyebilirsiniz).*

## 📌 Örnek Test Komutları (CURL)

### 1) Öğretmen Register
```bash
curl -X POST http://localhost:8090/api/auth/register \
-H "Content-Type: application/json" \
-d '{"name":"Ahmet Hoca","email":"ahmet@golearn.com","password":"123","role":"teacher"}'
```

### 2) Login
```bash
curl -X POST http://localhost:8090/api/auth/login \
-H "Content-Type: application/json" \
-d '{"email":"ahmet@golearn.com","password":"123"}'
```
> Çıktıdaki `token` değerini kopyalayıp aşağıdaki işlemler için Bearer Token olarak kullanın.

### 3) Kurs Oluşturma (Sadece Öğretmen)
```bash
curl -X POST http://localhost:8090/api/courses \
-H "Authorization: Bearer <TOKEN BURAYA>" \
-H "Content-Type: application/json" \
-d '{"title":"Go Programlama","description":"Baştan Sona Go","category":"Yazılım"}'
```

### 4) Quiz Ekleme
```bash
curl -X POST http://localhost:8090/api/lessons/1/quiz \
-H "Authorization: Bearer <TOKEN BURAYA>" \
-H "Content-Type: application/json" \
-d '{"questions": [{"text":"Go ne zaman çıktı?","option_a":"2005","option_b":"2009","option_c":"2012","option_d":"2015","correct":"B"}]}'
```

## 💬 WebSocket Kullanım Örneği
WebSocket odalarına bağlanmak için bir websocket client (Postman WebSocket tool, wscat veya basit bir HTML JS scripti) kullanabilirsiniz.

```bash
wscat -c ws://localhost:8090/ws/classroom/1?username=Ahmet
```
Bağlanınca `{"type":"join","username":"Ahmet","text":"Ahmet has joined the classroom","course_id":"1"}` şeklinde veri gelir. Yazdığınız mesajlar json olarak okutulup diğer üyelere iletilir.
