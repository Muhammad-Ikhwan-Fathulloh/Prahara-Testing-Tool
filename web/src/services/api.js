import axios from 'axios';

const API_BASE = '/api';

const api = axios.create({
    baseURL: API_BASE,
    headers: {
        'Content-Type': 'application/json'
    }
});

// Add token to requests if available
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('prahara_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Auth
export const login = (username, password) => api.post('/auth/login', { username, password });
export const register = (username, password, role) => api.post('/auth/register', { username, password, role });

// Scripts
export const getScripts = () => api.get('/scripts');
export const createScript = (data) => api.post('/scripts', data);
export const runScript = (id) => api.post(`/scripts/${id}/run`);

// URLs
export const getUrls = (category) => api.get('/urls', { params: category ? { category } : {} });
export const createUrl = (data) => api.post('/urls', data);
export const deleteUrl = (id) => api.delete(`/urls/${id}`);

// Dynamic Test (Quick Storm)
export const runDynamicTest = (data) => api.post('/run-dynamic', data);

// Metrics
export const getMetrics = (range) => api.get('/metrics', { params: { range } });

// Test Runs
export const getRuns = () => api.get('/runs');

export default api;
