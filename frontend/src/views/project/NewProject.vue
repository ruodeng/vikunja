<template>
	<CreateEdit
		v-model:loading="isSubmitting"
		:title="$t('project.create.header')"
		:primary-disabled="project.title === ''"
		@create="createProject()"
	>
		<div class="field">
			<label
				class="label"
				for="projectTitle"
			>{{ $t('project.title') }}</label>
			<div
				:class="{ 'is-loading': projectService.loading }"
				class="control"
			>
				<input
					v-model="project.title"
					v-focus
					:class="{ disabled: projectService.loading }"
					class="input"
					:placeholder="$t('project.create.titlePlaceholder')"
					type="text"
					name="projectTitle"
					@keyup.enter="createProject()"
					@keyup.esc="$router.back()"
				>
			</div>
		</div>
		<p
			v-if="showError && project.title === ''"
			class="help is-danger"
		>
			{{ $t('project.create.addTitleRequired') }}
		</p>
		<div
			v-if="projectStore.hasProjects"
			class="field"
		>
			<label class="label">{{ $t('project.parent') }}</label>
			<div class="control">
				<ProjectSearch v-model="parentProject" />
			</div>
		</div>
		<div class="field">
			<label class="label">{{ $t('project.color') }}</label>
			<div class="control">
				<ColorPicker v-model="project.hexColor" />
			</div>
		</div>
		
		<div class="columns">
			<div class="column field">
				<label class="label">{{ $t('task.assignee.label') }}</label>
				<div class="control">
					<Multiselect
						v-model="project.assignees"
						:multiple="true"
						:close-after-select="false"
						:search-results="foundUsers"
						label="name"
						:placeholder="$t('task.assignee.placeholder')"
						:select-placeholder="$t('task.assignee.selectPlaceholder')"
						@search="findUser"
					>
						<template #searchResult="{option: user}">
							<User
								:avatar-size="24"
								:show-username="true"
								:user="user"
							/>
						</template>
					</Multiselect>
				</div>
			</div>
			
			<div class="column field">
				<label class="label">{{ $t('task.dueDate') }}</label>
				<div class="control">
					<Datepicker
						v-model="project.endDate"
						:choose-date-label="$t('task.chooseDueDate')"
					/>
				</div>
			</div>
		</div>
	</CreateEdit>
</template>

<script setup lang="ts">
import {ref, reactive, shallowReactive, watch} from 'vue'
import {useI18n} from 'vue-i18n'

import ProjectService from '@/services/project'
import ProjectModel from '@/models/project'
import CreateEdit from '@/components/misc/CreateEdit.vue'
import ColorPicker from '@/components/input/ColorPicker.vue'
import Datepicker from '@/components/input/Datepicker.vue'
import Multiselect from '@/components/input/Multiselect.vue'
import User from '@/components/misc/User.vue'

import {success} from '@/message'
import {useTitle} from '@/composables/useTitle'
import {useProjectStore} from '@/stores/projects'
import ProjectSearch from '@/components/tasks/partials/ProjectSearch.vue'
import type {IProject} from '@/modelTypes/IProject'
import type {IUser} from '@/modelTypes/IUser'
import UserService from '@/services/user'
import {useAuthStore} from '@/stores/auth'

const props = defineProps<{
	parentProjectId?: number,
}>()

const {t} = useI18n({useScope: 'global'})

useTitle(() => t('project.create.header'))

const showError = ref(false)
const project = reactive(new ProjectModel())
const projectService = shallowReactive(new ProjectService())
const projectStore = useProjectStore()
const parentProject = ref<IProject | null>(null)
const isSubmitting = ref(false)

const authStore = useAuthStore()
const userService = new UserService()
const foundUsers = ref<IUser[]>([])

// Init default assignees
if (authStore.info) {
	project.assignees = [authStore.info]
}

async function findUser(query: string) {
	if (!query) return
	foundUsers.value = await userService.getAll(undefined, {s: query}) as IUser[]
}


watch(
	() => props.parentProjectId,
	parentProjectId => {
		if (parentProjectId && projectStore.projects[parentProjectId]) {
			parentProject.value = projectStore.projects[parentProjectId]
		}
	},
	{immediate: true},
)

async function createProject() {
	if (project.title === '') {
		showError.value = true
		return
	}
	showError.value = false

	if (isSubmitting.value) {
		return
	}

	isSubmitting.value = true

	if (parentProject.value) {
		project.parentProjectId = parentProject.value.id
	}

	try {
		await projectStore.createProject(project)
		success({message: t('project.create.createdSuccess')})
	} finally {
		isSubmitting.value = false
	}
}
</script>
