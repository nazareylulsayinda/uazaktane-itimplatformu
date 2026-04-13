import { AuthApi, CourseApi, LessonApi } from './api.js';

// --- State Management ---
function getUser() {
    try { 
        return JSON.parse(localStorage.getItem('user')); 
    } catch(e) { 
        return null; 
    }
}

// --- DOM elements ---
const contentDiv = document.getElementById('app-content');
const navLinks = document.getElementById('nav-links');

// --- Global App Logic ---
window.addEventListener('hashchange', () => {
    updateNav();
    router();
});

window.addEventListener('load', () => {
    updateNav();
    router();
});

function navigateTo(hash) {
    window.location.hash = hash;
}

function updateNav() {
    const user = getUser();
    if (user) {
        navLinks.innerHTML = `
            <span class="text-muted" style="margin-right: 10px;">Hoşgeldin, <b>${user.name}</b></span>
            <a href="#dashboard" class="btn btn-outline" style="padding: 0.5rem 1rem; margin-right: 10px;">Dashboard</a>
            <button id="logout-btn" class="btn btn-primary" style="padding: 0.5rem 1rem;">Çıkış Yap</button>
        `;
        document.getElementById('logout-btn').addEventListener('click', () => {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            updateNav();
            navigateTo('login');
        });
    } else {
        navLinks.innerHTML = `
            <a href="#login" class="btn btn-outline" style="padding: 0.5rem 1rem; margin-right:  10px;">Giriş Yap</a>
            <a href="#register" class="btn btn-primary" style="padding: 0.5rem 1rem;">Kayıt Ol</a>
        `;
    }
}

// --- Simple Toast Notification ---
window.showToast = (message, type = 'success') => {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `glass-panel animate-fade-in`;
    toast.style.padding = '1rem';
    toast.style.marginBottom = '0.5rem';
    toast.style.borderLeft = `4px solid var(--${type})`;
    toast.style.color = 'white';
    toast.style.backgroundColor = type === 'danger' ? 'rgba(239, 68, 68, 0.6)' : 'rgba(16, 185, 129, 0.6)';
    toast.innerHTML = message;
    container.appendChild(toast);
    setTimeout(() => {
        toast.style.opacity = '0';
        setTimeout(() => toast.remove(), 300);
    }, 4000);
}

// --- Router Logic ---
function router() {
    let hash = window.location.hash || '#';
    const user = getUser();
    
    // Auth protection
    if (!user && hash !== '#login' && hash !== '#register') {
        navigateTo('login');
        return;
    }

    // Redirect if already logged in
    if (user && (hash === '#login' || hash === '#register')) {
        navigateTo('dashboard');
        return;
    }
    
    // Route matching
    if (hash === '#login') return renderLogin();
    if (hash === '#register') return renderRegister();
    if (hash === '#dashboard' || hash === '#') return renderDashboard();
    
    if (hash.startsWith('#course/')) {
        const id = hash.split('/')[1];
        return renderCourseDetail(id);
    }
    
    if (hash.startsWith('#classroom/')) {
        const id = hash.split('/')[1];
        return renderClassroom(id);
    }

    contentDiv.innerHTML = '<div class="text-center mt-8"><h2>404 - Sayfa Bulunamadı</h2><a href="#dashboard" class="btn btn-primary mt-4">Dashboarda Dön</a></div>';
}

// --- Pages ---

function renderLogin() {
    contentDiv.innerHTML = `
        <div class="grid" style="place-items: center; min-height: 60vh;">
            <div class="card" style="width: 100%; max-width: 400px;">
                <h2 class="card-title text-center text-gradient mb-8">Giriş Yap</h2>
                <form id="login-form">
                    <div class="form-group">
                        <label>E-posta</label>
                        <input type="email" id="email" class="form-control" required placeholder="ornek@golearn.com">
                    </div>
                    <div class="form-group">
                        <label>Şifre</label>
                        <input type="password" id="password" class="form-control" required placeholder="••••••">
                    </div>
                    <button type="submit" class="btn btn-primary btn-block mb-4">Giriş</button>
                </form>
                <div class="text-center text-muted" style="font-size:0.9rem;">
                    Hesabın yok mu? <a href="#register">Kayıt Ol</a>
                </div>
            </div>
        </div>
    `;

    document.getElementById('login-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        
        try {
            const res = await AuthApi.login(email, password);
            localStorage.setItem('token', res.token);
            // Now we use the user object directly from the backend
            localStorage.setItem('user', JSON.stringify(res.user));
            
            updateNav();
            showToast('Başarıyla giriş yapıldı!');
            navigateTo('dashboard');
        } catch (err) {
            showToast(err.message, 'danger');
        }
    });
}

function renderRegister() {
    contentDiv.innerHTML = `
        <div class="grid" style="place-items: center; min-height: 60vh;">
            <div class="card" style="width: 100%; max-width: 400px;">
                <h2 class="card-title text-center text-gradient mb-8">Kayıt Ol</h2>
                <form id="register-form">
                    <div class="form-group">
                        <label>Ad Soyad</label>
                        <input type="text" id="name" class="form-control" required>
                    </div>
                    <div class="form-group">
                        <label>E-posta</label>
                        <input type="email" id="email" class="form-control" required>
                    </div>
                    <div class="form-group">
                        <label>Şifre</label>
                        <input type="password" id="password" class="form-control" required>
                    </div>
                    <div class="form-group">
                        <label>Rol</label>
                        <select id="role" class="form-control">
                            <option value="student">Öğrenci</option>
                            <option value="teacher">Öğretmen</option>
                        </select>
                    </div>
                    <button type="submit" class="btn btn-primary btn-block mb-4">Kayıt Ol</button>
                </form>
                <div class="text-center text-muted" style="font-size:0.9rem;">
                    Zaten üye misin? <a href="#login">Giriş Yap</a>
                </div>
            </div>
        </div>
    `;

    document.getElementById('register-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const role = document.getElementById('role').value;
        
        try {
            const res = await AuthApi.register(name, email, password, role);
            showToast('Kayıt başarılı! Şimdi giriş yapabilirsiniz.');
            navigateTo('login');
        } catch (err) {
            showToast(err.message, 'danger');
        }
    });
}

async function renderDashboard() {
    const user = getUser();
    contentDiv.innerHTML = `
        <div class="hero-section animate-fade-in">
            <div class="hero-content">
                <h1 class="mb-4">Merhaba, <span class="text-gradient">${user?.name}</span>! 👋</h1>
                <p class="text-muted" style="max-width: 600px;">GoLearn platformuna hoş geldin. Bugün yeni bir şeyler öğrenmek için harika bir gün!</p>
            </div>
        </div>

        <div id="dashboard-stats" class="stats-grid">
            <div class="stat-card">
                <span class="stat-value animate-fade-in" id="stat-courses">-</span>
                <span class="stat-label">Aktif Kurslar</span>
            </div>
            <div class="stat-card">
                <span class="stat-value animate-fade-in" id="stat-progress">-</span>
                <span class="stat-label">Ort. İlerleme</span>
            </div>
            <div class="stat-card">
                <span class="stat-value animate-fade-in" id="stat-extra">-</span>
                <span class="stat-label">${user?.role === 'teacher' ? 'Tamamlanan' : 'Sertifikalar'}</span>
            </div>
        </div>

        <div class="flex items-center justify-between mb-4">
            <h2 class="text-gradient">Kurs Keşfet</h2>
            ${user?.role === 'teacher' ? '<button id="btn-create-course" class="btn btn-primary">+ Yeni Kurs Ekle</button>' : ''}
        </div>

        <div class="category-pills">
            <div class="pill active" data-category="all">Hepsi</div>
            <div class="pill" data-category="Yazılım">Yazılım</div>
            <div class="pill" data-category="Tasarım">Tasarım</div>
            <div class="pill" data-category="İşletme">İşletme</div>
            <div class="pill" data-category="Kişisel Gelişim">Kişisel Gelişim</div>
        </div>

        <div id="course-list" class="grid md:grid-cols-3 mt-4">
            <div class="text-muted">Kurslar yükleniyor...</div>
        </div>
        
        <!-- Modal Form (Hidden) -->
        <div id="course-modal" class="hidden" style="position:fixed; top:0; left:0; width:100%; height:100%; background:rgba(0,0,0,0.8); z-index:1000; place-items:center; display:flex;">
            <div class="card" style="width:100%; max-width:500px; margin:auto;">
                <h3 class="mb-4">Yeni Kurs Oluştur</h3>
                <form id="create-course-form">
                    <div class="form-group"><label>Başlık</label><input type="text" id="new-course-title" class="form-control" required></div>
                    <div class="form-group"><label>Açıklama</label><input type="text" id="new-course-desc" class="form-control" required></div>
                    <div class="form-group"><label>Kategori</label>
                        <select id="new-course-cat" class="form-control" required>
                            <option value="Yazılım">Yazılım</option>
                            <option value="Tasarım">Tasarım</option>
                            <option value="İşletme">İşletme</option>
                            <option value="Kişisel Gelişim">Kişisel Gelişim</option>
                        </select>
                    </div>
                    <div class="flex gap-4" style="justify-content:flex-end;">
                        <button type="button" class="btn btn-outline" id="close-course-modal">İptal</button>
                        <button type="submit" class="btn btn-primary">Oluştur</button>
                    </div>
                </form>
            </div>
        </div>
    `;

    // Category Filter Logic
    const pills = document.querySelectorAll('.pill');
    pills.forEach(pill => {
        pill.addEventListener('click', () => {
            pills.forEach(p => p.classList.remove('active'));
            pill.classList.add('active');
            loadCourses(pill.dataset.category === 'all' ? {} : { category: pill.dataset.category });
        });
    });

    if(user?.role === 'teacher'){
        document.getElementById('btn-create-course').addEventListener('click', () => {
            document.getElementById('course-modal').classList.remove('hidden');
        });
        document.getElementById('close-course-modal').addEventListener('click', () => {
            document.getElementById('course-modal').classList.add('hidden');
        });
        document.getElementById('create-course-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            try {
                await CourseApi.create({
                    title: document.getElementById('new-course-title').value,
                    description: document.getElementById('new-course-desc').value,
                    category: document.getElementById('new-course-cat').value
                });
                showToast('Kurs oluşturuldu');
                document.getElementById('course-modal').classList.add('hidden');
                loadCourses();
            } catch(e) { showToast(e.message, 'danger') }
        });
    }

    loadStats();
    loadCourses();
}

async function loadStats() {
    const user = getUser();
    try {
        const coursesReq = await CourseApi.getAll();
        const courses = coursesReq.data || [];
        
        document.getElementById('stat-courses').innerText = courses.length;
        
        if (user.role === 'student') {
            const progressReq = await LessonApi.getMyProgress();
            const progress = progressReq || [];
            const avgProgress = progress.length > 0 ? 
                (progress.reduce((acc, curr) => acc + curr.percent, 0) / progress.length).toFixed(0) : 0;
            
            document.getElementById('stat-progress').innerText = `%${avgProgress}`;
            document.getElementById('stat-extra').innerText = '0';
        } else {
            document.getElementById('stat-progress').innerText = '-';
            document.getElementById('stat-extra').innerText = courses.filter(c => c.TeacherID === user.id).length;
        }
    } catch(e) { console.error(e); }
}

async function loadCourses(params = {}) {
    const list = document.getElementById('course-list');
    try {
        const req = await CourseApi.getAll(params);
        const courses = req.data || [];
        if(courses.length === 0){
            list.innerHTML = '<div class="text-muted text-center" style="grid-column: 1/-1; padding: 3rem;">Bu kategoride henüz kurs bulunmuyor.</div>';
            return;
        }
        list.innerHTML = courses.map(c => `
            <div class="card animate-fade-in">
                <div class="flex justify-between items-center mb-2">
                    <span class="badge badge-primary">${c.category}</span>
                </div>
                <h3 class="card-title">${c.title}</h3>
                <p class="card-desc">${c.description.length > 80 ? c.description.substring(0, 80) + '...' : c.description}</p>
                <div class="flex gap-4">
                     <a href="#course/${c.ID}" class="btn btn-primary" style="flex:1;">Eğitime Git</a>
                     <a href="#classroom/${c.ID}" class="btn btn-outline" title="Canlı Sınıf">💬</a>
                </div>
            </div>
        `).join('');
    } catch(e) {
        list.innerHTML = `<div class="text-danger text-center" style="grid-column: 1/-1;">${e.message}</div>`;
    }
}

async function renderCourseDetail(id) {
    const user = getUser();
    contentDiv.innerHTML = `<div class="text-center mt-8">Ders yükleniyor...</div>`;
    try {
        const course = await CourseApi.getById(id);
        const lessonsRes = await CourseApi.getLessons(id);
        const lessons = lessonsRes.data || [];
        
        contentDiv.innerHTML = `
            <div class="mb-8">
                <a href="#dashboard" class="text-muted mb-4" style="display:inline-block;">← Dashboarda Dön</a>
                <h1 class="text-gradient">${course.title}</h1>
                <p class="text-muted">${course.description}</p>
            </div>
            
            <div class="grid md:grid-cols-3">
                <div style="grid-column: span 2;">
                    ${user?.role === 'teacher' ? `
                        <div class="card mb-4" style="border: 1px dashed var(--primary); background: transparent;">
                            <h4>Yeni Ders Videosu/İçeriği Ekle</h4>
                            <form id="add-lesson-form" class="flex gap-4 mt-4 items-center">
                                <input type="text" id="lesson-title" class="form-control" placeholder="Ders Başlığı" required>
                                <input type="text" id="lesson-url" class="form-control" placeholder="Video URL" required>
                                <button type="submit" class="btn btn-primary">Ekle</button>
                            </form>
                        </div>
                    ` : ''}
                    
                    <h3>Ders İçerikleri</h3>
                    <div class="grid grid-cols-1 mt-4" id="lessons-container">
                        ${lessons.length === 0 ? '<p class="text-muted">Henüz ders yok.</p>' : 
                          lessons.map(l => `
                            <div class="card" style="flex-direction:row; justify-content:space-between; align-items:center;">
                                <div>
                                    <h4 style="margin:0">${l.title}</h4>
                                    <small class="text-muted">${l.content_url}</small>
                                </div>
                                <button class="btn btn-outline btn-complete-lesson" data-id="${l.ID}">Tamamla</button>
                            </div>
                          `).join('')
                        }
                    </div>
                </div>
                <div>
                     <div class="card text-center mb-4">
                         <div style="font-size:3rem; margin-bottom: 1rem;">💬</div>
                         <h3 class="mb-4">Canlı Sınıf</h3>
                         <p class="text-muted mb-4">Eğitmen ve öğrencilerle anlık sohbet edin.</p>
                         <a href="#classroom/${id}" class="btn btn-primary btn-block">Sınıfa Katıl</a>
                     </div>
                </div>
            </div>
        `;

        if(user?.role === 'teacher') {
            document.getElementById('add-lesson-form').addEventListener('submit', async(e) => {
                e.preventDefault();
                try {
                    await CourseApi.createLesson(id, {
                        title: document.getElementById('lesson-title').value,
                        content_url: document.getElementById('lesson-url').value,
                        order_index: lessons.length + 1
                    });
                    showToast('Ders eklendi.');
                    renderCourseDetail(id);
                } catch(e){ showToast(e.message, 'danger'); }
            });
        }

        const completeBtns = document.querySelectorAll('.btn-complete-lesson');
        completeBtns.forEach(btn => {
            btn.addEventListener('click', async(e) => {
                const lessonId = e.target.getAttribute('data-id');
                try {
                    await LessonApi.completeLesson(lessonId);
                    e.target.innerHTML = '✓ Tamamlandı';
                    e.target.classList.replace('btn-outline','btn-primary');
                    showToast('Ders tamamlandı işaretlendi!');
                } catch(err){ showToast(err.message, 'danger'); }
            })
        });

    } catch(err) {
        contentDiv.innerHTML = `<div class="text-danger text-center mt-8">Hata: ${err.message}</div>`;
    }
}

let ws = null;
function renderClassroom(courseId) {
    const user = getUser();
    contentDiv.innerHTML = `
        <div class="flex justify-between items-center mb-4 mt-4">
            <div>
                 <a href="#course/${courseId}" class="text-muted">← Kursa Dön</a>
                 <h2 class="text-gradient">Canlı Sınıf Sohbeti</h2>
            </div>
            <div id="ws-status" class="badge">Bağlanıyor...</div>
        </div>
        <div class="card">
             <div class="chat-window">
                  <div class="chat-messages" id="chat-messages">
                       <div class="text-center text-muted" style="margin-top:auto; margin-bottom:auto;">Sohbete katılınıyor...</div>
                  </div>
                  <div class="chat-input-area">
                       <input type="text" id="chat-input" class="form-control" placeholder="Bir mesaj yaz..." autocomplete="off">
                       <button id="chat-send" class="btn btn-primary">Gönder</button>
                  </div>
             </div>
        </div>
    `;

    const chatMessages = document.getElementById('chat-messages');
    const input = document.getElementById('chat-input');
    const sendBtn = document.getElementById('chat-send');
    const status = document.getElementById('ws-status');

    if(ws) ws.close();
    
    ws = new WebSocket(`ws://localhost:8090/ws/classroom/${courseId}?username=${user?.name || 'Anonim'}`);
    
    ws.onopen = () => {
        status.innerHTML = '● Bağlandı';
        status.style.color = '#10b981';
        status.style.background = 'rgba(16, 185, 129, 0.2)';
        chatMessages.innerHTML = '';
    };
    
    ws.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        const div = document.createElement('div');
        
        if (msg.type === 'join' || msg.type === 'leave') {
            div.className = 'text-center text-muted';
            div.style.fontSize = '0.8rem';
            div.innerText = msg.text;
        } else {
            const isMe = msg.username === user?.name;
            div.className = `message ${isMe ? 'sent' : 'received'}`;
            div.innerHTML = `
                <div class="message-meta">${msg.username}</div>
                <div>${msg.text}</div>
            `;
        }
        
        chatMessages.appendChild(div);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    };
    
    ws.onclose = () => {
        status.innerHTML = '○ Bağlantı Koptu';
        status.style.color = '#ef4444';
        status.style.background = 'rgba(239, 68, 68, 0.2)';
    };

    const sendMessage = () => {
        if(!input.value.trim() || ws.readyState !== WebSocket.OPEN) return;
        ws.send(JSON.stringify({ text: input.value, type: "message" }));
        input.value = '';
    };

    sendBtn.addEventListener('click', sendMessage);
    input.addEventListener('keypress', (e) => { if (e.key === 'Enter') sendMessage(); });
}
