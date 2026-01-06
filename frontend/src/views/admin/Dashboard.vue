<template>
	<div class="container user-view">
		<h1>{{ $t('admin.dashboard') }}</h1>
		
		<div
			v-if="loading"
			class="base-card loading"
		>
			<h3 class="has-text-centered">
				Loading...
			</h3>
		</div>
		
		<div
			v-else-if="error"
			class="base-card error"
		>
			<div class="notification is-danger">
				{{ error }}
			</div>
		</div>

		<div
			v-else
			class="dashboard-content"
		>
			<!-- Overview Cards -->
			<div class="columns is-multiline">
				<div class="column is-4">
					<div class="base-card overview-card">
						<h3>{{ $t('admin.total_users') }}</h3>
						<p class="stat-number">
							{{ stats.total_users }}
						</p>
					</div>
				</div>
				<div class="column is-4">
					<div class="base-card overview-card">
						<h3>{{ $t('admin.total_tasks') }}</h3>
						<p class="stat-number">
							{{ stats.total_tasks }}
						</p>
					</div>
				</div>
				<div class="column is-4">
					<div
						class="base-card overview-card"
						:class="{'has-text-danger': stats.overdue_tasks > 0}"
					>
						<h3>{{ $t('admin.overdue_tasks') }}</h3>
						<p class="stat-number">
							{{ stats.overdue_tasks }}
						</p>
					</div>
				</div>
			</div>

			<!-- User Statistics Table -->
			<div class="base-card">
				<h3>{{ $t('admin.user_performance') }}</h3>
				<div class="table-container">
					<table class="table is-fullwidth is-hoverable">
						<thead>
							<tr>
								<th>{{ $t('admin.user') }}</th>
								<th>{{ $t('admin.active_tasks') }}</th>
								<th>{{ $t('admin.completed_today') }}</th>
								<th>{{ $t('admin.overdue') }}</th>
							</tr>
						</thead>
						<tbody>
							<tr
								v-for="userStat in stats.user_stats"
								:key="userStat.user.id"
								class="is-clickable"
								@click="openUserDetail(userStat.user)"
							>
								<td>
									<div class="is-flex is-align-items-center">
										<div class="user-avatar-wrapper mr-2">
											<User
												:user="userStat.user"
												:avatar-size="36"
											/>
										</div>
										<span>{{ userStat.user.username }}</span>
										<span
											v-if="userStat.user.name"
											class="has-text-grey is-size-7 ml-2"
										>({{ userStat.user.name }})</span>
									</div>
								</td>
								<td>{{ userStat.active_tasks }}</td>
								<td class="has-text-success">
									{{ userStat.completed_today }}
								</td>
								<td :class="{'has-text-danger': userStat.overdue_tasks > 0}">
									{{ userStat.active_tasks > 0 ? userStat.overdue_tasks : '-' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>

			<!-- Debug Info -->
			<div class="base-card m-t-2">
				<h4 class="is-size-5 mb-2">
					{{ $t('admin.debug_info') }}
				</h4>
				<pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; overflow: auto;">
{{ $t('admin.api_url') }}: {{ apiUrl }}
{{ $t('admin.stats_raw') }}: {{ JSON.stringify(stats, null, 2) }}
{{ $t('admin.error') }}: {{ error }}
				</pre>
			</div>
		</div>
	</div>

	<!-- User Detail Modal -->
	<Modal
		v-if="showUserDetail"
		:title="$t('admin.user_details')"
		@close="closeUserDetail"
	>
		<div
			v-if="selectedUserDetails"
			class="user-detail-content"
		>
			<!-- Header -->
			<div class="user-header">
				<div class="user-avatar-large">
					<User
						:user="selectedUserDetails.user"
						:avatar-size="64"
					/>
				</div>
				<div class="user-info">
					<h3 class="title is-4 mb-0">
						{{ selectedUserDetails.user.username }}
					</h3>
					<p
						v-if="selectedUserDetails.user.email"
						class="subtitle is-6 has-text-grey"
					>
						{{ selectedUserDetails.user.email }}
					</p>
					<span
						v-if="selectedUserDetails.user.name"
						class="tag is-light mt-1"
					>
						{{ selectedUserDetails.user.name }}
					</span>
				</div>
			</div>

			<hr class="dropdown-divider my-4">

			<!-- Stats Grid -->
			<div class="stats-grid mb-5">
				<div class="stat-card">
					<div class="stat-icon is-primary-bg">
						<i class="fas fa-tasks" />
					</div>
					<div class="stat-content">
						<p class="heading">
							{{ $t('admin.total_tasks') }}
						</p>
						<p class="title is-4">
							{{ selectedUserDetails.total_tasks }}
						</p>
					</div>
				</div>
				
				<div class="stat-card">
					<div class="stat-icon is-danger-bg">
						<i class="fas fa-exclamation-circle" />
					</div>
					<div class="stat-content">
						<p class="heading has-text-danger">
							{{ $t('admin.overdue_tasks') }}
						</p>
						<p class="title is-4 has-text-danger">
							{{ selectedUserDetails.overdue_tasks }}
						</p>
					</div>
				</div>

				<div class="stat-card">
					<div class="stat-icon is-success-bg">
						<i class="fas fa-check-circle" />
					</div>
					<div class="stat-content">
						<p class="heading has-text-success">
							{{ $t('admin.active') }}
						</p>
						<p class="title is-4 has-text-success">
							{{ selectedUserDetails.active_tasks }}
						</p>
					</div>
				</div>
			</div>

			<!-- Task List -->
			<div class="tasks-section">
				<h4 class="title is-5 mb-3">
					{{ $t('task.task') }}
				</h4>
				<div
					class="table-container custom-scrollbar"
					style="max-height: 350px; overflow-y: auto;"
				>
					<table
						v-if="selectedUserDetails.tasks && selectedUserDetails.tasks.length > 0"
						class="table is-fullwidth is-hoverable is-striped custom-table"
					>
						<thead>
							<tr>
								<th style="width: 60px;">
									{{ $t('admin.id') }}
								</th>
								<th>{{ $t('admin.task_title') }}</th>
								<th>{{ $t('admin.due_date') }}</th>
								<th>{{ $t('admin.status') }}</th>
							</tr>
						</thead>
						<tbody>
							<tr
								v-for="task in selectedUserDetails.tasks"
								:key="task.id"
							>
								<td class="has-text-grey-light">
									#{{ task.id }}
								</td>
								<td class="has-text-weight-medium">
									<RouterLink
										:to="{ name: 'task.detail', params: { id: task.id } }"
										class="task-link"
										@click.stop="closeUserDetail"
									>
										{{ task.title }}
									</RouterLink>
								</td>
								<td>
									<span :class="{'has-text-danger': task.due_date && new Date(task.due_date).getFullYear() > 1000 && new Date(task.due_date) < new Date() && !task.done}">
										{{ task.due_date ? new Date(task.due_date).toLocaleDateString() : '-' }}
									</span>
								</td>
								<td>
									<span
										class="tag is-rounded"
										:class="task.done ? 'is-success is-light' : (task.due_date && new Date(task.due_date).getFullYear() > 1000 && new Date(task.due_date) < new Date() ? 'is-danger is-light' : 'is-info is-light')"
									>
										{{ task.done ? $t('admin.done') : (task.due_date && new Date(task.due_date).getFullYear() > 1000 && new Date(task.due_date) < new Date() ? $t('admin.overdue') : $t('admin.active')) }}
									</span>
								</td>
							</tr>
						</tbody>
					</table>
					<div
						v-else
						class="empty-state py-5"
					>
						<p class="has-text-centered has-text-grey">
							<i class="fas fa-check-double fa-2x mb-2" />
							<br>
							{{ $t('admin.no_tasks') }}
						</p>
					</div>
				</div>
			</div>
		</div>
		<div
			v-else
			class="loader-wrapper is-flex is-justify-content-center p-6"
		>
			<div class="loader is-loading" />
		</div>
	</Modal>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { AuthenticatedHTTPFactory } from '@/helpers/fetcher'
import User from '@/components/misc/User.vue'
import Modal from '@/components/misc/Modal.vue'

interface Task {
	id: number
	title: string
	done: boolean
	due_date?: string
}

interface UserStat {
	user: any
	active_tasks: number
	overdue_tasks: number
	completed_today: number
}

interface UserDetailStats {
	user: UserStat['user']
	total_tasks: number
	active_tasks: number
	overdue_tasks: number
	tasks: Task[]
}

interface DashboardStats {
	total_tasks: number
	overdue_tasks: number
	total_users: number
	user_stats: UserStat[]
}

const stats = ref<DashboardStats>({
	total_tasks: 0,
	overdue_tasks: 0,
	total_users: 0,
	user_stats: [],
})

const { t } = useI18n()
const loading = ref(true)
const error = ref('')
const apiUrl = window.API_URL

// User Detail Modal State
const showUserDetail = ref(false)
const selectedUserDetails = ref<UserDetailStats | null>(null)

const openUserDetail = async (user: UserStat['user']) => {
	showUserDetail.value = true
	selectedUserDetails.value = null // Reset previous data
	
	try {
		const api = AuthenticatedHTTPFactory()
		const response = await api.get(`/admin/users/${user.id}`)
		selectedUserDetails.value = response.data
	} catch (e: any) {
		console.error('Failed to load user details:', e)
		// Optionally show an error in the modal or toast
	}
}

const closeUserDetail = () => {
	showUserDetail.value = false
	selectedUserDetails.value = null
}

onMounted(async () => {
	try {
		loading.value = true
		console.log('Dashboard: Using API URL:', window.API_URL)
		const api = AuthenticatedHTTPFactory()
		const response = await api.get('/admin/dashboard')
		console.log('Dashboard: Response:', response)
		stats.value = response.data
	} catch (e: any) {
		console.error('Dashboard Error:', e)
		if (e.response && e.response.status === 403) {
			error.value = t('admin.not_authorized')
		} else {
			error.value = t('admin.failed_to_load') + ': ' + (e.message || 'Unknown error')
		}
	} finally {
		loading.value = false
	}
})
</script>

<style lang="scss" scoped>
.user-detail-content {
	background-color: var(--bg-surface, #ffffff);
	border-radius: var(--radius-large);
	padding: 2rem;
	box-shadow: var(--box-shadow);
	text-align: left;
	width: 100%;
	color: var(--text-main);
	
	@media screen and (min-width: 768px) {
		min-width: 600px;
	}
}

.overview-card {
	text-align: center;
	padding: 1.5rem;
	
	h3 {
		margin-bottom: 0.5rem;
		color: var(--text-medium);
		font-weight: bold;
	}
	
	.stat-number {
		font-size: 2.5rem;
		font-weight: bold;
		color: var(--text-main);
	}
}

.user-avatar-wrapper {
	width: 36px;
	height: 36px;
	overflow: hidden;
	border-radius: 50%;
	flex-shrink: 0;
	
	/* Ensure content (like alt text) doesn't break out */
	display: flex;
	align-items: center;
	justify-content: center;
}

.table-container {
	margin-top: 1rem;
}

.mr-2 {
	margin-right: 0.5rem;
}
.ml-2 {
	margin-left: 0.5rem;
}

.is-clickable {
	cursor: pointer;
	transition: background-color 0.2s;
	
	&:hover {
		background-color: var(--background-hover);
	}
}

/* User Detail Modal Styles */
.user-header {
	display: flex;
	align-items: center;
	
	.user-avatar-large {
		margin-right: 1.5rem;
	}
	
	.user-info {
		.title {
			line-height: 1.2;
		}
		.tag {
			font-size: 0.75rem;
		}
	}
}

.stats-grid {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 1rem;
	
	.stat-card {
		background: var(--background-light);
		border-radius: var(--radius-main);
		padding: 1rem;
		display: flex;
		align-items: center;
		transition: transform 0.2s, box-shadow 0.2s;
		
		&:hover {
			transform: translateY(-2px);
			box-shadow: 0 4px 12px rgba(0,0,0,0.05);
		}
		
		.stat-icon {
			width: 40px;
			height: 40px;
			border-radius: 10px;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-right: 1rem;
			font-size: 1.2rem;
			
			&.is-primary-bg {
				background: rgba(var(--primary-rgb), 0.1);
				color: var(--primary);
			}
			&.is-danger-bg {
				background: rgba(255, 56, 96, 0.1);
				color: #ff3860;
			}
			&.is-success-bg {
				background: rgba(35, 209, 96, 0.1);
				color: #23d160;
			}
		}
		
		.stat-content {
			.heading {
				margin-bottom: 0px !important;
				font-size: 0.75rem;
				font-weight: 600;
				letter-spacing: 0.5px;
				text-transform: uppercase;
				color: var(--text-light);
			}
			.title {
				margin-bottom: 0 !important;
				font-weight: 700;
			}
		}
	}
}

.custom-table {
	th {
		border-bottom-width: 2px;
		color: var(--text-medium);
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}
	
	td {
		vertical-align: middle;
	}
}

.task-link {
	color: var(--primary);
	text-decoration: none;
	transition: color 0.2s;
	
	&:hover {
		color: var(--primary-dark, #1a73e8);
		text-decoration: underline;
	}
}

.empty-state {
	background: var(--background-light);
	border-radius: 8px;
	opacity: 0.7;
}

@media screen and (max-width: 768px) {
	.stats-grid {
		grid-template-columns: 1fr;
	}
}
</style>
