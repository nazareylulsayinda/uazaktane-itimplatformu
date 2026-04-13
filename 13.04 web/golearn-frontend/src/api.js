const BASE_URL = 'http://localhost:8090/api';

function getToken() {
    return localStorage.getItem('token');
}

export async function fetchApi(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Content-Type': 'application/json',
        ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
        ...options.headers
    };

    const response = await fetch(`${BASE_URL}${endpoint}`, {
        ...options,
        headers
    });

    // Handle 401 Unauthorized globally, BUT skip for the login endpoint
    if (response.status === 401 && !endpoint.includes('/auth/login')) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.hash = '#login';
        throw new Error("Session expired, please login again.");
    }

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP Error ${response.status}`);
    }

    return response.json();
}

export const AuthApi = {
    login: (email, password) => fetchApi('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    }),
    register: (name, email, password, role) => fetchApi('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ name, email, password, role })
    })
};

export const CourseApi = {
    getAll: (params = {}) => {
        const query = new URLSearchParams(params).toString();
        return fetchApi(`/courses?${query}`);
    },
    getById: (id) => fetchApi(`/courses/${id}`),
    create: (data) => fetchApi('/courses', { method: 'POST', body: JSON.stringify(data) }),
    getLessons: (courseId) => fetchApi(`/courses/${courseId}/lessons`),
    createLesson: (courseId, data) => fetchApi(`/courses/${courseId}/lessons`, { method: 'POST', body: JSON.stringify(data) })
};

export const LessonApi = {
    getQuiz: (lessonId) => fetchApi(`/lessons/${lessonId}/quiz`),
    createQuiz: (lessonId, data) => fetchApi(`/lessons/${lessonId}/quiz`, { method: 'POST', body: JSON.stringify(data) }),
    submitQuiz: (quizId, score) => fetchApi(`/quiz/${quizId}/submit`, { method: 'POST', body: JSON.stringify({ score }) }),
    completeLesson: (lessonId) => fetchApi(`/lessons/${lessonId}/complete`, { method: 'POST' }),
    getMyProgress: () => fetchApi('/my/progress')
};
