<template>
  <div class="flex flex-col flex-grow">
    <div class="flex w-full min-h-0 flex-grow">
      <BuildProcList v-model:selected-proc-id="selectedProcId" :build="build" />

      <div class="flex flex-grow relative">
        <div v-if="error" class="flex flex-col p-4">
          <span class="text-red-400 font-bold text-xl mb-2">{{ $t('repo.build.execution_error') }}</span>
          <span class="text-red-400">{{ error }}</span>
        </div>

        <div v-else-if="build.status === 'blocked'" class="flex flex-col flex-grow justify-center items-center">
          <Icon name="status-blocked" class="w-32 h-32 text-color" />
          <p class="text-xl text-color">{{ $t('repo.build.protected.awaits') }}</p>
          <div v-if="repoPermissions.push" class="flex mt-2 space-x-4">
            <Button
              color="green"
              :text="$t('repo.build.protected.approve')"
              :is-loading="isApprovingBuild"
              @click="approveBuild"
            />
            <Button
              color="red"
              :text="$t('repo.build.protected.decline')"
              :is-loading="isDecliningBuild"
              @click="declineBuild"
            />
          </div>
        </div>

        <div v-else-if="build.status === 'declined'" class="flex flex-col flex-grow justify-center items-center">
          <Icon name="status-blocked" class="w-32 h-32 text-color" />
          <p class="text-xl text-color">{{ $t('repo.build.protected.declined') }}</p>
        </div>

        <BuildLog
          v-else-if="selectedProcId"
          v-model:proc-id="selectedProcId"
          :build="build"
          class="fixed top-0 left-0 w-full h-full md:absolute"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, inject, PropType, Ref, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';
import BuildLog from '~/components/repo/build/BuildLog.vue';
import BuildProcList from '~/components/repo/build/BuildProcList.vue';
import useApiClient from '~/compositions/useApiClient';
import { useAsyncAction } from '~/compositions/useAsyncAction';
import useNotifications from '~/compositions/useNotifications';
import { Build, BuildProc, Repo, RepoPermissions } from '~/lib/api/types';
import { findProc } from '~/utils/helpers';

export default defineComponent({
  name: 'Build',

  components: {
    Button,
    BuildProcList,
    Icon,
    BuildLog,
  },

  props: {
    procId: {
      type: String as PropType<string | null>,
      default: null,
    },
  },

  setup(props) {
    const apiClient = useApiClient();
    const router = useRouter();
    const route = useRoute();
    const notifications = useNotifications();
    const i18n = useI18n();

    const build = inject<Ref<Build>>('build');
    const repo = inject<Ref<Repo>>('repo');
    const repoPermissions = inject<Ref<RepoPermissions>>('repo-permissions');
    if (!repo || !repoPermissions || !build) {
      throw new Error('Unexpected: "repo", "repoPermissions" & "build" should be provided at this place');
    }

    const procId = toRef(props, 'procId');

    const defaultProcId = computed(() => {
      if (!build.value || !build.value.procs || !build.value.procs[0].children) {
        return null;
      }

      return build.value.procs[0].children[0].pid;
    });

    const selectedProcId = computed({
      get() {
        if (procId.value !== '' && procId.value !== null) {
          const id = parseInt(procId.value, 10);
          const proc = build.value?.procs?.reduce(
            (prev, p) => prev || p.children?.find((c) => c.pid === id),
            undefined as BuildProc | undefined,
          );
          if (proc) {
            return proc.pid;
          }

          // return fallback if proc-id is provided, but proc can not be found
          return defaultProcId.value;
        }

        // is opened on >= md-screen
        if (window.innerWidth > 768) {
          return defaultProcId.value;
        }

        return null;
      },
      set(_selectedProcId: number | null) {
        if (!_selectedProcId) {
          router.replace({ params: { ...route.params, procId: '' } });
          return;
        }

        router.replace({ params: { ...route.params, procId: `${_selectedProcId}` } });
      },
    });

    const selectedProc = computed(() => findProc(build.value.procs || [], selectedProcId.value || -1));
    const error = computed(() => build.value?.error || selectedProc.value?.error);

    const { doSubmit: approveBuild, isLoading: isApprovingBuild } = useAsyncAction(async () => {
      if (!repo) {
        throw new Error('Unexpected: Repo is undefined');
      }

      await apiClient.approveBuild(repo.value.owner, repo.value.name, `${build.value.number}`);
      notifications.notify({ title: i18n.t('repo.build.protected.approve_success'), type: 'success' });
    });

    const { doSubmit: declineBuild, isLoading: isDecliningBuild } = useAsyncAction(async () => {
      if (!repo) {
        throw new Error('Unexpected: Repo is undefined');
      }

      await apiClient.declineBuild(repo.value.owner, repo.value.name, `${build.value.number}`);
      notifications.notify({ title: i18n.t('repo.build.protected.decline_success'), type: 'success' });
    });

    return {
      repoPermissions,
      selectedProcId,
      build,
      error,
      isApprovingBuild,
      isDecliningBuild,
      approveBuild,
      declineBuild,
    };
  },
});
</script>
