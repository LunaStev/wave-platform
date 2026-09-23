import assert from 'node:assert/strict'
import test from 'node:test'
import { documentationCatalog, documentationPath, documentationProject, projectDocuments } from '../src/services/documentNavigation.ts'

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
