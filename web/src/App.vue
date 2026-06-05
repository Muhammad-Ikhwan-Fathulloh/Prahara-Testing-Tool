<template>
  <div class="min-h-screen flex flex-col bg-slate-950 text-slate-50">
    <!-- Sidebar / Nav -->
    <header class="h-16 border-b border-slate-800 flex items-center justify-between px-8 bg-slate-900/50 backdrop-blur-md sticky top-0 z-50">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 bg-indigo-600 rounded-lg flex items-center justify-center shadow-lg shadow-indigo-500/20">
          <ZapIcon class="w-6 h-6 text-white" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-white uppercase italic">Prahara</h1>
      </div>
      
      <nav v-if="user" class="flex gap-1 bg-slate-800/50 p-1 rounded-lg border border-slate-700">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          @click="currentTab = tab.id"
          :class="[
            'px-4 py-1.5 rounded-md text-sm font-medium transition-all duration-200',
            currentTab === tab.id ? 'bg-indigo-600 text-white shadow-md' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-700/50'
          ]"
        >
          {{ tab.name }}
        </button>
      </nav>

      <div class="flex items-center gap-4">
        <template v-if="user">
          <div class="flex flex-col items-end">
            <span class="text-sm font-medium text-slate-200">{{ user.username }}</span>
            <span class="text-xs text-slate-500 uppercase tracking-widest">{{ user.role }}</span>
          </div>
          <button @click="logout" class="p-2 hover:bg-slate-800 rounded-full transition-colors">
            <LogOutIcon class="w-5 h-5 text-slate-400" />
          </button>
        </template>
        <button v-else @click="showLogin = true" class="px-4 py-2 bg-indigo-600 rounded-lg font-medium hover:bg-indigo-500 transition-colors">
          Login
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 p-8 max-w-7xl mx-auto w-full">
      <div v-if="!user" class="flex flex-col items-center justify-center h-full py-20">
        <div class="text-center max-w-2xl px-4">
          <h2 class="text-5xl font-extrabold mb-6 bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
            Unleash the Storm on Your APIs
          </h2>
          <p class="text-slate-400 text-lg mb-8">
            Prahara is a lightweight, secure, and powerful testing suite for load and QA testing. 
            Built on k6 and InfluxDB for professional-grade performance monitoring.
          </p>
          <div class="flex gap-4 justify-center">
            <button @click="currentTab = 'dashboard'" class="px-8 py-3 bg-indigo-600 rounded-xl font-bold shadow-lg shadow-indigo-500/25 hover:scale-105 transition-transform">
              Get Started
            </button>
            <button class="px-8 py-3 bg-slate-800 rounded-xl font-bold border border-slate-700 hover:bg-slate-700 transition-colors">
              Documentation
            </button>
          </div>
        </div>
      </div>

      <div v-else>
        <Dashboard v-if="currentTab === 'dashboard'" />
        <URLManager v-if="currentTab === 'urls'" />
        <ScriptEditor v-if="currentTab === 'editor'" />
        <Settings v-if="currentTab === 'settings'" />
      </div>
    </main>

    <!-- Modal for Login -->
    <div v-if="showLogin" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md p-8 shadow-2xl">
        <h3 class="text-2xl font-bold mb-6">Login to Prahara</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Username</label>
            <input v-model="loginForm.username" type="text" class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Password</label>
            <input v-model="loginForm.password" type="password" class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all" />
          </div>
          <button @click="handleLogin" class="w-full bg-indigo-600 py-3 rounded-lg font-bold mt-4 hover:bg-indigo-500 transition-colors">
            Authorize
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { ZapIcon, LogOutIcon, BarChart3Icon, CodeIcon, SettingsIcon } from 'lucide-vue-next';
import Dashboard from './pages/Dashboard.vue';
import URLManager from './pages/URLManager.vue';
import ScriptEditor from './pages/ScriptEditor.vue';
import Settings from './pages/Settings.vue';

const user = ref(JSON.parse(localStorage.getItem('prahara_user')));
const currentTab = ref('dashboard');
const showLogin = ref(false);

const tabs = [
  { id: 'dashboard', name: 'Dashboard' },
  { id: 'urls', name: 'URL Registry' },
  { id: 'editor', name: 'Test Editor' },
  { id: 'settings', name: 'Configuration' }
];

const loginForm = reactive({ username: 'admin', password: 'password' });

const handleLogin = () => {
  // Mock login for now, will connect to API later
  const mockUser = { username: loginForm.username, role: 'admin' };
  user.value = mockUser;
  localStorage.setItem('prahara_user', JSON.stringify(mockUser));
  showLogin.value = false;
};

const logout = () => {
  user.value = null;
  localStorage.removeItem('prahara_user');
};
</script>

<style>
@tailwind base;
@tailwind components;
@tailwind utilities;
</style>
