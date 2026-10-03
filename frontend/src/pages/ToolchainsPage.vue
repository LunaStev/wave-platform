<script setup lang="ts">
import { computed, onMounted, ref, watchEffect } from 'vue'
import { Download, RefreshCw } from '@lucide/vue'
import { useI18n } from '../i18n'
import { applyPageSEO } from '../services/seo'
import { formatToolchainSize, parseToolchainCatalog, toolchainChecksumCommand, type ToolchainBundle } from '../services/toolchains'

const { locale } = useI18n()
const copy = computed(() => locale.value === 'ko' ? {
  title: 'LLVM 툴체인', intro: 'Wave를 빌드하는 데 사용하는 검증된 LLVM SDK입니다.',
  compiler: 'Wave 컴파일러 다운로드', scope: '각 SDK는 표시된 Linux 아키텍처에서 실행되며, LLVM 도구·LLD·공유 라이브러리·개발 헤더를 포함합니다. GCC와 시스템 라이브러리는 별도로 준비해야 합니다.',
  loading: '툴체인 목록을 불러오는 중입니다.', empty: '아직 게시된 툴체인이 없습니다. 검증을 마친 파일부터 이곳에 추가됩니다.',
  failed: '툴체인 목록을 불러오지 못했습니다.', retry: '다시 시도', download: '다운로드', checksum: '체크섬 파일',
  published: '게시일', verify: '다운로드 확인', verifyBody: '압축 파일과 체크섬 파일을 같은 디렉터리에 저장한 뒤 실행하세요.',
  versioning: '게시된 파일은 같은 주소에서 교체하지 않습니다. 빌드 구성이 바뀌면 새로운 revision으로 게시합니다.',
  license: 'LLVM 라이선스',
} : {
  title: 'LLVM toolchains', intro: 'Verified LLVM SDKs for building Wave.',
  compiler: 'Download the Wave compiler', scope: 'Each SDK runs on the listed Linux architecture and includes LLVM tools, LLD, shared libraries and development headers. GCC and system libraries are separate prerequisites.',
  loading: 'Loading toolchains…', empty: 'No toolchains have been published yet. Verified downloads will appear here.',
  failed: 'The toolchain catalog could not be loaded.', retry: 'Try again', download: 'Download', checksum: 'Checksum file',
  published: 'Published', verify: 'Verify your download', verifyBody: 'Save the archive and its checksum file in the same directory, then run:',
  versioning: 'Published files keep a fixed URL and contents. Build configuration changes are published as a new revision.',
  license: 'LLVM license',
})
const bundles = ref<ToolchainBundle[]>([])
const state = ref<'loading' | 'ready' | 'error'>('loading')
async function load() {
  state.value = 'loading'
  try {
    const response = await fetch('/downloads/toolchains/index.json', { cache: 'no-cache' })
    if (response.status === 404) bundles.value = []
    else {
      if (!response.ok) throw new Error('Catalog unavailable')
      bundles.value = parseToolchainCatalog(await response.json())
    }
    state.value = 'ready'
  } catch { state.value = 'error' }
}
watchEffect(() => applyPageSEO({ title: copy.value.title + ' · Wave', description: copy.value.intro, locale: locale.value, path: '/toolchains' }))
onMounted(load)
</script>

<template>
  <main class="portal-width toolchains-page">
    <header class="toolchains-header">
      <div><p class="toolchains-eyebrow">Wave · LLVM 21</p><h1>{{ copy.title }}</h1><p>{{ copy.intro }}</p></div>
      <a class="ui-button" href="https://github.com/wavefnd/Wave/releases">{{ copy.compiler }}</a>
    </header>
    <p class="toolchains-scope">{{ copy.scope }}</p>
    <p v-if="state === 'loading'" role="status">{{ copy.loading }}</p>
    <div v-else-if="state === 'error'" role="alert" class="toolchains-notice">
      <p>{{ copy.failed }}</p><button class="ui-button" type="button" @click="load"><RefreshCw :size="16" />{{ copy.retry }}</button>
    </div>
    <p v-else-if="!bundles.length" class="toolchains-notice">{{ copy.empty }}</p>
    <div v-else class="toolchains-grid">
      <article v-for="bundle in bundles" :key="bundle.url" class="toolchain-card">
        <p class="toolchains-eyebrow">LLVM {{ bundle.llvm_version }} · {{ bundle.revision }}</p>
        <h2>{{ bundle.target === 'linux-riscv64' ? 'Linux RISC-V64' : 'Linux LoongArch64' }}</h2>
        <p>{{ formatToolchainSize(bundle.size_bytes) }} · tar.xz</p>
        <p class="toolchain-date">{{ copy.published }} <time :datetime="bundle.published_at">{{ new Date(bundle.published_at).toLocaleDateString(locale) }}</time></p>
        <div class="toolchain-actions">
          <a class="ui-button" :href="bundle.url" download><Download :size="16" />{{ copy.download }}</a>
          <a :href="`${bundle.url}.sha256`" download>{{ copy.checksum }}</a>
        </div>
        <details><summary>SHA-256</summary><code class="toolchain-hash">{{ bundle.sha256 }}</code></details>
        <div class="toolchain-verification">
          <h3>{{ copy.verify }}</h3><p>{{ copy.verifyBody }}</p>
          <pre><code>{{ toolchainChecksumCommand(bundle) }}</code></pre>
          <a :href="`https://github.com/llvm/llvm-project/blob/llvmorg-${bundle.llvm_version}/llvm/LICENSE.TXT`">{{ copy.license }}</a>
        </div>
      </article>
    </div>
    <p class="toolchain-verification">{{ copy.versioning }}</p>
  </main>
</template>

<style scoped>
.toolchains-page { padding-block: 2.5rem 4rem; }
.toolchains-header { display: flex; align-items: center; justify-content: space-between; gap: 1.5rem; flex-wrap: wrap; }
.toolchains-header h1 { margin: .25rem 0 .75rem; font-size: 2rem; }
.toolchains-eyebrow { font-size: .8rem; letter-spacing: .04em; opacity: .7; }
.toolchains-scope { max-width: 72ch; line-height: 1.7; margin-block: 1.5rem 2rem; }
.toolchains-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 300px), 1fr)); gap: 1rem; }
.toolchain-card, .toolchains-notice { border: 1px solid var(--ui-border, #81818155); border-radius: .65rem; padding: 1.5rem; }
.toolchain-card { min-width: 0; }
.toolchain-card h2 { font-size: 1.2rem; }
.toolchain-date { font-size: .85rem; opacity: .7; }
.toolchain-actions { display: flex; align-items: center; gap: 1rem; flex-wrap: wrap; margin-block: 1.5rem; }
.toolchain-hash { display: block; overflow-wrap: anywhere; padding-block: .75rem; font-size: .8rem; }
.toolchain-verification { margin-top: 2.5rem; }
.toolchain-verification h3 { font-size: 1rem; }
.toolchain-verification pre { padding: 1rem; border: 1px solid var(--ui-border, #81818155); border-radius: .4rem; overflow-x: auto; }
.toolchain-card summary { cursor: pointer; }
</style>
