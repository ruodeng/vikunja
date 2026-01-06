<template>
	<div
		id="dashboard-comments"
		class="dashboard-comments is-max-width-desktop mbs-4"
	>
		<div class="is-flex is-align-items-center is-justify-content-space-between mbb-2">
			<h3 class="mbb-0">
				{{ $t('notification.title') }}
			</h3>
			<BaseButton
				v-if="notifications.length > 0"
				class="mis-auto"
				variant="tertiary"
				@click="markAllRead"
			>
				{{ $t('notification.markAllRead') }}
			</BaseButton>
		</div>

		<div
			v-if="notifications.length > 0"
			class="box"
		>
			<table class="table is-fullwidth is-hoverable">
				<thead>
					<tr>
						<th />
						<th>{{ $t('task.attributes.assignees') }}</th>
						<th>{{ $t('task.comment.title') }}</th>
						<th>{{ $t('task.title') }}</th>
						<th>{{ $t('project.title') }}</th>
						<th class="is-narrow" />
					</tr>
				</thead>
				<tbody>
					<tr
						v-for="n in notifications"
						:key="n.id"
					>
						<!-- Notification Type Icon (Optional enhancement) -->
						<td class="is-vcentered is-narrow">
							<!-- Could add icon here based on n.name -->
						</td>

						<!-- Who -->
						<td class="is-vcentered">
							<User
								v-if="getDoer(n)"
								:user="getDoer(n)"
								:show-username="true"
								:avatar-size="24"
							/>
						</td>

						<!-- Comment/Content -->
						<td class="is-vcentered is-wrap-anywhere">
							<!-- eslint-disable-next-line @typescript-eslint/no-explicit-any -->
							<span>{{ (n as any).toText(userInfo) }}</span>
							<span class="has-text-grey is-size-7 mis-2">
								{{ formatDisplayDate(n.created) }}
							</span>
						</td>

						<!-- Task -->
						<td class="is-vcentered">
							<RouterLink
								v-if="getTask(n)"
								:to="{ name: 'task.detail', params: { id: getTask(n).id } }"
								class="has-text-weight-bold"
								@click="markRead(n)"
							>
								{{ getTask(n).title }}
							</RouterLink>
						</td>

						<!-- Project -->
						<td class="is-vcentered">
							<RouterLink
								v-if="getProject(n)"
								:to="{ name: 'task.index', params: { projectId: getProject(n).id } }"
								@click="markRead(n)"
							>
								<span
									class="project-dot"
									:style="{ backgroundColor: getProject(n).hexColor }"
								/>
								{{ getProject(n).title }}
							</RouterLink>
							<!-- Fallback if project is nested in task -->
							<RouterLink
								v-else-if="getTask(n) && getTask(n).projectId" 
								:to="{ name: 'task.index', params: { projectId: getTask(n).projectId } }" 
								@click="markRead(n)"
							>
								<!-- We don't have project title here, show generic text or nothing -->
							</RouterLink>
						</td>

						<!-- Actions -->
						<td class="is-vcentered has-text-right">
							<BaseButton
								variant="transparent"
								icon="check"
								:title="$t('notification.markRead')"
								@click="markRead(n)"
							/>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
		
		<div
			v-else
			class="has-text-grey has-text-centered py-4 box"
		>
			{{ $t('notification.none') }}
		</div>
	</div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import NotificationService from '@/services/notification'
import { type INotification } from '@/modelTypes/INotification'
import { useAuthStore } from '@/stores/auth'
import { formatDisplayDate } from '@/helpers/time/formatDate'
import BaseButton from '@/components/base/BaseButton.vue'
import User from '@/components/misc/User.vue'
import { success } from '@/message'

const { t } = useI18n()
const authStore = useAuthStore()
const userInfo = computed(() => authStore.info)

const allNotifications = ref<INotification[]>([])

const notifications = computed(() => {
	// Filter unread and valid notifications
	return allNotifications.value
		.filter(n => n.name !== '' && n.readAt === null)
})

onMounted(() => {
	loadNotifications()
})

async function loadNotifications() {
	const notificationService = new NotificationService()
	allNotifications.value = await notificationService.getAll()
}

async function markRead(n: INotification) {
	const notificationService = new NotificationService()
	n.read = true
	// Optimistic update
	allNotifications.value = allNotifications.value.filter(notif => notif.id !== n.id)
	
	await notificationService.update(n)
}

async function markAllRead() {
	const notificationService = new NotificationService()
	await notificationService.markAllRead()
	success({ message: t('notification.markAllReadSuccess') })
	allNotifications.value = []
}

function getTask(n: INotification) {
	if ('task' in n.notification) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		return (n.notification as any).task
	}
	return undefined
}

function getProject(n: INotification) {
	if ('project' in n.notification) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		return (n.notification as any).project
	}
	return undefined
}

function getDoer(n: INotification) {
	if ('doer' in n.notification) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		return (n.notification as any).doer
	}
	return undefined
}

</script>

<style scoped lang="scss">
.dashboard-comments {
	margin-left: auto;
	margin-right: auto;
}

.table {
	background: transparent;
	
	th {
		border: none;
		color: var(--grey-500);
		font-size: 0.85rem;
		font-weight: 600;
	}

	td {
		border-bottom: 1px solid var(--grey-200);
		&:first-child {
			padding-left: 0.5rem;
		}
		&:last-child {
			padding-right: 0.5rem;
		}
	}

	tr:last-child td {
		border-bottom: none;
	}
}

.project-dot {
	display: inline-block;
	width: 10px;
	height: 10px;
	border-radius: 50%;
	margin-right: 0.5rem;
}

.is-wrap-anywhere {
	overflow-wrap: anywhere;
}
</style>
