<template>
  <div class="max-w-4xl space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <div>
      <h2 class="text-3xl font-bold mb-2">Configuration</h2>
      <p class="text-slate-500">Manage your Prahara system settings and RBAC</p>
    </div>

    <div class="grid grid-cols-1 gap-8">
      <!-- Database Settings -->
      <section class="bg-slate-900 border border-slate-800 p-8 rounded-3xl space-y-6">
        <h3 class="text-xl font-bold flex items-center gap-2">
          <DatabaseIcon class="w-5 h-5 text-indigo-400" /> Infrastructure
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-2">
            <label class="text-sm font-medium text-slate-400">InfluxDB URL</label>
            <input type="text" value="http://influxdb:8086" readonly class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-3 text-slate-300 font-mono text-sm" />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium text-slate-400">InfluxDB Bucket</label>
            <input type="text" value="prahara" readonly class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-3 text-slate-300 font-mono text-sm" />
          </div>
        </div>
      </section>

      <!-- RBAC Settings -->
      <section class="bg-slate-900 border border-slate-800 p-8 rounded-3xl space-y-6">
        <h3 class="text-xl font-bold flex items-center gap-2">
          <ShieldCheckIcon class="w-5 h-5 text-indigo-400" /> User Management
        </h3>
        <div class="overflow-hidden rounded-xl border border-slate-800">
          <table class="w-full text-left">
            <thead class="bg-slate-800 text-xs uppercase tracking-widest text-slate-500">
              <tr>
                <th class="px-6 py-3 font-semibold">User</th>
                <th class="px-6 py-3 font-semibold">Role</th>
                <th class="px-6 py-3 font-semibold">Last Login</th>
                <th class="px-6 py-3 font-semibold text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-if="user" class="hover:bg-slate-800/50 transition-colors">
                <td class="px-6 py-4 font-medium">{{ user.username }}</td>
                <td class="px-6 py-4 text-indigo-400 font-bold uppercase text-[10px] tracking-widest">{{ user.role }}</td>
                <td class="px-6 py-4 text-xs text-slate-500">Active Session</td>
                <td class="px-6 py-4 text-right">
                  <span class="text-[10px] bg-slate-800 px-2 py-1 rounded text-slate-400">CURRENT</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <button class="px-4 py-2 bg-slate-800 rounded-lg text-sm font-bold border border-slate-700 hover:bg-slate-700 transition-colors">
          + Add New User
        </button>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { DatabaseIcon, ShieldCheckIcon } from 'lucide-vue-next';

const user = ref(JSON.parse(localStorage.getItem('prahara_user')));
</script>
