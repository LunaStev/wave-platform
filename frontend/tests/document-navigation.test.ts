import assert from 'node:assert/strict'
import test from 'node:test'
import { canonicalDocumentPath, documentationCatalog, documentationPath, documentationProject, projectDocuments } from '../src/services/documentNavigation.ts'

test('URL selects the project and separates catalogue from document paths', () => {
  for (const path of ['', 'language/types', 'toolchain/whale-overview', 'whales/overview']) {
    assert.equal(documentationProject(path), 'wave')
    assert.equal(documentationPath(path), path)
  }
  assert.equal(documentationProject('whale'), 'whale')
  assert.equal(documentationPath('whale'), '')
  assert.equal(documentationProject('whale/overview'), 'whale')
  assert.equal(documentationPath('whale/overview'), 'whale/overview')
  for (const locale of ['en', 'ko', 'ja', 'zh', 'es', 'de', 'ru', 'id', 'vi']) {
    assert.equal(documentationCatalog(locale, 'wave'), `/docs/${locale}`)
    assert.equal(documentationCatalog(locale, 'whale'), `/docs/${locale}/whale`)
  }
})

test('navigation merges translations before sorting within the selected project', () => {
  const english = [
    { path: 'language/types', title: 'Wave', groupOrder: 1, order: 1 },
    { path: 'toolchain/whale-overview', title: 'Legacy Wave path', groupOrder: 4, order: 1 },
    { path: 'whale/symbols', title: 'Symbols', groupOrder: 5, order: 4 },
    { path: 'whale/overview', title: 'Overview', groupOrder: 5, order: 1 },
    { path: 'whale/alignment', title: 'Alignment', groupOrder: 5, order: 4 },
    { path: 'whale/other', title: 'Other group', groupOrder: 6, order: 1 },
  ]
  const translated = [{ ...english[2], title: '심볼' }]
  const whale = projectDocuments(translated, english, 'whale')
  assert.deepEqual(whale.map(item => item.path), ['whale/overview', 'whale/alignment', 'whale/symbols', 'whale/other'])
  assert.equal(whale[2].title, '심볼')
  assert.deepEqual(projectDocuments(translated, english, 'wave').map(item => item.path), ['language/types', 'toolchain/whale-overview'])
  // Search and previous/next consume this same scoped list.
  assert.equal(whale.filter(item => item.title.includes('Wave')).length, 0)
  assert.equal(whale[whale.findIndex(item => item.path === 'whale/symbols') - 1].path, 'whale/alignment')
  assert.deepEqual(projectDocuments([], [], 'whale'), [])
  assert.equal(english[2].title, 'Symbols')
})

test('standard library has an isolated catalogue, preserves old URLs and translation fallback', () => {
  const legacy = ['reference/standard-library', 'reference/string-and-bytes', 'reference/memory-and-buffer', 'reference/system-io-network-process']
  for (const path of [...legacy, 'stdlib/buffer']) {
    assert.equal(documentationProject(path), 'stdlib')
    assert.equal(documentationPath(path), path)
  }
  assert.equal(documentationProject('stdlib'), 'stdlib')
  assert.equal(documentationPath('stdlib'), '')
  assert.equal(documentationProject('stdlib-extra/page'), 'wave')
  for (const locale of ['ko', 'en', 'ja']) assert.equal(documentationCatalog(locale, 'stdlib'), `/docs/${locale}/stdlib`)
  const english = [
    { path: legacy[0], title: 'Library', groupOrder: 1, order: 1 },
    { path: 'stdlib/buffer', title: 'Buffer', groupOrder: 1, order: 2 },
    { path: 'language/types', title: 'Types', groupOrder: 3, order: 1 },
    { path: 'whale/overview', title: 'Whale', groupOrder: 5, order: 1 },
  ]
  const translated = [{ ...english[0], title: '표준 라이브러리' }]
  assert.deepEqual(projectDocuments(translated, english, 'stdlib').map(item => item.title), ['표준 라이브러리', 'Buffer'])
  assert.deepEqual(projectDocuments(translated, english, 'wave').map(item => item.path), ['language/types'])
  assert.deepEqual(projectDocuments([], [], 'stdlib'), [])
})


test('retired lesson and toolchain URLs point at the single canonical collection', () => {
  assert.equal(canonicalDocumentPath('learn/functions'), 'language/functions-and-generics')
  assert.equal(canonicalDocumentPath('learn/strings'), 'language/strings')
  assert.equal(canonicalDocumentPath('toolchain/build-link-targets'), 'whale/build-link-targets')
  assert.equal(canonicalDocumentPath('getting-started/design-goals'), 'getting-started/overview')
  assert.equal(canonicalDocumentPath('language/control-flow'), 'language/control-flow')
  assert.equal(documentationProject(canonicalDocumentPath('toolchain/whale-cli')), 'whale')
})
