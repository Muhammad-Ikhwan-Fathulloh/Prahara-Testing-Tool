<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Quick Storm Widget -->
    <QuickStorm @test-started="fetchData" />

    <!-- Header & Filter -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-slate-900 border border-slate-800 p-6 rounded-3xl shadow-lg">
      <div>
        <h2 class="text-2xl font-black text-white flex items-center gap-2">
          <ZapIcon class="w-6 h-6 text-indigo-400" />
          Analytics Dashboard
        </h2>
        <p class="text-slate-500 text-sm">Real-time load performance across your infra</p>
      </div>
      <div class="flex flex-wrap gap-2 bg-slate-800/50 p-1 rounded-2xl border border-slate-700">
        <button v-for="cat in categories" :key="cat"
          @click="selectedCategory = cat"
          :class="[
            'px-4 py-1.5 rounded-xl text-xs font-bold transition-all duration-300',
            selectedCategory === cat ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-slate-200'
          ]"
        >
          {{ cat }}
        </button>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="stat in stats" :key="stat.label" class="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl flex items-center gap-4">
        <div :class="[stat.color, 'w-12 h-12 rounded-xl flex items-center justify-center']">
          <component :is="stat.icon" class="w-6 h-6 text-white" />
        </div>
        <div>
          <p class="text-sm text-slate-500 font-medium uppercase tracking-wider">{{ stat.label }}</p>
          <p class="text-2xl font-bold">{{ stat.value }}</p>
        </div>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-2 bg-slate-900 border border-slate-800 p-8 rounded-3xl shadow-xl flex flex-col h-[500px]">
        <div class="flex items-center justify-between mb-8 flex-shrink-0">
          <h3 class="text-xl font-bold">Request Performance ({{ selectedCategory }})</h3>
          <div class="flex gap-2">
            <span class="flex items-center gap-1.5 text-xs text-indigo-400 font-medium px-2 py-1 bg-indigo-400/10 rounded-md">
              <span class="w-2 h-2 rounded-full bg-indigo-400"></span> Avg. Latency
            </span>
          </div>
        </div>
        <div class="flex-1 relative min-h-0">
          <canvas ref="lineChartCanvas"></canvas>
        </div>
      </div>

      <div class="bg-slate-900 border border-slate-800 p-8 rounded-3xl shadow-xl flex flex-col h-[500px]">
        <h3 class="text-xl font-bold mb-8 flex-shrink-0">Success Rate</h3>
        <div class="flex-1 flex items-center justify-center relative min-h-0">
          <canvas ref="doughnutChartCanvas"></canvas>
          <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
            <span class="text-3xl font-black text-white">{{ successRateDisplay }}%</span>
            <span class="text-xs text-slate-500 uppercase font-bold tracking-tighter">Healthy</span>
          </div>
        </div>
        <div class="mt-8 space-y-3 flex-shrink-0">
          <div class="flex justify-between text-sm">
            <span class="text-slate-400">Completed</span>
            <span class="font-bold text-emerald-400">{{ completedCount }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-400">Running</span>
            <span class="font-bold text-amber-400">{{ runningCount }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-400">Failed</span>
            <span class="font-bold text-rose-400">{{ failedCount }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Run Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-3xl overflow-hidden shadow-xl">
      <div class="p-6 border-b border-slate-800 flex justify-between items-center">
        <h3 class="font-bold text-lg">Recent Storms</h3>
        <span class="text-sm text-slate-500">{{ filteredRuns.length }} runs</span>
      </div>
      <table v-if="filteredRuns.length > 0" class="w-full text-left">
        <thead class="bg-slate-800/50 text-slate-400 text-xs uppercase tracking-widest">
          <tr>
            <th class="px-6 py-4 font-semibold">Test Name</th>
            <th class="px-6 py-4 font-semibold">Status</th>
            <th class="px-6 py-4 font-semibold">VUs</th>
            <th class="px-6 py-4 font-semibold">Duration</th>
            <th class="px-6 py-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="run in filteredRuns" :key="run.id" class="hover:bg-slate-800/30 transition-colors group">
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <span :class="['px-1.5 py-0.5 rounded text-[10px] font-bold border', getCategoryStyle(run.category)]">
                  {{ run.category || 'N/A' }}
                </span>
                <p class="font-bold text-slate-200 truncate max-w-[200px]">{{ run.name || 'Test #' + run.id }}</p>
              </div>
              <p v-if="run.target_url" class="text-xs text-slate-500 font-mono truncate max-w-[300px] mt-0.5">{{ run.target_url }}</p>
            </td>
            <td class="px-6 py-4">
              <span :class="[
                'px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-tighter',
                run.status === 'completed' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' :
                run.status === 'running' ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20 animate-pulse' :
                run.status === 'failed' ? 'bg-rose-500/10 text-rose-400 border border-rose-500/20' :
                'bg-amber-500/10 text-amber-400 border border-amber-500/20'
              ]">
                {{ run.status }}
              </span>
            </td>
            <td class="px-6 py-4 font-mono text-slate-300">{{ run.vus || '—' }}</td>
            <td class="px-6 py-4 text-slate-300">{{ run.duration || '—' }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openRunDetails(run)" class="p-2 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-slate-700 rounded-lg text-indigo-400" title="View Details">
                <ExternalLinkIcon class="w-4 h-4" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="py-16 flex flex-col items-center justify-center text-slate-600">
        <ZapIcon class="w-10 h-10 mb-3 opacity-20" />
        <p class="font-medium">No storms found for this category</p>
        <p class="text-sm text-slate-700">Try changing the filter or launch a new test</p>
      </div>
    </div>

    <!-- Run Details Modal -->
    <div v-if="showDetailsModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-lg overflow-hidden shadow-2xl animate-in zoom-in-95 duration-200">
        <div class="p-6 border-b border-slate-800 bg-gradient-to-r from-indigo-600/10 to-transparent flex justify-between items-center">
          <div>
            <h3 class="text-xl font-bold">Storm Report</h3>
            <p class="text-slate-500 text-xs">{{ selectedRun?.name }}</p>
          </div>
          <button @click="showDetailsModal = false" class="text-slate-500 hover:text-white">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        
        <div class="p-6 space-y-6">
          <!-- Modal Tabs -->
          <div class="flex gap-4 border-b border-slate-800">
            <button @click="modalTab = 'STATS'" :class="['pb-2 text-sm font-bold transition-all', modalTab === 'STATS' ? 'text-indigo-400 border-b-2 border-indigo-500' : 'text-slate-500']">Metrics Statistics</button>
            <button @click="modalTab = 'LOGS'" :class="['pb-2 text-sm font-bold transition-all', modalTab === 'LOGS' ? 'text-indigo-400 border-b-2 border-indigo-500' : 'text-slate-500']">Execution Logs</button>
          </div>

          <div v-if="modalTab === 'STATS'" class="space-y-6 animate-in fade-in duration-300">
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-slate-800/50 p-4 rounded-2xl border border-slate-700">
                <p class="text-[10px] uppercase font-bold text-slate-500 mb-1">Avg Latency</p>
                <p class="text-2xl font-black text-indigo-400">{{ Math.round(selectedRunMetrics?.avg_latency || 0) }}ms</p>
              </div>
              <div class="bg-slate-800/50 p-4 rounded-2xl border border-slate-700">
                <p class="text-[10px] uppercase font-bold text-slate-500 mb-1">Max Latency</p>
                <p class="text-2xl font-black text-rose-400">{{ Math.round(selectedRunMetrics?.max_latency || 0) }}ms</p>
              </div>
              <div class="bg-slate-800/50 p-4 rounded-2xl border border-slate-700">
                <p class="text-[10px] uppercase font-bold text-slate-500 mb-1">Total Requests</p>
                <p class="text-2xl font-black text-white">{{ selectedRunMetrics?.requests || 0 }}</p>
              </div>
              <div class="bg-slate-800/50 p-4 rounded-2xl border border-slate-700">
                <p class="text-[10px] uppercase font-bold text-slate-500 mb-1">Success / Failed</p>
                <p class="text-2xl font-black text-white">
                  <span class="text-emerald-400">{{ selectedRunMetrics?.success || 0 }}</span> / 
                  <span class="text-rose-400">{{ selectedRunMetrics?.failed || 0 }}</span>
                </p>
              </div>
              <div class="bg-slate-800/50 p-4 rounded-2xl border border-slate-700">
                <p class="text-[10px] uppercase font-bold text-slate-500 mb-1">Configuration</p>
                <p class="text-sm font-bold text-slate-300">{{ selectedRun?.vus }} VUs / {{ selectedRun?.duration }}</p>
              </div>
            </div>

            <div v-if="selectedRun?.target_url" class="bg-slate-950 p-4 rounded-xl border border-slate-800 font-mono text-xs break-all">
              <p class="text-slate-500 mb-1 uppercase font-bold tracking-tighter">Target URL</p>
              <p class="text-indigo-300">{{ selectedRun.method }} {{ selectedRun.target_url }}</p>
            </div>
          </div>

          <div v-else class="animate-in fade-in duration-300">
            <div class="bg-slate-950 rounded-xl border border-slate-800 p-4 font-mono text-[11px] h-[300px] overflow-y-auto space-y-1 text-slate-300 scrollbar-thin scrollbar-thumb-slate-800">
              <div v-if="!selectedRun?.logs" class="text-slate-600 italic">Waiting for k6 output...</div>
              <div v-for="(line, idx) in selectedRun?.logs?.split('\n')" :key="idx" class="break-all whitespace-pre-wrap">
                <span class="text-slate-600 mr-2">{{ idx + 1 }}</span> {{ line }}
              </div>
            </div>
            <div v-if="selectedRun?.status === 'running'" class="mt-2 flex items-center gap-2 text-xs text-amber-400 animate-pulse">
              <Loader2Icon class="w-3 h-3 animate-spin" />
              Monitoring live output...
            </div>
          </div>
        </div>

        <div class="p-6 bg-slate-800/30 border-t border-slate-800 flex justify-end">
          <button @click="showDetailsModal = false" class="px-6 py-2 bg-slate-800 rounded-xl font-bold hover:bg-slate-700 transition-colors">
            Close Report
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, computed, watch } from 'vue';
import { Chart, registerables } from 'chart.js';
import { ActivityIcon, UsersIcon, TimerIcon, ShieldAlertIcon, ZapIcon, ExternalLinkIcon } from 'lucide-vue-next';
import QuickStorm from '../components/QuickStorm.vue';
import { getMetrics, getRuns, getRunMetrics } from '../services/api';

Chart.register(...registerables);

const lineChartCanvas = ref(null);
const doughnutChartCanvas = ref(null);
let lineChart = null;
let doughnutChart = null;

const categories = ['ALL', 'API', 'BACKEND', 'FRONTEND', 'STAGING', 'PRODUCTION', 'CUSTOM'];
const selectedCategory = ref('ALL');

const showDetailsModal = ref(false);
const modalTab = ref('STATS'); // STATS or LOGS
const selectedRun = ref(null);
const selectedRunMetrics = ref(null);
let logPollTimer = null;

const stats = ref([
  { label: 'Total Requests', value: '0', icon: ActivityIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Peak VUs', value: '0', icon: UsersIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Avg Latency', value: '0ms', icon: TimerIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Error Rate', value: '0%', icon: ShieldAlertIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' }
]);

const recentRuns = ref([]);

const filteredRuns = computed(() => {
  if (selectedCategory.value === 'ALL') return recentRuns.value;
  return recentRuns.value.filter(r => r.category === selectedCategory.value);
});

const completedCount = computed(() => filteredRuns.value.filter(r => r.status === 'completed').length);
const runningCount = computed(() => filteredRuns.value.filter(r => r.status === 'running').length);
const failedCount = computed(() => filteredRuns.value.filter(r => r.status === 'failed').length);
const successRateDisplay = computed(() => {
  const total = completedCount.value + failedCount.value;
  if (total === 0) return 100;
  return Math.round((completedCount.value / total) * 100);
});

const formatDate = (dateStr) => {
  if (!dateStr) return '—';
  const d = new Date(dateStr);
  return d.toLocaleString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
};

const getCategoryStyle = (cat) => {
  const styles = {
    API: 'text-emerald-400 border-emerald-500/20 bg-emerald-500/5',
    BACKEND: 'text-purple-400 border-purple-500/20 bg-purple-500/5',
    FRONTEND: 'text-indigo-400 border-indigo-500/20 bg-indigo-500/5',
    STAGING: 'text-amber-400 border-amber-500/20 bg-amber-500/5',
    PRODUCTION: 'text-rose-400 border-rose-500/20 bg-rose-500/5',
    CUSTOM: 'text-slate-400 border-slate-500/20 bg-slate-500/5'
  };
  return styles[cat] || styles.CUSTOM;
};

const openRunDetails = async (run) => {
  selectedRun.value = run;
  showDetailsModal.value = true;
  modalTab.value = 'STATS';
  selectedRunMetrics.value = null;
  
  try {
    const res = await getRunMetrics(run.id);
    selectedRunMetrics.value = res.data;
    // Start polling logs
    pollLogs();
  } catch (err) {
    console.error('Failed to fetch run metrics:', err);
  }
};

const pollLogs = async () => {
  if (!showDetailsModal.value || !selectedRun.value) return;
  try {
    const res = await api.get(`/runs/${selectedRun.value.id}/logs`);
    selectedRun.value.logs = res.data.logs;
    
    // Continue polling if run is still running or pending
    if (selectedRun.value.status === 'running' || selectedRun.value.status === 'pending') {
      logPollTimer = setTimeout(pollLogs, 2000);
    }
  } catch (err) {
    console.error('Failed to poll logs:', err);
  }
};

watch(showDetailsModal, (val) => {
  if (!val && logPollTimer) {
    clearTimeout(logPollTimer);
  }
});

const initCharts = () => {
  lineChart = new Chart(lineChartCanvas.value, {
    type: 'line',
    data: {
      labels: [],
      datasets: [{
        label: 'Latency (ms)',
        data: [],
        borderColor: '#818cf8',
        backgroundColor: 'rgba(129, 140, 248, 0.1)',
        tension: 0.4,
        fill: true,
        borderWidth: 3,
        pointRadius: 0,
        pointHoverRadius: 6,
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: {
        y: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#64748b' } },
        x: { grid: { display: false }, ticks: { color: '#64748b' } }
      }
    }
  });

  doughnutChart = new Chart(doughnutChartCanvas.value, {
    type: 'doughnut',
    data: {
      datasets: [{
        data: [1, 0, 0],
        backgroundColor: ['#10b981', '#f59e0b', '#f43f5e'],
        borderWidth: 0,
        cutout: '85%'
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } }
    }
  });
};

let refreshTimer = null;

const fetchData = async () => {
  try {
    const [metricsRes, runsRes] = await Promise.all([
      getMetrics('-1h', selectedCategory.value),
      getRuns()
    ]);

    const data = metricsRes.data || [];
    recentRuns.value = runsRes.data || [];

    // Update Peak VUs & Error Rate from filtered runs
    if (filteredRuns.value.length > 0) {
      const peakVUs = Math.max(...filteredRuns.value.map(r => r.vus || 0));
      stats.value[1].value = peakVUs.toString();
      stats.value[3].value = `${100 - successRateDisplay.value}%`;
    } else {
      stats.value[1].value = '0';
      stats.value[3].value = '0%';
    }

    // Update doughnut chart
    if (doughnutChart) {
      doughnutChart.data.datasets[0].data = [
        completedCount.value || (filteredRuns.value.length === 0 ? 1 : 0),
        runningCount.value,
        failedCount.value
      ];
      doughnutChart.update();
    }

    // Update Line Chart (Latency)
    const durationData = data.filter(d => d.meas === 'http_req_duration');
    const reqsData = data.filter(d => d.meas === 'http_reqs');

    if (durationData.length > 0) {
      const displayData = durationData.slice(-20);
      lineChart.data.labels = displayData.map(d => new Date(d.time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }));
      lineChart.data.datasets[0].data = displayData.map(d => Math.round(d.value));
      lineChart.update();
      const avgLatency = Math.round(durationData.reduce((a, b) => a + b.value, 0) / durationData.length);
      stats.value[2].value = `${avgLatency}ms`;
    } else {
      lineChart.data.labels = [];
      lineChart.data.datasets[0].data = [];
      lineChart.update();
      stats.value[2].value = '0ms';
    }

    if (reqsData.length > 0) {
      stats.value[0].value = reqsData.length.toLocaleString();
    } else {
      stats.value[0].value = '0';
    }

  } catch (err) {
    console.error('Failed to fetch dashboard data:', err);
  } finally {
    // Schedule next fetch
    refreshTimer = setTimeout(fetchData, 5000);
  }
};

watch(selectedCategory, () => {
  clearTimeout(refreshTimer);
  fetchData();
});

onMounted(() => {
  initCharts();
  fetchData();
});

onUnmounted(() => {
  if (refreshTimer) clearTimeout(refreshTimer);
  if (lineChart) lineChart.destroy();
  if (doughnutChart) doughnutChart.destroy();
});
</script>
