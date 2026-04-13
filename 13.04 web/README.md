Öğrenci: Nazar Eylül Sayında
Okul No: 24080410213

# GoLearn - Full-Stack LMS (Learning Management System) 🎓

## 📝 Proje Tanımı
GoLearn, modern web teknolojileri ile geliştirilmiş bir öğrenme yönetim sistemidir. Öğrencilerin eğitim materyallerine erişebildiği, ilerleme takibi yapabildiği ve canlı sınıflarda etkileşime geçebildiği dinamik bir platform sunar. Bu proje, "Web Tabanlı Programlama" dersi "Finale Doğru" serisi kapsamında geliştirilmiştir.

---

## 🚀 Kurulum ve Çalıştırma

### 1. Ön Koşullar
- Docker ve Docker Compose yüklü olmalıdır.
- Node.js (v18+) yüklü olmalıdır.

### 2. Backend (API) Başlatma
Backend servisi Dockerize edilmiştir.
```bash
cd golearn
docker-compose up -d --build
```
- **API URL:** `http://localhost:8090`
- **Swagger Dokümantasyonu:** `http://localhost:8090/swagger/index.html`

### 3. Frontend Başlatma
```bash
cd golearn-frontend
npm install
npm run dev
```
- **Uygulama URL:** `http://localhost:5173`

---

## 🛠️ Teknik Özellikler & Fonksiyonellik
- **Backend:** Go (Gin Framework), GORM, SQLite.
- **Frontend:** HTML5, Vanilla CSS (Glassmorphism), JavaScript (ES6+), Vite.
- **Güvenlik:** JWT (JSON Web Token) Auth, RBAC (Role Based Access Control), API Rate Limiting.
- **Gerçek Zamanlı İletişim:** WebSocket (Classroom Chat).
- **Veritabanı Mimarisi:** User, Course, Lesson, Quiz, Progress modelleri ile tam ilişkisel yapı.

---

## 🗺️ API Endpoint Listesi (Önemli)

### Kimlik Doğrulama (Auth)
- `POST /api/auth/register` - Yeni kullanıcı kaydı.
- `POST /api/auth/login` - Kullanıcı girişi ve JWT üretimi.

### Kurs Yönetimi
- `GET /api/courses` - Tüm kursları listele (Filtreleme desteği).
- `GET /api/courses/:id` - Kurs detaylarını getir.
- `POST /api/courses` - Yeni kurs oluştur (Sadece Öğretmen).
- `GET /api/courses/:id/lessons` - Kursa ait dersleri listele.

### İlerleme & Etkileşim
- `POST /api/lessons/:id/complete` - Dersi tamamlandı olarak işaretle.
- `GET /api/my/progress` - Mevcut kullanıcının tüm kurslardaki ilerleme durumunu getir.
- `GET /ws/classroom/:courseId` - Canlı sınıf sohbet odası (WebSocket).

---

## 📂 Klasör Yapısı
- `golearn/` - Go Backend Projesi
  - `handlers/` - API Mantığı
  - `models/` - Veritabanı Modelleri
  - `middleware/` - Auth ve Güvenlik
- `golearn-frontend/` - Frontend Projesi
  - `src/` - JS Mantığı ve API Servisleri
  - `style.css` - UI Design System

---

**Nazar Eylül Sayında - 24080410213**
