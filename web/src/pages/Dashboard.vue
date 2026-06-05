<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Quick Storm Widget -->
    <QuickStorm />

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
      <div class="lg:col-span-2 bg-slate-900 border border-slate-800 p-8 rounded-3xl shadow-xl">
        <div class="flex items-center justify-between mb-8">
          <h3 class="text-xl font-bold">Request Performance</h3>
          <div class="flex gap-2">
            <span class="flex items-center gap-1.5 text-xs text-indigo-400 font-medium px-2 py-1 bg-indigo-400/10 rounded-md">
              <span class="w-2 h-2 rounded-full bg-indigo-400"></span> Avg. Latency
            </span>
          </div>
        </div>
        <canvas ref="lineChartCanvas" class="w-full h-80"></canvas>
      </div>

      <div class="bg-slate-900 border border-slate-800 p-8 rounded-3xl shadow-xl flex flex-col">
        <h3 class="text-xl font-bold mb-8">Success Rate</h3>
        <div class="flex-1 flex items-center justify-center relative">
          <canvas ref="doughnutChartCanvas" class="max-w-[200px]"></canvas>
          <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
            <span class="text-3xl font-black text-white">98.5%</span>
            <span class="text-xs text-slate-500 uppercase font-bold tracking-tighter">Healthy</span>
          </div>
        </div>
        <div class="mt-8 space-y-3">
          <div class="flex justify-between text-sm">
            <span class="text-slate-400">HTTP 2xx</span>
            <span class="font-bold text-emerald-400">9850</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-400">HTTP 4xx</span>
            <span class="font-bold text-amber-400">120</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-400">HTTP 5xx</span>
            <span class="font-bold text-rose-400">30</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Run Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-3xl overflow-hidden shadow-xl">
      <div class="p-6 border-b border-slate-800 flex justify-between items-center">
        <h3 class="font-bold text-lg">Recent Storms</h3>
        <button class="text-sm font-medium text-indigo-400 hover:text-indigo-300">View All Runs</button>
      </div>
      <table class="w-full text-left">
        <thead class="bg-slate-800/50 text-slate-400 text-xs uppercase tracking-widest">
          <tr>
            <th class="px-6 py-4 font-semibold">Test Name</th>
            <th class="px-6 py-4 font-semibold">Status</th>
            <th class="px-6 py-4 font-semibold">Peak VUs</th>
            <th class="px-6 py-4 font-semibold">Duration</th>
            <th class="px-6 py-4 font-semibold text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="run in recentRuns" :key="run.id" class="hover:bg-slate-800/30 transition-colors group">
            <td class="px-6 py-4">
              <p class="font-bold text-slate-200">{{ run.name }}</p>
              <p class="text-xs text-slate-500 italic">{{ run.date }}</p>
            </td>
            <td class="px-6 py-4">
              <span :class="[
                'px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-tighter',
                run.status === 'Completed' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
              ]">
                {{ run.status }}
              </span>
            </td>
            <td class="px-6 py-4 font-mono text-slate-300">{{ run.vus }}</td>
            <td class="px-6 py-4 text-slate-300">{{ run.duration }}</td>
            <td class="px-6 py-4 text-right">
              <button class="p-2 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-slate-700 rounded-lg">
                <ExternalLinkIcon class="w-4 h-4 text-indigo-400" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { Chart, registerables } from 'chart.js';
import { ActivityIcon, UsersIcon, TimerIcon, ShieldAlertIcon, ExternalLinkIcon } from 'lucide-vue-next';
import QuickStorm from '../components/QuickStorm.vue';

Chart.register(...registerables);

const lineChartCanvas = ref(null);
const doughnutChartCanvas = ref(null);

const stats = [
  { label: 'Total Requests', value: '1,245k', icon: ActivityIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Peak VUs', value: '500', icon: UsersIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Avg Latency', value: '142ms', icon: TimerIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' },
  { label: 'Error Rate', value: '0.12%', icon: ShieldAlertIcon, color: 'bg-indigo-600 shadow-indigo-600/20 shadow-lg' }
];

const recentRuns = [
  { id: 1, name: 'API Gateway Stress', status: 'Completed', vus: 200, duration: '10m', date: '2 mins ago' },
  { id: 2, name: 'Login Flow Soak', status: 'Running', vus: 50, duration: '2h', date: 'Ongoing' },
  { id: 3, name: 'Checkout Spike', status: 'Completed', vus: 1000, duration: '2m', date: '1 hour ago' }
];

onMounted(() => {
  // Line Chart
  new Chart(lineChartCanvas.value, {
    type: 'line',
    data: {
      labels: ['10:00', '10:05', '10:10', '10:15', '10:20', '10:25', '10:30'],
      datasets: [{
        label: 'Latency (ms)',
        data: [120, 140, 135, 170, 150, 160, 145],
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

  // Doughnut Chart
  new Chart(doughnutChartCanvas.value, {
    type: 'doughnut',
    data: {
      datasets: [{
        data: [985, 12, 3],
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
});
</script>
