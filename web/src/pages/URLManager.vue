<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Header -->
    <div class="flex justify-between items-center bg-slate-900 p-6 rounded-2xl border border-slate-800">
      <div>
        <h2 class="text-2xl font-bold">URL Registry</h2>
        <p class="text-slate-500 text-sm">Manage categorized endpoints for testing</p>
      </div>
      <button @click="showAddModal = true" class="px-4 py-2 bg-indigo-600 rounded-lg font-bold hover:bg-indigo-500 transition-colors flex items-center gap-2 shadow-lg shadow-indigo-600/20">
        <PlusIcon class="w-4 h-4" /> Add New URL
      </button>
    </div>

    <!-- Category Filters -->
    <div class="flex gap-2 bg-slate-900/50 p-1 rounded-xl border border-slate-800 w-fit">
      <button 
        v-for="cat in categories" 
        :key="cat"
        @click="filter = cat"
        :class="[
          'px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-widest transition-all',
          filter === cat ? 'bg-indigo-600 text-white' : 'text-slate-500 hover:text-slate-300'
        ]"
      >
        {{ cat }}
      </button>
    </div>

    <!-- URL Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="item in filteredUrls" :key="item.id" class="bg-slate-900 border border-slate-800 p-6 rounded-2xl group hover:border-indigo-500/50 transition-all shadow-xl">
        <div class="flex justify-between items-start mb-4">
          <span :class="[
            'px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-tighter border',
            getCategoryColor(item.category)
          ]">
            {{ item.category }}
          </span>
          <button @click="deleteUrl(item.id)" class="text-slate-600 hover:text-rose-400 transition-colors opacity-0 group-hover:opacity-100">
            <TrashIcon class="w-4 h-4" />
          </button>
        </div>
        <h3 class="font-bold text-lg text-white mb-1">{{ item.name }}</h3>
        <p class="text-xs font-mono text-slate-500 break-all mb-4">{{ item.url }}</p>
        <div class="flex gap-2">
          <button @click="copyToClipboard(item.url)" class="flex-1 py-2 bg-slate-800 rounded-lg text-xs font-bold hover:bg-slate-700 transition-colors flex items-center justify-center gap-2">
            <CopyIcon class="w-3.5 h-3.5" /> Copy URL
          </button>
          <a :href="item.url" target="_blank" class="p-2 bg-slate-800 rounded-lg hover:bg-slate-700 transition-colors">
            <ExternalLinkIcon class="w-4 h-4 text-indigo-400" />
          </a>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredUrls.length === 0" class="py-20 flex flex-col items-center justify-center border-2 border-dashed border-slate-800 rounded-3xl text-slate-600">
      <LinkIcon class="w-12 h-12 mb-4 opacity-20" />
      <p class="text-lg font-medium">No URLs found in this category</p>
    </div>

    <!-- Add Modal -->
    <div v-if="showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md p-8 shadow-2xl animate-in zoom-in-95 duration-300">
        <h3 class="text-2xl font-bold mb-6">Register New Endpoint</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Friendly Name</label>
            <input v-model="newUrl.name" type="text" placeholder="e.g. Auth Service Production" class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Target URL</label>
            <input v-model="newUrl.url" type="text" placeholder="https://api.example.com/v1" class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all font-mono text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Category</label>
            <select v-model="newUrl.category" class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all">
              <option v-for="cat in categories.slice(1)" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>
          <div class="flex gap-3 mt-6">
            <button @click="showAddModal = false" class="flex-1 py-3 bg-slate-800 rounded-lg font-bold hover:bg-slate-700 transition-colors">Cancel</button>
            <button @click="saveUrl" class="flex-1 py-3 bg-indigo-600 rounded-lg font-bold hover:bg-indigo-500 transition-colors">Save Endpoint</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { PlusIcon, TrashIcon, CopyIcon, ExternalLinkIcon, LinkIcon } from 'lucide-vue-next';

const filter = ref('ALL');
const showAddModal = ref(false);
const categories = ['ALL', 'FRONTEND', 'BACKEND', 'API', 'STAGING', 'PRODUCTION'];

const urls = ref([
  { id: 1, name: 'Main Frontend', url: 'https://prahara.example.com', category: 'FRONTEND' },
  { id: 2, name: 'Auth Microservice', url: 'https://auth.api.example.com', category: 'BACKEND' },
  { id: 3, name: 'Payment Gateway', url: 'https://pay.api.example.com', category: 'API' }
]);

const newUrl = ref({ name: '', url: '', category: 'FRONTEND' });

const filteredUrls = computed(() => {
  if (filter.value === 'ALL') return urls.value;
  return urls.value.filter(u => u.category === filter.value);
});

const getCategoryColor = (cat) => {
  const colors = {
    FRONTEND: 'bg-indigo-500/10 text-indigo-400 border-indigo-500/20',
    BACKEND: 'bg-purple-500/10 text-purple-400 border-purple-500/20',
    API: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
    STAGING: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
    PRODUCTION: 'bg-rose-500/10 text-rose-400 border-rose-500/20'
  };
  return colors[cat] || 'bg-slate-500/10 text-slate-400 border-slate-500/20';
};

const saveUrl = () => {
  if (!newUrl.value.name || !newUrl.value.url) return;
  urls.value.push({ ...newUrl.value, id: Date.now() });
  newUrl.value = { name: '', url: '', category: 'FRONTEND' };
  showAddModal.value = false;
};

const deleteUrl = (id) => {
  urls.value = urls.value.filter(u => u.id !== id);
};

const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text);
  // Ideally show a toast here
};
</script>
