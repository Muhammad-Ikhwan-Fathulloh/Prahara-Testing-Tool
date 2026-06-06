<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-slate-900 border border-slate-800 p-8 rounded-3xl shadow-lg">
      <div>
        <h2 class="text-3xl font-black text-white flex items-center gap-3">
          <UsersIcon class="w-8 h-8 text-indigo-400" />
          User Management
        </h2>
        <p class="text-slate-500 text-sm">Manage administrative access and tester permissions</p>
      </div>
      <button @click="showAddModal = true" class="px-6 py-3 bg-indigo-600 rounded-xl font-bold flex items-center gap-2 hover:bg-indigo-500 transition-all shadow-lg shadow-indigo-600/20 active:scale-95">
        <UserPlusIcon class="w-5 h-5" />
        Add New User
      </button>
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-3xl overflow-hidden shadow-xl">
      <table class="w-full text-left border-collapse">
        <thead class="bg-slate-800/50 text-slate-400 text-xs uppercase tracking-widest border-b border-white/5">
          <tr>
            <th class="px-8 py-5 font-semibold">User</th>
            <th class="px-8 py-5 font-semibold">Role</th>
            <th class="px-8 py-5 font-semibold">Last Login</th>
            <th class="px-8 py-5 text-right font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="user in users" :key="user.id" class="hover:bg-indigo-600/5 transition-colors group">
            <td class="px-8 py-5">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-slate-700 to-slate-800 flex items-center justify-center border border-white/10 group-hover:border-indigo-500/50 transition-colors">
                  <UserIcon class="w-5 h-5 text-slate-400" />
                </div>
                <div>
                  <p class="font-bold text-white">{{ user.username }}</p>
                  <p class="text-[10px] text-slate-500 uppercase tracking-tighter">ID: {{ user.id }}</p>
                </div>
              </div>
            </td>
            <td class="px-8 py-5">
              <span :class="[
                'px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border',
                user.role === 'admin' ? 'bg-indigo-500/10 text-indigo-400 border-indigo-500/20' : 'bg-slate-500/10 text-slate-400 border-slate-500/20'
              ]">
                {{ user.role }}
              </span>
            </td>
            <td class="px-8 py-5">
              <p class="text-sm font-medium" :class="user.last_login ? 'text-slate-300' : 'text-slate-600 italic'">
                {{ user.last_login ? formatDate(user.last_login) : 'Never logged in' }}
              </p>
            </td>
            <td class="px-8 py-5 text-right">
              <button @click="confirmDelete(user)" class="p-2.5 bg-rose-500/10 text-rose-500 rounded-xl opacity-0 group-hover:opacity-100 transition-all hover:bg-rose-500 hover:text-white" title="Delete User">
                <TrashIcon class="w-4 h-4" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Modal -->
    <div v-if="showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-in fade-in duration-200">
      <div class="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-md overflow-hidden shadow-2xl animate-in zoom-in-95 duration-200">
        <div class="p-6 bg-gradient-to-r from-indigo-600/10 to-transparent border-b border-white/5 flex justify-between items-center">
          <h3 class="text-xl font-bold">New User Account</h3>
          <button @click="showAddModal = false" class="text-slate-500 hover:text-white transition-colors">
            <XIcon class="w-6 h-6" />
          </button>
        </div>
        
        <div class="p-8 space-y-5">
          <div>
            <label class="block text-xs font-bold text-slate-500 uppercase tracking-widest mb-2">Username</label>
            <input v-model="newUser.username" type="text" class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 focus:ring-2 focus:ring-indigo-600 transition-all outline-none" placeholder="e.g. jdoe" />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-500 uppercase tracking-widest mb-2">Password</label>
            <input v-model="newUser.password" type="password" class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 focus:ring-2 focus:ring-indigo-600 transition-all outline-none" placeholder="••••••••" />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-500 uppercase tracking-widest mb-2">System Role</label>
            <div class="grid grid-cols-2 gap-3">
              <button 
                v-for="role in ['admin', 'tester']" :key="role"
                @click="newUser.role = role"
                :class="[
                  'px-4 py-3 rounded-xl text-xs font-bold uppercase transition-all border',
                  newUser.role === role ? 'bg-indigo-600 border-indigo-500 text-white shadow-lg' : 'bg-slate-950 border-slate-800 text-slate-500 hover:border-slate-600'
                ]"
              >
                {{ role }}
              </button>
            </div>
          </div>
        </div>

        <div class="p-6 bg-slate-800/30 border-t border-white/5 flex gap-3">
          <button @click="showAddModal = false" class="flex-1 px-4 py-3 bg-slate-800 rounded-xl font-bold hover:bg-slate-700 transition-colors text-slate-300">
            Cancel
          </button>
          <button @click="createUser" :disabled="isSaving" class="flex-1 px-4 py-3 bg-indigo-600 rounded-xl font-bold hover:bg-indigo-500 transition-all flex items-center justify-center gap-2 group shadow-lg shadow-indigo-600/20">
            <Loader2Icon v-if="isSaving" class="w-4 h-4 animate-spin" />
            {{ isSaving ? 'Creating...' : 'Create User' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { UsersIcon, UserPlusIcon, TrashIcon, UserIcon, XIcon, Loader2Icon } from 'lucide-vue-next';
import api from '../services/api';

const users = ref([]);
const showAddModal = ref(false);
const isSaving = ref(false);
const newUser = ref({ username: '', password: '', role: 'tester' });

const fetchUsers = async () => {
  try {
    const res = await api.get('/users');
    users.value = res.data;
  } catch (err) {
    console.error('Failed to fetch users:', err);
  }
};

const createUser = async () => {
  if (!newUser.value.username || !newUser.value.password) return;
  isSaving.value = true;
  try {
    await api.post('/users', newUser.value);
    await fetchUsers();
    showAddModal.value = false;
    newUser.value = { username: '', password: '', role: 'tester' };
  } catch (err) {
    alert('Failed to create user: ' + (err.response?.data?.error || err.message));
  } finally {
    isSaving.value = false;
  }
};

const confirmDelete = async (user) => {
  if (!confirm(`Are you sure you want to delete user "${user.username}"?`)) return;
  try {
    await api.delete(`/users/${user.id}`);
    await fetchUsers();
  } catch (err) {
    alert('Failed to delete user: ' + (err.response?.data?.error || err.message));
  }
};

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
};

onMounted(fetchUsers);
</script>
