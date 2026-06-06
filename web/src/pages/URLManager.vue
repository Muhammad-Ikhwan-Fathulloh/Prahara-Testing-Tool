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

    <!-- URL Grid / Table Container -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <table v-if="filteredUrls.length > 0" class="w-full text-left">
        <thead class="bg-slate-950/50 text-slate-500 text-xs uppercase font-bold">
          <tr>
            <th class="px-6 py-4">Endpoint</th>
            <th class="px-6 py-4 text-center">Category</th>
            <th class="px-6 py-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="url in filteredUrls" :key="url.id" class="hover:bg-slate-800/50 transition-colors">
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center flex-shrink-0">
                  <LinkIcon class="w-4 h-4 text-indigo-400" />
                </div>
                <div class="min-w-0">
                  <p class="font-medium truncate">{{ url.name }}</p>
                  <p class="text-xs text-slate-500 font-mono truncate">{{ url.url }}</p>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-center">
              <span :class="['px-2 py-1 rounded text-[10px] font-bold border whitespace-nowrap', getCategoryColor(url.category)]">
                {{ url.category }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <div class="flex justify-end gap-2">
                <button @click="openTestModal(url)" class="p-2 text-emerald-400 hover:bg-emerald-500/10 rounded-lg transition-colors" title="Quick Test">
                  <PlayIcon class="w-4 h-4" />
                </button>
                <button @click="copyToClipboard(url.url)" class="p-2 text-slate-400 hover:bg-slate-800 rounded-lg transition-colors">
                  <CopyIcon class="w-4 h-4" />
                </button>
                <button @click="deleteUrl(url.id)" class="p-2 text-rose-400 hover:bg-rose-500/10 rounded-lg transition-colors">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- Empty State -->
      <div v-if="filteredUrls.length === 0" class="py-20 flex flex-col items-center justify-center text-slate-600">
        <LinkIcon class="w-12 h-12 mb-4 opacity-20" />
        <p class="text-lg font-medium">No URLs found in this category</p>
      </div>
    </div>

    <!-- Add URL Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-md overflow-hidden shadow-2xl animate-in zoom-in-95 duration-200">
        <div class="p-6 border-b border-slate-800 bg-gradient-to-r from-indigo-600/10 to-transparent">
          <h3 class="text-xl font-bold">Register New Endpoint</h3>
          <p class="text-slate-500 text-sm">Add a new URL to your testing registry</p>
        </div>
        
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Friendly Name</label>
            <input v-model="newUrl.name" type="text" placeholder="e.g. Auth Service" class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all font-medium text-white" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Target URL</label>
            <input v-model="newUrl.url" type="url" placeholder="https://api.example.com" class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all font-mono text-white text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Category</label>
            <select v-model="newUrl.category" class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 focus:ring-2 focus:ring-indigo-500 transition-all font-medium text-white">
              <option v-for="cat in categories.filter(c => c !== 'ALL')" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>
        </div>

        <div class="p-6 bg-slate-800/30 border-t border-slate-800 flex gap-3">
          <button @click="showAddModal = false" class="flex-1 px-4 py-2 bg-slate-800 rounded-xl font-medium hover:bg-slate-700 transition-colors">
            Cancel
          </button>
          <button @click="saveUrl" :disabled="isSaving" class="flex-2 px-6 py-2 bg-indigo-600 rounded-xl font-bold hover:bg-indigo-500 transition-colors flex items-center justify-center gap-2 shadow-lg shadow-indigo-600/20">
            <Loader2Icon v-if="isSaving" class="w-4 h-4 animate-spin" />
            {{ isSaving ? 'Saving...' : 'Register URL' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Quick Test Modal -->
    <div v-if="showTestModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-md overflow-hidden shadow-2xl animate-in zoom-in-95 duration-200">
        <div class="p-6 border-b border-slate-800 bg-gradient-to-r from-emerald-600/10 to-transparent">
          <div class="flex items-center gap-3 mb-1">
            <ZapIcon class="w-5 h-5 text-emerald-400" />
            <h3 class="text-xl font-bold">Quick Storm Test</h3>
          </div>
          <p class="text-slate-500 text-sm">Launch an instant load test on this endpoint</p>
        </div>
        
        <div class="p-6 space-y-4">
          <div class="bg-slate-800/50 p-4 rounded-xl border border-slate-700">
            <p class="text-xs text-slate-500 uppercase font-bold mb-1">Target Endpoint</p>
            <p class="font-medium text-emerald-400 truncate">{{ selectedUrl?.url }}</p>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-1 flex items-center gap-2">
                <UsersIcon class="w-3.5 h-3.5" /> Virtual Users
              </label>
              <input v-model.number="testConfig.vus" type="number" class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 focus:ring-2 focus:ring-emerald-500 transition-all font-medium text-white" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-1 flex items-center gap-2">
                <TimerIcon class="w-3.5 h-3.5" /> Duration
              </label>
              <input v-model="testConfig.duration" type="text" placeholder="30s" class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 focus:ring-2 focus:ring-emerald-500 transition-all font-medium text-white" />
            </div>
          </div>
        </div>

        <div class="p-6 bg-slate-800/30 border-t border-slate-800 flex gap-3">
          <button @click="showTestModal = false" class="flex-1 px-4 py-2 bg-slate-800 rounded-xl font-medium hover:bg-slate-700 transition-colors">
            Cancel
          </button>
          <button @click="runUrlTest" :disabled="isTesting" class="flex-2 px-6 py-2 bg-emerald-600 rounded-xl font-bold hover:bg-emerald-500 transition-colors flex items-center justify-center gap-2 shadow-lg shadow-emerald-600/20">
            <PlayIcon v-if="!isTesting" class="w-4 h-4" />
            <Loader2Icon v-else class="w-4 h-4 animate-spin" />
            {{ isTesting ? 'Launching...' : 'Ignite Storm' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { 
  PlusIcon, TrashIcon, CopyIcon, LinkIcon, 
  Loader2Icon, PlayIcon, ZapIcon, TimerIcon, UsersIcon 
} from 'lucide-vue-next';
import { getUrls, createUrl, deleteUrl as removeUrl, runDynamicTest } from '../services/api';

const emit = defineEmits(['test-started']);
const filter = ref('ALL');
const showAddModal = ref(false);
const showTestModal = ref(false);
const isSaving = ref(false);
const isLoading = ref(false);
const isTesting = ref(false);
const categories = ['ALL', 'FRONTEND', 'BACKEND', 'API', 'STAGING', 'PRODUCTION'];

const urls = ref([]);
const newUrl = ref({ name: '', url: '', category: 'API' });
const selectedUrl = ref(null);
const testConfig = reactive({
  vus: 10,
  duration: '30s',
  method: 'GET'
});

const fetchUrls = async () => {
  isLoading.value = true;
  try {
    const response = await getUrls();
    urls.value = response.data || [];
  } catch (err) {
    console.error('Failed to fetch URLs:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchUrls);

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

const saveUrl = async () => {
  if (!newUrl.value.name || !newUrl.value.url) return;
  isSaving.value = true;
  try {
    await createUrl(newUrl.value);
    await fetchUrls();
    newUrl.value = { name: '', url: '', category: 'API' };
    showAddModal.value = false;
  } catch (err) {
    console.error('Failed to save URL:', err);
  } finally {
    isSaving.value = false;
  }
};

const deleteUrl = async (id) => {
  if (!confirm('Are you sure you want to remove this endpoint?')) return;
  try {
    await removeUrl(id);
    await fetchUrls();
  } catch (err) {
    console.error('Failed to delete URL:', err);
  }
};

const openTestModal = (url) => {
  selectedUrl.value = url;
  showTestModal.value = true;
};

const runUrlTest = async () => {
  if (!selectedUrl.value) return;
  isTesting.value = true;
  try {
    await runDynamicTest({
      url: selectedUrl.value.url,
      ...testConfig
    });
    showTestModal.value = false;
    emit('test-started');
  } catch (err) {
    alert('Failed to launch test: ' + (err.response?.data?.error || err.message));
  } finally {
    isTesting.value = false;
  }
};

const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text);
};
</script>
