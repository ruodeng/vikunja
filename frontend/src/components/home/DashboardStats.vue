<template>
	<div class="dashboard-stats is-max-width-desktop mbs-4">
		<h3 class="mbb-2">
			{{ $t('home.stats') }}
		</h3>
		<div class="stats-grid">
			<!-- In Progress -->
			<div 
				class="stat-card clickable"
				@click="scrollToSection('dashboard-tasks')"
			>
				<div class="stat-value">
					{{ inProgressCount }}
				</div>
				<div class="stat-label">
					{{ $t('task.stats.inProgress') }}
				</div>
			</div>

			<!-- Overdue -->
			<div
				class="stat-card clickable"
				:class="{ 'is-danger': overdueCount > 0 }"
				@click="scrollToSection('dashboard-tasks')"
			>
				<div class="stat-value">
					{{ overdueCount }}
				</div>
				<div class="stat-label">
					{{ $t('task.stats.overdue') }}
				</div>
			</div>

			<!-- Comments -->
			<div 
				class="stat-card clickable"
				@click="scrollToSection('dashboard-comments')"
			>
				<div class="stat-value">
					{{ commentsCount }}
				</div>
				<div class="stat-label">
					{{ $t('task.stats.comments') }}
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import TaskService from '@/services/task'
import NotificationService from '@/services/notification'
import type { INotification } from '@/modelTypes/INotification'

const inProgressCount = ref(0)
const overdueCount = ref(0)
const commentsCount = ref(0)

onMounted(() => {
	fetchStats()
})

async function fetchStats() {
	const taskService = new TaskService()
	const notificationService = new NotificationService()

	// Fetch all tasks for now (optimize later if needed)
	// Using filter parameters if API supports it would be better
	try {
		// Mocking counts basically by fetching first page or reasonable number, 
		// but since we need accurate counts, we might need a specific endpoint or fetch all.
		// For MVP, letting's assumes we can get a list.
		// NOTE: Fetching *all* might be heavy. 
		// Ideally we use filter='undone'
		 
		const allTasks = await taskService.getAll({
			filter: 'undone',
			sort_by: 'due_date',
			order_by: 'asc',
		} as any)
		
		const now = new Date()
		
		// In Progress: Undone tasks
		inProgressCount.value = allTasks.length

		// Overdue: Due date exists, is before now, and not done
		overdueCount.value = allTasks.filter(t => {
			if (!t.dueDate) return false
			const dueDate = new Date(t.dueDate)
			return dueDate.getFullYear() > 1000 && dueDate < now
		}).length

	} catch (e) {
		console.error('Failed to fetch task stats', e)
	}

	try {
		// Comments: Unread notifications that are mentions or comments
		const notifications = await notificationService.getAll()
		commentsCount.value = notifications.filter((n: INotification) => n.readAt === null).length
	} catch (e) {
		console.error('Failed to fetch notification stats', e)
	}
}

function scrollToSection(id: string) {
	const element = document.getElementById(id)
	if (element) {
		element.scrollIntoView({ behavior: 'smooth' })
	}
}
</script>

<style scoped lang="scss">
.dashboard-stats {
	margin-left: auto;
	margin-right: auto;
}

.stats-grid {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 1rem;
}

.stat-card {
	background: var(--white);
	border-radius: $radius;
	box-shadow: var(--shadow-sm);
	padding: 1.5rem;
	text-align: center;
	border: 1px solid var(--grey-200);
	transition: transform 0.2s, box-shadow 0.2s;

	&.clickable {
		cursor: pointer;
		&:hover {
			transform: translateY(-2px);
			box-shadow: var(--shadow);
		}
	}

	&.is-danger {
		border-color: var(--danger);
		color: var(--danger);
		.stat-value {
			color: var(--danger);
		}
	}
}

.stat-value {
	font-size: 2rem;
	font-weight: bold;
	color: var(--primary);
	line-height: 1.2;
}

.stat-label {
	color: var(--grey-500);
	font-size: 0.9rem;
	margin-top: 0.5rem;
}
</style>
