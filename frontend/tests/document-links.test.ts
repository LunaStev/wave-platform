import assert from 'node:assert/strict'
import { test } from 'node:test'
import { renderDocumentMarkdown } from '../src/services/documentMarkdown.ts'
import { validateDocumentLinks, type DocumentLinkEntry } from '../src/services/documentLinks.ts'
const locales = new Set(['en', 'ko', 'pt'])
const entry = (locale: string, path: string, markdown: string): DocumentLinkEntry => {
  const result = renderDocumentMarkdown(markdown, { playgrounds: false, slotPrefix: 'test', renderCode: () => '' })
  return { locale, path, anchors: new Set(result.headings.map(item => item.anchor)), links: result.links }
}

test('checks local headings and reference links using rendered multilingual and duplicate IDs, ignoring code and external URLs', () => {
  const source = entry('ko', 'language/first', '## ¿언제 `var`?\n\n## Same\n\n## Same\n\n[local](#언제-var) [second](#same-1) [ref][target]\n\n[target]: /docs/ko/language/second#arrays\n\n```wave\n[not a link](#missing)\n```\n\n[external](https://example.org/docs/en/missing) [raw source](/api/v1/raw?path=README.md)')
  assert.deepEqual(validateDocumentLinks([source, entry('ko', 'language/second', '## Arrays')], locales, {}), [])
})

test('resolves English fallback, moved paths, encoded anchors, and catalogues', () => {
  const source = entry('pt', 'language/first', '[fallback](/docs/pt/language/missing-translation#title) [moved](/docs/en/learn/old#title) [unicode](/docs/ko/language/local#%ED%95%9C%EA%B8%80) [catalogue](/docs/pt/whale)')
  assert.deepEqual(validateDocumentLinks([source, entry('en', 'language/missing-translation', '## Title'), entry('ko', 'language/local', '## 한글')], locales, { 'learn/old': 'language/missing-translation' }), [])
})

test('reports broken paths, fragments, unsupported locales, duplicate paths, and malformed encoding', () => {
  const source = entry('en', 'first', '[path](/docs/en/missing) [anchor](#missing) [locale](/docs/xx) [invalid](#%ZZ) [source](first.md#title) [encoded](/docs/en/first%2Emd#title)')
  const failures = validateDocumentLinks([source, source], locales, {})
  for (const reason of ['duplicate document path', 'document not found', 'heading not found', 'unsupported locale', 'invalid URL', 'Markdown file link']) assert.ok(failures.some(item => item.includes(reason)), reason)
})
