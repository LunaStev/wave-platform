import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createDocumentLoader } from '../src/services/documentLoader.ts'
import type { DocumentLocale, DocumentSummary, DocumentView } from '../src/services/http.ts'

const item = (path = 'language/first', locale: DocumentLocale = 'en') => ({ id: `${locale}/${path}`, path, locale, title: path, summary: '', group: 'language', groupOrder: 1, order: 1 })
const view = (locale: DocumentLocale) => ({ ...item('language/first', locale), markdown: '## Title', blocks: [], translations: [], updatedAt: '' }) satisfies DocumentView
const missing = { status: 404 }
const unavailable = { status: 503 }
const isNotFound = (error: unknown) => error === missing

test('catalogue calls are coalesced across concurrent pages, cached for one minute, and fresh retries bypass the cache', async () => {
  let now = 0
  const calls: string[] = []
  let release: (items: DocumentSummary[]) => void
  const pending = new Promise<DocumentSummary[]>(resolve => { release = resolve })
  const loader = createDocumentLoader({ list: locale => { calls.push(locale); return calls.length <= 2 ? pending : Promise.resolve([item('', locale)]) }, get: async (_path, locale) => view(locale), isNotFound }, () => now)
  const a = loader.load('', 'ko')
  const b = loader.load('language/first', 'ko')
  await Promise.resolve()
  assert.deepEqual(calls, ['ko', 'en'])
  release!([item()])
  await Promise.all([a, b])
  await loader.load('', 'ko')
  assert.equal(calls.length, 2)
  now = 60_000
  await loader.load('', 'ko')
  assert.equal(calls.length, 4)
  await loader.load('', 'ko', true)
  assert.equal(calls.length, 6)
})

test('a failed catalogue is not cached and the next request can recover', async () => {
  let attempts = 0
  const loader = createDocumentLoader({ list: async () => { if (!attempts++) throw unavailable; return [] }, get: async (_path, locale) => view(locale), isNotFound })
  await assert.rejects(loader.load('', 'en'), error => error === unavailable)
  assert.deepEqual((await loader.load('', 'en')).translated, [])
  assert.equal(attempts, 2)
})

test('only a translated 404 triggers English fallback; other failures and missing English remain errors', async () => {
  const calls: string[] = []
  let failure: unknown = missing
  const loader = createDocumentLoader({ list: async () => [], get: async (_path, locale) => { calls.push(locale); if (locale !== 'en') throw failure; return view(locale) }, isNotFound })
  const result = await loader.load('language/first', 'ko')
  assert.equal(result.englishFallback, true)
  assert.equal(result.document?.locale, 'en')
  assert.deepEqual(calls, ['ko', 'en'])
  failure = unavailable
  await assert.rejects(loader.load('language/first', 'ko'), error => error === unavailable)
  assert.deepEqual(calls, ['ko', 'en', 'ko'])
  const noEnglish = createDocumentLoader({ list: async () => [], get: async () => { throw missing }, isNotFound })
  await assert.rejects(noEnglish.load('missing', 'en'), error => error === missing)
  await assert.rejects(noEnglish.load('missing', 'ko'), error => error === missing)
})

test('catalogue cache stays separate by locale and keeps English fallback available', async () => {
  const english = [item(), item('whale/first'), item('stdlib/first')]
  const loader = createDocumentLoader({ list: async locale => locale === 'en' ? english : [item('language/first', locale)], get: async (_path, locale) => view(locale), isNotFound })
  const korean = await loader.load('', 'ko')
  assert.equal(korean.translated[0].locale, 'ko')
  assert.deepEqual(korean.english, english)
  const japanese = await loader.load('', 'ja')
  assert.equal(japanese.translated[0].locale, 'ja')
  assert.equal(japanese.english, korean.english)
  assert.equal(english.length, 3)
})
