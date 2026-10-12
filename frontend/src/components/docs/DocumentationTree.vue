<script setup lang="ts">
import { ChevronRight, Search } from '@lucide/vue'
import { useI18n } from '../../i18n'
import type { DocumentSummary } from '../../services/http'

const props = defineProps<{
  groups: Array<{ id: string; title: string; items: DocumentSummary[] }>
  base: string
  catalog: string
  projectName: string
  currentPath: string
  activeGroup?: string
}>()
const query = defineModel<string>({ required: true })
const { t } = useI18n()
</script>

<template>
  <div class="docs-tree-heading">{{ projectName }} {{ t('docs.title') }}</div>
  <label class="docs-search">
    <Search :size="16" aria-hidden="true" />
    <input v-model="query" type="search" :aria-label="t('docs.search')" :placeholder="t('docs.filter')" />
  </label>
  <nav class="docs-tree-links" :aria-label="t('docs.navigation')">
    <RouterLink class="docs-tree-overview" :to="catalog" :aria-current="!currentPath ? 'page' : undefined">{{ t('docs.overview') }}</RouterLink>
    <details v-for="group in groups" :key="`${group.id}:${query.trim()}`" :open="!!query.trim() || !activeGroup || group.id === activeGroup">
      <summary><ChevronRight :size="14" aria-hidden="true" /><span>{{ group.title }}</span></summary>
      <div class="docs-tree-group">
        <RouterLink v-for="item in group.items" :key="item.path" :to="`${base}/${item.path}`" :aria-current="item.path === currentPath ? 'page' : undefined">{{ item.title }}</RouterLink>
      </div>
    </details>
    <p v-if="query.trim() && !props.groups.length" class="docs-tree-empty" role="status">{{ t('docs.noResults') }}</p>
  </nav>
</template>
