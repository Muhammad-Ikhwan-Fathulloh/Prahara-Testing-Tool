<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Quick Storm Widget -->
    <QuickStorm @test-started="fetchData" />

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
          <h3 class="text-xl font-bold">Request Performance</h3>
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
        <span class="text-sm text-slate-500">{{ recentRuns.length }} runs</span>
      </div>
      <table v-if="recentRuns.length > 0" class="w-full text-left">
        <thead class="bg-slate-800/50 text-slate-400 text-xs uppercase tracking-widest">
          <tr>
            <th class="px-6 py-4 font-semibold">Test Name</th>
            <th class="px-6 py-4 font-semibold">Status</th>
            <th class="px-6 py-4 font-semibold">VUs</th>
            <th class="px-6 py-4 font-semibold">Duration</th>
            <th class="px-6 py-4 font-semibold text-right">Started</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="run in recentRuns" :key="run.id" class="hover:bg-slate-800/30 transition-colors group">
            <td class="px-6 py-4">
              <p class="font-bold text-slate-200 truncate max-w-[300px]">{{ run.name || 'Test #' + run.id }}</p>
              <p v-if="run.target_url" class="text-xs text-slate-500 font-mono truncate max-w-[300px]">{{ run.target_url }}</p>
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
            <td class="px-6 py-4 text-right text-sm text-slate-400">{{ formatDate(run.started_at) }}</td>
          </tr>
        </tbody>
      </table>
      <div v-else class="py-16 flex flex-col items-center justify-center text-slate-600">
        <ZapIcon class="w-10 h-10 mb-3 opacity-20" />
        <p class="font-medium">No storms launched yet</p>
        <p class="text-sm text-slate-700">Use Quick Storm above to run your first test</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, computed } from 'vue';
import { Chart, registerables } from 'chart.js';
import { ActivityIcon, UsersIcon, TimerIcon, ShieldAlertIcon, ZapIcon } from 'lucide-vue-next';
import QuickStorm from '../components/QuickStorm.vue';
import { getMetrics, getRuns } from '../services/api';

Chart.register(...registerables);

const lineChartCanvas = ref(null);
const doughnutChartCanvas = ref(null);
let lineChart = null;
let doughnutChart = null;

const stats = ref([
  { label: 'Total Requests', value: '0', icon: ActivityIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Peak VUs', value: '0', icon: UsersIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Avg Latency', value: '0ms', icon: TimerIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Error Rate', value: '0%', icon: ShieldAlertIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' }
]);

const recentRuns = ref([]);

const completedCount = computed(() => recentRuns.value.filter(r => r.status === 'completed').length);
const runningCount = computed(() => recentRuns.value.filter(r => r.status === 'running').length);
const failedCount = computed(() => recentRuns.value.filter(r => r.status === 'failed').length);
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
        data: [100, 0, 0],
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

const fetchData = async () => {
  try {
    const [metricsRes, runsRes] = await Promise.all([
      getMetrics('-1h'),
      getRuns()
    ]);

    const data = metricsRes.data || [];
    recentRuns.value = runsRes.data || [];

    // Update Peak VUs & Error Rate from runs
    if (recentRuns.value.length > 0) {
      const peakVUs = Math.max(...recentRuns.value.map(r => r.vus || 0));
      if (peakVUs > 0) stats.value[1].value = peakVUs.toString();
      stats.value[3].value = `${100 - successRateDisplay.value}%`;
    }

    // Update doughnut chart with real run status counts
    if (doughnutChart) {
      doughnutChart.data.datasets[0].data = [
        completedCount.value || 1,
        runningCount.value,
        failedCount.value
      ];
      doughnutChart.update();
    }

    if (data.length === 0) return;

    const durationData = data.filter(d => d.meas === 'http_req_duration');
    const reqsData = data.filter(d => d.meas === 'http_reqs');

    if (durationData.length > 0) {
      const displayData = durationData.slice(-20);
      lineChart.data.labels = displayData.map(d => new Date(d.time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }));
      lineChart.data.datasets[0].data = displayData.map(d => Math.round(d.value));
      lineChart.update();
      const avgLatency = Math.round(durationData.reduce((a, b) => a + b.value, 0) / durationData.length);
      stats.value[2].value = `${avgLatency}ms`;
    }

    if (reqsData.length > 0) {
      stats.value[0].value = reqsData.length.toLocaleString();
    }

  } catch (err) {
    console.error('Failed to fetch dashboard data:', err);
  }
};

let refreshInterval = null;

onMounted(() => {
  initCharts();
  fetchData();
  refreshInterval = setInterval(fetchData, 5000);
});

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval);
  if (lineChart) lineChart.destroy();
  if (doughnutChart) doughnutChart.destroy();
});
</script>
