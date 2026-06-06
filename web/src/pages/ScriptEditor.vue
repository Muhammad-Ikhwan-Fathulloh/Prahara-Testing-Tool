<template>
  <div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <div class="flex justify-between items-center bg-slate-900 p-6 rounded-2xl border border-slate-800">
      <div>
        <h2 class="text-2xl font-bold">Storm Architect</h2>
        <p class="text-slate-500 text-sm">Design your k6 load testing scenarios</p>
      </div>
      <div class="flex gap-3">
        <button @click="saveScript" :disabled="isSaving" class="px-4 py-2 bg-slate-800 rounded-lg font-medium hover:bg-slate-700 transition-colors flex items-center gap-2">
          <SaveIcon v-if="!isSaving" class="w-4 h-4" />
          <Loader2Icon v-else class="w-4 h-4 animate-spin" />
          {{ isSaving ? 'Saving...' : 'Save Script' }}
        </button>
        <button @click="executeTest" :disabled="isRunning" class="px-6 py-2 bg-indigo-600 rounded-lg font-bold hover:bg-indigo-500 transition-colors flex items-center gap-2 shadow-lg shadow-indigo-600/20">
          <PlayIcon v-if="!isRunning" class="w-4 h-4" />
          <Loader2Icon v-else class="w-4 h-4 animate-spin" />
          {{ isRunning ? 'Storming...' : 'Launch Storm' }}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-4">
        <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
          <div class="bg-slate-800 px-4 py-2 border-b border-slate-700 flex items-center justify-between">
            <span class="text-xs font-mono text-slate-400">script.js</span>
            <div class="flex gap-2">
              <span class="w-3 h-3 rounded-full bg-rose-500/50"></span>
              <span class="w-3 h-3 rounded-full bg-amber-500/50"></span>
              <span class="w-3 h-3 rounded-full bg-emerald-500/50"></span>
            </div>
          </div>
          <textarea 
            v-model="scriptContent" 
            class="w-full h-[500px] bg-slate-950 text-indigo-300 font-mono p-6 outline-none resize-none spellcheck-false"
            spellcheck="false"
          ></textarea>
        </div>
      </div>

      <div class="space-y-6">
        <div class="bg-slate-900 border border-slate-800 p-6 rounded-2xl">
          <h3 class="font-bold mb-4 flex items-center gap-2 text-indigo-400">
            <InfoIcon class="w-4 h-4" /> Quick Info
          </h3>
          <div class="space-y-4 text-sm text-slate-400">
            <p>Prahara uses standard <strong>k6</strong> syntax. You can use JS modules and k6/http tools.</p>
            <div class="bg-indigo-500/5 p-4 rounded-xl border border-indigo-500/10 text-xs leading-relaxed">
              <p class="font-bold text-indigo-300 mb-1">PRO TIP:</p>
              Use <code>options.thresholds</code> to define QA pass/fail criteria for your CI/CD pipelines.
            </div>
          </div>
        </div>

        <div class="bg-slate-900 border border-slate-800 p-6 rounded-2xl">
          <h3 class="font-bold mb-4 flex items-center gap-2 text-indigo-400">
            <ZapIcon class="w-4 h-4" /> Templates
          </h3>
          <div class="space-y-2">
            <button v-for="t in templates" :key="t.name" @click="scriptContent = t.content" class="w-full text-left p-3 rounded-xl hover:bg-slate-800 transition-colors border border-transparent hover:border-slate-700">
              <p class="font-bold text-slate-200">{{ t.name }}</p>
              <p class="text-xs text-slate-500">{{ t.desc }}</p>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { SaveIcon, PlayIcon, InfoIcon, ZapIcon, Loader2Icon } from 'lucide-vue-next';
import { getScripts, createScript, runScript } from '../services/api';

const emit = defineEmits(['test-started']);
const isRunning = ref(false);
const isSaving = ref(false);
const scripts = ref([]);
const currentScriptId = ref(null);
const scriptContent = ref('');
const scriptName = ref('New Performance Script');

const fetchScripts = async () => {
  try {
    const response = await getScripts();
    scripts.value = response.data || [];
    if (scripts.value.length > 0 && !currentScriptId.value) {
      loadScript(scripts.value[0]);
    }
  } catch (err) {
    console.error('Failed to fetch scripts:', err);
  }
};

const loadScript = (s) => {
  currentScriptId.value = s.id;
  scriptContent.value = s.content;
  scriptName.value = s.name;
};

const saveScript = async () => {
  isSaving.value = true;
  try {
    await createScript({
      name: scriptName.value,
      content: scriptContent.value
    });
    await fetchScripts();
    alert('Script saved successfully!');
  } catch (err) {
    alert('Failed to save script: ' + (err.response?.data?.error || err.message));
  } finally {
    isSaving.value = false;
  }
};

const executeTest = async () => {
  if (!currentScriptId.value) {
    alert('Please save the script first!');
    return;
  }
  
  isRunning.value = true;
  try {
    const response = await runScript(currentScriptId.value);
    emit('test-started');
  } catch (err) {
    alert('Storm failed to ignite: ' + (err.response?.data?.error || err.message));
  } finally {
    isRunning.value = false;
  }
};

const templates = [
  { 
    name: 'Get Request', 
    desc: 'Simple load test for a single endpoint', 
    content: `import http from 'k6/http';\nimport { sleep } from 'k6';\n\nexport default function() {\n  http.get('https://test.k6.io');\n  sleep(1);\n}` 
  },
  { 
    name: 'Stress Test', 
    desc: 'Ramp-up VUs to find breaking point', 
    content: `import http from 'k6/http';\n\nexport const options = {\n  stages: [\n    { duration: '2m', target: 100 },\n    { duration: '5m', target: 100 },\n    { duration: '2m', target: 0 },\n  ],\n};\n\nexport default function() {\n  http.get('https://test.k6.io');\n}` 
  }
];

onMounted(fetchScripts);
</script>
