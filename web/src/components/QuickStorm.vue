<template>
  <div class="bg-indigo-600 rounded-3xl p-8 text-white relative overflow-hidden shadow-2xl shadow-indigo-600/20 group">
    <!-- Background Decor -->
    <div class="absolute -right-20 -top-20 w-64 h-64 bg-white/10 rounded-full blur-3xl group-hover:scale-110 transition-transform duration-700"></div>
    <div class="absolute -left-20 -bottom-20 w-48 h-48 bg-indigo-400/20 rounded-full blur-2xl group-hover:scale-110 transition-transform duration-700"></div>

    <div class="relative z-10 flex flex-col lg:flex-row gap-8 items-center">
      <div class="lg:flex-1 space-y-4">
        <div class="flex items-center gap-3">
          <div class="p-2 bg-white/20 rounded-lg backdrop-blur-md">
            <ZapIcon class="w-6 h-6" />
          </div>
          <h2 class="text-2xl font-black uppercase italic tracking-wider">Quick Storm</h2>
        </div>
        <p class="text-indigo-100 font-medium leading-relaxed">
          Need a quick load check? Launch a dynamic storm on any URL instantly without writing a script.
        </p>
      </div>

      <div class="w-full lg:w-3/5 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="md:col-span-2 relative">
          <input 
            v-model="config.url" 
            type="text" 
            placeholder="https://api.example.com/check" 
            class="w-full bg-white/10 border border-white/20 rounded-2xl px-6 py-4 placeholder:text-white/40 focus:bg-white/20 focus:ring-2 focus:ring-white/50 transition-all outline-none font-mono text-sm"
          />
        </div>
        
        <div class="flex bg-white/10 border border-white/20 rounded-2xl overflow-hidden">
          <select v-model="config.method" class="w-full bg-transparent px-4 py-4 outline-none appearance-none font-bold text-center cursor-pointer hover:bg-white/5 transition-colors">
            <option class="text-slate-900" v-for="m in methods" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="bg-white/10 border border-white/20 rounded-2xl p-2 flex flex-col items-center justify-center">
            <span class="text-[10px] uppercase font-bold text-indigo-200">VUs</span>
            <input v-model.number="config.vus" type="number" class="w-full bg-transparent text-center font-black text-lg outline-none" />
          </div>
          <div class="bg-white/10 border border-white/20 rounded-2xl p-2 flex flex-col items-center justify-center">
            <span class="text-[10px] uppercase font-bold text-indigo-200">Duration</span>
            <input v-model="config.duration" type="text" class="w-full bg-transparent text-center font-black text-lg outline-none" />
          </div>
        </div>
        
        <button 
          @click="launchStorm" 
          :disabled="isLaunching"
          class="md:col-span-2 bg-white text-indigo-600 py-4 rounded-2xl font-black uppercase tracking-widest hover:bg-indigo-50 hover:scale-[1.02] active:scale-95 transition-all shadow-xl shadow-black/10 flex items-center justify-center gap-3 disabled:opacity-50 disabled:scale-100"
        >
          <template v-if="isLaunching">
            <ActivityIcon class="w-5 h-5 animate-spin" /> Igniting...
          </template>
          <template v-else>
            <PlayIcon class="w-5 h-5" /> Launch Dynamic Storm
          </template>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { ZapIcon, PlayIcon, ActivityIcon } from 'lucide-vue-next';
import { runDynamicTest } from '../services/api';

const emit = defineEmits(['test-started']);
const isLaunching = ref(false);
const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'];

const config = reactive({
  url: '',
  method: 'GET',
  vus: 10,
  duration: '30s'
});

const launchStorm = async () => {
  if (!config.url) return;
  isLaunching.value = true;
  
  try {
    await runDynamicTest(config);
    emit('test-started');
    config.url = '';
  } catch (err) {
    alert('Failed to ignite storm: ' + (err.response?.data?.error || err.message));
  } finally {
    isLaunching.value = false;
  }
};
</script>
