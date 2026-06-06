<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Header -->
    <div class="bg-indigo-600 rounded-3xl p-8 text-white relative overflow-hidden shadow-2xl shadow-indigo-600/20">
      <div class="relative z-10 max-w-2xl">
        <h2 class="text-3xl font-black mb-3">Learning Center</h2>
        <p class="text-indigo-100 opacity-90 leading-relaxed">
          Master load testing with Prahara. From basic URL checks to complex user scenarios, 
          use these guides and templates to ensure your infrastructure can handle the storm.
        </p>
      </div>
      <ZapIcon class="absolute right-[-20px] bottom-[-20px] w-64 h-64 text-white/10 rotate-12" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Sidebar Navigation -->
      <div class="space-y-4">
        <div class="bg-slate-900 border border-slate-800 p-6 rounded-3xl sticky top-24">
          <h3 class="font-bold mb-4 flex items-center gap-2">
            <BookOpenIcon class="w-4 h-4 text-indigo-400" />
            Jump to Section
          </h3>
          <nav class="space-y-1">
            <button v-for="section in sections" :key="section.id"
              @click="scrollTo(section.id)"
              class="w-full text-left px-4 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 hover:bg-slate-800 text-slate-400 hover:text-white"
            >
              {{ section.name }}
            </button>
          </nav>
        </div>
      </div>

      <!-- Content Area -->
      <div class="lg:col-span-2 space-y-12 pb-20">
        <!-- User Guide -->
        <section id="guide" class="space-y-6">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center">
              <span class="font-bold text-sm">1</span>
            </div>
            <h3 class="text-2xl font-bold">Panduan Penggunaan (Tutorial)</h3>
          </div>
          
          <div class="space-y-4 bg-slate-900 border border-slate-800 p-8 rounded-3xl leading-relaxed">
            <div v-for="tip in guideSteps" :key="tip.title" class="group">
              <h4 class="font-bold text-white mb-2 flex items-center gap-2">
                <div class="w-1.5 h-1.5 rounded-full bg-indigo-500"></div>
                {{ tip.title }}
              </h4>
              <p class="text-slate-400 text-sm pl-4 border-l border-slate-800 group-hover:border-indigo-500/50 transition-colors">
                {{ tip.desc }}
              </p>
            </div>
          </div>
        </section>

        <!-- Script Templates -->
        <section id="templates" class="space-y-6">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center">
              <span class="font-bold text-sm">2</span>
            </div>
            <h3 class="text-2xl font-bold">k6 Script Templates</h3>
          </div>

          <div v-for="tpl in templates" :key="tpl.name" class="bg-slate-900 border border-slate-800 rounded-3xl overflow-hidden group hover:border-slate-700 transition-colors shadow-lg">
            <div class="p-6 border-b border-slate-800 flex justify-between items-center bg-slate-800/30">
              <div>
                <h4 class="font-bold text-indigo-300 group-hover:text-indigo-200">{{ tpl.name }}</h4>
                <p class="text-xs text-slate-500 mt-1">{{ tpl.desc }}</p>
              </div>
              <button @click="copyCode(tpl.code)" class="p-2 hover:bg-slate-700 rounded-lg transition-colors text-slate-400 hover:text-white" title="Copy Template">
                <CopyIcon class="w-4 h-4" />
              </button>
            </div>
            <div class="bg-slate-950 p-6">
              <pre class="text-xs font-mono text-indigo-200/80 leading-relaxed overflow-x-auto whitespace-pre"><code>{{ tpl.code }}</code></pre>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ZapIcon, BookOpenIcon, CopyIcon } from 'lucide-vue-next';

const sections = [
  { id: 'guide', name: 'User Guide' },
  { id: 'templates', name: 'Script Library' }
];

const guideSteps = [
  { title: 'Quick Storm', desc: 'Cara tercepat untuk mengetes satu URL langsung dari Dashboard. Masukkan URL, VUs, dan Durasi.' },
  { title: 'URL Registry', desc: 'Daftarkan endpoint infra Anda di sini. Kategorikan (API, Frontend, dll) agar laporan di Dashboard lebih rapi.' },
  { title: 'Script Editor', desc: 'Gunakan Script Editor untuk skenario yang lebih kompleks seperti simulasi Login atau Checkout flow.' },
  { title: 'Storm Report', desc: 'Klik ikon "Detail" di tabel History untuk melihat statistik mendalam (Avg, Max, Requests) tiap run.' }
];

const templates = [
  {
    name: 'Basic Load Test',
    desc: 'Cek kestabilan dengan jumlah user konstan.',
    code: `import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  http.get('https://example.com');
  sleep(1);
}`
  },
  {
    name: 'Ramping Up (Stress Test)',
    desc: 'Naikkan jumlah user secara bertahap untuk mencari batas server.',
    code: `import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m', target: 20 },
    { duration: '30s', target: 0 },
  ],
};

export default function () {
  http.get('https://example.com');
  sleep(1);
}`
  },
  {
    name: 'API POST Payload',
    desc: 'Simulasi pengiriman data JSON ke API.',
    code: `import http from 'k6/http';

export default function () {
  const url = 'https://example.com/api/login';
  const payload = JSON.stringify({
    user: 'test', pass: 'secret'
  });
  const params = {
    headers: { 'Content-Type': 'application/json' }
  };
  http.post(url, payload, params);
}`
  }
];

const scrollTo = (id) => {
  const el = document.getElementById(id);
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
};

const copyCode = (text) => {
  navigator.clipboard.writeText(text);
};
</script>
